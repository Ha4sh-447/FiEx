package ui

import (
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Ha4sh-447/FiEx/internal"
	"github.com/Ha4sh-447/FiEx/internal/cache"
	"github.com/Ha4sh-447/FiEx/internal/ui/custom"
	pkg "github.com/Ha4sh-447/FiEx/pkg"
	files "github.com/Ha4sh-447/FiEx/pkg/files"
	"github.com/shirou/gopsutil/disk"
)

// FUNCTIONS :
/*

InitScreenData: Creates a new instance of ScreenData type

InitScreen: Populates the screen with initial data

ScreenUpdate: Update the content of the screen to render new files and folders

backButton: Returns a button for the functionality to move up a directory

Volume: Lists the drive's present in system

DriveItem: Render's the drive

FileItem: Render's file based on type

*/

type ScreenData struct {
	files       binding.UntypedList
	fContainer  *fyne.Container
	topbar      *fyne.Container
	widget_Data *widget.Label
	working_dir binding.ExternalString
	searchQuery binding.String
	screen      *fyne.Container
	searchRes   []string
}

// Creates a new instance of ScreenData struct
func InitScreenData(cwd binding.ExternalString) *ScreenData {
	return &ScreenData{
		files:       binding.NewUntypedList(),
		fContainer:  container.NewVBox(),
		topbar:      container.NewHBox(),
		widget_Data: widget.NewLabelWithData(cwd),
		working_dir: cwd,
		searchQuery: binding.NewString(),
		screen:      container.NewVBox(),
	}
}

func InitScreen(w fyne.Window) (fyne.Window, error) {
	fpath := ""
	cwd := binding.BindString(&fpath)
	sd := InitScreenData(cwd)
	cachePath := internal.GetCachePath()
	cache, err := cache.GetCache(cachePath)

	if err != nil {
		slog.Error("Can't find local cache: ", "error", err)
	}

	cwd.AddListener(binding.NewDataListener(func() {
		ScreenUpdate(w, sd, false)
	}))

	wd_widget := widget.NewLabel("Current Directory:")
	wd_widget_data := sd.widget_Data

	back_button := backButton(sd)

	sd.topbar.Add(back_button)
	sd.topbar.Add(wd_widget)
	sd.topbar.Add(wd_widget_data)

	search := widget.NewEntryWithData(sd.searchQuery)
	search.SetPlaceHolder("Search files...")

	// searchSync, _ := debounce.Debounce(func() {
	searchSync := func() {
		query, _ := sd.searchQuery.Get()
		dir := sd.widget_Data.Text

		if query == "" {
			ScreenUpdate(w, sd, false)
			return
		}

		cacheRes := pkg.SearchInCache(query, cache)

		// fmt.Println(cacheRes)

		if cacheRes == nil {

			res := pkg.Search(dir, query)

			if len(res) == 0 {
				sd.fContainer.RemoveAll()
				sd.fContainer.Add(widget.NewLabel("No result found"))
				return
			}

			sd.searchRes = append(sd.searchRes, res...)
			cache.Store[query] = sd.searchRes
			cache.WriteToFile(cachePath)
		} else {
			slog.Info("found in cache")
			sd.searchRes = append(sd.searchRes, cacheRes...)
		}
		ScreenUpdate(w, sd, true)
		sd.searchRes = nil
		// }, 200*time.Millisecond)
	}

	search.OnSubmitted = func(_ string) {
		sd.fContainer.RemoveAll()
		w := widget.NewLabel("Searching...")
		w.Move(fyne.NewPos(400, 300))
		sd.fContainer.Add(w)
		searchSync()
	}
	search.FocusLost()

	sd.screen.Add(sd.topbar)
	sd.screen.Add(search)
	sd.screen.Add(widget.NewSeparator())
	sd.screen.Add(sd.fContainer)

	scrollContainer := container.NewScroll(sd.fContainer)
	scrollContainer.SetMinSize(fyne.NewSize(500, 600))
	scrollContainer.ScrollToTop()

	screen := container.NewBorder(sd.topbar, nil, nil, nil, container.NewVBox(search, scrollContainer))

	w.SetContent(screen)
	return w, nil
}

func backButton(sd *ScreenData) *widget.Button {
	back_button := widget.NewButton("Back", func() {
		wd, err := sd.working_dir.Get()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}

		if wd == "" || wd == "\\" {
			return // if at the root, do nothing
		}

		new_path := filepath.Dir(wd)
		if new_path == wd {
			new_path = ""
		}

		sd.working_dir.Set(new_path)
	})
	return back_button
}

func ScreenUpdate(w fyne.Window, sd *ScreenData, isSearchRes bool) {

	q, _ := sd.searchQuery.Get()
	if isSearchRes && q != "" {
		sd.fContainer.RemoveAll()

		for _, r := range sd.searchRes {

			item := widget.NewButtonWithIcon(r, theme.FileIcon(), func() {
				files.OpenFile(r)
			})
			item.Alignment = widget.ButtonAlignLeading

			sd.fContainer.Add(FileItem(r, w, sd, true))
			// sd.fContainer.Add(item)
		}
		return
	} else {

		// Reload the value from the external string binding
		err := sd.working_dir.Reload()
		if err != nil {
			fyne.LogError("Error reloading working directory", err)
			return
		}

		// Get the current value of the working directory
		fpath, err := sd.working_dir.Get()
		if err != nil {
			fyne.LogError("Error getting working directory value", err)
			return
		}

		if fpath == "" {
			Volume(w, sd)
			return
		}

		sd.files.Set(nil)
		sd.fContainer.RemoveAll()
		sd.searchQuery.Set("")
		// Otherwise, list the contents of the directory
		_, f, err := files.Files(fpath)
		if err != nil {
			fyne.LogError("Can't load files", err)
			return
		}

		for _, file := range f {
			sd.fContainer.Add(FileItem(file.Name(), w, sd, false))
			sd.files.Append(file)
		}
	}

	w.Content().Refresh()
}
func Volume(w fyne.Window, sd *ScreenData) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		slog.Error("Failed to get disk partitions: ", "error", err)
	}

	sd.files.Set(nil)
	sd.fContainer.RemoveAll()

	for _, partition := range partitions {
		mountpoint := partition.Mountpoint
		usage, err := disk.Usage(mountpoint)
		if err != nil {
			slog.Error("Failed to get disk usage", "error", err)
			continue
		}
		// Round of usage and other numbers to 100's
		sd.fContainer.Add(widget.NewSeparator())
		sd.fContainer.Add(DriveItem(mountpoint, usage, w, sd))
	}

	w.Content().Refresh()
}

func DriveItem(drive string, usage *disk.UsageStat, w fyne.Window, sd *ScreenData) *fyne.Container {
	d, _ := strings.CutSuffix(drive, ":")
	vol := widget.NewButtonWithIcon(fmt.Sprintf(" %s ", d), theme.StorageIcon(), func() {
		sd.working_dir.Set(drive + "\\")
	})
	vol.Resize(fyne.NewSize(45, 15))
	info := fmt.Sprintf("%d GB free of %d GB \t\t\t\t\t", uint64(usage.Free/1000000000), uint64(usage.Total/1000000000))
	infoLabel := widget.NewLabel(info)

	progress := widget.NewProgressBar()
	progress.Min = 0
	progress.Max = 100
	progress.SetValue(float64(usage.UsedPercent))

	box := container.NewVBox(
		progress,
		infoLabel,
	)

	driveItem := container.NewHBox(
		vol,
		widget.NewSeparator(),
		box,
	)
	return driveItem
}

func FileItem(file string, w fyne.Window, sd *ScreenData, isSearchRes bool) *fyne.Container {
	dir, err := sd.working_dir.Get()
	if err != nil {
		fyne.LogError("Error getting cwd", err)
	}

	var f fs.FileInfo
	if !isSearchRes {

		var fullPath string
		if strings.HasPrefix(file, dir) {
			fullPath = file
		} else {
			fullPath = filepath.Join(dir, file)
		}
		// fmt.Println("Converted path: ", fullPath)
		f, err = os.Lstat(fullPath)
		if err != nil {
			slog.Warn("Lstat error: ", "error", err)
		}
	} else {
		f, err = os.Lstat(file)
		if err != nil {
			slog.Warn("Lstat error: ", "error", err)
		}
	}
	icon := fileExtType(file, f.IsDir())

	if isSearchRes {

		c := custom.NewCustomButton(icon, f.Name(), file, func() {
			// fpath := filepath.Join(dir, file)
			if f.IsDir() {
				sd.working_dir.Set(file)
			} else {
				files.OpenFile(file)
			}
		})

		fileItem := container.NewHBox(
			c,
		)
		return fileItem
	}
	clickable := widget.NewButtonWithIcon(file, icon, func() {

		fpath := filepath.Join(dir, file)
		if f.IsDir() {
			sd.working_dir.Set(fpath)
		} else {
			files.OpenFile(fpath)
		}
	})
	clickable.Alignment = widget.ButtonAlignLeading

	fileItem := container.NewHBox(
		clickable,
	)
	return fileItem
}

func fileExtType(file string, isDir bool) fyne.Resource {
	var icon fyne.Resource

	if isDir {
		icon = theme.FolderIcon()
		return icon
	}

	switch filepath.Ext(file) {
	case ".txt":
		icon = theme.FileTextIcon()
	case ".png", ".jpg", ".jpeg", ".gif":
		icon = theme.FileImageIcon()
	case ".mp4", ".mkv", ".avi":
		icon = theme.FileVideoIcon()
	case ".mp3", ".wav", ".flac":
		icon = theme.FileAudioIcon()
	case ".pdf":
		icon = theme.FileTextIcon()
	default:
		icon = theme.FileIcon()
	}
	return icon
}

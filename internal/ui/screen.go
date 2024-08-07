package ui

import (
	"fmt"
	"image/color"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
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

	bg := canvas.NewRectangle(&color.RGBA{R: 0, G: 0, B: 0, A: 255})
	bg.FillColor = color.Black
	bg.Resize(fyne.NewSize(1200, 600))

	// Main screen content setup
	wd_widget := widget.NewLabel("Current Directory:")
	wd_widget_data := sd.widget_Data

	back_button := backButton(sd)

	sd.topbar.Add(back_button)
	sd.topbar.Add(wd_widget)
	sd.topbar.Add(wd_widget_data)

	search := widget.NewEntryWithData(sd.searchQuery)
	search.SetPlaceHolder("Search files...")

	// Declare sysCache here so it's available in the searchSync function
	var sysCache *cache.SearchCache

	searchSync := func() {
		query, _ := sd.searchQuery.Get()
		dir := filepath.Clean(sd.widget_Data.Text)

		if query == "" {
			ScreenUpdate(w, sd, false)
			return
		}

		var dirRes []string
		if sysCache != nil {
			dirRes = pkg.SearchInCache(dir, sysCache)
		}

		if len(dirRes) == 0 {
			slog.Info("Not found in cache", "Traversing directory", dir)
			dirRes, err := pkg.TraverseDir(dir)
			if err != nil {
				slog.Error("Can't traverse directory", "failed", err)
			}
			if sysCache != nil {
				sysCache.SyncStore.Store(dir, dirRes)
			}
		}

		res := pkg.Search(dir, query, dirRes)

		if len(res) == 0 {
			sd.fContainer.RemoveAll()
			sd.fContainer.Add(widget.NewLabel("No result found"))
			return
		}

		sd.searchRes = append(sd.searchRes, res...)
		ScreenUpdate(w, sd, true)
		sd.searchRes = nil
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
	content := container.NewStack(bg, screen)

	// Update the window content to the actual application
	w.SetContent(content)

	// Start a goroutine to load or create the cache
	go func() {
		var err error
		sysCache, err = cache.GetCache_msg(internal.GetCachePath())
		if err != nil {
			slog.Error("Can't find local cache: ", "error", err)
			// Show a small box indicating cache creation
			cacheBox := widget.NewLabel("Creating cache...")
			cacheContainer := container.NewStack(bg, container.NewVBox(cacheBox))
			w.SetContent(container.NewStack(content, cacheContainer))

			// Create cache file
			slog.Info("Creating local cache", "CACHE", "")
			sysCache = cache.CreateSysCache()
			slog.Info("Created local cache", "CACHE", "COMPLETED")

			// Remove the cache creation box
			w.SetContent(content)
		}

		slog.Info("INFO", "Loaded cache", len(sysCache.Store))

		cwd.AddListener(binding.NewDataListener(func() {
			ScreenUpdate(w, sd, false)
		}))
	}()

	return w, nil
}
func detailsTile() *fyne.Container {
	// Create header labels with centered text alignment
	nameLabel := widget.NewLabelWithStyle("Name", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	dateModifiedLabel := widget.NewLabelWithStyle("Date modified", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	typeLabel := widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	// sizeLabel := widget.NewLabelWithStyle("Size", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Header container with equal column width
	header := container.NewGridWithColumns(3, // Adjust the number based on the number of labels
		nameLabel,
		dateModifiedLabel,
		typeLabel,
		// sizeLabel,
	)

	// Set the header to fill the entire width of the window
	header.Resize(fyne.NewSize(1200, 100)) // Assuming 1200 is the window width

	// Return the header container
	return header
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

// func ScreenUpdate(w fyne.Window, sd *ScreenData, isSearchRes bool) {

// 	q, _ := sd.searchQuery.Get()
// 	if isSearchRes && q != "" {
// 		sd.fContainer.RemoveAll()

// 		for _, r := range sd.searchRes {

// 			item := widget.NewButtonWithIcon(r, theme.FileIcon(), func() {
// 				files.OpenFile(r)
// 			})
// 			item.Alignment = widget.ButtonAlignLeading

// 			sd.fContainer.Add(FileItem(r, w, sd, true))
// 			// sd.fContainer.Add(item)
// 		}
// 		return
// 	} else {

// 		// Reload the value from the external string binding
// 		err := sd.working_dir.Reload()
// 		if err != nil {
// 			fyne.LogError("Error reloading working directory", err)
// 			return
// 		}

// 		// Get the current value of the working directory
// 		fpath, err := sd.working_dir.Get()
// 		if err != nil {
// 			fyne.LogError("Error getting working directory value", err)
// 			return
// 		}

// 		if fpath == "" {
// 			ShowVolumes(w, sd)
// 			return
// 		}

// 		sd.files.Set(nil)
// 		sd.fContainer.RemoveAll()
// 		sd.searchQuery.Set("")
// 		// Otherwise, list the contents of the directory
// 		_, f, err := files.Files(fpath)
// 		if err != nil {
// 			fyne.LogError("Can't load files", err)
// 			return
// 		}
// 		if len(f) > 0 {
// 			// Add the detailsTile if files are present
// 			sd.screen.Add(detailsTile())
// 			sd.screen.Add(widget.NewSeparator())

// 			// Add files to the screen
// 			for _, file := range f {
// 				sd.fContainer.Add(FileItem(file.Name(), w, sd, false))
// 				sd.files.Append(file)
// 			}
// 		} else {
// 			sd.fContainer.Add(widget.NewLabel("No files found"))
// 		}

// 		// for _, file := range f {
// 		// 	sd.fContainer.Add(FileItem(file.Name(), w, sd, false))
// 		// 	sd.files.Append(file)
// 		// }
// 	}

//		w.Content().Refresh()
//	}
func ScreenUpdate(w fyne.Window, sd *ScreenData, isSearchRes bool) {
	sd.screen.RemoveAll()

	// Top bar with current directory and search bar
	sd.screen.Add(sd.topbar)

	q, _ := sd.searchQuery.Get()
	if isSearchRes && q != "" {
		sd.screen.Add(widget.NewSeparator())
		sd.fContainer.RemoveAll()

		for _, r := range sd.searchRes {
			item := widget.NewButtonWithIcon(r, theme.FileIcon(), func() {
				files.OpenFile(r)
			})
			item.Alignment = widget.ButtonAlignLeading

			sd.fContainer.Add(FileItem(r, w, sd, true))
		}
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
			ShowVolumes(w, sd)
			return
		}

		sd.files.Set(nil)
		sd.fContainer.RemoveAll()
		sd.searchQuery.Set("")

		// Get directory contents
		_, f, err := files.Files(fpath)
		if err != nil {
			fyne.LogError("Can't load files", err)
			dialog.ShowError(err, w)
			// return
		}

		if len(f) > 0 {
			// Add the detailsTile if files are present
			sd.screen.Add(detailsTile())
			sd.screen.Add(widget.NewSeparator())

			// Add files to the screen
			for _, file := range f {
				sd.fContainer.Add(FileItem(file.Name(), w, sd, false))
				sd.files.Append(file)
			}
		} else {
			sd.fContainer.Add(widget.NewLabel("No files found"))
		}
	}

	// Add the file container to the screen
	scrollContainer := container.NewScroll(sd.fContainer)
	scrollContainer.SetMinSize(fyne.NewSize(500, 600))
	scrollContainer.ScrollToTop()

	sd.screen.Add(scrollContainer)
	w.Content().Refresh()
}

func ShowVolumes(w fyne.Window, sd *ScreenData) {
	diskInfoList, err := files.GetDiskUsage()
	if err != nil {
		slog.Error("Failed to get disk partitions: ", "error", err)
		return
	}

	Volume(w, sd, diskInfoList)
}

func Volume(w fyne.Window, sd *ScreenData, diskInfoList []files.DiskInfo) {
	sd.files.Set(nil)
	sd.fContainer.RemoveAll()

	for _, diskInfo := range diskInfoList {
		sd.fContainer.Add(widget.NewSeparator())
		sd.fContainer.Add(DriveItem(diskInfo.Mountpoint, diskInfo.Usage, w, sd))
	}

	w.Content().Refresh()
}

func ShowError(w fyne.Window, err error) {
	dialog.ShowError(err, w)
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
	var fullPath string
	if !isSearchRes {
		if strings.HasPrefix(file, dir) {
			fullPath = file
		} else {
			fullPath = filepath.Join(dir, file)
		}
	} else {
		fullPath = file
	}
	f, err = os.Lstat(fullPath)
	if err != nil {
		slog.Warn("Lstat error: ", "error", err)
	}

	iconResource := fileExtType(file, f.IsDir())
	icon := widget.NewIcon(iconResource)
	icon.Resize(fyne.NewSize(128, 128)) // Increase icon size

	var content fyne.CanvasObject
	if isSearchRes {
		c := custom.NewCustomButton(iconResource, f.Name(), file, func() {
			if f.IsDir() {
				sd.working_dir.Set(file)
			} else {
				files.OpenFile(file)
			}
		})
		content = c
	} else {
		fName := f.Name()
		if len(f.Name()) > 40 {
			fName = strings.Join([]string{fName[:41], "..."}, "")
		}
		fileName := widget.NewLabel(fName)
		fileName.Alignment = fyne.TextAlignLeading
		fileName.TextStyle = fyne.TextStyle{Bold: true}

		fileDateTime := widget.NewLabel(fmt.Sprintf("Date Modified: %s", f.ModTime().Format("02-01-2006 15:04")))
		fileDateTime.Alignment = fyne.TextAlignLeading

		var fileExt string
		if f.IsDir() {
			fileExt = "Folder"
		} else {
			fileExt = strings.TrimPrefix(filepath.Ext(file), ".")
			if fileExt == "" {
				fileExt = "File"
			}
		}
		fileExtLabel := widget.NewLabel(fmt.Sprintf("Type: %s", fileExt))
		fileExtLabel.Alignment = fyne.TextAlignLeading

		newLabel := widget.NewLabel(fmt.Sprintf("%s \nDateModified: %s\t|\tType:%s", file, f.ModTime().Format("02-01-2006 15:04"), fileExt))

		// details := container.NewHBox(
		// 	fileName,
		// 	container.NewHBox(fileDateTime, widget.NewLabel("|"), fileExtLabel),
		// )
		// details.Resize(fyne.NewSize(details.Size().Width, 64)) // Match height with icon

		iconAndDetails := container.NewHBox(icon, newLabel)

		clickable := widget.NewButton("", func() {
			if f.IsDir() {
				sd.working_dir.Set(fullPath)
			} else {
				files.OpenFile(fullPath)
			}
		})
		clickable.Importance = widget.LowImportance

		content = container.NewMax(clickable, iconAndDetails)
	}

	bg := canvas.NewRectangle(&color.RGBA{R: 0, G: 0, B: 0, A: 255})

	fileItem := container.NewStack(bg, content)

	// Add padding around the file item
	paddedFileItem := container.NewPadded(fileItem)

	return paddedFileItem
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

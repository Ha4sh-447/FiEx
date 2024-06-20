package ui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	pkg "github.com/Ha4sh-447/fileExp/pkg/files"
	"github.com/shirou/gopsutil/disk"
)

type ScreenData struct {
	files       binding.UntypedList
	fContainer  *fyne.Container
	widget_Data *widget.Label
	working_dir binding.ExternalString
	searchQuery binding.String
}

// creates new instance of ScreenData struct
func InitScreenData(cwd binding.ExternalString) *ScreenData {
	return &ScreenData{
		binding.NewUntypedList(),
		container.NewVBox(),
		widget.NewLabelWithData(cwd),
		cwd,
		binding.NewString(),
	}
}

func InitScreen(w fyne.Window) (fyne.Window, error) {
	// _, err := user.Current()
	// if err != nil {
	// 	fyne.LogError("Error getting user", err)
	// 	return nil, err
	// }

	fpath := ""

	cwd := binding.BindString(&fpath)
	sd := InitScreenData(cwd)
	cwd.AddListener(binding.NewDataListener(func() {
		ScreenUpdate(w, sd)
		fmt.Println("Screen updater called")
	}))

	wd_widget := widget.NewLabel("Current Directory:")
	wd_widget_data := sd.widget_Data

	back_button := backButton(sd)

	wd_container := container.NewHBox(
		back_button,
		wd_widget,
		wd_widget_data,
	)

	search := widget.NewEntryWithData(sd.searchQuery)
	search.SetPlaceHolder("Search files...")

	sd.searchQuery.AddListener(binding.NewDataListener(func() {
		query, _ := sd.searchQuery.Get()
		fmt.Println(query)
	}))

	ScreenUpdate(w, sd)

	scrollContainer := container.NewScroll(sd.fContainer)
	scrollContainer.SetMinSize(fyne.NewSize(500, 600))
	scrollContainer.ScrollToTop()

	screen := container.NewBorder(wd_container, nil, nil, nil, container.NewVBox(search, scrollContainer))

	w.SetContent(screen)
	return w, nil
}

func backButton(sd *ScreenData) *widget.Button {
	back_button := widget.NewButton("Back", func() {
		wd, err := sd.working_dir.Get()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}
		fmt.Println("Back button pressed")

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

func ScreenUpdate(w fyne.Window, sd *ScreenData) {
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

	// If the working directory is empty, list the drives
	if fpath == "" {
		Volume(w, sd)
		return
	}

	// Otherwise, list the contents of the directory
	_, f, err := pkg.Files(fpath)
	if err != nil {
		fyne.LogError("Can't load files", err)
		return
	}

	// Update the binding list and container
	sd.files.Set(nil)
	sd.fContainer.RemoveAll()

	for _, file := range f {
		sd.files.Append(file)
		sd.fContainer.Add(FileItem(file, w, sd))
	}

	w.Content().Refresh()
}

func Volume(w fyne.Window, sd *ScreenData) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Fatalf("Failed to get disk partitions: %v", err)
	}

	sd.files.Set(nil)
	sd.fContainer.RemoveAll()

	seen := make(map[string]bool)
	for _, partition := range partitions {
		mountpoint := partition.Mountpoint
		if !seen[mountpoint] {
			seen[mountpoint] = true
			fmt.Printf("Listing files in %s:\n", mountpoint)
			sd.fContainer.Add(DriveItem(mountpoint, w, sd))
		}
	}

	w.Content().Refresh()
}

func DriveItem(drive string, w fyne.Window, sd *ScreenData) *fyne.Container {
	clickable := widget.NewButton(drive, func() {
		sd.working_dir.Set(drive + "\\")
	})

	driveItem := container.NewHBox(
		clickable,
	)
	return driveItem
}

func FileItem(file os.DirEntry, w fyne.Window, sd *ScreenData) *fyne.Container {
	var icon fyne.Resource

	if file.IsDir() {
		icon = theme.FolderIcon()
	} else {
		switch filepath.Ext(file.Name()) {
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
	}

	clickable := widget.NewButtonWithIcon(file.Name(), icon, func() {
		dir, err := sd.working_dir.Get()
		if err != nil {
			fyne.LogError("Error getting cwd", err)
		}

		fpath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			// Change directory path
			sd.working_dir.Set(fpath)
		} else {
			pkg.OpenFile(fpath)
		}
	})
	clickable.Alignment = widget.ButtonAlignLeading // Align text to the left

	fileItem := container.NewHBox(
		clickable,
	)
	return fileItem
}

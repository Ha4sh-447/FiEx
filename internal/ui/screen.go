package ui

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	pkg "github.com/Ha4sh-447/fileExp/pkg/files"
)

type ScreenData struct {
	files       binding.UntypedList
	fContainer  *fyne.Container
	widget_Data *widget.Label
	working_dir binding.ExternalString
}

// creates new instance of ScreenData struct
func InitScreenData(cwd binding.ExternalString) *ScreenData {
	return &ScreenData{
		binding.NewUntypedList(),
		container.NewVBox(),
		widget.NewLabelWithData(cwd),
		cwd,
	}
}

func InitScreen(w fyne.Window) (fyne.Window, error) {
	// a := app.New()
	// w := a.NewWindow("File Explorer")
	usr, err := user.Current()
	cwd := binding.BindString(&usr.HomeDir)
	sd := InitScreenData(cwd)

	if err != nil {
		fyne.LogError("Error getting user", err)
		return nil, err
	}

	wd_widget := widget.NewLabel("Current Directory:")
	wd_widget_data := sd.widget_Data

	wd_container := container.NewHBox(
		wd_widget,
		wd_widget_data,
	)

	wd, err := cwd.Get()
	if err != nil {
		fyne.LogError("Couldn't get current directory", err)
		return nil, err
	}

	// populate fContainer
	updateScreenData(wd, sd)

	scrollContainer := container.NewScroll(sd.fContainer)
	scrollContainer.ScrollToTop()

	screen := container.NewBorder(wd_container, nil, nil, nil, scrollContainer)

	w.SetContent(screen)
	return w, nil
}

func updateScreenData(fpath string, sd *ScreenData) {
	// fmt.Println(fpath)
	_, f, err := pkg.Files(fpath)
	sd.working_dir.Set(fpath)

	if err != nil {
		fyne.LogError("Can't load file", err)
	}

	sd.files.Set(nil)
	sd.fContainer.RemoveAll()

	for _, file := range f {
		sd.files.Append(file)
	}

	for _, file := range f {
		sd.fContainer.Add(FileItem(file, sd.working_dir, sd))
	}

	scrollContainer := container.NewScroll(sd.fContainer)
	scrollContainer.Refresh()
	// scrollContainer.ScrollToTop()
}

func FileItem(file os.DirEntry, cwd binding.ExternalString, sd *ScreenData) *fyne.Container {
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

	// label := widget.NewLabel(file.Name())

	clickable := widget.NewButtonWithIcon(file.Name(), icon, func() {
		dir, err := cwd.Get()
		if err != nil {
			fyne.LogError("Error getting cwd", err)
		}

		fpath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			// Change directory path
			cwd.Set(fpath)
			updateScreenData(fpath, sd)
			sd.fContainer.Refresh()
		} else {
			pkg.OpenFile(fpath)
		}
	})
	clickable.Alignment = widget.ButtonAlignLeading // Align text to the left

	fileItem := container.NewHBox(
		clickable,
	)
	fmt.Println(sd.files.Length())

	return fileItem
}

package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Ha4sh-447/fileExp/internal"
)

var files binding.UntypedList = binding.NewUntypedList()

func updateScreenData(files binding.UntypedList, fpath string) {
	_, f, err := internal.Files(fpath)
	if err != nil {
		fyne.LogError("Can't load file", err)
	}

	for _, file := range f {
		files.Append(file)
	}
	fmt.Println("Files appended")
}

func main() {
	a := app.New()
	w := a.NewWindow("File Explorer")
	usr, err := user.Current()

	if err != nil {
		fyne.LogError("Error getting user", err)
		return
	}

	cwd := binding.BindString(&usr.HomeDir)

	wd, err := cwd.Get()
	if err != nil {
		fyne.LogError("Couldn't get current directory", err)
		return
	}

	_, f, err := internal.Files(wd)

	for _, file := range f {
		files.Append(file)
	}
	fmt.Println(files.Length())

	wd_widget := widget.NewLabel("Current Directory:")
	wd_container := container.NewHBox(
		wd_widget,
		widget.NewLabelWithData(cwd),
	)

	if err != nil {
		fyne.LogError("Failed to load directory files", err)
		return
	}

	fContainer := container.NewVBox()
	for _, file := range f {
		fContainer.Add(fileItem(file, cwd))
	}

	scrollContainer := container.NewScroll(fContainer)
	screen := container.NewBorder(wd_container, nil, nil, nil, scrollContainer)

	w.SetContent(screen)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

func updateFileContainer(f []fs.DirEntry, cwd binding.ExternalString) *fyne.Container {
	cont := container.NewVBox()
	for _, file := range f {
		cont.Add(fileItem(file, cwd))
	}

	return cont
}

func fileItem(file os.DirEntry, cwd binding.ExternalString) *fyne.Container {
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
			updateScreenData(files, fpath)
		} else {
			openFile(fpath)
		}
	})
	clickable.Alignment = widget.ButtonAlignLeading // Align text to the left

	fileItem := container.NewHBox(
		clickable,
	)
	fmt.Println(files.Length())

	return fileItem
}

func openFile(path string) {
	switch filepath.Ext(path) {
	case ".txt":
		exec.Command("notepad", path).Start()
	default:
		// Handle other file types or open with the default application
		exec.Command("cmd", "/C", "start", path).Start()
	}
}

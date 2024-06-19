package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	pkg "github.com/Ha4sh-447/fileExp/pkg/files"
)

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

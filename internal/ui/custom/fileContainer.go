package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FileContainer struct {
	widget.BaseWidget
	Content fyne.CanvasObject
	onClick func()
}

func NewFileContainer(content fyne.CanvasObject, onClick func()) *FileContainer {
	fCont := &FileContainer{
		Content: content,
		onClick: onClick,
	}
	fCont.ExtendBaseWidget(fCont)
	return fCont
}

func (fc *FileContainer) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewBorder(nil, nil, nil, nil, fc.Content)
	return widget.NewSimpleRenderer(c)
}

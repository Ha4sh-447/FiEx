package custom

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CustomButton struct {
	widget.BaseWidget
	Icon    fyne.Resource
	Name    string
	Path    string
	onClick func()
}

func NewCustomButton(icon fyne.Resource, name string, path string, onClick func()) *CustomButton {
	cb := &CustomButton{
		Icon:    icon,
		Name:    name,
		Path:    path,
		onClick: onClick,
	}
	cb.ExtendBaseWidget(cb)
	return cb
}

func (cb *CustomButton) CreateRenderer() fyne.WidgetRenderer {
	icon := canvas.NewImageFromResource(cb.Icon)
	icon.SetMinSize(fyne.NewSize(32, 32))
	name := widget.NewLabel(cb.Name)
	name.TextStyle = fyne.TextStyle{Bold: true}
	path := widget.NewLabel(cb.Path)

	border := canvas.NewRectangle(color.Black)
	border.StrokeColor = color.White
	border.CornerRadius = 1
	border.StrokeWidth = 0.2

	objects := []fyne.CanvasObject{border, icon, name, path}
	return &CustomButtonRenderer{
		customButton: cb,
		border:       border,
		icon:         icon,
		name:         name,
		path:         path,
		objects:      objects,
	}
}

func (cb *CustomButton) Tapped(event *fyne.PointEvent) {
	if cb.onClick != nil {
		cb.onClick()
	}
}

func (cb *CustomButton) TappedSecondary(event *fyne.PointEvent) {
	// Handle right-click if needed
}

type CustomButtonRenderer struct {
	customButton *CustomButton
	border       *canvas.Rectangle
	icon         *canvas.Image
	name         *widget.Label
	path         *widget.Label
	objects      []fyne.CanvasObject
}

func (r *CustomButtonRenderer) Layout(size fyne.Size) {
	padding := theme.Padding()
	iconSize := r.icon.MinSize()
	r.icon.Resize(iconSize)
	r.icon.Move(fyne.NewPos(padding*2, (size.Height-iconSize.Height)/2))

	nameSize := r.name.MinSize()
	pathSize := r.path.MinSize()

	totalHeight := nameSize.Height + pathSize.Height + padding
	startY := (size.Height - totalHeight) / 2

	r.name.Resize(nameSize)
	r.name.Move(fyne.NewPos(iconSize.Width+2*padding, startY))

	r.path.Resize(pathSize)
	r.path.Move(fyne.NewPos(iconSize.Width+2*padding, startY+nameSize.Height+padding))

	// Adjust the border width to be as long as the path
	borderWidth := iconSize.Width + pathSize.Width + 2*padding
	r.border.Resize(fyne.NewSize(borderWidth, size.Height))
	r.border.Move(fyne.NewPos(padding, 0))
}

func (r *CustomButtonRenderer) MinSize() fyne.Size {
	iconSize := r.icon.MinSize()
	nameSize := r.name.MinSize()
	pathSize := r.path.MinSize()
	width := iconSize.Width + pathSize.Width + 3*theme.Padding()
	height := fyne.Max(iconSize.Height, nameSize.Height+pathSize.Height+theme.Padding())
	return fyne.NewSize(width, height)
}

func (r *CustomButtonRenderer) Refresh() {
	r.border.Refresh()
	r.icon.Refresh()
	r.name.Refresh()
	r.path.Refresh()
}

func (r *CustomButtonRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *CustomButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *CustomButtonRenderer) Destroy() {}

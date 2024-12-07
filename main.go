package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	EditWidget *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile fyne.URI
	SaveMenuItem *fyne.MenuItem
}

var cfg Config

func (c *Config) makeUI() (*widget.Entry, *widget.RichText) {
	var entry = widget.NewMultiLineEntry()
	var preview = widget.NewRichTextFromMarkdown("")

	c.EditWidget = entry
	c.PreviewWidget = preview

	entry.OnChanged = preview.ParseMarkdown

	return entry, preview
}

func main() {
	var _app = app.New()
	var win = _app.NewWindow("Markdown Editor")

	var entry, preview = cfg.makeUI()

	win.SetContent(container.NewHSplit(entry, preview))
	win.Resize(fyne.Size{Width: 600, Height: 500})
	win.CenterOnScreen()
	win.ShowAndRun()
}

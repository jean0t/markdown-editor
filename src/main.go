package main

import (
	"io"
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	EditWidget *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile fyne.URI
	SaveMenuItem *fyne.MenuItem
}

// allows only the markdown files to be showed
// really important!!!
var filter storage.FileFilter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

var cfg Config

func (c *Config) makeUI() (*widget.Entry, *widget.RichText) {
	var entry = widget.NewMultiLineEntry()
	var preview = widget.NewRichTextFromMarkdown("")

	c.EditWidget = entry
	c.PreviewWidget = preview

	return entry, preview
}

func (c *Config) makeMenu(w fyne.Window) {

	var openFile = fyne.NewMenuItem("Open File", c.openFile(w))
	var saveFile = fyne.NewMenuItem("Save", c.save(w))
	c.SaveMenuItem = saveFile
	c.SaveMenuItem.Disabled = true
	var saveAsFile = fyne.NewMenuItem("Save as...", c.saveAs(w))

	var fileMenu = fyne.NewMenu("File", openFile, saveFile, saveAsFile)
	var mainMenu = fyne.NewMainMenu(fileMenu)

	w.SetMainMenu(mainMenu)
}

func (c *Config) openFile(w fyne.Window) func() {
	return func() {
		var openDialog = dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			if read == nil {
				return
			}
			defer read.Close()

			contents, err := io.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			c.EditWidget.Text = string(contents)
			c.CurrentFile = read.URI()
			w.SetTitle("Markdown - " + read.URI().Name())

			c.SaveMenuItem.Disabled = false
		}, w)

		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

func (c *Config) saveAs(w fyne.Window) func() {
	return func() {
		var saveDialog = dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			if write == nil { //action cancelled
				return
			}
			defer write.Close()

			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Use a .md extension", w)
				return
			}

			write.Write([]byte(c.EditWidget.Text))
			c.CurrentFile = write.URI()

			w.SetTitle("Markdown - " + write.URI().Name())
			c.SaveMenuItem.Disabled = false
		}, w)

		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}

func (c *Config) save(w fyne.Window) func() {
	return func() {
		if c.CurrentFile != nil {
			var write, err = storage.Writer(c.CurrentFile)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			defer write.Close()

			write.Write([]byte(c.EditWidget.Text))
			w.SetTitle("Markdown - " + write.URI().Name())
			dialog.ShowInformation("Success", "File was saved", w)
		}
	}
}

func main() {
	var _app = app.New()
	var win = _app.NewWindow("Markdown Editor")

	var entry, preview = cfg.makeUI()
	cfg.makeMenu(win)

	entry.OnChanged = func(content string) {
		preview.ParseMarkdown(content)
		if cfg.CurrentFile == nil {
			return
		}
		if !strings.HasSuffix(win.Title(), "*") {
			win.SetTitle(win.Title() + "*")
		}
	}


	win.SetContent(container.NewHSplit(entry, preview))
	win.Resize(fyne.Size{Width: 600, Height: 500})
	win.CenterOnScreen()
	win.ShowAndRun()
}

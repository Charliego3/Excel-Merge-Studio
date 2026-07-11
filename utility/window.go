package utility

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func NewDialog(options application.WebviewWindowOptions) *application.WebviewWindow {
	state := State()
	dialog := state.App.Window.NewWithOptions(options)
	state.MainWindow.AttachModal(dialog)
	state.MainWindow.SetEnabled(false)
	dialog.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		state.MainWindow.SetEnabled(true)
		state.MainWindow.Focus()
	})
	return dialog
}

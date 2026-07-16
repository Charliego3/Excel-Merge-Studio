package utility

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func NewDialog(options application.WebviewWindowOptions) *application.WebviewWindow {
	dialog := State.App.Window.NewWithOptions(options)
	State.MainWindow.AttachModal(dialog)
	State.MainWindow.SetEnabled(false)
	dialog.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		State.MainWindow.SetEnabled(true)
		State.MainWindow.Focus()
	})
	return dialog
}

func ShowWarning(message string) {
	State.App.Dialog.Warning().
		AttachToWindow(State.MainWindow).
		SetTitle("🚧 提醒").
		SetMessage(message).
		Show()
}

func Confirm(title, message string, callback func(bool)) {
	dialog := State.App.Dialog.Question().
		SetTitle(title).
		AttachToWindow(State.MainWindow).
		SetMessage(message)

	cancelBtn := dialog.AddButton("取消")
	cancelBtn.OnClick(func() {
		callback(false)
	})
	overrideBtn := dialog.AddButton("确定")
	overrideBtn.OnClick(func() {
		callback(true)
	})
	dialog.SetDefaultButton(overrideBtn)
	dialog.SetCancelButton(cancelBtn)
	dialog.Show()
}

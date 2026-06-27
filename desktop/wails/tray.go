package main

import (
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func configureWindowHooks(app *application.App, window *application.WebviewWindow) {
	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window.Hide()
		e.Cancel()
	})

	if runtime.GOOS == "darwin" {
		app.Event.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(event *application.ApplicationEvent) {
			window.Show()
		})
	}
}

func configureTray(app *application.App, desktop *DesktopApp, window *application.WebviewWindow) {
	tray := app.SystemTray.New()
	if runtime.GOOS == "darwin" {
		tray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	menu := app.NewMenu()
	menu.Add("Open Settings").OnClick(func(ctx *application.Context) {
		window.Show()
	})
	menu.Add("Connect Now").OnClick(func(ctx *application.Context) {
		_, _ = desktop.ConnectNow()
	})
	menu.Add("Disconnect").OnClick(func(ctx *application.Context) {
		_, _ = desktop.DisconnectNow()
	})
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	tray.SetMenu(menu)
	tray.OnClick(func() {
		window.Show()
	})
}

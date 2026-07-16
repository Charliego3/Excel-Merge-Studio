package main

import (
	"embed"
	"fmt"
	"log"

	"merger/services"
	"merger/utility"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	utility.State.App = application.New(application.Options{
		Name:        "Excel Merge Studio",
		Description: "Merge Excel Workbooks by a selected key column.",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Services: []application.Service{
			application.NewService(&services.Reader{}),
			application.NewService(&services.Workbook{}),
			application.NewService(&services.Setting{}),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		OnShutdown: func() {
			/// Store window bounds before shutdown
			var bounds = utility.State.MainWindow.Bounds()
			fmt.Printf("%+v\n", bounds)
		},
	})

	utility.State.MainWindow = utility.State.App.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Excel Merge Studio",
		URL:   "/",

		/// TODO: Restore window bounds on startup
		X:      0,
		Y:      400,
		Width:  900,
		Height: 718,

		InitialPosition: application.WindowXY,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,

			Appearance: application.NSAppearanceNameAqua,
			TitleBar:   application.MacTitleBarHiddenInset,
		},
	})

	if err := utility.State.App.Run(); err != nil {
		log.Fatal(err)
	}
}

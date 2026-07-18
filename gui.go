package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetsaver"

	"dial/internal/gui"
	"dial/internal/store"
)

//go:embed frontend/dist
var assets embed.FS

func launchGUI() error {
	dbPath, err := store.DefaultPath()
	if err != nil {
		return err
	}
	db, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	app := gui.NewApp(db)

	return wails.Run(&options.App{
		Title: "Dial",
		Width: 480,
		Height: 640,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Bind: []interface{}{app},
	})
}

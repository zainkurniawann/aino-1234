package main

import (
	"document/database"
	"github.com/fatih/color"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"os"
)

func main() {
	color.Yellow("Connecting to Database...")

	db := database.Connection()
	if db == nil {
		color.Red("Failed to connect to PostgreSQL")
		os.Exit(1)
	}
	defer db.Close()

	sqlDB := db.DB

	// FIX: Deklarasikan err dengan :=
	err := os.Chdir("database")
	if err != nil {
		color.Red("Failed to change directory to databases:", err)
		os.Exit(1)
	}

	color.Yellow("Running migrations...")

	// Get param from command line
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "down":
			color.Yellow("Rolling back migrations...")
			err = goose.Down(sqlDB, "./migration")
		case "redo":
			color.Yellow("Redoing migrations...")
			err = goose.Redo(sqlDB, "./migration")
		case "status":
			color.Yellow("Checking migration status...")
			err = goose.Status(sqlDB, "./migration")
		case "downall":
			color.Yellow("Rolling back all migrations...")
			err = goose.DownTo(sqlDB, "./migration", 0)
		case "version":
			color.Yellow("Checking migration version...")
			err = goose.Version(sqlDB, "./migration")
		default:
			color.Red("Unknown command:", os.Args[1])
			os.Exit(1)
		}
	} else {
		color.Yellow("Running migrations...")
		err = goose.Up(sqlDB, "./migration")
	}

	if err != nil {
		color.Red("Migration failed:", err)
		os.Exit(1)
	}

	color.Green("Migration successful")
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charlieroth/godo/internal/cli"
	"github.com/charlieroth/godo/internal/store"
)

func main() {
	jsonStore := store.NewJsonStore("db.json")
	if err := jsonStore.Load(); err != nil {
		fmt.Println("Error loading database:", err)
		os.Exit(1)
	}

	app := cli.NewApp(jsonStore)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

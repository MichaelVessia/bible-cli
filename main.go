package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bible-cli",
		Usage: "A CLI for browsing the Bible",
		Action: func(*cli.Context) error {
			fmt.Println("Let there be light.")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

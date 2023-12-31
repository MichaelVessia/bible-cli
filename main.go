package main

import (
	"bible-cli/database"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", "./bibleReadings.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &cli.App{
		Name:  "bible-cli",
		Usage: "A CLI for browsing the Bible",

		// Define commands
		Commands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Usage:   "Fetch readings for a date. Supports full dates like 2023-11-28 or partial dates like 2023-11",
				Action: func(c *cli.Context) error {
					return fetchReadings(c, db)
				},
			},
			{
				Name:    "seed-db",
				Aliases: []string{"sd"},
				Usage:   "Seed the database with readings from CSV files in the lectionary directory",
				Action:  seedDatabase,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func fetchReadings(c *cli.Context, db *sql.DB) error {
	// Get the date argument
	date := c.Args().First()
	if date == "" {
		fmt.Println("Please provide a date. Usage: bible-cli fetch <date>")
		return nil
	}

	// Call the function to fetch readings from the database
	readingsList, err := database.FetchReadingsForDate(db, date)
	if err != nil {
		log.Fatal(err)
	}

	if len(readingsList) == 0 {
		fmt.Println("No readings found for", date)
		return nil
	}

	for _, readings := range readingsList {
		fmt.Printf("Readings for %s (%s):\n", readings.CalendarDate, readings.LiturgicalDate)
		fmt.Println("First Reading:", readings.FirstReading)
		fmt.Println("Psalm Reading:", readings.PsalmReading)
		fmt.Println("Second Reading:", readings.SecondReading)
		fmt.Println("Gospel Reading:", readings.GospelReading)
		fmt.Println()
	}

	return nil
}

func seedDatabase(c *cli.Context) error {
	rowsInserted, err := database.SeedDb()
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Database seeded without error, added %v rows.\n", rowsInserted)
	return nil
}

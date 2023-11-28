package database

import (
	"database/sql"
	"encoding/csv"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
)

func SeedDb() error {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", "./bibleReadings.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Setup the database schema
	if err := CreateDb(db); err != nil {
		log.Fatal(err)
	}

	// Open the CSV file
	csvFile, err := os.Open("./Year_C_2021-2022.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	// Create a new CSV reader from the file
	reader := csv.NewReader(csvFile)

	// Read the CSV data
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		// Insert the data into the SQLite database
		// Adjust the SQL statement according to your table's schema
		_, err = db.Exec("INSERT INTO readings (liturgical_date, calendar_date, first_reading, psalm_reading, second_reading, gospel_reading) VALUES (?, ?, ?, ?, ?, ?)", row[0], row[1], row[2], row[3], row[4], row[5])
		if err != nil {
			log.Fatal(err)
		}
	}
	return err
}

func CreateDb(db *sql.DB) error {
	// SQL statement to create the table
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS readings (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        liturgical_date TEXT,
        calendar_date TEXT,
        first_reading TEXT,
        psalm_reading TEXT,
        second_reading TEXT,
        gospel_reading TEXT
    );`

	_, err := db.Exec(createTableSQL)
	return err

}

type Mass struct {
	LiturgicalDate string
	CalendarDate   string
	FirstReading   string
	PsalmReading   string
	SecondReading  string
	GospelReading  string
}

func FetchReadingsForDate(db *sql.DB, date string) ([]Mass, error) {
	var readingsList []Mass

	// Query the database
	rows, err := db.Query("SELECT liturgical_date, calendar_date, first_reading, psalm_reading, second_reading, gospel_reading FROM readings WHERE calendar_date = ?", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var masses Mass
		if err := rows.Scan(&masses.LiturgicalDate, &masses.CalendarDate, &masses.FirstReading, &masses.PsalmReading, &masses.SecondReading, &masses.GospelReading); err != nil {
			return nil, err
		}
		readingsList = append(readingsList, masses)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return readingsList, nil
}

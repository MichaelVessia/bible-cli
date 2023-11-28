package database

import (
	"database/sql"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func SeedDb() (int, error) {
	db, err := sql.Open("sqlite3", "./bibleReadings.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	if err := CreateDb(db); err != nil {
		return 0, err
	}

	directory := "./lectionary/"

	// Read all files in the directory
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return 0, err
	}

	rowsInserted := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".csv") {
			csvFilePath := filepath.Join(directory, file.Name())
			seedCount, err := seedFromFile(db, csvFilePath)
			if err != nil {
				log.Printf("Error seeding from file %s: %v", csvFilePath, err)
			}
			rowsInserted += seedCount
		}
	}

	return rowsInserted, nil
}

func seedFromFile(db *sql.DB, csvFilePath string) (int, error) {
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return 0, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	// Toss the first row since it's header info
	reader.Read()
	count := 0
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		_, err = db.Exec("INSERT INTO readings (liturgical_date, calendar_date, first_reading, psalm_reading, second_reading, gospel_reading) VALUES (?, ?, ?, ?, ?, ?)", row[0], row[1], row[2], row[3], row[4], row[5])
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
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
	rows, err := db.Query("SELECT liturgical_date, calendar_date, first_reading, psalm_reading, second_reading, gospel_reading FROM readings WHERE calendar_date LIKE ?", date+"%")
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

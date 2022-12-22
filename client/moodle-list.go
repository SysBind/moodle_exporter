package client

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

// List of Moodle instances
type MoodleList struct {
	poolConfig *pgxpool.Config
	moodles    []Moodle
}

func NewMoodleList(hostname string, username string, password string) (moodles MoodleList, err error) {
	connconf := pgx.ConnConfig{Host: hostname, User: username, Password: password}
	conn, err := pgx.Connect(connconf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close()

	var rows *pgx.Rows
	if rows, err = conn.Query("SELECT datname FROM pg_database"); err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var datname string
		if err = rows.Scan(&datname); err != nil {
			fmt.Fprintf(os.Stderr, " failed: %v\n", err)
			return
		}
		fmt.Printf("Checking database %s for moodle tables", datname)
	}

	return
}

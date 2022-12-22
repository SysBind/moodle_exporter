package client

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

// List of Moodle instances
type MoodleList struct {
	Moodles []*Moodle
}

func NewMoodleList(hostname string, username string, password string) (list MoodleList, err error) {
	connconf := pgx.ConnConfig{Host: hostname, User: username, Password: password}
	conn, err := pgx.Connect(connconf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close()

	list = MoodleList{Moodles: []*Moodle{}}

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
		if datname == "postgres" || datname == "template0" || datname == "template1" {
			continue
		}

		var moodle *Moodle
		if moodle, err = NewMoodle(hostname, username, password, datname); err != nil {
			fmt.Printf("Skipping database %s, contains no moodle tables\n", datname)
			err = nil
			continue
		}
		fmt.Printf("Moodle list: adding %s\n", moodle)
		list.Moodles = append(list.Moodles, moodle)
	}

	return
}

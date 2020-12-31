package client

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

type Moodle struct {
	Connection *pgx.Conn
	Shortname  string
}

func New(hostname string, username string, database string) (moodle *Moodle, err error) {
	moodle = nil
	conn, err := pgx.Connect(
		pgx.ConnConfig{Host: hostname, User: username, Database: database})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close()

	var shortname string
	err = conn.QueryRow("SELECT shortname FROM mdl_course WHERE id=1").Scan(&shortname)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return
	}
	moodle = &Moodle{Connection: conn, Shortname: shortname}

	return
}

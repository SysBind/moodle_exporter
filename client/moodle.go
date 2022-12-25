package client

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Moodle struct represents single moodle instance
type Moodle struct {
	poolconfig *pgxpool.Config
	shortname  string
}

func (m Moodle) String() string {
	return fmt.Sprintf("Moodle Instance: %s", m.shortname)
}

func NewMoodle(hostname string, username string, password string, database string) (moodle *Moodle, err error) {
	moodle = nil
	var poolconf *pgxpool.Config
	connstr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		username, password, hostname, database)
	poolconf, err = pgxpool.ParseConfig(connstr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed configuring connection pool: %v\n", err)
		return
	}

	ctx := context.Background()
	var conn *pgxpool.Pool
	conn, err = pgxpool.NewWithConfig(ctx, poolconf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close()

	var shortname string
	err = conn.QueryRow(ctx, "SELECT shortname FROM mdl_course WHERE id=1").Scan(&shortname)
	if err != nil {
		fmt.Printf("QueryRow failed: %v\n", err)
		return
	}
	moodle = &Moodle{poolconfig: poolconf, shortname: shortname}

	return
}

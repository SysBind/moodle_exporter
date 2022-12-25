package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// List of Moodle instances
type MoodleList struct {
	moodles []*Moodle
}

func (list MoodleList) String() (str string) {
	str = "Moodles:"
	for i, moodle := range list.moodles {
		str = str + fmt.Sprintf(" -- %d: %s", i, moodle)
	}
	return
}

func NewMoodleList(hostname string, username string, password string) (list MoodleList, err error) {
	var poolconf *pgxpool.Config
	connstr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		username, password, hostname, "postgres")
	poolconf, err = pgxpool.ParseConfig(connstr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed configuring connection pool: %v\n", err)
		return
	}

	ctx := context.Background()
	attempt := 0
	var conn *pgxpool.Pool
	for {
		conn, err = pgxpool.NewWithConfig(ctx, poolconf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			if attempt < 5 {
				attempt++
				fmt.Printf("Sleeping %d seconds before retry..\n", attempt*2)
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return
		}
		break
	}
	defer conn.Close()

	list = MoodleList{moodles: []*Moodle{}}

	var rows pgx.Rows
	if rows, err = conn.Query(ctx, "SELECT datname FROM pg_database"); err != nil {
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
		fmt.Printf("Moodle list: checking datname %s\n", datname)
		if moodle, err = NewMoodle(hostname, username, password, datname); err != nil {
			fmt.Printf("Skipping database %s, contains no moodle tables\n", datname)
			err = nil
			continue
		}
		fmt.Printf("Moodle list: adding %s\n", moodle)
		list.moodles = append(list.moodles, moodle)
	}

	return
}

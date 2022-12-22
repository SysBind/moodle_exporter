package client

import "github.com/jackc/pgx/v4/pgxpool"

// List of Moodle instances
type MoodleList struct {
	poolConfig *pgxpool.Config
	moodles    []Moodle
}

package client

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx"
)

type Moodle struct {
	Connection *pgx.Conn // defer Connection.close()
	Shortname  string
}

type UserStats struct {
	LiveUsers                        int
	ExpectedUpcomingExamParticipants int
}

func New(hostname string, username string, password string, database string) (moodle *Moodle, err error) {
	moodle = nil
	conn, err := pgx.Connect(
		pgx.ConnConfig{Host: hostname, User: username, Password: password, Database: database})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}

	var shortname string
	err = conn.QueryRow("SELECT shortname FROM mdl_course WHERE id=1").Scan(&shortname)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return
	}
	moodle = &Moodle{Connection: conn, Shortname: shortname}

	return
}

func (m *Moodle) GetUserStats() (stats *UserStats, err error) {
	stats = &UserStats{LiveUsers: 0, ExpectedUpcomingExamParticipants: 0}

	before5minutes := strconv.Itoa(int(time.Now().Unix() - 300))
	err = m.Connection.QueryRow("SELECT COUNT(id) FROM mdl_user WHERE lastaccess >" + before5minutes).Scan(&stats.LiveUsers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return
	}
	more10minutes := strconv.Itoa(int(time.Now().Unix() + 600))
	err = m.Connection.QueryRow("SELECT COUNT (u.id) FROM mdl_quiz q JOIN mdl_course c ON c.id = q.course JOIN mdl_enrol e ON e.courseid = c.id JOIN mdl_user_enrolments ue ON ue.enrolid = e.id JOIN mdl_user u ON u.id = ue.userid WHERE q.timeopen <" + more10minutes + " AND q.timeclose >" + more10minutes).Scan(&stats.ExpectedUpcomingExamParticipants)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return
	}
	return
}

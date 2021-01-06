package client

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Moodle struct {
	PoolConfig *pgxpool.Config
	Shortname  string
}

type UserStats struct {
	LiveUsers                        int
	ExpectedUpcomingExamParticipants int
}

func New(hostname string, username string, password string, database string) (moodle *Moodle, err error) {
	moodle = nil
	connconf := pgx.ConnConfig{Host: hostname, User: username, Password: password, Database: database}
	conn, err := pgx.Connect(connconf)
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

	connstr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		username, password, hostname, database)
	poolconf, err := pgxpool.ParseConfig(connstr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed configuring connection pool: %v\n", err)
		return
	}
	moodle = &Moodle{PoolConfig: poolconf, Shortname: shortname}

	return
}

func GetLiveUsers(ctx context.Context, conn *pgxpool.Pool,
	stats *UserStats, wg *sync.WaitGroup) {

	wg.Add(1)
	go func() {
		defer wg.Done()
		before5minutes := strconv.Itoa(int(time.Now().Unix() - 300))
		var count int
		if err := conn.QueryRow(ctx, "SELECT COUNT(id) FROM mdl_user WHERE lastaccess >"+before5minutes).Scan(&count); err != nil {
			fmt.Fprintf(os.Stderr, "GetExpectedUpcomingExamParticipants failed: %v\n", err)
			return
		}
		stats.LiveUsers = count
	}()
}

func GetExpectedUpcomingExamParticipants(ctx context.Context,
	conn *pgxpool.Pool, stats *UserStats, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		more10minutes := strconv.Itoa(int(time.Now().Unix() + 600))
		var count int

		err := conn.QueryRow(ctx, "SELECT COUNT (u.id) FROM mdl_quiz q JOIN mdl_course c ON c.id = q.course JOIN mdl_enrol e ON e.courseid = c.id JOIN mdl_user_enrolments ue ON ue.enrolid = e.id JOIN mdl_user u ON u.id = ue.userid WHERE (q.timeclose - q.timeopen) < 60*60*10 AND q.timeopen <"+more10minutes+" AND q.timeclose >"+more10minutes).Scan(&count)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetExpectedUpcomingExamParticipants failed: %v\n", err)
			return
		}
		stats.ExpectedUpcomingExamParticipants = count
	}()
}

func (m *Moodle) GetUserStats() (stats *UserStats, err error) {
	stats = &UserStats{LiveUsers: 0, ExpectedUpcomingExamParticipants: 0}

	ctx := context.Background()
	conn, err := pgxpool.ConnectConfig(ctx, m.PoolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup

	GetLiveUsers(ctx, conn, stats, &wg)
	GetExpectedUpcomingExamParticipants(ctx, conn, stats, &wg)

	wg.Wait()

	return
}

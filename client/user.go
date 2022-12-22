package client

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStats struct {
	LiveUsers                        int
	ExpectedUpcomingExamParticipants int
}

func (m *Moodle) GetUserStats() (stats *UserStats, err error) {
	stats = &UserStats{LiveUsers: 0, ExpectedUpcomingExamParticipants: 0}

	ctx := context.Background()
	conn, err := pgxpool.NewWithConfig(ctx, m.poolconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup

	getLiveUsers(ctx, conn, stats, &wg)
	getExpectedUpcomingExamParticipants(ctx, conn, stats, &wg)

	wg.Wait()

	return
}

func getLiveUsers(ctx context.Context, conn *pgxpool.Pool,
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

func getExpectedUpcomingExamParticipants(ctx context.Context,
	conn *pgxpool.Pool, stats *UserStats, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		before5minutes := strconv.Itoa(int(time.Now().Unix() - 300))
		more20minutes := strconv.Itoa(int(time.Now().Unix() + 1200))
		var count int

		err := conn.QueryRow(ctx, "SELECT COUNT (DISTINCT(u.id)) FROM mdl_quiz q JOIN mdl_course c ON c.id = q.course JOIN mdl_enrol e ON e.courseid = c.id JOIN mdl_user_enrolments ue ON ue.enrolid = e.id JOIN mdl_user u ON u.id = ue.userid WHERE u.lastaccess <"+before5minutes+" AND (q.timeclose - q.timeopen) < 60*60*5 AND q.timeopen <"+more20minutes+" AND q.timeclose >"+more20minutes).Scan(&count)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetExpectedUpcomingExamParticipants failed: %v\n", err)
			return
		}
		stats.ExpectedUpcomingExamParticipants = count
	}()
}

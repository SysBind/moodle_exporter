package client

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type StorageStats struct {
	BytesAssignSubmission map[int]int
	BytesBackup           map[int]int
	BytesBackupAuto       map[int]int
	BytesAll              int64
}

func (m *Moodle) GetStorageStats() (stats *StorageStats, err error) {
	stats = &StorageStats{
		BytesAssignSubmission: make(map[int]int),
		BytesBackup:           make(map[int]int),
		BytesBackupAuto:       make(map[int]int),
		BytesAll:              0,
	}

	ctx := context.Background()
	conn, err := pgxpool.ConnectConfig(ctx, m.PoolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup

	getBytesAssignSubmission(ctx, conn, stats, &wg)
	getBytesBackup(ctx, conn, stats, &wg)
	getBytesBackupAuto(ctx, conn, stats, &wg)
	getBytesAll(ctx, conn, stats, &wg)

	wg.Wait()

	return
}

func getBytesAssignSubmission(ctx context.Context, conn *pgxpool.Pool,
	stats *StorageStats, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		if rows, err := conn.Query(ctx,
			`SELECT cm.course, SUM(f.filesize) FROM mdl_files f
			 JOIN mdl_context ctx ON ctx.id = f.contextid 
                         JOIN mdl_course_modules cm ON cm.id = ctx.instanceid 
                         WHERE component IN ('assignsubmission_file', 'assignfeedback_editpdf') 
                         AND ctx.contextlevel = 70 GROUP BY cm.course`); err != nil {
			fmt.Fprintf(os.Stderr, "getBytesAssignSubmission failed: %v\n", err)
		} else {
			defer rows.Close()

			for rows.Next() {
				var course, bytes int
				err := rows.Scan(&course, &bytes)
				if err != nil {
					fmt.Fprintf(os.Stderr, "getBytesAssignSubmission failed: %v\n", err)
					return
				}
				stats.BytesAssignSubmission[course] = bytes
			}
		}
	}()
}

func getBytesBackup(ctx context.Context, conn *pgxpool.Pool,
	stats *StorageStats, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		if rows, err := conn.Query(ctx,
			`SELECT c.id, SUM(f.filesize) FROM mdl_files f
			 JOIN mdl_context ctx ON ctx.id = f.contextid
                         JOIN mdl_course c ON c.id = ctx.instanceid
                         WHERE component = 'backup' AND filearea='course'
                         AND ctx.contextlevel = 50 GROUP BY c.id`); err != nil {
			fmt.Fprintf(os.Stderr, "getBytesBackup failed: %v\n", err)
		} else {
			defer rows.Close()

			for rows.Next() {
				var course, bytes int
				err := rows.Scan(&course, &bytes)
				if err != nil {
					fmt.Fprintf(os.Stderr, "getBytesBackup failed: %v\n", err)
					return
				}
				stats.BytesBackup[course] = bytes
			}
		}
	}()
}

func getBytesBackupAuto(ctx context.Context, conn *pgxpool.Pool,
	stats *StorageStats, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		if rows, err := conn.Query(ctx,
			`SELECT c.id, SUM(f.filesize) FROM mdl_files f
			 JOIN mdl_context ctx ON ctx.id = f.contextid
                         JOIN mdl_course c ON c.id = ctx.instanceid
                         WHERE component = 'backup' AND filearea='automated'
                         AND ctx.contextlevel = 50 GROUP BY c.id`); err != nil {
			fmt.Fprintf(os.Stderr, "getBytesBackupAuto failed: %v\n", err)
		} else {
			defer rows.Close()

			for rows.Next() {
				var course, bytes int
				err := rows.Scan(&course, &bytes)
				if err != nil {
					fmt.Fprintf(os.Stderr, "getBytesBackupAuto failed: %v\n", err)
					return
				}
				stats.BytesBackupAuto[course] = bytes
			}
		}
	}()
}

func getBytesAll(ctx context.Context, conn *pgxpool.Pool,
	stats *StorageStats, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		var bytes int64

		if err := conn.QueryRow(ctx,
			`SELECT SUM(filesize) FROM mdl_files`).Scan(&bytes); err != nil {
			fmt.Fprintf(os.Stderr, "getBytesAll failed: %v\n", err)
		} else {
			stats.BytesAll = bytes
		}
	}()
}

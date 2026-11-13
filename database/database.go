package database

import (
	"database/sql"
	"fmt"
	"kube-api/config"
	"time"

	_ "github.com/lib/pq"
)

type Solves struct {
	ID                  int
	CUBE_TYPE           string
	SCRAMBLE            string
	DURATION_NS         int
	TIMESTAMP_UNIX_NANO int
}

func Connect(cfg config.Config) *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DATABASE.HOSTNAME, cfg.DATABASE.PORT, cfg.DATABASE.USERNAME, cfg.DATABASE.PASSWORD, cfg.DATABASE.DATABASE)
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		fmt.Println("Failed to connect")
		fmt.Println(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to PING: %v\n", err)
	}
	fmt.Printf("Connected to database %s\n", cfg.DATABASE.HOSTNAME)

	return db
}

func GetFastestTime(database *sql.DB, cubeType string) (int, error) {
	stmt, err := database.Prepare(`
		SELECT id, cube_type, scramble, duration_ns, timestamp_unix_nano
		FROM cube_timer.solves
		WHERE cube_type = '?'
		ORDER BY duration_ns
		LIMIT 1
	`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var solve Solves
	err = stmt.QueryRow(cubeType).Scan(&solve.ID, &solve.CUBE_TYPE, &solve.SCRAMBLE, &solve.DURATION_NS, &solve.TIMESTAMP_UNIX_NANO)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, err
		}
	}

	return solve.DURATION_NS, nil
}

func GetLastTenSolves(database *sql.DB) ([]Solves, error) {
	var solves []Solves

	query := `
		SELECT id, cube_type, scramble, duration_ns, timestamp_unix_nano
		FROM cube_timer.solves
		ORDER BY timestamp_unix_nano DESC
		LIMIT 10
	`

	row, err := database.Query(query)
	if err != nil {
		return solves, err
	}
	defer row.Close()

	for row.Next() {
		var s Solves
		err := row.Scan(&s.ID, &s.CUBE_TYPE, &s.SCRAMBLE, &s.DURATION_NS, &s.TIMESTAMP_UNIX_NANO)
		if err != nil {
			if row.Err() != nil {
				return solves, err
			}
		}

		solves = append(solves, s)
	}

	return solves, nil
}

func GetTodayAverage(database *sql.DB, cubeType string) (int, error) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).UnixNano()

	fmt.Println(today)

	stmt, err := database.Prepare(`
		SELECT AVG(duration_ns)::bigint
		FROM cube_timer.solves
		WHERE cube_type = $1
		AND timestamp_unix_nano >= $2
	`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var average sql.NullInt64
	err = stmt.QueryRow(cubeType, today).Scan(&average)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, err
		}
	}

	if !average.Valid {
		return -1, fmt.Errorf("No valid average for today")
	}

	return int(average.Int64), nil
}

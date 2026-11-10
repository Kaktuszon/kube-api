package database

import (
	"database/sql"
	"fmt"
	"kube-api/config"

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
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to PING: %v\n", err)
	}
	fmt.Printf("Connected to database %s\n", cfg.DATABASE.HOSTNAME)

	return db
}

func GetFastestTime(database *sql.DB, cubeType string) (int, error) {
	var solve []Solves

	query := fmt.Sprintf(`
		SELECT id, cube_type, scramble, duration_ns, timestamp_unix_nano
		FROM cube_timer.solves
		WHERE cube_type = '%s'
		ORDER BY duration_ns
		LIMIT 1
	`, cubeType)

	row, err := database.Query(query)
	if err != nil {
		return 0, err
	}

	for row.Next() {
		var s Solves
		err := row.Scan(&s.ID, &s.CUBE_TYPE, &s.SCRAMBLE, &s.DURATION_NS, &s.TIMESTAMP_UNIX_NANO)
		if err != nil {
			return 0, err
		}

		if row.Err() != nil {
			return 0, err
		}

		solve = append(solve, s)
	}

	return solve[0].DURATION_NS, nil
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

	for row.Next() {
		var s Solves
		err := row.Scan(&s.ID, &s.CUBE_TYPE, &s.SCRAMBLE, &s.DURATION_NS, &s.TIMESTAMP_UNIX_NANO)
		if err != nil {
			return solves, err
		}

		if row.Err() != nil {
			return solves, err
		}

		solves = append(solves, s)
	}

	return solves, nil
}

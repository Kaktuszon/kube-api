package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var host string = os.Getenv("DB_HOSTNAME")
var port int = 5432
var user string = os.Getenv("DB_USERNAME")
var password string = os.Getenv("DB_PASSWORD")
var dbname string = os.Getenv("DB_DBNAME")
var serverhost string = os.Getenv("SERVER_HOSTNAME")
var serverport string = "8080"
var serveraddress string

type Solves struct {
	ID                  int
	CUBE_TYPE           string
	SCRAMBLE            string
	DURATION_NS         int
	TIMESTAMP_UNIX_NANO int
}

func Connect() *sql.DB {
	if host == "" {
		host = "localhost"
	}
	if serverhost == "" {
		serverhost = "localhost"
	}
	serveraddress = serverhost + ":" + serverport
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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
	fmt.Printf("Connected to database %s\n", host)

	return db
}

func GetFastestTime(database *sql.DB, cubeType string) (int, error) {
	var solve []Solves

	query := `
		SELECT id, cube_type, scramble, duration_ns, timestamp_unix_nano
		FROM cube_timer.solves
		ORDER BY duration_ns
		LIMIT 1
	`

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

		solve = append(solve, s)
	}

	return solve[0].DURATION_NS, nil
}

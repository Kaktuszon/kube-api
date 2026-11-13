package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"kube-api/config"
	"kube-api/database"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var SERVER_HOSTNAME = os.Getenv("SERVER_HOSTNAME")
var SERVER_PORT = os.Getenv("SERVER_PORT")

func main() {
	// Load config, parse the -path flag
	configPath := flag.String("path", "", "Path to config file, if empty ~/.kopierare/config.yaml")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("Could not load config, exiting...")
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	db := database.Connect(cfg)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the kube-timer API")
	})

	mux.HandleFunc("GET /api/fastest/{cube_type}", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		cubeType := path[len("/api/fastest/"):]
		time, _ := database.GetFastestTime(db, cubeType)
		fmt.Println(time)
		fmt.Fprintf(w, "%d", time)
	})

	mux.HandleFunc("GET /api/lasttensolves", func(w http.ResponseWriter, r *http.Request) {
		solves, err := database.GetLastTenSolves(db)
		if err != nil {
			fmt.Println("Could not get last ten solves")
		}

		jsondata, err := json.Marshal(solves)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, "%s", string(jsondata))
	})

	mux.HandleFunc("GET /api/todayaverage/{cube_type}", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		cubeType := path[len("/api/todayaverage/"):]
		average, err := database.GetTodayAverage(db, cubeType)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, "%d", average)
	})

	serveraddress := fmt.Sprintf("%s:%d", cfg.SERVER.HOSTNAME, cfg.SERVER.PORT)
	if err := http.ListenAndServe(serveraddress, mux); err != nil {
		fmt.Println(err)
	}

	db.Close()
}

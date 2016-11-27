package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"net/http"
)

func main() {
	config := pgx.ConnConfig{
		Host:     "localhost",
		Database: "mememarket",
		User:     "postgres",
	}
	conn, err := pgx.Connect(config)
	if err != nil {
		panic(err)
	}

	// Create and start the server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// TODO: Do this the right way, likely by creating multiple callbacks
		if path == "/" {
			handleRoot(w, r, conn)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	var memeCount int
	err := conn.QueryRow("SELECT COUNT(*) FROM memes").Scan(&memeCount)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Root. Memes: %d", memeCount)
}

package main

import (
	"fmt"
	"github.com/jackc/pgx"
	// "html"
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
	fmt.Println("IT WORKED")

	// Create and start the server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var memeCount int
		err = conn.QueryRow("SELECT COUNT(*) FROM memes").Scan(&memeCount)
		fmt.Fprintf(w, "Memes: %d", memeCount)
		// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func server() {

}

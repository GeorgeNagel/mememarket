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
			handleRoot(w, conn)
		} else if path == "/memes/add" {
			addMeme(w, conn)
		} else if path == "/memes" {
			listMemes(w, conn)
		} else if path == "/accounts/add" {
			addAccount(w, conn)
		} else if path == "/accounts" {
			listAccounts(w, conn)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, conn *pgx.Conn) {
	var memeCount int
	err := conn.QueryRow("SELECT COUNT(*) FROM memes").Scan(&memeCount)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Root. Memes: %d", memeCount)
}

func addAccount(w http.ResponseWriter, conn *pgx.Conn) {
	_, err := conn.Exec("insert into accounts (username, settled) VALUES ('test', 1000)")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Added account successfully.")
}

func listAccounts(w http.ResponseWriter, conn *pgx.Conn) {
	rows, err := conn.Query("select username, settled from accounts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	accounts := ""
	for rows.Next() {
		var username string
		var settled int
		err = rows.Scan(&username, &settled)
		if err != nil {
			// TODO: Clean up error propagation
			panic(err)
		}
		accounts = fmt.Sprintf("%s\n%s: $%d free", accounts, username, settled)
	}
	fmt.Fprint(w, accounts)
}

func listMemes(w http.ResponseWriter, conn *pgx.Conn) {
	rows, err := conn.Query("select name, price from memes")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	memes := ""
	for rows.Next() {
		var name string
		var price int
		err = rows.Scan(&name, &price)
		if err != nil {
			// TODO: Clean up error propagation
			panic(err)
		}
		memes = fmt.Sprintf("%s\n%s: $%d/share", memes, name, price)
	}
	fmt.Fprint(w, memes)
}

func addMeme(w http.ResponseWriter, conn *pgx.Conn) {
	_, err := conn.Exec("insert into memes (name, price) VALUES ('test', 12)")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Added meme successfully.")
}

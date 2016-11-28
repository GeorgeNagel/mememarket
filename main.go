package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Create and start the server
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/memes/add", addMeme)
	http.HandleFunc("/memes", listMemes)
	http.HandleFunc("/accounts/add", addAccount)
	http.HandleFunc("/accounts", listAccounts)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func connectPG() (*pgx.Conn, error) {
	config := pgx.ConnConfig{
		Host:     "localhost",
		Database: "mememarket",
		User:     "postgres",
	}
	conn, err := pgx.Connect(config)
	return conn, err
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Root.")
}

func addAccount(w http.ResponseWriter, r *http.Request) {
	conn, err := connectPG()
	if err != nil {
		fmt.Fprint(w, "Problem connecting to Postgres")
		return
	}
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles("accounts_new.html")
		t.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			// TODO: Handle this better
			panic(err)
		}
		formData := r.Form
		username := formData["username"][0]

		_, err = conn.Exec("insert into accounts (username, settled) VALUES ($1, 0)", username)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "Added account successfully.")
	}
}

func listAccounts(w http.ResponseWriter, r *http.Request) {
	conn, err := connectPG()
	if err != nil {
		fmt.Fprint(w, "Problem connecting to Postgres")
		return
	}

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

func listMemes(w http.ResponseWriter, r *http.Request) {
	conn, err := connectPG()
	if err != nil {
		fmt.Fprint(w, "Problem connecting to Postgres")
		return
	}

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

func addMeme(w http.ResponseWriter, r *http.Request) {
	conn, err := connectPG()
	if err != nil {
		fmt.Fprint(w, "Problem connecting to Postgres")
		return
	}

	_, err = conn.Exec("insert into memes (name, price) VALUES ('test', 12)")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Added meme successfully.")
}

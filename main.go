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
	http.HandleFunc("/purchase", purchaseShare)

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
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles("templates/accounts_new.html")
		t.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		conn, err := connectPG()
		if err != nil {
			fmt.Fprint(w, "Problem connecting to Postgres")
			return
		}

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

func purchaseShare(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles("templates/purchase_share.html")
		t.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		conn, err := connectPG()
		if err != nil {
			fmt.Fprint(w, "Problem connecting to Postgres")
			return
		}

		err = r.ParseForm()
		if err != nil {
			// TODO: Handle this better
			panic(err)
		}
		formData := r.Form
		account := formData["account"][0]
		meme := formData["meme"][0]
		amount := formData["amount"][0]

		_, err = conn.Exec(
			"insert into accounts_memes_map (account_id, meme_id, amount) VALUES ($1, $2, $3)",
			account, meme, amount)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "Purchase successful.")
	}
}

func listAccounts(w http.ResponseWriter, r *http.Request) {
	conn, err := connectPG()
	if err != nil {
		fmt.Fprint(w, "Problem connecting to Postgres")
		return
	}

	rows, err := conn.Query("select id, username, settled from accounts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	accounts := ""
	for rows.Next() {
		var id int
		var username string
		var settled int
		err = rows.Scan(&id, &username, &settled)
		if err != nil {
			// TODO: Clean up error propagation
			panic(err)
		}
		accounts = fmt.Sprintf("%s\n%d. %s: $%d free", accounts, id, username, settled)
	}
	fmt.Fprint(w, accounts)
}

func listMemes(w http.ResponseWriter, r *http.Request) {
	conn, err := connectPG()
	if err != nil {
		fmt.Fprint(w, "Problem connecting to Postgres")
		return
	}

	rows, err := conn.Query("select id, name, price from memes")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	memes := ""
	for rows.Next() {
		var id int
		var name string
		var price int
		err = rows.Scan(&id, &name, &price)
		if err != nil {
			// TODO: Clean up error propagation
			panic(err)
		}
		memes = fmt.Sprintf("%s\n%d. %s: $%d/share", memes, id, name, price)
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

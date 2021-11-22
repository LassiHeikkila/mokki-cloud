package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/LassiHeikkila/mokki-cloud/server/auth"
	"github.com/LassiHeikkila/mokki-cloud/server/auth/internal"
)

func main() {
	dbName := "auth.db"
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println("failed to open:", dbName)
		return
	}
	defer db.Close()

	auth.RegisterDatabase(db)
	if err := auth.InitializeDatabase(); err != nil {
		fmt.Println("failed to initialize database", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println()
		fmt.Println("create (c) a user?")
		fmt.Println("update (u) a user's password?")
		fmt.Println("remove (r) a user?")
		fmt.Println("check password (p) for user?")
		fmt.Println("or quit (q) to exit")
		fmt.Println()
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		switch input {
		case "create", "c":
			handleCreate(reader, db)
		case "remove", "r":
			handleRemove(reader, db)
		case "update", "u":
			handleUpdate(reader, db)
		case "password", "p":
			handleValidate(reader, db)
		case "quit", "exit", "q":
			fmt.Println("exiting!")
			return
		default:
			fmt.Println("unknown operation!")
		}
	}
}

func handleCreate(reader *bufio.Reader, db *sql.DB) {
	fmt.Println("enter username:")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Println("enter password:")
	password, _ := reader.ReadString('\n')
	password = strings.TrimRight(password, "\n")
	hashedPassword := auth.HashPassword(password)

	if err := internal.InsertUser(db, username, hashedPassword); err != nil {
		fmt.Println("error adding user to database:", err)
		return
	}
}

func handleRemove(reader *bufio.Reader, db *sql.DB) {
	fmt.Println("enter username of user to remove:")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if err := internal.RemoveUser(db, username); err != nil {
		fmt.Println("error adding user to database:", err)
		return
	}
}

func handleUpdate(reader *bufio.Reader, db *sql.DB) {
	fmt.Println("unimplemented")
}

func handleValidate(reader *bufio.Reader, db *sql.DB) {
	fmt.Println("enter username:")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Println("enter password:")
	password, _ := reader.ReadString('\n')
	password = strings.TrimRight(password, "\n")

	if auth.IsAuthorizedUser(username, password) {
		fmt.Println("password matches database!")
	} else {
		fmt.Println("password does not match database!")
	}
}

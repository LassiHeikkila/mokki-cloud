package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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
		fmt.Println("create (c) or revoke (r) a token?")
		fmt.Println("or quit (q) to exit")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		switch input {
		case "c":
			handleCreate(reader, db)
		case "r":
			handleRevoke(reader, db)
		case "quit", "exit", "q":
			fmt.Println("exiting!")
			return
		default:
			fmt.Println("unknown operation!")
		}
	}
}

func handleCreate(reader *bufio.Reader, db *sql.DB) {
	fmt.Println("how long should token be valid for?")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	dur, err := parseDuration(input)
	if err != nil {
		fmt.Println("error parsing duration:", err)
		return
	}
	now := time.Now()
	expiryDate := now.Add(dur)

	token, err := auth.GenerateToken()
	if err != nil {
		fmt.Println("failed to generate token:", err)
		return
	}

	if err := internal.InsertToken(db, token, now, expiryDate, false); err != nil {
		fmt.Println("failed to add token to database:", err)
		return
	}
}

func handleRevoke(reader *bufio.Reader, db *sql.DB) {
	fmt.Println("enter token to revoke:")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	err := internal.RevokeToken(db, input)
	if err != nil {
		fmt.Println("failed to revoke token:", err)
		return
	}
}

func parseDuration(input string) (time.Duration, error) {
	if d, err := time.ParseDuration(input); err == nil {
		return d, nil
	}

	re := regexp.MustCompile(`^([0-9]+)([a-zA-Z])$`)
	if !re.MatchString(input) {
		return 0, errors.New("input does not match expected format")
	}

	submatches := re.FindStringSubmatch(input)
	if len(submatches) != 3 {
		return 0, errors.New("input does not match expected format")
	}

	// [0] is full match
	// [1] is first group, i.e. integer value'
	// [2] is second group, i.e. unit

	dur, err := strconv.ParseInt(submatches[1], 10, 64)
	if err != nil {
		return 0, errors.New("input does not match expected format")
	}

	const (
		day  = 24 * time.Hour
		week = 7 * day
		year = 365 * day
	)

	var multiplier time.Duration // hours
	switch submatches[2] {
	case "h":
		multiplier = time.Hour
	case "d":
		multiplier = day
	case "w":
		multiplier = week
	case "y":
		multiplier = year
	default:
		return 0, errors.New("unknown time unit: " + submatches[2])
	}

	return time.Duration(dur) * multiplier, nil
}

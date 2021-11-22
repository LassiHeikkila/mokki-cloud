package internal

import (
	"database/sql"
	"errors"
	"time"
)

const (
	// ISO8601 looks like "2016-01-01 10:20:05.123"
	iso8601 = `2006-01-02 03:04:05.000`
)

func InsertUser(db *sql.DB, username, hashedPassword string) error {
	if db == nil {
		return errors.New("no database registered")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO credentials
		(username, password)
		VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(username, hashedPassword); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func RemoveUser(db *sql.DB, username string) error {
	if db == nil {
		return errors.New("no database registered")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`DELETE FROM credentials
		WHERE username == ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(username); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func GetUsersHashedPassword(db *sql.DB, username string) (string, error) {
	if db == nil {
		return "", errors.New("no database registered")
	}

	stmt, err := db.Prepare(`SELECT password FROM credentials
		WHERE username == ?`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var hashedPw string
		if err := rows.Scan(&hashedPw); err != nil {
			return "", err
		}
		return hashedPw, nil
	}

	return "", errors.New("no matching user found")
}

func InsertToken(db *sql.DB, token string, validFrom time.Time, validTo time.Time, revoked bool) error {
	if db == nil {
		return errors.New("no database registered")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO tokens
		(token, validFrom, validTo, revoked)
		VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	fromStr := validFrom.Format(iso8601)
	toStr := validTo.Format(iso8601)
	rev := 0
	if revoked {
		rev = 1
	}

	if _, err := stmt.Exec(token, fromStr, toStr, rev); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func RevokeToken(db *sql.DB, token string) error {
	if db == nil {
		return errors.New("no database registered")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`UPDATE tokens
		SET revoked = 1
		WHERE token = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(token); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func ContainsValidToken(db *sql.DB, token string) bool {
	if db == nil {
		return false
	}

	stmt, err := db.Prepare(`SELECT token, validFrom, validTo, revoked FROM tokens
		WHERE token == ?`)
	if err != nil {
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(token)
	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var t, validFrom, validTo string
		var revoked int
		if err := rows.Scan(&t, &validFrom, &validTo, &revoked); err != nil {
			return false
		}
		if t != token {
			continue
		}
		from, err := time.Parse(iso8601, validFrom)
		if err != nil {
			continue
		}
		to, err := time.Parse(iso8601, validTo)
		if err != nil {
			continue
		}
		now := time.Now()
		if now.After(from) && now.Before(to) && revoked == 0 {
			return true
		}
		return false
	}

	return false
}

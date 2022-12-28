package auth

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/LassiHeikkila/mokki-cloud/server/auth/internal"
)

var (
	databaseHandle *sql.DB = nil
)

func RegisterDatabase(db *sql.DB) {
	databaseHandle = db
}

func InitializeDatabase() error {
	if databaseHandle == nil {
		return errors.New("no database registered")
	}

	if _, err := databaseHandle.Exec(internal.CredentialsTableInitStmt); err != nil {
		return err
	}
	if _, err := databaseHandle.Exec(internal.TokensTableInitStmt); err != nil {
		return err
	}

	return nil
}

// AuthorizedUser returns true if user and password are valid,
// otherwise false.
func IsAuthorizedUser(user string, password string) bool {
	h, err := internal.GetUsersHashedPassword(databaseHandle, user)
	if err != nil {
		log.Println("did not find matching user")
		// no user found
		return false
	}
	return comparePasswordAndHash(password, h)
}

func TokenIsValid(token string) bool {
	return internal.ContainsValidToken(databaseHandle, token)
}

// Generates a new token which will be valid for given dur, or 4 weeks if dur is zero.
func GenerateToken(dur time.Duration) (string, error) {
	if dur == 0 {
		dur = time.Hour * 24 * 28
	}
	tok, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	validFrom := time.Now().UTC()
	validTo := time.Now().UTC().Add(dur)
	err = internal.InsertToken(databaseHandle, tok.String(), validFrom, validTo, false)
	if err != nil {
		return "", err
	}
	return tok.String(), nil
}

func HashPassword(password string) string {
	const bcryptCost = 14
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return ""
	}
	return string(b)
}

func comparePasswordAndHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

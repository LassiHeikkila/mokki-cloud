package internal

const CredentialsTableInitStmt = `
CREATE TABLE IF NOT EXISTS "credentials"
(
	id INTEGER NOT NULL,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	PRIMARY KEY (id)
);
`

const TokensTableInitStmt = `
CREATE TABLE IF NOT EXISTS "tokens"
(
	id INTEGER NOT NULL,
	token TEXT NOT NULL UNIQUE,
	validFrom TEXT NOT NULL,
	validTo TEXT NOT NULL,
	revoked INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (id)
);
`

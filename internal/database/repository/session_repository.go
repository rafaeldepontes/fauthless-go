package repository

import (
	"database/sql"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

type Session struct {
	Id           string
	Username     string
	IsRevoked    bool
	RefreshToken string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

// NewSessionRepository initialize a new SessionRepository containing
// a database connection, it returns a pointer to the new SessionRepository.
func NewSessionRepository(conn *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: conn,
	}
}

// CreateSession creates a new session, it expects a session object
// and returns an error if any
func (r *SessionRepository) CreateSession(session *Session) (string, error) {
	query := `
	INSERT INTO sessions (id, username, is_revoked, refresh_token, expires_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	var stmt *sql.Stmt
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var id string
	err = stmt.QueryRow(session.Id, session.Username, session.IsRevoked, session.RefreshToken, session.ExpiresAt).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// FindSessionById searchs for a session based on its identifier, it
// expects the session identifier, returns the session and an error if
// any.
func (r *SessionRepository) FindSessionById(id string) (*Session, error) {
	query := `
	SELECT id, username, is_revoked, refresh_token, created_at, expires_at 
	FROM sessions
	WHERE id = $1
	`
	var stmt *sql.Stmt
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var session Session
	err = stmt.QueryRow(id).Scan(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// RevokeSession revokes a session by its identifier, expects the identifier
// and returns an error if any.
func (r *SessionRepository) RevokeSession(id string) error {
	query := `
	UPDATE sessions
	SET is_revoked = true
	WHERE id = $1
	`

	var stmt *sql.Stmt
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSession removes a session from the database by its identifier,
// expects the id and return an error if any.
func (r *SessionRepository) DeleteSession(id string) error {
	query := `
	DELETE FROM sessions WHERE id = $1
	`

	var stmt *sql.Stmt
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

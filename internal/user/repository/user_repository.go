package repository

import (
	"database/sql"
	"errors"

	"github.com/rafaeldepontes/fauthless-go/internal/domain"
	"github.com/rafaeldepontes/fauthless-go/internal/errorhandler"
	"github.com/rafaeldepontes/fauthless-go/internal/user"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository initialize a new UserRepository containing
// a database connection, it returns a pointer to the new UserRepository.
func NewUserRepository(conn *sql.DB) user.Repository {
	return &userRepository{
		db: conn,
	}
}

// FindAllUsers search for all the users without a filter
// using a pagination method based on cursor, has a pointer
// to the next record from the database.
//
// returns a slice of User, the total of records in the database
// and an Error if any.
func (repo *userRepository) FindAllUsersCursor(cursor int64, size int) ([]domain.User, int64, error) {
	query := `
		SELECT id, username, age
		FROM users
		WHERE id >= $1
		ORDER BY id ASC
		LIMIT $2;
	`
	rows, err := repo.db.Query(query, cursor, size)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]domain.User, 0, size)
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.Id, &u.Username, &u.Age); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var nextCursor int64
	if len(users) == size {
		nextCursor = *(users[len(users)-1].Id)
		users = users[:len(users)-1]
	} else {
		nextCursor = 0
	}

	return users, nextCursor, nil
}

// FindAllUsers search for all the users without a filter
// using a pagination method based on offset and limit
// returns a slice of User, the total of records in the database
// and an Error if any.
func (repo *userRepository) FindAllUsers(size, page int) ([]domain.User, int, error) {
	queryCount := `SELECT COUNT(id) FROM users;`

	var total int
	if err := repo.db.QueryRow(queryCount).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, username, age FROM users ORDER BY id ASC LIMIT $1 OFFSET $2;`
	offset := (page - 1) * size

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.Id, &u.Username, &u.Age); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindUserById search for an user by his id
// returns the user and an error if any.
func (repo *userRepository) FindUserById(id int64) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, age FROM users WHERE id = $1;`
	row := repo.db.QueryRow(query, id)
	if err := row.Scan(&user.Id, &user.Username, &user.Age); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorhandler.ErrUserNotFound // or return custom NotFound error
		}
		return nil, err
	}
	return &user, nil
}

// FindUserByUsername search for an user by his username
// returns the user and an error if any.
func (repo *userRepository) FindUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, password, username, age FROM users WHERE username = $1;`

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Id, &user.HashedPassword, &user.Username, &user.Age)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// RegisterUser saves a new user into the database, expects a
// pointer to a user and returns an error if any.
func (repo *userRepository) RegisterUser(u *domain.User) error {
	query := `
		INSERT INTO users (username, password, age)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	err := repo.db.QueryRow(query, u.Username, u.HashedPassword, u.Age).Scan(&u.Id)
	return err
}

// SetUserToken save the user token into the database, expects
// both csrft and session token and also the Id related to the user
// returns an error if any.
func (repo *userRepository) SetUserToken(token, csrfToken string, userId int64) error {
	query := `
	UPDATE users 
	SET session_token = $1, 
	csrf_token = $2 
	WHERE id = $3;
	`

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(token, csrfToken, userId)

	return nil
}

// UpdateUserDetails updates the user details in the ddatabase
// with the new username and/or age, returns an erro if any.
func (repo *userRepository) UpdateUserDetails(user *domain.User) error {
	query := `
	UPDATE users
	SET age = $1
	WHERE id = $2
	`

	var stmt *sql.Stmt
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Age, user.Id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAccount deletes the user account from the database
// based on the username, returns an error if any.
func (repo *userRepository) DeleteAccount(username string) error {
	_, err := repo.db.Exec(`DELETE FROM users WHERE username = $1`, username)
	return err
}

package repository

import (
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

type User struct {
	Username       *string `json:"username"`
	HashedPassword *string `json:"password,omitempty"`
	Id             *uint   `json:"id,omitempty"`
	Age            *uint   `json:"age,omitempty"`
}

// NewUserRepository initialize a new UserRepository containing
// a database connection, it returns a pointer to the new UserRepository.
func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{
		db: conn,
	}
}

// FindAllUsers search for all the users without a filter
// using a pagination method based on offset and limit
// returns a slice of User, the total of records in the database
// and an Error if any.
func (repo *UserRepository) FindAllUsers(size, page int) ([]User, int, error) {
	queryCount := `SELECT COUNT(id) FROM users;`

	var total int
	if err := repo.db.QueryRow(queryCount).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, username, age FROM users LIMIT $1 OFFSET $2;`
	offset := (page - 1) * size

	rows, err := repo.db.Query(query, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
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
func (repo *UserRepository) FindUserById(id uint) (*User, error) {
	var user User
	query := `SELECT id, username, age FROM users WHERE id = $1;`

	err := repo.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Age)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

// FindUserByUsername search for an user by his username
// returns the user and an error if any.
func (repo *UserRepository) FindUserByUsername(username string) (*User, error) {
	var user User
	query := `SELECT id, password, username FROM users WHERE username = $1;`

	err := repo.db.QueryRow(query, username).Scan(&user.Id, &user.HashedPassword, &user.Username)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

// RegisterUser saves a new user into the database, expects a
// pointer to a user and returns an error if any.
func (repo *UserRepository) RegisterUser(user *User) error {
	query := `
	INSERT INTO users (username, password, age)
	VALUES ($1, $2, $3) 
	RETURNING id;
	`

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return err
	}

	stmt.Exec(user.Username, user.HashedPassword, user.Age)

	return nil
}

// SetUserToken save the user token into the database, expects
// both csrft and session token and also the Id related to the user
// returns an error if any.
func (repo *UserRepository) SetUserToken(token, csrfToken string, userId uint) error {
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

	stmt.Exec(token, csrfToken, userId)

	return nil
}

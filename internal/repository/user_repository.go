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
// a database connection, it returns a pointer to the new UserRepository
func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{
		db: conn,
	}
}

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

func (repo *UserRepository) FindUserById(id uint) (User, error) {
	var user User
	query := `SELECT id, username, age FROM users WHERE id = $1`

	err := repo.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Age)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

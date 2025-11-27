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

func (repo *UserRepository) FindAllUsers() ([]User, error) {
	query := `SELECT id, username, age FROM users;`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usersMemAddress := new([]User)
	users := *usersMemAddress

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Username, &u.Age); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

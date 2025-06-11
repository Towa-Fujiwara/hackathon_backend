package dao

import (
	"database/sql"
	"fmt"
	"hackathon/model"
	"log"
)

type UserDao interface {
	FindById(id string) (*model.User, error)
	Create(user *model.User) error
	SearchByName(query string) ([]model.User, error)
}

type userDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) UserDao {
	return &userDao{db: db}
}

func (d *userDao) FindById(Id string) (*model.User, error) {
	row := d.db.QueryRow("SELECT id, name, age, password, display_name, bio, icon_url FROM users WHERE id = ?", Id)
	
	user := &model.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Age, &user.Password, &user.Profile.DisplayName, &user.Profile.Bio, &user.Profile.IconUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	
	return user, nil
}

func (d *userDao) Create(user *model.User) error {
	_, err := d.db.Exec(
		"INSERT INTO users (id, name, age, password, display_name, bio, icon_url, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.Id,
		user.Name,
		user.Age,
		user.Password,
		user.Profile.DisplayName,
		user.Profile.Bio,
		user.Profile.IconUrl,
		user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}


func (d *userDao) SearchByName(query string) ([]model.User, error) {
	query = "%" + query + "%"

	SQL := "SELECT id, name, age, password, display_name, bio, icon_url FROM users WHERE name LIKE ?"
	rows, err := d.db.Query(SQL, query)
	if err != nil {
		log.Printf("ERROR: Failed to search users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Password, &user.Profile.DisplayName, &user.Profile.Bio, &user.Profile.IconUrl)
		if err != nil {
			log.Printf("ERROR: Failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

package dao

import (
	"database/sql"
	"fmt"
	"hackathon/model"
	"log"
)

type UserDao interface {
	FindById(id string) (*model.User, error)
	FindByFirebaseUID(firebaseUID string) (*model.User, error)
	Create(user *model.User) error
	SearchByName(query string) ([]model.User, error)
}

type userDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) UserDao {
	return &userDao{db: db}
}

func (d *userDao) FindByFirebaseUID(firebaseUID string) (*model.User, error) {
	row := d.db.QueryRow("SELECT userId, firebaseUid, name, bio, iconUrl FROM users WHERE firebaseUid = ?", firebaseUID)
	user := &model.User{}
	err := row.Scan(&user.UserId, &user.FirebaseUID, &user.Name, &user.Bio, &user.IconUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return user, nil
}

func (d *userDao) FindById(Id string) (*model.User, error) {
	row := d.db.QueryRow("SELECT userId, firebaseUid, name, bio, iconUrl FROM users WHERE userId = ?", Id)
	
	user := &model.User{}
	err := row.Scan(&user.UserId, &user.FirebaseUID, &user.Name, &user.Bio, &user.IconUrl)
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	
	return user, nil
}

func (d *userDao) Create(user *model.User) error {
	_, err := d.db.Exec(
		"INSERT INTO users (userId, firebaseUid, name, bio, iconUrl, createdAt) VALUES (?, ?, ?, ?, ?, ?)",
		user.UserId,
		user.FirebaseUID,
		user.Name,
		user.Bio,
		user.IconUrl,
		user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}


func (d *userDao) SearchByName(query string) ([]model.User, error) {
	query = "%" + query + "%"

	SQL := "SELECT userId, firebaseUid, name, bio, iconUrl FROM users WHERE name LIKE ?"
	rows, err := d.db.Query(SQL, query)
	if err != nil {
		log.Printf("ERROR: Failed to search users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.UserId, &user.FirebaseUID, &user.Name, &user.Bio, &user.IconUrl)
		if err != nil {
			log.Printf("ERROR: Failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

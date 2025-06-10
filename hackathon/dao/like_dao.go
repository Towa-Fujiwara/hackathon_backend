package dao

import (
	"database/sql"
	"fmt"
	"hackathon/model"
)

type LikeDao interface {
	FindById(userId string, postId string) (*model.Like, error)
	Create(like *model.Like) error
	Delete(id string) error
}

type likeDao struct {
	db *sql.DB
}

func NewLikeDao(db *sql.DB) LikeDao {
	return &likeDao{db: db}
}
func (d *likeDao) FindById(userId string, postId string) (*model.Like, error) {
	row := d.db.QueryRow("SELECT id, user_id, post_id, created_at FROM like WHERE user_id = ? AND post_id = ?", userId, postId)
	var like model.Like
	if err := row.Scan(&like.Id, &like.UserId, &like.PostId, &like.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan like: %w", err)
	}
	return &like, nil
}

//Like作成
func (d *likeDao) Create(like *model.Like) error {
	_, err := d.db.Exec("INSERT INTO like (id, user_id, post_id, created_at) VALUES (?, ?, ?, ?)",
		like.Id, like.UserId, like.PostId, like.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}
	return nil
}
//Like削除
func (d *likeDao) Delete(id string) error {
	result, err := d.db.Exec("DELETE FROM like WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("like id %s not found or already deleted", id)
	}
	return nil
}
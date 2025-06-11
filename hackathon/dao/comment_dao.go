package dao

import (
	"hackathon/model"
	"database/sql"
	"log"
)

type CommentDao interface {
	CreateComment(comment *model.Comment) error
	GetCommentsByPostId(postId string) ([]model.Comment, error)
}

type commentDao struct {
	db *sql.DB
}

func NewCommentDao(db *sql.DB) CommentDao {
	return &commentDao{db: db}
}

func (c *commentDao) CreateComment(comment *model.Comment) error {
	query := "INSERT INTO comments (id, post_id, user_id, text, created_at) VALUES (?, ?, ?, ?, ?)"
	_, err := c.db.Exec(query, comment.Id, comment.PostId, comment.UserId, comment.Text, comment.CreatedAt)
	return err
}

func (c *commentDao) GetCommentsByPostId(postId string) ([]model.Comment, error) {
	query := "%" + postId + "%"

	SQL := "SELECT id, post_id, user_id, text, created_at FROM comments WHERE post_id LIKE ?"
	rows, err := c.db.Query(SQL, query)
	if err != nil {
		log.Printf("ERROR: Failed to search users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(&comment.Id, &comment.PostId, &comment.UserId, &comment.Text, &comment.CreatedAt)
		if err != nil {
			log.Printf("ERROR: Failed to scan user: %v", err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
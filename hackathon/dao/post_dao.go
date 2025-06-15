package dao

import (
	"database/sql"
	"fmt"
	"hackathon/model" 
)

type PostDao interface {
    FindAll() ([]model.Post, error)
    Create(post *model.Post) error
	FindById(id string) (*model.Post, error)
	Update(post *model.Post) error
	Delete(id string) error
	FindAllByUserId(uid string) ([]model.Post, error)
}
type postDao struct {
	db *sql.DB
}

func NewPostDao(db *sql.DB) PostDao {
	return &postDao{db: db}
}

func (d *postDao) FindById(Id string) (*model.Post, error) {
	row := d.db.QueryRow("SELECT id, userId, text, image, createdAt FROM posts WHERE id = ?", Id)

	post := &model.Post{}
	err := row.Scan(&post.Id, &post.UserId, &post.Text, &post.Image, &post.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to scan post: %w", err)
	}
	return post, nil
}



func (d *postDao) FindAll() ([]model.Post, error) {
	const selectPosts = "SELECT id, userId, text, image, createdAt FROM posts ORDER BY createdAt DESC"

	rows, err := d.db.Query(selectPosts)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for all posts: %w", err)
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.Id, &post.UserId, &post.Text, &post.Image, &post.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post row: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("an error occurred during rows iteration: %w", err)
	}

	return posts, nil
}

func (d *postDao) FindAllByUserId(uid string) ([]model.Post, error) {
	const selectPosts = "SELECT id, userId, text, image, createdAt FROM posts WHERE userId = ? ORDER BY createdAt DESC"

	rows, err := d.db.Query(selectPosts, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for all posts: %w", err)
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.Id, &post.UserId, &post.Text, &post.Image, &post.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post row: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("an error occurred during rows iteration: %w", err)
	}

	return posts, nil
}

func (d *postDao) Create(post *model.Post) error {
	if post == nil {
		return fmt.Errorf("post is empty")
	}
	if post.UserId == "" {
		return fmt.Errorf("user id is not set")
	}
	_, err := d.db.Exec("INSERT INTO posts (id, userId, text, image, createdAt) VALUES (?, ?, ?, ?, ?)",
		post.Id, post.UserId, post.Text, post.Image, post.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

func (d *postDao) Update(post *model.Post) error {
	if post == nil {
		return fmt.Errorf("post is empty")
	}
	_, err := d.db.Exec("UPDATE posts SET text = ?, image = ? WHERE id = ?", post.Text, post.Image, post.Id)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

func (d *postDao) Delete(Id string) error {
	result, err := d.db.Exec("DELETE FROM posts WHERE id = ?", Id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("post id %s not found or already deleted", Id)
	}

	return nil
}

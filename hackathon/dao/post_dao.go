package dao

import (
	"database/sql"
	"fmt"
	"hackathon/model"
	"log"
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
	const selectPost = `
		SELECT
			p.id,
			p.userId,
			u.name AS userName,
			p.text,
			p.image,
			p.createdAt,
			COUNT(DISTINCT l.id) AS likeCount,
			COUNT(DISTINCT c.id) AS commentCount
		FROM
			posts AS p
		LEFT JOIN
			users AS u ON p.userId = u.userId
		LEFT JOIN
			likes AS l ON p.id = l.postId
		LEFT JOIN
			comments AS c ON p.id = c.postId
		WHERE
			p.id = ?
		GROUP BY
			p.id
	`
	row := d.db.QueryRow(selectPost, Id)

	post := &model.Post{}
	err := row.Scan(&post.Id, &post.UserId, &post.UserName, &post.Text, &post.Image, &post.CreatedAt, &post.LikeCount, &post.CommentCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan post: %w", err)
	}
	return post, nil
}

func (d *postDao) FindAll() ([]model.Post, error) {
	const selectPosts = `
        SELECT
            p.id,
            p.userId,
            u.name AS userName,
            p.text,
            p.image,
            p.createdAt,
            COUNT(DISTINCT l.id) AS likeCount,
            COUNT(DISTINCT c.id) AS commentCount
        FROM
            posts AS p
        LEFT JOIN
            users AS u ON p.userId = u.userId
        LEFT JOIN
            likes AS l ON p.id = l.postId
        LEFT JOIN
            comments AS c ON p.id = c.postId
        GROUP BY
            p.id
        ORDER BY
            p.createdAt DESC
    `

	rows, err := d.db.Query(selectPosts)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for all posts: %w", err)
	}
	defer rows.Close()

	posts := []model.Post{}
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.Id, &post.UserId, &post.UserName, &post.Text, &post.Image, &post.CreatedAt, &post.LikeCount, &post.CommentCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post row: %w", err)
		}
		
		// 投稿データをログ出力
		log.Printf("取得した投稿: ID=%s, UserID=%s, UserName=%s, Text=%s, Image=%s, CreatedAt=%s, LikeCount=%d, CommentCount=%d",
			post.Id, post.UserId, post.UserName, post.Text, post.Image, post.CreatedAt, post.LikeCount, post.CommentCount)
		
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("an error occurred during rows iteration: %w", err)
	}

	log.Printf("合計 %d 件の投稿を取得しました", len(posts))
	return posts, nil
}

func (d *postDao) FindAllByUserId(uid string) ([]model.Post, error) {
	const selectPosts = `
		SELECT
			p.id,
			p.userId,
			u.name AS userName,
			p.text,
			p.image,
			p.createdAt,
			COUNT(DISTINCT l.id) AS likeCount,
			COUNT(DISTINCT c.id) AS commentCount
		FROM
			posts AS p
		LEFT JOIN
			users AS u ON p.userId = u.userId
		LEFT JOIN
			likes AS l ON p.id = l.postId
		LEFT JOIN
			comments AS c ON p.id = c.postId
		WHERE
			p.userId = ?
		GROUP BY
			p.id
		ORDER BY
			p.createdAt DESC
	`

	rows, err := d.db.Query(selectPosts, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for all posts by user: %w", err)
	}
	defer rows.Close()

	posts := []model.Post{}
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.Id, &post.UserId, &post.UserName, &post.Text, &post.Image, &post.CreatedAt, &post.LikeCount, &post.CommentCount)
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
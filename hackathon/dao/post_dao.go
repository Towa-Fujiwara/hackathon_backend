package dao

import (
	"database/sql"
	"fmt"
	"hackathon/model" // modelパッケージのインポートを追加
)

// var db *sql.DB // 重複定義なのでコメントアウト（または削除）

type PostDao interface {
    // DBから取得・スキャンして []model.Post を返す
    FindAll() ([]model.Post, error)
    // DBにINSERTする
    Create(post *model.Post) error
	FindById(id string) (*model.Post, error)
	Update(post *model.Post) error
	Delete(id string) error
}
type postDao struct {
	db *sql.DB
}
//後で変更？
func NewPostDao(db *sql.DB) PostDao {
	return &postDao{db: db}
}

//ポストIDによるポスト取得
func (d *postDao) FindById(Id string) (*model.Post, error) {
	// "content" を "text, image" に修正
	row := d.db.QueryRow("SELECT id, user_id, text, image, created_at FROM post WHERE id = ?", Id)

	post := &model.Post{}
	// モデルに合わせて UserId, Text, Image に修正
	err := row.Scan(&post.Id, &post.UserId, &post.Text, &post.Image, &post.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to scan post: %w", err)
	}
	return post, nil
}



func (d *postDao) FindAll() ([]model.Post, error) {
	// 投稿を取得するためのSQLクエリ（実際には定数として定義）
	const selectPosts = "SELECT id, user_id, text, image, created_at FROM post ORDER BY created_at DESC"

	rows, err := d.db.Query(selectPosts)
	if err != nil {
		// panic(err) の代わりに、エラーをラップして呼び出し元に返します
		return nil, fmt.Errorf("failed to execute query for all posts: %w", err)
	}
	defer rows.Close()

	// 投稿の一覧を格納する配列
	var posts []model.Post

	// 取得した投稿を一つずつ取りだして、配列に格納する
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.Id, &post.UserId, &post.Text, &post.Image, &post.CreatedAt)
		if err != nil {
			// ここも panic ではなく、エラーを返します
			return nil, fmt.Errorf("failed to scan post row: %w", err)
		}
		posts = append(posts, post)
	}

	// ループ処理中に発生したエラーがないか最終チェックします
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("an error occurred during rows iteration: %w", err)
	}

	return posts, nil
}


//ポスト作成
func (d *postDao) Create(post *model.Post) error {
	if post == nil {
		return fmt.Errorf("post is empty")
	}
	// モデルに合わせて UserId, Text に修正
	if post.UserId == "" {
		return fmt.Errorf("user id is not set")
	}
	_, err := d.db.Exec("INSERT INTO post (id, user_id, text, image, created_at) VALUES (?, ?, ?, ?, ?)",
		post.Id, post.UserId, post.Text, post.Image, post.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

//ポスト更新
func (d *postDao) Update(post *model.Post) error {
	if post == nil {
		return fmt.Errorf("post is empty")
	}
	_, err := d.db.Exec("UPDATE post SET text = ?, image = ? WHERE id = ?", post.Text, post.Image, post.Id)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

//ポスト削除
func (d *postDao) Delete(Id string) error {

	result, err := d.db.Exec("DELETE FROM post WHERE id = ?", Id)
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

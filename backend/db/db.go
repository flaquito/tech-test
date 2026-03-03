package db

import (
	"database/sql"
	"encoding/json"
	"example.com/backend/post"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const pageSize = 10

func GetPostsPage(db *sql.DB, page int) ([]post.Post, error) {
	offset := (page - 1) * pageSize

	rows, err := db.Query(`
	SELECT * FROM posts
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []post.Post

	for rows.Next() {
		var p post.Post
		var tagsRaw string

		if err := rows.Scan(&p.ID,
			&p.ImageURL,
			&p.Text,
			&tagsRaw,
			&p.CreatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(tagsRaw), &p.Tags); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func CreatePost(
	db *sql.DB,
	imageDir string,
	text string,
	tagsJSON []byte,
	imageBytes []byte,
	ext string,
) (post.Post, error) {

	tx, err := db.Begin()
	if err != nil {
		return post.Post{}, err
	}
	defer tx.Rollback()

	now := time.Now()

	res, err := tx.Exec(`
		INSERT INTO posts (image_url, text, tags, created_at)
		VALUES (?, ?, ?, ?)
	`, "", text, tagsJSON, now)
	if err != nil {
		return post.Post{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return post.Post{}, err
	}

	filename := fmt.Sprintf("%d%s", id, ext)
	path := filepath.Join(imageDir, filename)

	if err := os.WriteFile(path, imageBytes, 0644); err != nil {
		return post.Post{}, err
	}

	_, err = tx.Exec(`
		UPDATE posts SET image_url = ?
		WHERE id = ?
	`, filename, id)
	if err != nil {
		return post.Post{}, err
	}

	if err := tx.Commit(); err != nil {
		return post.Post{}, err
	}

	var tags []string
	if err := json.Unmarshal(tagsJSON, &tags); err != nil {
		tags = []string{}
	}
	if tags == nil {
		tags = []string{}
	}

	return post.Post{
		ID:        int(id),
		ImageURL:  filename,
		Text:      text,
		Tags:      tags,
		CreatedAt: now,
	}, nil
}

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath     = "../../data/app.db"
	imageDir   = "../../data/images"
	postCount  = 25
	randomSeed = 42
)

var availableTags = []string{
	"cat", "dog", "bear", "nature",
	"sun", "moon", "flowers", "friday",
}

func main() {
	rand.New(rand.NewSource(randomSeed))

	if err := os.MkdirAll(imageDir, 0755); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTable(db); err != nil {
		log.Fatal(err)
	}

	alreadySeeded, err := hasPosts(db)
	if err != nil {
		log.Fatal(err)
	}
	if alreadySeeded {
		log.Println("Database already seeded. Skipping.")
		return
	}

	for i := 1; i <= postCount; i++ {
		tags := randomTags()
		text := generateText(tags)
		imagePath := fmt.Sprintf("%s/%d.jpg", imageDir, i)

		if err := downloadImage(i, imagePath); err != nil {
			log.Fatal(err)
		}

		tagsJSON, _ := json.Marshal(tags)

		_, err := db.Exec(`
			INSERT INTO posts (image_url, text, tags, created_at)
			VALUES (?, ?, ?, ?)
		`,
			fmt.Sprintf("%d.jpg", i), // URL path for serving
			text,
			string(tagsJSON),
			time.Now(),
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Seeding completed successfully.")
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			image_url TEXT NOT NULL,
			text TEXT NOT NULL,
			tags TEXT NOT NULL,
			created_at DATETIME NOT NULL
		)
	`)
	return err
}

func hasPosts(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM posts`).Scan(&count)
	return count > 0, err
}

func randomTags() []string {
	n := rand.Intn(4) // 0–3 tags
	if n == 0 {
		return []string{}
	}

	shuffled := rand.Perm(len(availableTags))
	var tags []string
	for i := range n {
		tags = append(tags, availableTags[shuffled[i]])
	}
	return tags
}

func generateText(tags []string) string {
	if len(tags) == 0 {
		return "A quiet moment in nature."
	}

	switch len(tags) {
	case 1:
		return fmt.Sprintf("Enjoying some %s vibes.", tags[0])
	case 2:
		return fmt.Sprintf("A day with %s and %s.", tags[0], tags[1])
	default:
		return fmt.Sprintf("Exploring %s, %s and %s.", tags[0], tags[1], tags[2])
	}
}

func downloadImage(seed int, path string) error {
	url := fmt.Sprintf("https://picsum.photos/seed/%d/400/300", seed)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: %s", resp.Status)
	}

	file, err := os.Create(filepath.Clean(path))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

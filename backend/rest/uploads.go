package rest

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	database "example.com/backend/db"
	"example.com/backend/service/image"
	"example.com/backend/websocket"
	"github.com/gin-gonic/gin"
)

func abortWithError(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error": msg,
	})
}

func PostUploadsHandler(db *sql.DB, imageDir string, baseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			abortWithError(c, http.StatusBadRequest, "file required")
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			abortWithError(c, http.StatusBadRequest, "invalid file")
			return
		}
		defer file.Close()

		text := c.PostForm("text")
		rawTags := c.PostForm("tags")

		tagSlice := []string{}
		if rawTags != "" {
			for t := range strings.SplitSeq(rawTags, ",") {
				if trimmed := strings.TrimSpace(t); trimmed != "" {
					tagSlice = append(tagSlice, trimmed)
				}
			}
		}
		tagsJSON, err := json.Marshal(tagSlice)
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, "failed to process tags")
			return
		}

		normalized, ext, err := image.Normalize(file)
		if err != nil {
			abortWithError(c, http.StatusBadRequest, "invalid image")
			return
		}

		createdPost, err := database.CreatePost(db, imageDir, text, tagsJSON, normalized, ext)
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, "failed to create post")
			return
		}
		resp := createdPost.ToResponse(baseURL)
		websocket.Broadcast(resp)

		c.Status(http.StatusCreated)
	}
}

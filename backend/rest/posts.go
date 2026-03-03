package rest

import (
	"database/sql"
	"net/http"
	"strconv"

	database "example.com/backend/db"
	"example.com/backend/post"
	"github.com/gin-gonic/gin"
)

func GetPostsHandler(db *sql.DB, baseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := 1

		if p := c.Query("page"); p != "" {
			if n, err := strconv.Atoi(p); err == nil && n > 0 {
				page = n
			}
		}

		posts, err := database.GetPostsPage(db, page)
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, "failed to fetch posts")
			return
		}

		responses := make([]post.Response, 0, len(posts))

		for _, p := range posts {
			responses = append(responses, p.ToResponse(baseURL))
		}

		c.IndentedJSON(http.StatusOK, responses)
	}
}

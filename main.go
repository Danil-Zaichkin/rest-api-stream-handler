package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Repo interface {
	SaveOperation(ctx context.Context, p Package) (int, error)
}

type PackagePostBody struct {
	PackageID string     `json:"packageId"`
	Operation Opertation `json:"operation"`
	Value     int        `json:"value"`
}

func InitHandler(repo Repo) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/api/v1/package", func(c *gin.Context) {
		streamID := c.Query("streamId")
		jsonBody, _ := io.ReadAll(c.Request.Body)

		var body PackagePostBody
		_ = json.Unmarshal(jsonBody, &body)

		value, err := repo.SaveOperation(c, Package{
			Value:     body.Value,
			PackageID: body.PackageID,
			StreamID:  streamID,
			Op:        body.Operation,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": fmt.Sprintf("invalid request: %s", err.Error()),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"value": value,
		})
	})

	return r
}

func main() {
	repo := NewRepo()

	r := InitHandler(repo)
	r.Run()
}

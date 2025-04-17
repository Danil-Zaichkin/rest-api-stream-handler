package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Danil-Zaichkin/rest-api-stream-handler/internal/entity"
	"github.com/gin-gonic/gin"
)

type CalcUC interface {
	ApplyOperation(ctx context.Context, p entity.Package) (int, error)
}

type Handler struct {
	calcUC CalcUC
}

func New(calcUC CalcUC) *gin.Engine {
	h := &Handler{calcUC: calcUC}
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/api/v1/package", h.HandlePostPackage)

	return r
}

func (h *Handler) HandlePostPackage(ctx *gin.Context) {
	streamID := ctx.Query("streamId")
	jsonBody, _ := io.ReadAll(ctx.Request.Body)

	var body PackagePostBody
	_ = json.Unmarshal(jsonBody, &body)

	req := convertOperationRequest(body, streamID)

	res, err := h.calcUC.ApplyOperation(ctx, req)
	if err != nil {
		log.Printf("error: %v", err)
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": fmt.Sprintf("invalid request: %s", err.Error()),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"value": res,
	})
}

func convertOperationRequest(raw PackagePostBody, streamID string) entity.Package {
	return entity.Package{
		Value:     int(raw.Value),
		Op:        entity.Opertation(raw.Operation),
		PackageID: raw.PackageID,
		StreamID:  streamID,
	}
}

package api

import "github.com/Danil-Zaichkin/rest-api-stream-handler/internal/entity"

type PackagePostBody struct {
	PackageID string            `json:"packageId"`
	Operation entity.Opertation `json:"operation"`
	Value     int               `json:"value"`
}

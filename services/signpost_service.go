package services

import (
	"context"

	"github.com/99ridho/metrickit-backend/models"
)

type SignpostService interface {
	Store(ctx context.Context, signpost *models.AppSignpost, metadata *models.AppMetadata) (int64, error)
}

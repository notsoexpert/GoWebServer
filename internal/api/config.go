package api

import (
	"sync/atomic"

	"github.com/notsoexpert/gowebserver/internal/database"
)

type APIConfig struct {
	DBQueries      *database.Queries
	fileserverHits atomic.Int32
	Platform       string
	Secret         string
	PolkaKey       string
}

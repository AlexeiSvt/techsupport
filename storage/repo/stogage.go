package repo

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Storage struct {
	Pool *pgxpool.Pool
	mu   sync.RWMutex
}
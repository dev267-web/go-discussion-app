package health

import (
	"database/sql"
	"time"
)

type HealthStatus struct {
	Status    string            `json:"status"`   // "ok" or "fail"
	Checks    map[string]string `json:"checks"`   // e.g. { "database": "ok" }
	Timestamp time.Time         `json:"timestamp"`
}

type HealthService struct {
	db *sql.DB
}

func NewHealthService(db *sql.DB) *HealthService {
	return &HealthService{db: db}
}

func (hs *HealthService) CheckHealth() HealthStatus {
	checks := make(map[string]string)

	// Check DB connection
	if err := hs.db.Ping(); err != nil {
		checks["database"] = "fail"
		return HealthStatus{
			Status:    "fail",
			Checks:    checks,
			Timestamp: time.Now().UTC(),
		}
	}

	checks["database"] = "ok"

	return HealthStatus{
		Status:    "ok",
		Checks:    checks,
		Timestamp: time.Now().UTC(),
	}
}

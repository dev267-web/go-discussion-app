// repository.go 
package tag

import (
    "context"
    "database/sql"

    "go-discussion-app/models"
)

// TagRepository defines methods to interact with the tags table.
type TagRepository interface {
    // GetAll returns all tags in the database.
    GetAll(ctx context.Context) ([]models.Tag, error)
}

type repo struct {
    db *sql.DB
}

// NewRepository constructs a TagRepository backed by *sql.DB.
func NewRepository(db *sql.DB) TagRepository {
    return &repo{db: db}
}

func (r *repo) GetAll(ctx context.Context) ([]models.Tag, error) {
    const q = `
      SELECT id, name, created_at
      FROM tags
      ORDER BY name;
    `
    rows, err := r.db.QueryContext(ctx, q)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tags []models.Tag
    for rows.Next() {
        var t models.Tag
        if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
            return nil, err
        }
        tags = append(tags, t)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return tags, nil
}

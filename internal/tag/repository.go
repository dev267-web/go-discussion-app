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
    GetByName(ctx context.Context, name string) (*models.Tag, error)
    Create(ctx context.Context, name string) (int, error)
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
func (r *repo) GetByName(ctx context.Context, name string) (*models.Tag, error) {
    const q = `SELECT id, name, created_at FROM tags WHERE name = $1;`
    row := r.db.QueryRowContext(ctx, q, name)
    var t models.Tag
    if err := row.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &t, nil
}

func (r *repo) Create(ctx context.Context, name string) (int, error) {
    const q = `
        INSERT INTO tags (name, created_at)
        VALUES ($1, NOW())
        RETURNING id;
    `
    var id int
    err := r.db.QueryRowContext(ctx, q, name).Scan(&id)
    return id, err
}
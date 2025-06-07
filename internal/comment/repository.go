// repository.go 
package comment

import (
    "context"
    "database/sql"

    "go-discussion-app/models"
)

type Repository interface {
    Create(ctx context.Context, c *models.Comment) (int, error)
    ListByDiscussion(ctx context.Context, discussionID int) ([]models.Comment, error)
}

type repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
    return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, c *models.Comment) (int, error) {
    const q = `
      INSERT INTO comments (discussion_id, user_id, content, created_at)
      VALUES ($1, $2, $3, $4)
      RETURNING id;
    `
    var id int
    err := r.db.QueryRowContext(ctx, q,
        c.DiscussionID, c.UserID, c.Content, c.CreatedAt,
    ).Scan(&id)
    return id, err
}

func (r *repository) ListByDiscussion(ctx context.Context, discussionID int) ([]models.Comment, error) {
    const q = `
      SELECT id, discussion_id, user_id, content, created_at
      FROM comments
      WHERE discussion_id = $1
      ORDER BY created_at ASC;
    `
    rows, err := r.db.QueryContext(ctx, q, discussionID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []models.Comment
    for rows.Next() {
        var c models.Comment
        if err := rows.Scan(&c.ID, &c.DiscussionID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
            return nil, err
        }
        comments = append(comments, c)
    }
    return comments, rows.Err()
}

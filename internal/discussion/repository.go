// repository.go 
package discussion

import (
    "context"
    "database/sql"
    
    "time"

    "go-discussion-app/models"
)

type Repository interface {
    Create(ctx context.Context, d *models.Discussion) (int, error)
    GetAll(ctx context.Context) ([]models.Discussion, error)
    GetByID(ctx context.Context, id int) (*models.Discussion, error)
    Update(ctx context.Context, d *models.Discussion) error
    Delete(ctx context.Context, id int) error

    GetByUser(ctx context.Context, userID int) ([]models.Discussion, error)
    GetByTag(ctx context.Context, tag string) ([]models.Discussion, error)
    AddTags(ctx context.Context, discussionID int, tagIDs []int) error
}

type repo struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
    return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, d *models.Discussion) (int, error) {
    const q = `
      INSERT INTO discussions (user_id, title, content, scheduled_at, created_at, updated_at)
      VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;
    `
    var id int
    err := r.db.QueryRowContext(ctx, q,
        d.UserID, d.Title, d.Content, d.ScheduledAt, d.CreatedAt, d.UpdatedAt,
    ).Scan(&id)
    return id, err
}

func (r *repo) GetAll(ctx context.Context) ([]models.Discussion, error) {
    const q = `
      SELECT id, user_id, title, content, scheduled_at, created_at, updated_at
      FROM discussions
      ORDER BY created_at DESC;
    `
    rows, err := r.db.QueryContext(ctx, q)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ds []models.Discussion
    for rows.Next() {
        var d models.Discussion
        if err := rows.Scan(&d.ID, &d.UserID, &d.Title, &d.Content, &d.ScheduledAt, &d.CreatedAt, &d.UpdatedAt); err != nil {
            return nil, err
        }
        ds = append(ds, d)
    }
    return ds, rows.Err()
}

func (r *repo) GetByID(ctx context.Context, id int) (*models.Discussion, error) {
    const q = `
      SELECT id, user_id, title, content, scheduled_at, created_at, updated_at
      FROM discussions WHERE id=$1;
    `
    row := r.db.QueryRowContext(ctx, q, id)
    var d models.Discussion
    if err := row.Scan(&d.ID, &d.UserID, &d.Title, &d.Content, &d.ScheduledAt, &d.CreatedAt, &d.UpdatedAt); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &d, nil
}

func (r *repo) Update(ctx context.Context, d *models.Discussion) error {
    const q = `
      UPDATE discussions
      SET title=$1, content=$2, scheduled_at=$3, updated_at=$4
      WHERE id=$5;
    `
    _, err := r.db.ExecContext(ctx, q,
        d.Title, d.Content, d.ScheduledAt, time.Now().UTC(), d.ID,
    )
    return err
}

func (r *repo) Delete(ctx context.Context, id int) error {
    _, err := r.db.ExecContext(ctx, `DELETE FROM discussions WHERE id=$1`, id)
    return err
}

func (r *repo) GetByUser(ctx context.Context, userID int) ([]models.Discussion, error) {
    const q = `
      SELECT id, user_id, title, content, scheduled_at, created_at, updated_at
      FROM discussions WHERE user_id=$1 ORDER BY created_at DESC;
    `
    rows, err := r.db.QueryContext(ctx, q, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ds []models.Discussion
    for rows.Next() {
        var d models.Discussion
        if err := rows.Scan(&d.ID, &d.UserID, &d.Title, &d.Content, &d.ScheduledAt, &d.CreatedAt, &d.UpdatedAt); err != nil {
            return nil, err
        }
        ds = append(ds, d)
    }
    return ds, rows.Err()
}

func (r *repo) GetByTag(ctx context.Context, tag string) ([]models.Discussion, error) {
    const q = `
      SELECT d.id, d.user_id, d.title, d.content, d.scheduled_at, d.created_at, d.updated_at
      FROM discussions d
      JOIN discussion_tags dt ON d.id = dt.discussion_id
      JOIN tags t ON dt.tag_id = t.id
      WHERE t.name = $1
      ORDER BY d.created_at DESC;
    `
    rows, err := r.db.QueryContext(ctx, q, tag)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ds []models.Discussion
    for rows.Next() {
        var d models.Discussion
        if err := rows.Scan(&d.ID, &d.UserID, &d.Title, &d.Content, &d.ScheduledAt, &d.CreatedAt, &d.UpdatedAt); err != nil {
            return nil, err
        }
        ds = append(ds, d)
    }
    return ds, rows.Err()
}

func (r *repo) AddTags(ctx context.Context, discussionID int, tagIDs []int) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    stmt, err := tx.PrepareContext(ctx, `
      INSERT INTO discussion_tags (discussion_id, tag_id)
      VALUES ($1, $2) ON CONFLICT DO NOTHING;
    `)
    if err != nil {
        tx.Rollback()
        return err
    }
    defer stmt.Close()

    for _, tagID := range tagIDs {
        if _, err := stmt.ExecContext(ctx, discussionID, tagID); err != nil {
            tx.Rollback()
            return err
        }
    }
    return tx.Commit()
}

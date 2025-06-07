// repository.go 
// repository.go 
package user

import (
    "context"
    "database/sql"
    "time"

    "go-discussion-app/models"
)

type UserRepository interface {
    Create(ctx context.Context, u *models.User) (int, error)
    GetByID(ctx context.Context, id int) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, u *models.User) (sql.Result, error)
    Delete(ctx context.Context, id int) (sql.Result, error)
}

type userRepo struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) UserRepository {
    return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *models.User) (int, error) {
    const q = `
      INSERT INTO users
        (username, email, password_hash, full_name, bio, created_at, updated_at)
      VALUES ($1,$2,$3,$4,$5,$6,$7)
      RETURNING id;`
    var id int
    err := r.db.QueryRowContext(ctx, q,
        u.Username, u.Email, u.PasswordHash, u.FullName, u.Bio,
        u.CreatedAt, u.UpdatedAt,
    ).Scan(&id)
    return id, err
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
    const q = `
      SELECT id, username, email, password_hash, full_name, bio, created_at, updated_at
      FROM users WHERE id=$1;`
    row := r.db.QueryRowContext(ctx, q, id)
    var u models.User
    if err := row.Scan(
        &u.ID, &u.Username, &u.Email, &u.PasswordHash,
        &u.FullName, &u.Bio, &u.CreatedAt, &u.UpdatedAt,
    ); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &u, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    const q = `
      SELECT id, username, email, password_hash, full_name, bio, created_at, updated_at
      FROM users WHERE email=$1;`
    row := r.db.QueryRowContext(ctx, q, email)
    var u models.User
    if err := row.Scan(
        &u.ID, &u.Username, &u.Email, &u.PasswordHash,
        &u.FullName, &u.Bio, &u.CreatedAt, &u.UpdatedAt,
    ); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &u, nil
}

func (r *userRepo) Update(ctx context.Context, u *models.User) (sql.Result, error) {
    const q = `
      UPDATE users SET
        username=$1, email=$2, password_hash=$3, full_name=$4, bio=$5, updated_at=$6
      WHERE id=$7;`
    return r.db.ExecContext(ctx, q,
        u.Username, u.Email, u.PasswordHash, u.FullName, u.Bio,
        time.Now().UTC(), u.ID,
    )
}

func (r *userRepo) Delete(ctx context.Context, id int) (sql.Result, error) {
    const q = `DELETE FROM users WHERE id=$1;`
    return r.db.ExecContext(ctx, q, id)
}

package user

import (
	"context"
	"database/sql"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) *userRepositoryImpl {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) Insert(ctx context.Context, user entity.User) (err error) {
	statement := "INSERT INTO users (id, username, email, name, password, role) VALUES ($1, $2, $3, $4, $5, $6);"

	result, dbErr := u.db.ExecContext(ctx, statement, user.ID, user.Username, user.Email, user.Name, user.Password, user.Role)
	if dbErr != nil {
		err = repository.ErrDatabase
		return
	}

	count, dbErr := result.RowsAffected()
	if dbErr != nil || count < 1 {
		err = repository.ErrDatabase
		return
	}

	return
}

func (u *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (user entity.User, err error) {
	statement := `SELECT id,
  					     username,
       					 email,
       					 name,
       					 password,
       					 role,
       					 is_active,
       					 created_at,
       					 updated_at
							 FROM users
							 WHERE username = $1;`

	row := u.db.QueryRow(statement, username)

	switch dbErr := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	); dbErr {
	case sql.ErrNoRows:
		{
			err = repository.ErrRecordNotFound
			return
		}
	case nil:
		{
			return
		}
	default:
		{
			err = repository.ErrDatabase
			return
		}
	}
}

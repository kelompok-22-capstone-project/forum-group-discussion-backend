package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
	"github.com/lib/pq"
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
		switch e := dbErr.(type) {
		case *pq.Error:
			{
				if e.Code == "23505" {
					err = repository.ErrRecordAlreadyExists
					return
				}
				err = repository.ErrDatabase
				return
			}
		default:
			{
				log.Println(dbErr)
				err = repository.ErrDatabase
				return
			}
		}
	}

	count, dbErr := result.RowsAffected()
	if dbErr != nil {
		err = repository.ErrDatabase
		return
	}

	if count < 1 {
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

	row := u.db.QueryRowContext(ctx, statement, username)

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
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
	}
}

func (u *userRepositoryImpl) FindAllWithStatusAndPagination(
	ctx context.Context,
	accessorUserID string,
	orderBy entity.UserOrderBy,
	userStatus entity.UserStatus,
	pageInfo entity.PageInfo,
) (pagination entity.Pagination[entity.User], err error) {
	var userOrderBy string
	if orderBy == entity.Ranking {
		userOrderBy = "total_follower"
	} else {
		userOrderBy = "u.created_at"
	}

	statement := fmt.Sprintf(`SELECT u.id,
       u.username,
       u.email,
       u.name,
       u.role,
       u.is_active,
       u.created_at,
       u.updated_at,
       (SELECT count(t.id) FROM threads t WHERE t.creator_id = u.id)           AS total_thread,
       (SELECT count(uf.id) FROM user_follows uf WHERE uf.following_id = u.id) AS total_follower,
       (SELECT CASE WHEN count(uf.id) > 0 THEN true ELSE false END
        FROM user_follows uf
        WHERE uf.user_id = $1
          AND uf.following_id = u.id)                                          AS is_followed
FROM users u
WHERE is_active = $2 AND u.role = 'user'
ORDER BY %s
OFFSET $3 LIMIT $4;`, userOrderBy)

	rows, dbErr := u.db.QueryContext(ctx, statement, accessorUserID, userStatus, (pageInfo.Page-1)*pageInfo.Limit, pageInfo.Limit*1)
	if dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	defer func(rows *sql.Rows) {
		if dbErr := rows.Close(); dbErr != nil {
			log.Println(dbErr)
		}
	}(rows)

	pagination.List = make([]entity.User, 0)
	for rows.Next() {
		var user entity.User
		if dbErr := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Name,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.TotalThread,
			&user.TotalFollower,
			&user.IsFollowed,
		); dbErr != nil {
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
		pagination.List = append(pagination.List, user)
	}

	if dbErr := rows.Err(); dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	countStatement := "SELECT count(u.id) FROM users u WHERE is_active = $1 AND u.role = 'user';"

	row := u.db.QueryRowContext(ctx, countStatement, userStatus)

	var count uint
	switch dbErr := row.Scan(&count); dbErr {
	case sql.ErrNoRows:
		{
			err = repository.ErrRecordNotFound
			return
		}
	case nil:
		{
			pagination.PageInfo.Limit = pageInfo.Limit
			pagination.PageInfo.Page = pageInfo.Page
			pagination.PageInfo.PageTotal = uint(math.Ceil(float64(count) / float64(pageInfo.Limit)))
			pagination.PageInfo.Total = count
			return
		}
	default:
		{
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
	}
}

func (u *userRepositoryImpl) FindByUsernameWithAccessor(
	ctx context.Context,
	accessorUserID string,
	username string,
) (user entity.User, err error) {
	statement := `SELECT u.id,
       u.username,
       u.email,
       u.name,
       u.role,
       u.is_active,
       u.created_at,
       u.updated_at,
       (SELECT count(t.id) FROM threads t WHERE t.creator_id = u.id)           AS total_thread,
       (SELECT count(uf.id) FROM user_follows uf WHERE uf.following_id = u.id) AS total_follower,
       (SELECT CASE WHEN count(uf.id) > 0 THEN true ELSE false END
        FROM user_follows uf
        WHERE uf.user_id = $1
          AND uf.following_id = u.id)                                          AS is_followed
FROM users u
WHERE u.username = $2
  AND role = 'user';`

	row := u.db.QueryRowContext(ctx, statement, accessorUserID, username)

	switch dbErr := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.TotalThread,
		&user.TotalFollower,
		&user.IsFollowed,
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
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
	}
}

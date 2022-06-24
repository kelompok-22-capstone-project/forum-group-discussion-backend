package admin

import (
	"context"
	"database/sql"
	"log"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
)

type adminRepositoryImpl struct {
	db *sql.DB
}

func NewAdminRepositoryImpl(db *sql.DB) *adminRepositoryImpl {
	return &adminRepositoryImpl{db: db}
}

func (a *adminRepositoryImpl) FindDashboardInfo(ctx context.Context) (info entity.DashboardInfo, err error) {
	statement := `SELECT (SELECT count(u.id)
        FROM users u)        AS total_user,
       (SELECT count(t.id)
        FROM threads t)      AS total_thread,
       (SELECT count(m.id)
        FROM moderators m)   AS total_moderator,
       (SELECT count(r.id)
        FROM user_banneds r) AS total_report;`

	row := a.db.QueryRowContext(ctx, statement)

	switch dbErr := row.Scan(
		&info.TotalUser,
		&info.TotalThread,
		&info.TotalModerator,
		&info.TotalReport,
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

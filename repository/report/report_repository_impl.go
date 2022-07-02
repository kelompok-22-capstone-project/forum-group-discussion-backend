package report

import (
	"context"
	"database/sql"
	"log"
	"math"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
)

type reportRepositoryImpl struct {
	db *sql.DB
}

func NewReportRepositoryImpl(db *sql.DB) *reportRepositoryImpl {
	return &reportRepositoryImpl{db: db}
}

func (r *reportRepositoryImpl) GetReportsWithPagination(
	ctx context.Context,
	pageInfo entity.PageInfo,
	reportStatus entity.ReportStatus,
) (pagination entity.Pagination[entity.UserBanned], err error) {
	statement := `SELECT t.id           AS report_id,
       t.moderator_id AS moderator_id,
       m.user_id      AS moderator_user_id,
       u1.username    AS moderator_username,
       u1.name        AS moderator_name,
       m.thread_id    AS moderator_thread_id,
       m.created_at   AS moderator_created_at,
       m.updated_at   AS moderator_updated_at,
       t.user_id      AS reported_user_id,
       u2.username    AS reported_username,
       u2.name        AS reported_name,
       th.id          AS thread_id,
       th.title       AS thread_title,
       c.id           AS comment_id,
       c.comment      AS comment,
       c.created_at   AS comment_created_at,
       c.updated_at   AS comment_updated_at,
       t.reason       AS reason,
       t.status       AS status,
       t.created_at   AS created_at,
       t.updated_at   AS updated_at
FROM user_banneds t
         INNER JOIN moderators m on m.id = t.moderator_id
         INNER JOIN users u1 on m.user_id = u1.id
         INNER JOIN users u2 on t.user_id = u2.id
         INNER JOIN comments c on c.id = t.comment_id
         INNER JOIN threads th on c.thread_id = th.id
WHERE t.status = $1
OFFSET $2 LIMIT $3;`

	var status string
	if reportStatus == entity.Accepted {
		status = "accepted"
	} else {
		status = "review"
	}

	rows, dbErr := r.db.QueryContext(ctx, statement, status, (pageInfo.Page-1)*pageInfo.Limit, pageInfo.Limit*1)
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

	pagination.List = make([]entity.UserBanned, 0)
	for rows.Next() {
		var userBanned entity.UserBanned
		if dbErr := rows.Scan(
			&userBanned.ID,
			&userBanned.Moderator.ID,
			&userBanned.Moderator.User.ID,
			&userBanned.Moderator.User.Username,
			&userBanned.Moderator.User.Name,
			&userBanned.Moderator.ThreadID,
			&userBanned.Moderator.CreatedAt,
			&userBanned.Moderator.UpdatedAt,
			&userBanned.User.ID,
			&userBanned.User.Username,
			&userBanned.User.Name,
			&userBanned.Thread.ID,
			&userBanned.Thread.Title,
			&userBanned.Comment.ID,
			&userBanned.Comment.Comment,
			&userBanned.Comment.CreatedAt,
			&userBanned.Comment.UpdatedAt,
			&userBanned.Reason,
			&userBanned.Status,
			&userBanned.CreatedAt,
			&userBanned.UpdatedAt,
		); dbErr != nil {
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
		pagination.List = append(pagination.List, userBanned)
	}

	if dbErr := rows.Err(); dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	countStatement := "SELECT count(t.id) FROM user_banneds t WHERE status = $1;"

	row := r.db.QueryRowContext(ctx, countStatement, status)

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

func (r *reportRepositoryImpl) Insert(
	ctx context.Context,
	ID,
	moderatorID,
	userID,
	commentID,
	reason string,
) (err error) {
	statement := "INSERT INTO user_banneds (id, moderator_id, user_id, comment_id, reason) VALUES ($1, $2, $3, $4, $5);"

	result, dbErr := r.db.ExecContext(ctx, statement, ID, moderatorID, userID, commentID, reason)
	if dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	count, dbErr := result.RowsAffected()
	if dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	if count < 1 {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	return
}

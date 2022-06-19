package thread

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
)

type threadRepositoryImpl struct {
	db *sql.DB
}

func NewThreadRepositoryImpl(db *sql.DB) *threadRepositoryImpl {
	return &threadRepositoryImpl{db: db}
}

func (t *threadRepositoryImpl) Insert(ctx context.Context, thread entity.Thread) (err error) {
	statement := "INSERT INTO threads(id, title, description, creator_id, category_id) VALUES ($1, $2, $3, $4, $5);"

	result, dbErr := t.db.ExecContext(ctx, statement, thread.ID, thread.Title, thread.Description, thread.Creator.ID, thread.Category.ID)
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

func (t *threadRepositoryImpl) FindAllWithQueryAndPagination(
	ctx context.Context,
	accessorUserID string,
	query string,
	pageInfo entity.PageInfo,
) (pagination entity.Pagination[entity.Thread], err error) {
	var rows *sql.Rows
	var dbErr error

	if query == "" {
		statement := `SELECT t.id                                                                        as thread_id,
       t.title                                                                                     as thread_title,
       t.description                                                                               as thread_description,
       t.total_viewer,
       t.creator_id,
       u.username                                                                                  as creator_username,
       u.email                                                                                     as creator_email,
       u.name                                                                                      as creator_name,
       u.password                                                                                  as creator_password,
       u.role                                                                                      as creator_rule,
       u.is_active                                                                                 as creator_is_active,
       u.created_at                                                                                as creator_created_at,
       u.updated_at                                                                                as creator_updated_at,
       t.category_id,
       c.name                                                                                      as category_name,
       c.description                                                                               as category_description,
       c.created_at                                                                                as category_created_at,
       c.updated_at                                                                                as category_updated_at,
       t.created_at                                                                                as thread_created_at,
       t.updated_at                                                                                as thread_updated_at,
       (SELECT CASE WHEN count(likes.id) > 0 THEN true ELSE false END
        FROM likes
        WHERE likes.user_id = $1
          AND likes.thread_id = t.id)                                                              as is_liked,
       (SELECT CASE WHEN count(thread_follows.id) > 0 THEN true ELSE false END
        FROM thread_follows
        WHERE thread_follows.user_id = $1
          AND thread_follows.thread_id = t.id)                                                     as is_followed,
       (SELECT count(thread_follows.id) FROM thread_follows WHERE thread_follows.thread_id = t.id) as total_follower,
       (SELECT count(likes.id) FROM likes WHERE likes.thread_id = t.id)                            as total_like,
       (SELECT count(comments.id) FROM comments WHERE comments.thread_id = t.id)                   as total_comment
FROM threads as t
         INNER JOIN categories c
                    on c.id = t.category_id
         INNER JOIN users u on t.creator_id = u.id
ORDER BY t.created_at DESC
OFFSET $2 LIMIT $3;`

		rows, dbErr = t.db.QueryContext(ctx, statement, accessorUserID, (pageInfo.Page-1)*pageInfo.Limit, pageInfo.Limit*1)
		if dbErr != nil {
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
	} else {
		statement := `SELECT t.id                                                                        as thread_id,
       t.title                                                                                     as thread_title,
       t.description                                                                               as thread_description,
       t.total_viewer,
       t.creator_id,
       u.username                                                                                  as creator_username,
       u.email                                                                                     as creator_email,
       u.name                                                                                      as creator_name,
       u.password                                                                                  as creator_password,
       u.role                                                                                      as creator_rule,
       u.is_active                                                                                 as creator_is_active,
       u.created_at                                                                                as creator_created_at,
       u.updated_at                                                                                as creator_updated_at,
       t.category_id,
       c.name                                                                                      as category_name,
       c.description                                                                               as category_description,
       c.created_at                                                                                as category_created_at,
       c.updated_at                                                                                as category_updated_at,
       t.created_at                                                                                as thread_created_at,
       t.updated_at                                                                                as thread_updated_at,
       (SELECT CASE WHEN count(likes.id) > 0 THEN true ELSE false END
        FROM likes
        WHERE likes.user_id = $1
          AND likes.thread_id = t.id)                                                              as is_liked,
       (SELECT CASE WHEN count(thread_follows.id) > 0 THEN true ELSE false END
        FROM thread_follows
        WHERE thread_follows.user_id = $1
          AND thread_follows.thread_id = t.id)                                                     as is_followed,
       (SELECT count(thread_follows.id) FROM thread_follows WHERE thread_follows.thread_id = t.id) as total_follower,
       (SELECT count(likes.id) FROM likes WHERE likes.thread_id = t.id)                            as total_like,
       (SELECT count(comments.id) FROM comments WHERE comments.thread_id = t.id)                   as total_comment
FROM threads as t
         INNER JOIN categories c
                    on c.id = t.category_id
         INNER JOIN users u on t.creator_id = u.id
WHERE t.title ILIKE $4
ORDER BY t.created_at DESC
OFFSET $2 LIMIT $3;`

		rows, dbErr = t.db.QueryContext(ctx, statement, accessorUserID, (pageInfo.Page-1)*pageInfo.Limit, pageInfo.Limit*1, fmt.Sprintf("%%%s%%", query))
		if dbErr != nil {
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
	}

	defer func(rows *sql.Rows) {
		if dbErr := rows.Close(); dbErr != nil {
			log.Println(dbErr)
		}
	}(rows)

	pagination.List = make([]entity.Thread, 0)
	for rows.Next() {
		var thread entity.Thread
		if dbErr := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Description,
			&thread.TotalViewer,
			&thread.Creator.ID,
			&thread.Creator.Username,
			&thread.Creator.Email,
			&thread.Creator.Name,
			&thread.Creator.Password,
			&thread.Creator.Role,
			&thread.Creator.IsActive,
			&thread.Creator.CreatedAt,
			&thread.Creator.UpdatedAt,
			&thread.Category.ID,
			&thread.Category.Name,
			&thread.Category.Description,
			&thread.Category.CreatedAt,
			&thread.Category.UpdatedAt,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.IsLiked,
			&thread.IsFollowed,
			&thread.TotalFollower,
			&thread.TotalLike,
			&thread.TotalComment,
		); dbErr != nil {
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
		pagination.List = append(pagination.List, thread)
	}

	if dbErr := rows.Err(); dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	countStatement := "SELECT count(threads.id) FROM threads WHERE title ILIKE $1;"

	row := t.db.QueryRowContext(ctx, countStatement, fmt.Sprintf("%%%s%%", query))

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

func (t *threadRepositoryImpl) FindAllByCategoryIDWithPagination(
	ctx context.Context,
	accessorUserID string,
	categoryID string,
	pageInfo entity.PageInfo,
) (pagination entity.Pagination[entity.Thread], err error) {
	statement := `SELECT t.id                                                                        as thread_id,
       t.title                                                                                     as thread_title,
       t.description                                                                               as thread_description,
       t.total_viewer,
       t.creator_id,
       u.username                                                                                  as creator_username,
       u.email                                                                                     as creator_email,
       u.name                                                                                      as creator_name,
       u.password                                                                                  as creator_password,
       u.role                                                                                      as creator_rule,
       u.is_active                                                                                 as creator_is_active,
       u.created_at                                                                                as creator_created_at,
       u.updated_at                                                                                as creator_updated_at,
       t.category_id,
       c.name                                                                                      as category_name,
       c.description                                                                               as category_description,
       c.created_at                                                                                as category_created_at,
       c.updated_at                                                                                as category_updated_at,
       t.created_at                                                                                as thread_created_at,
       t.updated_at                                                                                as thread_updated_at,
       (SELECT CASE WHEN count(likes.id) > 0 THEN true ELSE false END
        FROM likes
        WHERE likes.user_id = $1
          AND likes.thread_id = t.id)                                                              as is_liked,
       (SELECT CASE WHEN count(thread_follows.id) > 0 THEN true ELSE false END
        FROM thread_follows
        WHERE thread_follows.user_id = $1
          AND thread_follows.thread_id = t.id)                                                     as is_followed,
       (SELECT count(thread_follows.id) FROM thread_follows WHERE thread_follows.thread_id = t.id) as total_follower,
       (SELECT count(likes.id) FROM likes WHERE likes.thread_id = t.id)                            as total_like,
       (SELECT count(comments.id) FROM comments WHERE comments.thread_id = t.id)                   as total_comment
FROM threads as t
         INNER JOIN categories c
                    on c.id = t.category_id
         INNER JOIN users u on t.creator_id = u.id
WHERE t.category_id = $4
ORDER BY t.created_at DESC
OFFSET $2 LIMIT $3;`

	rows, dbErr := t.db.QueryContext(ctx, statement, accessorUserID, (pageInfo.Page-1)*pageInfo.Limit, pageInfo.Limit*1, categoryID)
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

	pagination.List = make([]entity.Thread, 0)
	for rows.Next() {
		var thread entity.Thread
		if dbErr := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Description,
			&thread.TotalViewer,
			&thread.Creator.ID,
			&thread.Creator.Username,
			&thread.Creator.Email,
			&thread.Creator.Name,
			&thread.Creator.Password,
			&thread.Creator.Role,
			&thread.Creator.IsActive,
			&thread.Creator.CreatedAt,
			&thread.Creator.UpdatedAt,
			&thread.Category.ID,
			&thread.Category.Name,
			&thread.Category.Description,
			&thread.Category.CreatedAt,
			&thread.Category.UpdatedAt,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.IsLiked,
			&thread.IsFollowed,
			&thread.TotalFollower,
			&thread.TotalLike,
			&thread.TotalComment,
		); dbErr != nil {
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
		pagination.List = append(pagination.List, thread)
	}

	if dbErr := rows.Err(); dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	countStatement := "SELECT count(threads.id) FROM threads WHERE category_id = $1;"

	row := t.db.QueryRowContext(ctx, countStatement, categoryID)

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

func (t *threadRepositoryImpl) FindAllByUserIDWithPagination(
	ctx context.Context,
	accessorUserID string,
	userID string,
	pageInfo entity.PageInfo,
) (pagination entity.Pagination[entity.Thread], err error) {
	statement := `SELECT t.id                                                                        as thread_id,
       t.title                                                                                     as thread_title,
       t.description                                                                               as thread_description,
       t.total_viewer,
       t.creator_id,
       u.username                                                                                  as creator_username,
       u.email                                                                                     as creator_email,
       u.name                                                                                      as creator_name,
       u.password                                                                                  as creator_password,
       u.role                                                                                      as creator_rule,
       u.is_active                                                                                 as creator_is_active,
       u.created_at                                                                                as creator_created_at,
       u.updated_at                                                                                as creator_updated_at,
       t.category_id,
       c.name                                                                                      as category_name,
       c.description                                                                               as category_description,
       c.created_at                                                                                as category_created_at,
       c.updated_at                                                                                as category_updated_at,
       t.created_at                                                                                as thread_created_at,
       t.updated_at                                                                                as thread_updated_at,
       (SELECT CASE WHEN count(likes.id) > 0 THEN true ELSE false END
        FROM likes
        WHERE likes.user_id = $1
          AND likes.thread_id = t.id)                                                              as is_liked,
       (SELECT CASE WHEN count(thread_follows.id) > 0 THEN true ELSE false END
        FROM thread_follows
        WHERE thread_follows.user_id = $1
          AND thread_follows.thread_id = t.id)                                                     as is_followed,
       (SELECT count(thread_follows.id) FROM thread_follows WHERE thread_follows.thread_id = t.id) as total_follower,
       (SELECT count(likes.id) FROM likes WHERE likes.thread_id = t.id)                            as total_like,
       (SELECT count(comments.id) FROM comments WHERE comments.thread_id = t.id)                   as total_comment
FROM threads as t
         INNER JOIN categories c
                    on c.id = t.category_id
         INNER JOIN users u on t.creator_id = u.id
WHERE t.creator_id = $4
ORDER BY t.created_at DESC
OFFSET $2 LIMIT $3;`

	rows, dbErr := t.db.QueryContext(ctx, statement, accessorUserID, (pageInfo.Page-1)*pageInfo.Limit, pageInfo.Limit*1, userID)
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

	pagination.List = make([]entity.Thread, 0)
	for rows.Next() {
		var thread entity.Thread
		if dbErr := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Description,
			&thread.TotalViewer,
			&thread.Creator.ID,
			&thread.Creator.Username,
			&thread.Creator.Email,
			&thread.Creator.Name,
			&thread.Creator.Password,
			&thread.Creator.Role,
			&thread.Creator.IsActive,
			&thread.Creator.CreatedAt,
			&thread.Creator.UpdatedAt,
			&thread.Category.ID,
			&thread.Category.Name,
			&thread.Category.Description,
			&thread.Category.CreatedAt,
			&thread.Category.UpdatedAt,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.IsLiked,
			&thread.IsFollowed,
			&thread.TotalFollower,
			&thread.TotalLike,
			&thread.TotalComment,
		); dbErr != nil {
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
		pagination.List = append(pagination.List, thread)
	}

	if dbErr := rows.Err(); dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	countStatement := "SELECT count(id) FROM threads WHERE creator_id = $1;"

	row := t.db.QueryRowContext(ctx, countStatement, userID)

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

func (t *threadRepositoryImpl) FindByID(
	ctx context.Context,
	accessorUserID string,
	ID string,
) (thread entity.Thread, err error) {
	statement := `SELECT t.id                                                                        as thread_id,
       t.title                                                                                     as thread_title,
       t.description                                                                               as thread_description,
       t.total_viewer,
       t.creator_id,
       u.username                                                                                  as creator_username,
       u.email                                                                                     as creator_email,
       u.name                                                                                      as creator_name,
       u.password                                                                                  as creator_password,
       u.role                                                                                      as creator_rule,
       u.is_active                                                                                 as creator_is_active,
       u.created_at                                                                                as creator_created_at,
       u.updated_at                                                                                as creator_updated_at,
       t.category_id,
       c.name                                                                                      as category_name,
       c.description                                                                               as category_description,
       c.created_at                                                                                as category_created_at,
       c.updated_at                                                                                as category_updated_at,
       t.created_at                                                                                as thread_created_at,
       t.updated_at                                                                                as thread_updated_at,
       (SELECT CASE WHEN count(likes.id) > 0 THEN true ELSE false END
        FROM likes
        WHERE likes.user_id = $1
          AND likes.thread_id = t.id)                                                              as is_liked,
       (SELECT CASE WHEN count(thread_follows.id) > 0 THEN true ELSE false END
        FROM thread_follows
        WHERE thread_follows.user_id = $1
          AND thread_follows.thread_id = t.id)                                                     as is_followed,
       (SELECT count(thread_follows.id) FROM thread_follows WHERE thread_follows.thread_id = t.id) as total_follower,
       (SELECT count(likes.id) FROM likes WHERE likes.thread_id = t.id)                            as total_like,
       (SELECT count(comments.id) FROM comments WHERE comments.thread_id = t.id)                   as total_comment
FROM threads as t
         INNER JOIN categories c
                    on c.id = t.category_id
         INNER JOIN users u on t.creator_id = u.id
WHERE t.id = $2;`

	row := t.db.QueryRowContext(ctx, statement, accessorUserID, ID)

	switch dbErr := row.Scan(
		&thread.ID,
		&thread.Title,
		&thread.Description,
		&thread.TotalViewer,
		&thread.Creator.ID,
		&thread.Creator.Username,
		&thread.Creator.Email,
		&thread.Creator.Name,
		&thread.Creator.Password,
		&thread.Creator.Role,
		&thread.Creator.IsActive,
		&thread.Creator.CreatedAt,
		&thread.Creator.UpdatedAt,
		&thread.Category.ID,
		&thread.Category.Name,
		&thread.Category.Description,
		&thread.Category.CreatedAt,
		&thread.Category.UpdatedAt,
		&thread.CreatedAt,
		&thread.UpdatedAt,
		&thread.IsLiked,
		&thread.IsFollowed,
		&thread.TotalFollower,
		&thread.TotalLike,
		&thread.TotalComment,
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
			log.Println(repository.ErrDatabase)
			err = repository.ErrDatabase
			return
		}
	}
}

func (t *threadRepositoryImpl) Update(
	ctx context.Context,
	ID string,
	thread entity.Thread,
) (err error) {
	statement := `UPDATE threads
SET title       = $2,
    description = $3,
    category_id = $4,
    updated_at  = current_timestamp
WHERE id = $1;`

	result, dbErr := t.db.ExecContext(ctx, statement, ID, thread.Title, thread.Description, thread.Category.ID)
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
		err = repository.ErrRecordNotFound
		return
	}

	return
}

func (t *threadRepositoryImpl) Delete(
	ctx context.Context,
	ID string,
) (err error) {
	statement := "DELETE FROM threads WHERE id = $1;"

	result, dbErr := t.db.ExecContext(ctx, statement, ID)
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
		err = repository.ErrRecordNotFound
		return
	}

	return
}

func (t *threadRepositoryImpl) FindAllModeratorByThreadID(
	ctx context.Context,
	threadID string,
) (moderators []entity.Moderator, err error) {
	statement := `SELECT m.id as moderator_id,
       u.id         as user_id,
       m.thread_id  as thread_id,
       m.created_at as moderator_created_at,
       m.updated_at as moderator_updated_at,
       u.username   as user_username,
       u.email      as user_email,
       u.name       as user_name,
       u.role       as user_role,
       u.is_active  as user_is_active,
       u.created_at as user_created_at,
       u.updated_at as user_updated_at
FROM moderators as m
         INNER JOIN users u on u.id = m.user_id
WHERE m.thread_id = $1;`

	rows, dbErr := t.db.QueryContext(ctx, statement, threadID)
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

	moderators = make([]entity.Moderator, 0)
	for rows.Next() {
		var moderator entity.Moderator
		if dbErr := rows.Scan(
			&moderator.ID,
			&moderator.User.ID,
			&moderator.ThreadID,
			&moderator.CreatedAt,
			&moderator.UpdatedAt,
			&moderator.User.Username,
			&moderator.User.Email,
			&moderator.User.Name,
			&moderator.User.Role,
			&moderator.User.IsActive,
			&moderator.User.CreatedAt,
			&moderator.User.UpdatedAt,
		); dbErr != nil {
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
		moderators = append(moderators, moderator)
	}

	if dbErr := rows.Err(); dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	return
}

func (t *threadRepositoryImpl) FindAllCommentByThreadID(
	ctx context.Context,
	threadID string,
	pageInfo entity.PageInfo,
) (pagination entity.Pagination[entity.Comment], err error) {
	statement := `SELECT c.id as comment_id,
       c.user_id   as user_id,
       t.id        as thread_id,
       c.comment,
       c.created_at,
       c.updated_at,
       u.username  as user_username,
       u.email     as user_email,
       u.name      as user_name,
       u.role      as user_role,
       u.is_active as user_is_active
FROM comments c
         INNER JOIN users u on c.user_id = u.id
         INNER JOIN threads t on t.id = c.thread_id
WHERE c.thread_id = $1
ORDER BY c.created_at DESC
OFFSET $2 LIMIT $3;`

	rows, dbErr := t.db.QueryContext(ctx, statement, threadID, (pageInfo.Page-1)*pageInfo.Limit, pageInfo.Limit*1)
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

	pagination.List = make([]entity.Comment, 0)
	for rows.Next() {
		var comment entity.Comment
		if dbErr := rows.Scan(
			&comment.ID,
			&comment.User.ID,
			&comment.Thread.ID,
			&comment.Comment,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.User.Username,
			&comment.User.Email,
			&comment.User.Name,
			&comment.User.Role,
			&comment.User.IsActive,
		); dbErr != nil {
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
		pagination.List = append(pagination.List, comment)
	}

	if dbErr := rows.Err(); dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	countStatement := "SELECT count(c.id) FROM comments c WHERE c.thread_id = $1;"

	row := t.db.QueryRowContext(ctx, countStatement, threadID)

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

func (t *threadRepositoryImpl) InsertComment(
	ctx context.Context,
	comment entity.Comment,
) (err error) {
	statement := "INSERT INTO comments(id, user_id, thread_id, comment) VALUES ($1, $2, $3, $4);"

	result, dbErr := t.db.ExecContext(ctx, statement, comment.ID, comment.User.ID, comment.Thread.ID, comment.Comment)
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
		err = repository.ErrDatabase
		return
	}

	return
}

func (t *threadRepositoryImpl) InsertFollowThread(
	ctx context.Context,
	threadFollow entity.ThreadFollow,
) (err error) {
	statement := "INSERT INTO thread_follows(id, user_id, thread_id) VALUES ($1, $2, $3);"

	result, dbErr := t.db.ExecContext(ctx, statement, threadFollow.ID, threadFollow.User.ID, threadFollow.Thread.ID)
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
		err = repository.ErrDatabase
		return
	}

	return
}

func (t *threadRepositoryImpl) DeleteFollowThread(
	ctx context.Context,
	threadFollow entity.ThreadFollow,
) (err error) {
	statement := "DELETE FROM thread_follows WHERE user_id = $1  AND thread_id = $2;"

	result, dbErr := t.db.ExecContext(ctx, statement, threadFollow.User.ID, threadFollow.Thread.ID)
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
		err = repository.ErrRecordNotFound
		return
	}

	return
}

func (t *threadRepositoryImpl) InsertLike(
	ctx context.Context,
	like entity.Like,
) (err error) {
	statement := "INSERT INTO likes(id, user_id, thread_id) VALUES ($1, $2, $3);"

	result, dbErr := t.db.ExecContext(ctx, statement, like.ID, like.User.ID, like.Thread.ID)
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
		err = repository.ErrDatabase
		return
	}

	return
}

func (t *threadRepositoryImpl) DeleteLike(
	ctx context.Context,
	like entity.Like,
) (err error) {
	statement := "DELETE FROM likes WHERE user_id = $1 AND thread_id = $2;"

	result, dbErr := t.db.ExecContext(ctx, statement, like.User.ID, like.Thread.ID)
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
		err = repository.ErrRecordNotFound
		return
	}

	return
}

func (t *threadRepositoryImpl) InsertModerator(
	ctx context.Context,
	moderator entity.Moderator,
) (err error) {
	statement := "INSERT INTO moderators(id, user_id, thread_id) VALUES ($1, $2, $3);"

	result, dbErr := t.db.ExecContext(ctx, statement, moderator.ID, moderator.User.ID, moderator.ThreadID)
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
		err = repository.ErrDatabase
		return
	}

	return
}

func (t *threadRepositoryImpl) DeleteModerator(
	ctx context.Context,
	moderator entity.Moderator,
) (err error) {
	statement := "DELETE FROM moderators WHERE user_id = $1 AND thread_id = $2;"

	result, dbErr := t.db.ExecContext(ctx, statement, moderator.User.ID, moderator.ThreadID)
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
		err = repository.ErrRecordNotFound
		return
	}

	return
}

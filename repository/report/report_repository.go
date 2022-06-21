package report

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
)

type ReportRepository interface {
	GetReportsWithPagination(
		ctx context.Context,
		pageInfo entity.PageInfo,
		reportStatus entity.ReportStatus,
	) (pagination entity.Pagination[entity.UserBanned], err error)

	Insert(
		ctx context.Context,
		ID,
		moderatorID,
		userID,
		reason string,
	) (err error)
}

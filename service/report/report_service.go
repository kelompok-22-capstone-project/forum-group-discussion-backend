package report

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
)

type ReportService interface {
	GetAll(
		ctx context.Context,
		accessorRole,
		status string,
		page uint,
		limit uint,
	) (r response.Pagination[response.Report], err error)

	Create(
		ctx context.Context,
		accessorUserID string,
		p payload.CreateReport,
	) (id string, err error)
}

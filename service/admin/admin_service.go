package admin

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
)

type AdminService interface {
	GetDashboardInfo(ctx context.Context, accessorRole string) (r response.DashboardInfo, err error)
}

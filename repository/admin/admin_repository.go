package admin

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
)

type AdminRepository interface {
	FindDashboardInfo(ctx context.Context) (info entity.DashboardInfo, err error)
}

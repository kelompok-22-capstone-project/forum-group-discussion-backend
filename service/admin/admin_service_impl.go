package admin

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/admin"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
)

type adminServiceImpl struct {
	adminRepository admin.AdminRepository
}

func NewAdminServiceImpl(adminRepository admin.AdminRepository) *adminServiceImpl {
	return &adminServiceImpl{adminRepository: adminRepository}
}

func (a *adminServiceImpl) GetDashboardInfo(ctx context.Context, accessorRole string) (r response.DashboardInfo, err error) {
	if accessorRole != "admin" {
		err = service.ErrAccessForbidden
		return
	}

	info, repoErr := a.adminRepository.FindDashboardInfo(ctx)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	r.TotalUser = info.TotalUser
	r.TotalReport = info.TotalReport
	r.TotalModerator = info.TotalModerator
	r.TotalThread = info.TotalThread

	return
}

package admin

import (
	"context"
	"fmt"
	"testing"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/admin/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDashboardInfo(t *testing.T) {
	mockAdminRepo := &mocks.AdminRepository{}
	var adminService AdminService = NewAdminServiceImpl(mockAdminRepo)

	testCases := []struct {
		name                  string
		inputAccessorRole     string
		expectedError         error
		expectedDashboardInfo response.DashboardInfo
		mockBehaviour         func()
	}{
		{
			name:                  "it should return service.ErrAccessForbidden, if accessorRole is not admin",
			inputAccessorRole:     "user",
			expectedError:         service.ErrAccessForbidden,
			expectedDashboardInfo: response.DashboardInfo{},
			mockBehaviour:         func() {},
		},
		{
			name:                  "it should return service.ErrRepository, when admin repository return repository.ErrrDatabase",
			inputAccessorRole:     "admin",
			expectedError:         service.ErrRepository,
			expectedDashboardInfo: response.DashboardInfo{},
			mockBehaviour: func() {
				mockAdminRepo.On(
					"FindDashboardInfo",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
				).Return(
					func(ctx context.Context) entity.DashboardInfo {
						return entity.DashboardInfo{}
					},
					func(ctx context.Context) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:              "it should return valid dashboard info, when no error is returned",
			inputAccessorRole: "admin",
			expectedError:     nil,
			expectedDashboardInfo: response.DashboardInfo{
				TotalUser:      1592,
				TotalThread:    95231,
				TotalModerator: 442139842,
				TotalReport:    125,
			},
			mockBehaviour: func() {
				mockAdminRepo.On(
					"FindDashboardInfo",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
				).Return(
					func(ctx context.Context) entity.DashboardInfo {
						return entity.DashboardInfo{
							TotalUser:      1592,
							TotalThread:    95231,
							TotalModerator: 442139842,
							TotalReport:    125,
						}
					},
					func(ctx context.Context) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			gotDashboardInfo, gotError := adminService.GetDashboardInfo(
				context.Background(),
				testCase.inputAccessorRole,
			)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotError, testCase.expectedError)
			} else {
				assert.Equal(t, testCase.expectedDashboardInfo, gotDashboardInfo)
			}
		})
	}
}

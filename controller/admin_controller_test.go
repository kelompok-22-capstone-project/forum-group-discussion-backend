package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	mas "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/admin/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	mtg "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteAdmin(t *testing.T) {
	mockAdminService := &mas.AdminService{}
	mockTokenGen := &mtg.TokenGenerator{}
	controller := NewAdminController(mockAdminService, mockTokenGen)
	g := echo.New().Group("/api/v1")
	controller.Route(g)
	assert.NotNil(t, controller)
}

func TestGetInfo(t *testing.T) {
	mockAdminService := &mas.AdminService{}
	mockTokenGen := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		dummyInfoResp := response.DashboardInfo{
			TotalUser:      30,
			TotalThread:    13,
			TotalModerator: 2,
			TotalReport:    3,
		}

		dummyResp := model.NewResponse("Success", "Get admin dashboard info successful.", dummyInfoResp)

		mockTokenGen.On(
			"ExtractToken",
			mock.AnythingOfType("*echo.context"),
		).Return(
			func(c echo.Context) generator.TokenPayload {
				return generator.TokenPayload{
					ID:       "a-abcd",
					Username: "sarifaturr",
					Role:     "admin",
					IsActive: true,
				}
			},
		).Once()

		mockAdminService.On(
			"GetDashboardInfo",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(ctx context.Context, accessorRole string) response.DashboardInfo {
				return dummyInfoResp
			},
			func(ctx context.Context, accessorRole string) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewAdminController(mockAdminService, mockTokenGen)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.getInfo(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)

				body := rec.Body.String()

				gotResponse := model.NewResponse("", "", response.DashboardInfo{})

				if err := json.Unmarshal([]byte(body), &gotResponse); assert.NoError(t, err) {
					assert.Equal(t, dummyResp.Data, gotResponse.Data)
				}
			}
		})
	})

	t.Run("failed scenario", func(t *testing.T) {
		testCases := []struct {
			name                 string
			expectedStatusCode   int
			expectedErrorMessage string
			mockBehaviours       func()
		}{
			{
				name:                 "it should return 500 status code, when error happened",
				expectedStatusCode:   http.StatusInternalServerError,
				expectedErrorMessage: "Something went wrong.",
				mockBehaviours: func() {
					mockTokenGen.On(
						"ExtractToken",
						mock.AnythingOfType("*echo.context"),
					).Return(
						func(c echo.Context) generator.TokenPayload {
							return generator.TokenPayload{
								ID:       "a-abcd",
								Username: "sarifaturr",
								Role:     "admin",
								IsActive: true,
							}
						},
					).Once()

					mockAdminService.On(
						"GetDashboardInfo",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
					).Return(
						func(ctx context.Context, accessorRole string) response.DashboardInfo {
							return response.DashboardInfo{}
						},
						func(ctx context.Context, accessorRole string) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewAdminController(mockAdminService, mockTokenGen)

				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				gotErr := controller.getInfo(c)
				if assert.Error(t, gotErr) {
					if echoHTTPError, ok := gotErr.(*echo.HTTPError); assert.Equal(t, true, ok) {
						assert.Equal(t, testCase.expectedStatusCode, echoHTTPError.Code)
						assert.Equal(t, testCase.expectedErrorMessage, echoHTTPError.Message)
					}
				}
			})
		}
	})
}

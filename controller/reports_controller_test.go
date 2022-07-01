package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	mrs "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/report/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	mtg "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteReports(t *testing.T) {
	mockReportService := &mrs.ReportService{}
	mockTokenGenerator := &mtg.TokenGenerator{}
	controller := NewReportsController(mockReportService, mockTokenGenerator)
	g := echo.New().Group("/api/v1")
	controller.Route(g)
	assert.NotNil(t, controller)
}

func TestPostCreateReport(t *testing.T) {
	mockReportService := &mrs.ReportService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		dummyReq := payload.CreateReport{
			Username: "sarifaturr",
			ThreadID: "t-123",
			Reason:   "spam",
		}

		dummyID := "u-123"
		dummyIDResponse := map[string]any{"ID": dummyID}
		dummyResp := model.NewResponse("success", "Create report successful.", dummyIDResponse)

		mockTokenGenerator.On(
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

		mockReportService.On(
			"Create",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", payload.CreateReport{})),
		).Return(
			func(ctx context.Context, accessorUserID string, p payload.CreateReport) string {
				return dummyID
			},
			func(ctx context.Context, accessorUserID string, p payload.CreateReport) error {
				return nil
			},
		).Once()

		t.Run("it should return 201 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewReportsController(mockReportService, mockTokenGenerator)
			requestBody, err := json.Marshal(dummyReq)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/reports", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.postCreateReport(c)) {
				assert.Equal(t, http.StatusCreated, rec.Code)

				body := rec.Body.String()

				gotResponse := make(map[string]any)

				if err := json.Unmarshal([]byte(body), &gotResponse); assert.NoError(t, err) {
					gotID := gotResponse["data"].(map[string]any)["ID"].(string)
					assert.Equal(t, dummyResp.Data["ID"], gotID)
				}
			}
		})
	})

	t.Run("failed scenario", func(t *testing.T) {})
	dummyReq := payload.CreateReport{
		Username: "sarifaturr",
		ThreadID: "t-123",
		Reason:   "spam",
	}

	testCases := []struct {
		name                 string
		inputPayload         payload.CreateReport
		expectedStatusCode   int
		expectedErrorMessage string
		mockBehaviours       func()
	}{
		{
			name:                 "it should return 500 status code, when error happened",
			inputPayload:         dummyReq,
			expectedStatusCode:   http.StatusInternalServerError,
			expectedErrorMessage: "Something went wrong.",
			mockBehaviours: func() {
				mockTokenGenerator.On(
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

				mockReportService.On(
					"Create",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", payload.CreateReport{})),
				).Return(
					func(ctx context.Context, accessorUserID string, p payload.CreateReport) string {
						return ""
					},
					func(ctx context.Context, accessorUserID string, p payload.CreateReport) error {
						return service.ErrRepository
					},
				).Once()
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviours()

			controller := NewReportsController(mockReportService, mockTokenGenerator)
			requestBody, err := json.Marshal(testCase.inputPayload)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/reports", strings.NewReader(string(requestBody)))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			gotErr := controller.postCreateReport(c)
			if assert.Error(t, gotErr) {
				if echoHTTPError, ok := gotErr.(*echo.HTTPError); assert.Equal(t, true, ok) {
					assert.Equal(t, testCase.expectedStatusCode, echoHTTPError.Code)
					assert.Equal(t, testCase.expectedErrorMessage, echoHTTPError.Message)
				}
			}
		})
	}
}

func TestGetReports(t *testing.T) {
	mockReportService := &mrs.ReportService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		now := time.Now()
		dummyReports := []response.Report{
			{
				ID:                "r-xyz",
				ModeratorID:       "m-abc",
				ModeratorUsername: "gaga12",
				ModeratorName:     "gaga",
				UserID:            "u-jkl",
				Username:          "sarifaturr",
				Name:              "sari faturr",
				Reason:            "spam",
				Status:            "review",
				ThreadID:          "t-cvb",
				ThreadTitle:       "Hello, world",
				ReportedOn:        now.Format(time.RFC822),
			},
		}

		dummyPagination := response.Pagination[response.Report]{
			List: dummyReports,
			PageInfo: response.PageInfo{
				Limit:     20,
				Page:      1,
				PageTotal: 1,
				Total:     1,
			},
		}
		dummyResp := model.NewResponse("success", "Get reports successful.", dummyPagination)

		mockTokenGenerator.On(
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

		mockReportService.On(
			"GetAll",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
		).Return(
			func(ctx context.Context,
				accessorRole string,
				status string,
				page uint,
				limit uint,
			) response.Pagination[response.Report] {
				return dummyPagination
			},
			func(ctx context.Context,
				accessorRole string,
				status string,
				page uint,
				limit uint,
			) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewReportsController(mockReportService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/reports", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.getReports(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)

				body := rec.Body.String()

				gotResponse := model.NewResponse("", "", response.Pagination[response.Report]{})

				if err := json.Unmarshal([]byte(body), &gotResponse); assert.NoError(t, err) {
					assert.Equal(t, dummyResp.Data.PageInfo, gotResponse.Data.PageInfo)
					assert.ElementsMatch(t, dummyResp.Data.List, gotResponse.Data.List)
				}
			}
		})
	})

	t.Run("failed scenario", func(t *testing.T) {})
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
				mockTokenGenerator.On(
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

				mockReportService.On(
					"GetAll",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
					mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
				).Return(
					func(ctx context.Context,
						accessorRole string,
						status string,
						page uint,
						limit uint,
					) response.Pagination[response.Report] {
						return response.Pagination[response.Report]{}
					},
					func(ctx context.Context,
						accessorRole string,
						status string,
						page uint,
						limit uint,
					) error {
						return service.ErrRepository
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		testCase.mockBehaviours()

		controller := NewReportsController(mockReportService, mockTokenGenerator)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/reports", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		gotErr := controller.getReports(c)
		if assert.Error(t, gotErr) {
			if echoHTTPError, ok := gotErr.(*echo.HTTPError); assert.Equal(t, true, ok) {
				assert.Equal(t, testCase.expectedStatusCode, echoHTTPError.Code)
				assert.Equal(t, testCase.expectedErrorMessage, echoHTTPError.Message)
			}
		}
	}
}

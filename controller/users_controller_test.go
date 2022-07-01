package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	mus "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/user/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	mtg "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteUsers(t *testing.T) {
	mockUserService := &mus.UserService{}
	mockTokenGenerator := &mtg.TokenGenerator{}
	controller := NewUsersController(mockUserService, mockTokenGenerator)
	g := echo.New().Group("/api/v1")
	controller.Route(g)
	assert.NotNil(t, controller)
}

func TestGetUsers(t *testing.T) {
	mockUserService := &mus.UserService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		now := time.Now()
		dummyUsers := []response.User{
			{
				UserID:         "u-xyz",
				Username:       "sarifaturr",
				Email:          "sarifaturr@gmail.com",
				Name:           "sari faturr",
				Role:           "user",
				IsActive:       true,
				RegisteredOn:   now.Format(time.RFC822),
				TotalThread:    10,
				TotalFollower:  500,
				TotalFollowing: 324,
				IsFollowed:     true,
			},
		}
		dummyPagination := response.Pagination[response.User]{
			List: dummyUsers,
			PageInfo: response.PageInfo{
				Limit:     20,
				Page:      1,
				PageTotal: 1,
				Total:     1,
			},
		}
		dummyResp := model.NewResponse("success", "Get users successful.", dummyPagination)

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

		mockUserService.On(
			"GetAll",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(ctx context.Context,
				accessorUserID string,
				orderBy string,
				status string,
				page uint,
				limit uint,
				keyword string,
			) response.Pagination[response.User] {
				return dummyPagination
			},
			func(ctx context.Context,
				accessorUserID string,
				orderBy string,
				status string,
				page uint,
				limit uint,
				keyword string,
			) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewUsersController(mockUserService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.getUsers(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)

				body := rec.Body.String()

				gotResponse := model.NewResponse("", "", response.Pagination[response.User]{})

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

				mockUserService.On(
					"GetAll",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
					mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context,
						accessorUserID string,
						orderBy string,
						status string,
						page uint,
						limit uint,
						keyword string,
					) response.Pagination[response.User] {
						return response.Pagination[response.User]{}
					},
					func(ctx context.Context,
						accessorUserID string,
						orderBy string,
						status string,
						page uint,
						limit uint,
						keyword string,
					) error {
						return service.ErrRepository
					},
				).Once()
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviours()

			controller := NewUsersController(mockUserService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/threads", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			gotErr := controller.getUsers(c)
			if assert.Error(t, gotErr) {
				if echoHTTPError, ok := gotErr.(*echo.HTTPError); assert.Equal(t, true, ok) {
					assert.Equal(t, testCase.expectedStatusCode, echoHTTPError.Code)
					assert.Equal(t, testCase.expectedErrorMessage, echoHTTPError.Message)
				}
			}
		})
	}
}

func TestGetMe(t *testing.T) {
	mockUserService := &mus.UserService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		now := time.Now()
		dummyUser := response.User{
			UserID:         "u-xyz",
			Username:       "sarifaturr",
			Email:          "sarifaturr@gmail.com",
			Name:           "sari faturr",
			Role:           "user",
			IsActive:       true,
			RegisteredOn:   now.Format(time.RFC822),
			TotalThread:    10,
			TotalFollower:  500,
			TotalFollowing: 324,
			IsFollowed:     true,
		}
		dummyResp := model.NewResponse("success", "Get user successful.", dummyUser)

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

		mockUserService.On(
			"GetOwn",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(ctx context.Context, accessorUserID string, accessorUsername string) response.User {
				return dummyUser
			},
			func(ctx context.Context, accessorUserID string, accessorUsername string) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewUsersController(mockUserService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.getMe(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)

				body := rec.Body.String()

				gotResponse := model.NewResponse("", "", response.User{})

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

					mockUserService.On(
						"GetOwn",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
					).Return(
						func(ctx context.Context, accessorUserID string, accessorUsername string) response.User {
							return response.User{}
						},
						func(ctx context.Context, accessorUserID string, accessorUsername string) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewUsersController(mockUserService, mockTokenGenerator)

				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				gotErr := controller.getMe(c)
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

func TestGetUserByUsername(t *testing.T) {
	mockUserService := &mus.UserService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		now := time.Now()
		dummyUser := response.User{
			UserID:         "u-xyz",
			Username:       "sarifaturr",
			Email:          "sarifaturr@gmail.com",
			Name:           "sari faturr",
			Role:           "user",
			IsActive:       true,
			RegisteredOn:   now.Format(time.RFC822),
			TotalThread:    10,
			TotalFollower:  500,
			TotalFollowing: 324,
			IsFollowed:     true,
		}
		dummyResp := model.NewResponse("success", "Get user successful.", dummyUser)

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

		mockUserService.On(
			"GetByUsername",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(ctx context.Context, accessorUserID string, username string) response.User {
				return dummyUser
			},
			func(ctx context.Context, accessorUserID string, username string) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewUsersController(mockUserService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.getUserByUsername(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)

				body := rec.Body.String()

				gotResponse := model.NewResponse("", "", response.User{})

				if err := json.Unmarshal([]byte(body), &gotResponse); assert.NoError(t, err) {
					assert.Equal(t, dummyResp.Data, gotResponse.Data)
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

				mockUserService.On(
					"GetByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, username string) response.User {
						return response.User{}
					},
					func(ctx context.Context, accessorUserID string, username string) error {
						return service.ErrRepository
					},
				).Once()
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviours()

			controller := NewUsersController(mockUserService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			gotErr := controller.getUserByUsername(c)
			if assert.Error(t, gotErr) {
				if echoHTTPError, ok := gotErr.(*echo.HTTPError); assert.Equal(t, true, ok) {
					assert.Equal(t, testCase.expectedStatusCode, echoHTTPError.Code)
					assert.Equal(t, testCase.expectedErrorMessage, echoHTTPError.Message)
				}
			}
		})
	}
}

func TestGetUserThreads(t *testing.T){
	mockUserService := &mus.UserService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T){
		now := time.Now().Format(time.RFC822)
		dummyThreads := []response.ManyThread{
			{
				ID:              "t-abc",
				Title:           "Go Programming Going Hype",
				CategoryID:      "c-xyz",
				CategoryName:    "Tech",
				PublishedOn:     now,
				IsLiked:         true,
				IsFollowed:      true,
				Description:     "Currently Go Programming going hype because it's popularity",
				TotalViewer:     324,
				TotalLike:       90,
				TotalFollower:   25,
				TotalComment:    42,
				CreatorID:       "u-abcdefg",
				CreatorUsername: "erikrios",
				CreatorName:     "Erik Rio Setiawan",
			},
		}

		dummyPagination := response.Pagination[response.ManyThread]{
			List: dummyThreads,
			PageInfo: response.PageInfo{
				Limit:     20,
				Page:      1,
				PageTotal: 1,
				Total:     15,
			},
		}
		dummyResp := model.NewResponse("success", "Get threads by username successful.", dummyPagination)

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

		mockUserService.On(
			"GetAllThreadByUsername",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
		).Return(
			func(ctx context.Context,
				accessorUserID string,
				username string,
				page uint,
				limit uint,
			) response.Pagination[response.ManyThread] {
				return dummyPagination
			},
			func(ctx context.Context,
				accessorUserID string,
				username string,
				page uint,
				limit uint,
			) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewUsersController(mockUserService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.getUserThreads(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)

				body := rec.Body.String()

				gotResponse := model.NewResponse("", "", response.Pagination[response.ManyThread]{})

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

				mockUserService.On(
					"GetAllThreadByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
					mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
				).Return(
					func(ctx context.Context,
						accessorUserID string,
						username string,
						page uint,
						limit uint,
					) response.Pagination[response.ManyThread] {
						return response.Pagination[response.ManyThread]{}
					},
					func(ctx context.Context,
						accessorUserID string,
						username string,
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
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviours()

			controller := NewUsersController(mockUserService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			gotErr := controller.getUserThreads(c)
			if assert.Error(t, gotErr) {
				if echoHTTPError, ok := gotErr.(*echo.HTTPError); assert.Equal(t, true, ok) {
					assert.Equal(t, testCase.expectedStatusCode, echoHTTPError.Code)
					assert.Equal(t, testCase.expectedErrorMessage, echoHTTPError.Message)
				}
			}
		})
	}
}

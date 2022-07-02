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
	mts "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/thread/mocks"
	mus "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/user/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteGuest(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockUserService := &mus.UserService{}
	controller := NewGuestController(mockThreadService, mockUserService)
	g := echo.New().Group("/api/v1")
	controller.Route(g)
	assert.NotNil(t, controller)
}

func TestGuestGetThreads(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockUserService := &mus.UserService{}

	t.Run("success scenario", func(t *testing.T) {
		now := time.Now().Format(time.RFC822)
		dummyThreads := []response.ManyThread{
			{
				ID:              "t-eG4HE",
				Title:           "Go Programming Going Hype",
				CategoryID:      "c-xyz",
				CategoryName:    "Tech",
				PublishedOn:     now,
				IsLiked:         true,
				IsFollowed:      false,
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
		dummyResp := model.NewResponse("success", "Get threads successful.", dummyPagination)

		mockThreadService.On(
			"GetAll",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(
				ctx context.Context,
				accessorUserID string,
				page uint,
				limit uint,
				query string,
			) response.Pagination[response.ManyThread] {
				return dummyPagination
			},
			func(
				ctx context.Context,
				accessorUserID string,
				page uint,
				limit uint,
				query string,
			) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/threads", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.getThreads(c)) {
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
				mockThreadService.On(
					"GetAll",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
					mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						page uint,
						limit uint,
						query string,
					) response.Pagination[response.ManyThread] {
						return response.Pagination[response.ManyThread]{}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						page uint,
						limit uint,
						query string,
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

			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/threads", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			gotErr := controller.getThreads(c)
			if assert.Error(t, gotErr) {
				if echoHTTPError, ok := gotErr.(*echo.HTTPError); assert.Equal(t, true, ok) {
					assert.Equal(t, testCase.expectedStatusCode, echoHTTPError.Code)
					assert.Equal(t, testCase.expectedErrorMessage, echoHTTPError.Message)
				}
			}
		})
	}
}

func TestGuestGetThread(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockUserService := &mus.UserService{}

	t.Run("success scenario", func(t *testing.T) {
		now := time.Now().Format(time.RFC822)
		dummyThread := response.Thread{
			ID:           "t-eG4HE",
			Title:        "Go Programming Going Hype",
			CategoryID:   "c-xyz",
			CategoryName: "Tech",
			PublishedOn:  now,
			IsLiked:      true,
			IsFollowed:   false,
			Moderators: []response.Moderator{
				{
					ID:           "m-xyz",
					UserID:       "u-abcdefg",
					Username:     "erikrios",
					Email:        "erikriosetiawan@gmail.com",
					Name:         "Erik Rio Setiawan",
					Role:         "user",
					IsActive:     true,
					RegisteredOn: now,
				},
			},
			Description:     "Currently Go Programming going hype because it's popularity",
			TotalViewer:     324,
			TotalLike:       90,
			TotalFollower:   25,
			TotalComment:    42,
			CreatorID:       "u-abcdefg",
			CreatorUsername: "erikrios",
			CreatorName:     "Erik Rio Setiawan",
		}
		dummyResp := model.NewResponse("success", "Get threads successful.", dummyThread)

		mockThreadService.On(
			"GetByID",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(
				ctx context.Context,
				accessorUserID string,
				ID string,
			) response.Thread {
				return dummyThread
			},
			func(
				ctx context.Context,
				accessorUserID string,
				ID string,
			) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/threads", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(dummyThread.ID)

			if assert.NoError(t, controller.getThread(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)

				body := rec.Body.String()

				gotResponse := model.NewResponse("", "", response.Thread{})

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
				mockThreadService.On(
					"GetByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						ID string,
					) response.Thread {
						return response.Thread{}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						ID string,
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

			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/threads", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues("t-XyzAbc")

			gotErr := controller.getThread(c)
			if assert.Error(t, gotErr) {
				if echoHTTPError, ok := gotErr.(*echo.HTTPError); assert.Equal(t, true, ok) {
					assert.Equal(t, testCase.expectedStatusCode, echoHTTPError.Code)
					assert.Equal(t, testCase.expectedErrorMessage, echoHTTPError.Message)
				}
			}
		})
	}
}

func TestGuestGetThreadComments(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockUserService := &mus.UserService{}

	t.Run("success scenario", func(t *testing.T) {
		now := time.Now().Format(time.RFC822)
		dummyComments := []response.Comment{
			{
				ID:          "c-EtILo81",
				UserID:      "u-abcdefg",
				Username:    "erikrios",
				Name:        "Erik Rio Setiawan",
				Comment:     "Nice thread, good job.",
				PublishedOn: now,
			},
		}
		dummyPagination := response.Pagination[response.Comment]{
			List: dummyComments,
			PageInfo: response.PageInfo{
				Limit:     20,
				Page:      1,
				PageTotal: 1,
				Total:     15,
			},
		}
		dummyResp := model.NewResponse("success", "Get comments successful", dummyPagination)

		mockThreadService.On(
			"GetComments",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
		).Return(
			func(
				ctx context.Context,
				threadID string,
				page uint,
				limit uint,
			) response.Pagination[response.Comment] {
				return dummyPagination
			},
			func(
				ctx context.Context,
				threadID string,
				page uint,
				limit uint,
			) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/threads", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id/comments")
			c.SetParamNames("id")
			c.SetParamValues("t-XyzAbc")

			if assert.NoError(t, controller.getThreadComments(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)

				body := rec.Body.String()

				gotResponse := model.NewResponse("", "", response.Pagination[response.Comment]{})

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
				mockThreadService.On(
					"GetComments",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
					mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
				).Return(
					func(
						ctx context.Context,
						threadID string,
						page uint,
						limit uint,
					) response.Pagination[response.Comment] {
						return response.Pagination[response.Comment]{}
					},
					func(
						ctx context.Context,
						threadID string,
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

			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/threads", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id/comments")
			c.SetParamNames("id")
			c.SetParamValues("t-XyzAbc")

			gotErr := controller.getThreadComments(c)
			if assert.Error(t, gotErr) {
				if echoHTTPError, ok := gotErr.(*echo.HTTPError); assert.Equal(t, true, ok) {
					assert.Equal(t, testCase.expectedStatusCode, echoHTTPError.Code)
					assert.Equal(t, testCase.expectedErrorMessage, echoHTTPError.Message)
				}
			}
		})
	}
}

func TestGuestGetUsers(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockUserService := &mus.UserService{}

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
			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/users", nil)
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

			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/threads", nil)
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

func TestGuestGetUserByUsername(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockUserService := &mus.UserService{}

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
			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/users", nil)
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

			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/users", nil)
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

func TestGuestGetUserThreads(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockUserService := &mus.UserService{}

	t.Run("success scenario", func(t *testing.T) {
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
			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/users", nil)
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

			controller := NewGuestController(mockThreadService, mockUserService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/guest/users", nil)
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

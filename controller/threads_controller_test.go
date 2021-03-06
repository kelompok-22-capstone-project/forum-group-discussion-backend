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
	mts "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/thread/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	mtg "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteThreads(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}
	controller := NewThreadsController(mockThreadService, mockTokenGenerator)
	g := echo.New().Group("/api/v1")
	controller.Route(g)
	assert.NotNil(t, controller)
}

func TestGetThreads(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

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

		mockTokenGenerator.On(
			"ExtractToken",
			mock.AnythingOfType("*echo.context"),
		).Return(
			func(c echo.Context) generator.TokenPayload {
				return generator.TokenPayload{
					ID:       "u-abcdefg",
					Username: "erikrios",
					Role:     "user",
					IsActive: true,
				}
			},
		).Once()

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
			controller := NewThreadsController(mockThreadService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/threads", nil)
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
				mockTokenGenerator.On(
					"ExtractToken",
					mock.AnythingOfType("*echo.context"),
				).Return(
					func(c echo.Context) generator.TokenPayload {
						return generator.TokenPayload{
							ID:       "u-abcdefg",
							Username: "erikrios",
							Role:     "user",
							IsActive: true,
						}
					},
				).Once()

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

			controller := NewThreadsController(mockThreadService, mockTokenGenerator)

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

func TestPostCreateThread(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		dummyReq := payload.CreateThread{
			Title:       "Go Programming Going Hype",
			Description: "Currently Go Programming going hype because it's popularity",
			CategoryID:  "c-xyz",
		}

		dummyID := "t-XyzAbc"
		dummyIDResponse := map[string]any{"ID": dummyID}
		dummyResp := model.NewResponse("success", "Create thread successful.", dummyIDResponse)

		mockTokenGenerator.On(
			"ExtractToken",
			mock.AnythingOfType("*echo.context"),
		).Return(
			func(c echo.Context) generator.TokenPayload {
				return generator.TokenPayload{
					ID:       "u-abcdefg",
					Username: "erikrios",
					Role:     "user",
					IsActive: true,
				}
			},
		).Once()

		mockThreadService.On(
			"Create",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", payload.CreateThread{})),
		).Return(
			func(ctx context.Context, accessorUserID string, p payload.CreateThread) string {
				return dummyID
			},
			func(ctx context.Context, accessorUserID string, p payload.CreateThread) error {
				return nil
			},
		).Once()

		t.Run("it should return 201 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewThreadsController(mockThreadService, mockTokenGenerator)
			requestBody, err := json.Marshal(dummyReq)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/threads", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.postCreateThread(c)) {
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

	t.Run("failed scenario", func(t *testing.T) {
		dummyReq := payload.CreateThread{
			Title:       "Go Programming Going Hype",
			Description: "Currently Go Programming going hype because it's popularity",
			CategoryID:  "c-xyz",
		}

		testCases := []struct {
			name                 string
			inputPayload         payload.CreateThread
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
								ID:       "u-abcdefg",
								Username: "erikrios",
								Role:     "erikrios",
								IsActive: true,
							}
						},
					).Once()

					mockThreadService.On(
						"Create",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", payload.CreateThread{})),
					).Return(
						func(ctx context.Context, accessorUserID string, p payload.CreateThread) string {
							return ""
						},
						func(ctx context.Context, accessorUserID string, p payload.CreateThread) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewThreadsController(mockThreadService, mockTokenGenerator)
				requestBody, err := json.Marshal(dummyReq)
				assert.NoError(t, err)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/api/v1/threads", strings.NewReader(string(requestBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				gotErr := controller.postCreateThread(c)
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

func TestGetThread(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

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

		mockTokenGenerator.On(
			"ExtractToken",
			mock.AnythingOfType("*echo.context"),
		).Return(
			func(c echo.Context) generator.TokenPayload {
				return generator.TokenPayload{
					ID:       "u-abcdefg",
					Username: "erikrios",
					Role:     "user",
					IsActive: true,
				}
			},
		).Once()

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
			controller := NewThreadsController(mockThreadService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/threads", nil)
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
				mockTokenGenerator.On(
					"ExtractToken",
					mock.AnythingOfType("*echo.context"),
				).Return(
					func(c echo.Context) generator.TokenPayload {
						return generator.TokenPayload{
							ID:       "u-abcdefg",
							Username: "erikrios",
							Role:     "user",
							IsActive: true,
						}
					},
				).Once()

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

			controller := NewThreadsController(mockThreadService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/threads", nil)
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

func TestPutUpdateThread(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		dummyReq := payload.UpdateThread{
			Title:       "Go Programming Going Hype",
			CategoryID:  "c-xyz",
			Description: "Currently Go Programming going hype because it's popularity",
		}

		mockTokenGenerator.On(
			"ExtractToken",
			mock.AnythingOfType("*echo.context"),
		).Return(
			func(c echo.Context) generator.TokenPayload {
				return generator.TokenPayload{
					ID:       "u-abcdefg",
					Username: "admin",
					Role:     "admin",
					IsActive: true,
				}
			},
		).Once()

		mockThreadService.On(
			"Update",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", payload.UpdateThread{})),
		).Return(
			func(ctx context.Context, accessorUserID string, id string, p payload.UpdateThread) error {
				return nil
			},
		).Once()

		t.Run("it should return 204 status code, when there is no error", func(t *testing.T) {
			controller := NewThreadsController(mockThreadService, mockTokenGenerator)
			requestBody, err := json.Marshal(dummyReq)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/v1/threads", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues("c-xyz")

			if assert.NoError(t, controller.putUpdateThread(c)) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
			}
		})
	})

	t.Run("failed scenario", func(t *testing.T) {
		dummyReq := payload.UpdateThread{
			Title:       "Go Programming Going Hype",
			CategoryID:  "c-xyz",
			Description: "Currently Go Programming going hype because it's popularity",
		}

		testCases := []struct {
			name                 string
			inputPayload         payload.UpdateThread
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
								ID:       "u-abcdefg",
								Username: "admin",
								Role:     "admin",
								IsActive: true,
							}
						},
					).Once()

					mockThreadService.On(
						"Update",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", payload.UpdateThread{})),
					).Return(
						func(ctx context.Context, accessorUserID string, id string, p payload.UpdateThread) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewThreadsController(mockThreadService, mockTokenGenerator)
				requestBody, err := json.Marshal(dummyReq)
				assert.NoError(t, err)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPut, "/api/v1/threads", strings.NewReader(string(requestBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/:id")
				c.SetParamNames("id")
				c.SetParamValues("c-xyz")

				gotErr := controller.putUpdateThread(c)
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

func TestDeleteThread(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		mockTokenGenerator.On(
			"ExtractToken",
			mock.AnythingOfType("*echo.context"),
		).Return(
			func(c echo.Context) generator.TokenPayload {
				return generator.TokenPayload{
					ID:       "u-abcdefg",
					Username: "admin",
					Role:     "admin",
					IsActive: true,
				}
			},
		).Once()

		mockThreadService.On(
			"Delete",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(ctx context.Context, accessorUserID string, role string, id string) error {
				return nil
			},
		).Once()

		t.Run("it should return 204 status code, when there is no error", func(t *testing.T) {
			controller := NewThreadsController(mockThreadService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/v1/categories", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues("c-xyz")

			if assert.NoError(t, controller.deleteThread(c)) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
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
								ID:       "u-abcdefg",
								Username: "admin",
								Role:     "admin",
								IsActive: true,
							}
						},
					).Once()

					mockThreadService.On(
						"Delete",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
					).Return(
						func(ctx context.Context, accessorUserID string, role string, id string) error {
							return service.ErrRepository
						},
					).Once()

				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewThreadsController(mockThreadService, mockTokenGenerator)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPut, "/api/v1/threads", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/:id")
				c.SetParamNames("id")
				c.SetParamValues("c-xyz")

				gotErr := controller.deleteThread(c)
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

func TestGetThreadComments(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

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
			controller := NewThreadsController(mockThreadService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/threads", nil)
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

			controller := NewThreadsController(mockThreadService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/threads", nil)
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

func TestPostCreateThreadComments(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		dummyReq := payload.CreateComment{
			Comment: "Nice thread, good job.",
		}

		dummyID := "c-XyzAbce"
		dummyIDResponse := map[string]any{"ID": dummyID}
		dummyResp := model.NewResponse("success", "Create comment successful.", dummyIDResponse)

		mockTokenGenerator.On(
			"ExtractToken",
			mock.AnythingOfType("*echo.context"),
		).Return(
			func(c echo.Context) generator.TokenPayload {
				return generator.TokenPayload{
					ID:       "u-abcdefg",
					Username: "erikrios",
					Role:     "user",
					IsActive: true,
				}
			},
		).Once()

		mockThreadService.On(
			"CreateComment",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", payload.CreateComment{})),
		).Return(
			func(ctx context.Context, threadID, accessorUserID string, p payload.CreateComment) string {
				return dummyID
			},
			func(ctx context.Context, threadID, accessorUserID string, p payload.CreateComment) error {
				return nil
			},
		).Once()

		t.Run("it should return 201 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewThreadsController(mockThreadService, mockTokenGenerator)
			requestBody, err := json.Marshal(dummyReq)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/threads", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id/comments")
			c.SetParamNames("id")
			c.SetParamValues("t-XyzAbc")

			if assert.NoError(t, controller.postCreateThreadComments(c)) {
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

	t.Run("failed scenario", func(t *testing.T) {
		dummyReq := payload.CreateComment{
			Comment: "Nice thread, good job.",
		}

		testCases := []struct {
			name                 string
			inputPayload         payload.CreateComment
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
								ID:       "u-abcdefg",
								Username: "erikrios",
								Role:     "user",
								IsActive: true,
							}
						},
					).Once()

					mockThreadService.On(
						"CreateComment",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", payload.CreateComment{})),
					).Return(
						func(ctx context.Context, threadID, accessorUserID string, p payload.CreateComment) string {
							return ""
						},
						func(ctx context.Context, threadID, accessorUserID string, p payload.CreateComment) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewThreadsController(mockThreadService, mockTokenGenerator)
				requestBody, err := json.Marshal(dummyReq)
				assert.NoError(t, err)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/api/v1/threads", strings.NewReader(string(requestBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/:id/comments")
				c.SetParamNames("id")
				c.SetParamValues("t-XyzAbc")

				gotErr := controller.postCreateThreadComments(c)
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

func TestPutThreadLike(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		mockTokenGenerator.On(
			"ExtractToken",
			mock.AnythingOfType("*echo.context"),
		).Return(
			func(c echo.Context) generator.TokenPayload {
				return generator.TokenPayload{
					ID:       "u-abcdefg",
					Username: "erikrios",
					Role:     "user",
					IsActive: true,
				}
			},
		).Once()

		mockThreadService.On(
			"ChangeLikeState",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(ctx context.Context, threadID, accessorUserID string) error {
				return nil
			},
		).Once()

		t.Run("it should return 204 status code, when there is no error", func(t *testing.T) {
			controller := NewThreadsController(mockThreadService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/v1/threads", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id/like")
			c.SetParamNames("id")
			c.SetParamValues("c-xyz")

			if assert.NoError(t, controller.putThreadLike(c)) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
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
								ID:       "u-abcdefg",
								Username: "erikrios",
								Role:     "user",
								IsActive: true,
							}
						},
					).Once()

					mockThreadService.On(
						"ChangeLikeState",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
					).Return(
						func(ctx context.Context, threadID, accessorUserID string) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewThreadsController(mockThreadService, mockTokenGenerator)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPut, "/api/v1/threads", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/:id/like")
				c.SetParamNames("id")
				c.SetParamValues("c-xyz")

				gotErr := controller.putThreadLike(c)
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

func TestPutThreadFollow(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		mockTokenGenerator.On(
			"ExtractToken",
			mock.AnythingOfType("*echo.context"),
		).Return(
			func(c echo.Context) generator.TokenPayload {
				return generator.TokenPayload{
					ID:       "u-abcdefg",
					Username: "erikrios",
					Role:     "user",
					IsActive: true,
				}
			},
		).Once()

		mockThreadService.On(
			"ChangeFollowingState",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(ctx context.Context, threadID, accessorUserID string) error {
				return nil
			},
		).Once()

		t.Run("it should return 204 status code, when there is no error", func(t *testing.T) {
			controller := NewThreadsController(mockThreadService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/v1/threads", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id/follow")
			c.SetParamNames("id")
			c.SetParamValues("t-XyzAbc")

			if assert.NoError(t, controller.putThreadFollow(c)) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
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
								ID:       "u-abcdefg",
								Username: "erikrios",
								Role:     "user",
								IsActive: true,
							}
						},
					).Once()

					mockThreadService.On(
						"ChangeFollowingState",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
					).Return(
						func(ctx context.Context, threadID, accessorUserID string) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewThreadsController(mockThreadService, mockTokenGenerator)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPut, "/api/v1/threads", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/:id/follow")
				c.SetParamNames("id")
				c.SetParamValues("t-XyzAbc")

				gotErr := controller.putThreadFollow(c)
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

func TestPutThreadAddModerator(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		dummyReq := payload.AddRemoveModerator{
			Username: "naruto",
		}

		mockTokenGenerator.On(
			"ExtractToken",
			mock.AnythingOfType("*echo.context"),
		).Return(
			func(c echo.Context) generator.TokenPayload {
				return generator.TokenPayload{
					ID:       "u-abcdefg",
					Username: "erikrios",
					Role:     "user",
					IsActive: true,
				}
			},
		).Once()

		mockThreadService.On(
			"AddModerator",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", payload.AddRemoveModerator{})),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(ctx context.Context, p payload.AddRemoveModerator, threadID, accessorUserID string) error {
				return nil
			},
		).Once()

		t.Run("it should return 204 status code, when there is no error", func(t *testing.T) {
			controller := NewThreadsController(mockThreadService, mockTokenGenerator)
			requestBody, err := json.Marshal(dummyReq)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/v1/threads", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id/moderators/add")
			c.SetParamNames("id")
			c.SetParamValues("t-XyzAbc")

			if assert.NoError(t, controller.putThreadAddModerator(c)) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
			}
		})
	})

	t.Run("failed scenario", func(t *testing.T) {
		dummyReq := payload.UpdateThread{
			Title:       "Go Programming Going Hype",
			CategoryID:  "c-xyz",
			Description: "Currently Go Programming going hype because it's popularity",
		}

		testCases := []struct {
			name                 string
			inputPayload         payload.UpdateThread
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
								ID:       "u-abcdefg",
								Username: "erikrios",
								Role:     "user",
								IsActive: true,
							}
						},
					).Once()

					mockThreadService.On(
						"AddModerator",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", payload.AddRemoveModerator{})),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
					).Return(
						func(ctx context.Context, p payload.AddRemoveModerator, threadID, accessorUserID string) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewThreadsController(mockThreadService, mockTokenGenerator)
				requestBody, err := json.Marshal(dummyReq)
				assert.NoError(t, err)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPut, "/api/v1/threads", strings.NewReader(string(requestBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/:id/moderators/add")
				c.SetParamNames("id")
				c.SetParamValues("t-XyzAbc")

				gotErr := controller.putThreadAddModerator(c)
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

func TestPutThreadRemoveModerator(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenario", func(t *testing.T) {
		dummyReq := payload.AddRemoveModerator{
			Username: "naruto",
		}

		mockTokenGenerator.On(
			"ExtractToken",
			mock.AnythingOfType("*echo.context"),
		).Return(
			func(c echo.Context) generator.TokenPayload {
				return generator.TokenPayload{
					ID:       "u-abcdefg",
					Username: "erikrios",
					Role:     "user",
					IsActive: true,
				}
			},
		).Once()

		mockThreadService.On(
			"RemoveModerator",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", payload.AddRemoveModerator{})),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(ctx context.Context, p payload.AddRemoveModerator, threadID, accessorUserID string) error {
				return nil
			},
		).Once()

		t.Run("it should return 204 status code, when there is no error", func(t *testing.T) {
			controller := NewThreadsController(mockThreadService, mockTokenGenerator)
			requestBody, err := json.Marshal(dummyReq)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/v1/threads", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id/moderators/remove")
			c.SetParamNames("id")
			c.SetParamValues("t-XyzAbc")

			if assert.NoError(t, controller.putThreadRemoveModerator(c)) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
			}
		})
	})

	t.Run("failed scenario", func(t *testing.T) {
		dummyReq := payload.UpdateThread{
			Title:       "Go Programming Going Hype",
			CategoryID:  "c-xyz",
			Description: "Currently Go Programming going hype because it's popularity",
		}

		testCases := []struct {
			name                 string
			inputPayload         payload.UpdateThread
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
								ID:       "u-abcdefg",
								Username: "erikrios",
								Role:     "user",
								IsActive: true,
							}
						},
					).Once()

					mockThreadService.On(
						"RemoveModerator",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", payload.AddRemoveModerator{})),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
					).Return(
						func(ctx context.Context, p payload.AddRemoveModerator, threadID, accessorUserID string) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewThreadsController(mockThreadService, mockTokenGenerator)
				requestBody, err := json.Marshal(dummyReq)
				assert.NoError(t, err)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPut, "/api/v1/threads", strings.NewReader(string(requestBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/:id/moderators/remove")
				c.SetParamNames("id")
				c.SetParamValues("t-XyzAbc")

				gotErr := controller.putThreadRemoveModerator(c)
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

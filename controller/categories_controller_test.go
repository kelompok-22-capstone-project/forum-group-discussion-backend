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
	mcs "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/category/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	mtg "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteCategories(t *testing.T) {
	mockCategoryService := &mcs.CategoryService{}
	mockTokenGenerator := &mtg.TokenGenerator{}
	controller := NewCategoriesController(mockCategoryService, mockTokenGenerator)
	g := echo.New().Group("/api/v1")
	controller.Route(g)
	assert.NotNil(t, controller)
}

func TestPostCreateCategory(t *testing.T) {
	mockCategoryService := &mcs.CategoryService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenarion", func(t *testing.T) {
		dummyReq := payload.CreateCategory{
			Name:        "Tech",
			Description: "Technology is used to make everything easier.",
		}

		dummyID := "c-XyzAbc"
		dummyIDResponse := map[string]any{"ID": dummyID}
		dummyResp := model.NewResponse("success", "Create category successful.", dummyIDResponse)

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

		mockCategoryService.On(
			"Create",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", payload.CreateCategory{})),
		).Return(
			func(ctx context.Context, accessorRole string, p payload.CreateCategory) string {
				return dummyID
			},
			func(ctx context.Context, accessorRole string, p payload.CreateCategory) error {
				return nil
			},
		).Once()

		t.Run("it should return 201 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewCategoriesController(mockCategoryService, mockTokenGenerator)
			requestBody, err := json.Marshal(dummyReq)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.postCreateCategory(c)) {
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
		dummyReq := payload.CreateCategory{
			Name:        "Tech",
			Description: "Technology is used to make everything easier.",
		}

		testCases := []struct {
			name                 string
			inputPayload         payload.CreateCategory
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

					mockCategoryService.On(
						"Create",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", payload.CreateCategory{})),
					).Return(
						func(ctx context.Context, accessorRole string, p payload.CreateCategory) string {
							return ""
						},
						func(ctx context.Context, accessorRole string, p payload.CreateCategory) error {
							return service.ErrRepository
						},
					).Once()

				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewCategoriesController(mockCategoryService, mockTokenGenerator)
				requestBody, err := json.Marshal(dummyReq)
				assert.NoError(t, err)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", strings.NewReader(string(requestBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				gotErr := controller.postCreateCategory(c)
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

func TestGetCategories(t *testing.T) {
	mockCategoryService := &mcs.CategoryService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenarion", func(t *testing.T) {
		now := time.Now().Format(time.RFC822)
		dummyCategories := []response.Category{
			{
				ID:          "c-xyz",
				Name:        "Tech",
				Description: "Technology is used to make everything easier.",
				CreatedOn:   now,
			},
		}
		dummyIDResponse := map[string][]response.Category{"categories": dummyCategories}
		dummyResp := model.NewResponse("success", "Get categories successful.", dummyIDResponse)

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

		mockCategoryService.On(
			"GetAll",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
		).Return(
			func(ctx context.Context) []response.Category {
				return dummyCategories
			},
			func(ctx context.Context) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewCategoriesController(mockCategoryService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.getCategories(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)

				body := rec.Body.String()

				gotResponse := make(map[string]any)

				if err := json.Unmarshal([]byte(body), &gotResponse); assert.NoError(t, err) {
					gotCategories := gotResponse["data"].(map[string]any)["categories"].([]any)
					assert.Equal(t, len(dummyResp.Data["categories"]), len(gotCategories))

					for i, gotCategory := range gotCategories {
						gotID := gotCategory.(map[string]any)["ID"]
						gotName := gotCategory.(map[string]any)["name"]
						gotDescription := gotCategory.(map[string]any)["description"]
						gotCreatedOn := gotCategory.(map[string]any)["createdOn"]
						assert.Equal(t, dummyResp.Data["categories"][i].ID, gotID)
						assert.Equal(t, dummyResp.Data["categories"][i].Name, gotName)
						assert.Equal(t, dummyResp.Data["categories"][i].Description, gotDescription)
						assert.Equal(t, dummyResp.Data["categories"][i].CreatedOn, gotCreatedOn)
					}
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
								ID:       "u-abcdefg",
								Username: "erikrios",
								Role:     "user",
								IsActive: true,
							}
						},
					).Once()

					mockCategoryService.On(
						"GetAll",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					).Return(
						func(ctx context.Context) []response.Category {
							return []response.Category{}
						},
						func(ctx context.Context) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewCategoriesController(mockCategoryService, mockTokenGenerator)

				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				gotErr := controller.getCategories(c)
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

func TestPutUpdateCategory(t *testing.T) {
	mockCategoryService := &mcs.CategoryService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenarion", func(t *testing.T) {
		dummyReq := payload.UpdateCategory{
			Name:        "Tech",
			Description: "Technology is used to make everything easier.",
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

		mockCategoryService.On(
			"Update",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", payload.UpdateCategory{})),
		).Return(
			func(ctx context.Context, accessorRole string, id string, p payload.UpdateCategory) error {
				return nil
			},
		).Once()

		t.Run("it should return 204 status code, when there is no error", func(t *testing.T) {
			controller := NewCategoriesController(mockCategoryService, mockTokenGenerator)
			requestBody, err := json.Marshal(dummyReq)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/v1/categories", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues("c-xyz")

			if assert.NoError(t, controller.putUpdateCategory(c)) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
			}
		})
	})

	t.Run("failed scenario", func(t *testing.T) {
		dummyReq := payload.UpdateCategory{
			Name:        "Tech",
			Description: "Technology is used to make everything easier.",
		}

		testCases := []struct {
			name                 string
			inputPayload         payload.UpdateCategory
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

					mockCategoryService.On(
						"Update",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", payload.UpdateCategory{})),
					).Return(
						func(ctx context.Context, accessorRole string, id string, p payload.UpdateCategory) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewCategoriesController(mockCategoryService, mockTokenGenerator)
				requestBody, err := json.Marshal(dummyReq)
				assert.NoError(t, err)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPut, "/api/v1/categories", strings.NewReader(string(requestBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/:id")
				c.SetParamNames("id")
				c.SetParamValues("c-xyz")

				gotErr := controller.putUpdateCategory(c)
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

func TestDeleteCategory(t *testing.T) {
	mockCategoryService := &mcs.CategoryService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenarion", func(t *testing.T) {
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

		mockCategoryService.On(
			"Delete",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
		).Return(
			func(ctx context.Context, accessorRole string, id string) error {
				return nil
			},
		).Once()

		t.Run("it should return 204 status code, when there is no error", func(t *testing.T) {
			controller := NewCategoriesController(mockCategoryService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/v1/categories", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues("c-xyz")

			if assert.NoError(t, controller.deleteCategory(c)) {
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

					mockCategoryService.On(
						"Delete",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
					).Return(
						func(ctx context.Context, accessorRole string, id string) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviours()

				controller := NewCategoriesController(mockCategoryService, mockTokenGenerator)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPut, "/api/v1/categories", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/:id")
				c.SetParamNames("id")
				c.SetParamValues("c-xyz")

				gotErr := controller.deleteCategory(c)
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

func TestGetCategoryThreads(t *testing.T) {
	mockCategoryService := &mcs.CategoryService{}
	mockTokenGenerator := &mtg.TokenGenerator{}

	t.Run("success scenarion", func(t *testing.T) {
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
		dummyResp := model.NewResponse("success", "Get threads by category successful.", dummyPagination)

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

		mockCategoryService.On(
			"GetAllByCategory",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", "")),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
			mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
		).Return(
			func(
				ctx context.Context,
				accessorID string,
				categoryID string,
				page uint,
				limit uint,
			) response.Pagination[response.ManyThread] {
				return dummyPagination
			},
			func(
				ctx context.Context,
				accessorID string,
				categoryID string,
				page uint,
				limit uint,
			) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewCategoriesController(mockCategoryService, mockTokenGenerator)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues("c-xyz")

			if assert.NoError(t, controller.getCategoryThreads(c)) {
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

					mockCategoryService.On(
						"GetAllByCategory",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", "")),
						mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
						mock.AnythingOfType(fmt.Sprintf("%T", uint(0))),
					).Return(
						func(
							ctx context.Context,
							accessorID string,
							categoryID string,
							page uint,
							limit uint,
						) response.Pagination[response.ManyThread] {
							return response.Pagination[response.ManyThread]{}
						},
						func(
							ctx context.Context,
							accessorID string,
							categoryID string,
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

				controller := NewCategoriesController(mockCategoryService, mockTokenGenerator)

				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetPath("/:id")
				c.SetParamNames("id")
				c.SetParamValues("c-xyz")

				gotErr := controller.getCategoryThreads(c)
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

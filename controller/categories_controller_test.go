package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
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

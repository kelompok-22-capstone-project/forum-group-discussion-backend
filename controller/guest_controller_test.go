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
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/user/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteGuest(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockUserService := &mocks.UserService{}
	controller := NewGuestController(mockThreadService, mockUserService)
	g := echo.New().Group("/api/v1")
	controller.Route(g)
	assert.NotNil(t, controller)
}

func TestGuestGetThreads(t *testing.T) {
	mockThreadService := &mts.ThreadService{}
	mockUserService := &mocks.UserService{}

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

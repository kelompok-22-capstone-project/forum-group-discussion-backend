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
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/user/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteRegister(t *testing.T) {
	mockUserService := &mocks.UserService{}
	controller := NewRegisterController(mockUserService)
	g := echo.New().Group("/api/v1")
	controller.Route(g)
	assert.NotNil(t, controller)
}

func TestPostRegister(t *testing.T) {
	mockUserService := &mocks.UserService{}

	t.Run("success scenario", func(t *testing.T) {
		dummyReq := payload.Register{
			Username: "erikrios",
			Email:    "erikriosetiawan15@gmail.com",
			Name:     "Erik Rio Setiawan",
			Password: "erikriosetiawan",
		}

		dummyID := "u-abcdef"
		dummyIDResponse := map[string]any{"userID": dummyID}
		dummyResp := model.NewResponse("success", "Register successful.", dummyIDResponse)

		mockUserService.On(
			"Register",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", payload.Register{})),
		).Return(
			func(ctx context.Context, p payload.Register) string {
				return dummyID
			},
			func(ctx context.Context, p payload.Register) error {
				return nil
			},
		).Once()

		t.Run("it should return 201 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewRegisterController(mockUserService)
			requestBody, err := json.Marshal(dummyReq)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.postRegister(c)) {
				assert.Equal(t, http.StatusCreated, rec.Code)

				body := rec.Body.String()

				gotResponse := make(map[string]any)

				if err := json.Unmarshal([]byte(body), &gotResponse); assert.NoError(t, err) {
					gotID := gotResponse["data"].(map[string]any)["userID"].(string)
					assert.Equal(t, dummyResp.Data["userID"], gotID)
				}
			}
		})
	})

	t.Run("failed scenario", func(t *testing.T) {
		dummyReq := payload.Register{
			Username: "erikrios",
			Email:    "erikriosetiawan15@gmail.com",
			Name:     "Erik Rio Setiawan",
			Password: "erikriosetiawan",
		}

		testCases := []struct {
			name                 string
			inputPayload         payload.Register
			expectedStatusCode   int
			expectedErrorMessage string
			mockBehaviour        func()
		}{
			{
				name:                 "it should return 500 status code, when error happened",
				inputPayload:         dummyReq,
				expectedStatusCode:   http.StatusInternalServerError,
				expectedErrorMessage: "Something went wrong.",
				mockBehaviour: func() {
					mockUserService.On(
						"Register",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", payload.Register{})),
					).Return(
						func(ctx context.Context, p payload.Register) string {
							return ""
						},
						func(ctx context.Context, p payload.Register) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviour()

				controller := NewRegisterController(mockUserService)
				requestBody, err := json.Marshal(dummyReq)
				assert.NoError(t, err)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(string(requestBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				gotErr := controller.postRegister(c)
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

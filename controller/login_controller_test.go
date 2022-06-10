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
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/user/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteLogin(t *testing.T) {
	mockUserService := &mocks.UserService{}
	controller := NewLoginController(mockUserService)
	g := echo.New().Group("/api/v1")
	controller.Route(g)
	assert.NotNil(t, controller)
}

func TestPostLogin(t *testing.T) {
	mockUserService := &mocks.UserService{}

	t.Run("success scenario", func(t *testing.T) {
		dummyReq := payload.Login{
			Username: "erikrios",
			Password: "erikriosetiawan",
		}

		dummyToken := "generatedtoken"
		dummyRole := "user"
		dummyLoginResponse := map[string]any{"token": dummyToken, "role": dummyRole}
		dummyResp := model.NewResponse("success", "Login successful.", dummyLoginResponse)

		mockUserService.On(
			"Login",
			mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
			mock.AnythingOfType(fmt.Sprintf("%T", payload.Login{})),
		).Return(
			func(ctx context.Context, p payload.Login) response.Login {
				return response.Login{
					Token: "generatedtoken",
					Role:  "user",
				}
			},
			func(ctx context.Context, p payload.Login) error {
				return nil
			},
		).Once()

		t.Run("it should return 200 status code with valid response, when there is no error", func(t *testing.T) {
			controller := NewLoginController(mockUserService)
			requestBody, err := json.Marshal(dummyReq)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, controller.postLogin(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)

				body := rec.Body.String()

				gotResponse := make(map[string]any)

				if err := json.Unmarshal([]byte(body), &gotResponse); assert.NoError(t, err) {
					token := gotResponse["data"].(map[string]any)["token"].(string)
					role := gotResponse["data"].(map[string]any)["role"].(string)
					assert.Equal(t, dummyResp.Data["token"], token)
					assert.Equal(t, dummyResp.Data["role"], role)
				}
			}
		})
	})

	t.Run("failed scenario", func(t *testing.T) {
		dummyReq := payload.Login{
			Username: "erikrios",
			Password: "erikriosetiawan",
		}

		testCases := []struct {
			name                 string
			inputPayload         payload.Login
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
						"Login",
						mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
						mock.AnythingOfType(fmt.Sprintf("%T", payload.Login{})),
					).Return(
						func(ctx context.Context, p payload.Login) response.Login {
							return response.Login{}
						},
						func(ctx context.Context, p payload.Login) error {
							return service.ErrRepository
						},
					).Once()
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				testCase.mockBehaviour()

				controller := NewLoginController(mockUserService)
				requestBody, err := json.Marshal(dummyReq)
				assert.NoError(t, err)

				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(string(requestBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				gotErr := controller.postLogin(c)
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

package user

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
	mur "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/user/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	mig "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator/mocks"
	mpg "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator/mocks"
	mtg "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	mockUserRepository := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}
	mockPwdGen := &mpg.PasswordGenerator{}
	mockTokenGen := &mtg.TokenGenerator{}

	var userService UserService = NewUserServiceImpl(
		mockUserRepository,
		mockIDGen,
		mockPwdGen,
		mockTokenGen,
	)

	testCases := []struct {
		name           string
		inputPayload   payload.Register
		expectedID     string
		expectedError  error
		mockBehaviours func()
	}{
		{
			name: "it should return service.ErrInvalidPayload error, when payload is invalid",
			inputPayload: payload.Register{
				Username: "",
				Email:    "erikriosetiawan15@gmail.com",
				Name:     "",
				Password: "",
			},
			expectedID:     "",
			expectedError:  service.ErrInvalidPayload,
			mockBehaviours: func() {},
		},
		{
			name: "it should return service.ErrRepository error, when generate user ID return an error",
			inputPayload: payload.Register{
				Username: "erikrios",
				Email:    "erikriosetiawan15@gmail.com",
				Name:     "Erik Rio Setiawan",
				Password: "erikriosetiawan",
			},
			expectedID:    "",
			expectedError: service.ErrRepository,
			mockBehaviours: func() {
				mockIDGen.On(
					"GenerateUserID",
				).Return(
					func() string {
						return ""
					},
					func() error {
						return errors.New("failed to generate user ID")
					},
				).Once()
			},
		},
		{
			name: "it should return service.ErrRepository error, when generate from password return an error",
			inputPayload: payload.Register{
				Username: "erikrios",
				Email:    "erikriosetiawan15@gmail.com",
				Name:     "Erik Rio Setiawan",
				Password: "erikriosetiawan",
			},
			expectedID:    "",
			expectedError: service.ErrRepository,
			mockBehaviours: func() {
				mockIDGen.On(
					"GenerateUserID",
				).Return(
					func() string {
						return "u-abcdef"
					},
					func() error {
						return nil
					},
				).Once()

				mockPwdGen.On(
					"GenerateFromPassword",
					mock.AnythingOfType(fmt.Sprintf("%T", []byte{})),
					mock.AnythingOfType(fmt.Sprintf("%T", 0)),
				).Return(
					func(p []byte, cost int) []byte {
						return []byte{}
					},
					func(p []byte, cost int) error {
						return errors.New("failed to generate password")
					},
				).Once()
			},
		},
		{
			name: "it should return service.ErrRepository error, when repository return an error",
			inputPayload: payload.Register{
				Username: "erikrios",
				Email:    "erikriosetiawan15@gmail.com",
				Name:     "Erik Rio Setiawan",
				Password: "erikriosetiawan",
			},
			expectedID:    "",
			expectedError: service.ErrRepository,
			mockBehaviours: func() {
				mockIDGen.On(
					"GenerateUserID",
				).Return(
					func() string {
						return "u-abcdef"
					},
					func() error {
						return nil
					},
				).Once()

				mockPwdGen.On(
					"GenerateFromPassword",
					mock.AnythingOfType(fmt.Sprintf("%T", []byte{})),
					mock.AnythingOfType(fmt.Sprintf("%T", 0)),
				).Return(
					func(p []byte, cost int) []byte {
						return []byte("generatedpassword")
					},
					func(p []byte, cost int) error {
						return nil
					},
				).Once()

				mockUserRepository.On(
					"Insert",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.User{})),
				).Return(
					func(ctx context.Context, user entity.User) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name: "it should return a valid ID, when no error is returned",
			inputPayload: payload.Register{
				Username: "erikrios",
				Email:    "erikriosetiawan15@gmail.com",
				Name:     "Erik Rio Setiawan",
				Password: "erikriosetiawan",
			},
			expectedID:    "u-abcdef",
			expectedError: nil,
			mockBehaviours: func() {
				mockIDGen.On(
					"GenerateUserID",
				).Return(
					func() string {
						return "u-abcdef"
					},
					func() error {
						return nil
					},
				).Once()

				mockPwdGen.On(
					"GenerateFromPassword",
					mock.AnythingOfType(fmt.Sprintf("%T", []byte{})),
					mock.AnythingOfType(fmt.Sprintf("%T", 0)),
				).Return(
					func(p []byte, cost int) []byte {
						return []byte("generatedpassword")
					},
					func(p []byte, cost int) error {
						return nil
					},
				).Once()

				mockUserRepository.On(
					"Insert",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.User{})),
				).Return(
					func(ctx context.Context, user entity.User) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviours()

			gotID, gotErr := userService.Register(context.Background(), testCase.inputPayload)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotErr, testCase.expectedError)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, testCase.expectedID, gotID)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	mockUserRepository := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}
	mockPwdGen := &mpg.PasswordGenerator{}
	mockTokenGen := &mtg.TokenGenerator{}

	var userService UserService = NewUserServiceImpl(
		mockUserRepository,
		mockIDGen,
		mockPwdGen,
		mockTokenGen,
	)

	testCases := []struct {
		name             string
		inputPayload     payload.Login
		expectedResponse response.Login
		expectedError    error
		mockBehaviours   func()
	}{
		{
			name: "it should return service.ErrInvalidPayload error, when payload is invalid",
			inputPayload: payload.Login{
				Username: "",
				Password: "erikriosetiawan",
			},
			expectedResponse: response.Login{},
			expectedError:    service.ErrInvalidPayload,
			mockBehaviours:   func() {},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviours()

			gotResponse, gotErr := userService.Login(context.Background(), testCase.inputPayload)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotErr, testCase.expectedError)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, testCase.expectedResponse, gotResponse)
			}
		})
	}
}

package user

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
	mur "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/user/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
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
			name: "it should return service.ErrInvalidPayload error, when email is invalid",
			inputPayload: payload.Register{
				Username: "erikrios",
				Email:    "erikriosetiawan15",
				Name:     "Erik Rio Setiawan",
				Password: "erikriosetiawan",
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
		{
			name: "it should return service.ErrUsernameNotFound error, when repository return an ErrRecordNotFound error",
			inputPayload: payload.Login{
				Username: "erikrios",
				Password: "erikriosetiawan",
			},
			expectedResponse: response.Login{},
			expectedError:    service.ErrUsernameNotFound,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, username string) entity.User {
						return entity.User{}
					},
					func(ctx context.Context, username string) error {
						return repository.ErrRecordNotFound
					},
				).Once()
			},
		},
		{
			name: "it should return service.ErrRepository error, when repository return an ErrDatabase error",
			inputPayload: payload.Login{
				Username: "erikrios",
				Password: "erikriosetiawan",
			},
			expectedResponse: response.Login{},
			expectedError:    service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, username string) entity.User {
						return entity.User{}
					},
					func(ctx context.Context, username string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name: "it should return service.ErrUsernameNotFound error, when user is inactive",
			inputPayload: payload.Login{
				Username: "erikrios",
				Password: "erikriosetiawan",
			},
			expectedResponse: response.Login{},
			expectedError:    service.ErrUsernameNotFound,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, username string) entity.User {
						return entity.User{
							ID:        "u-abcdef",
							Username:  "erikrios",
							Email:     "erikriosetiawan15@gmail.com",
							Name:      "Erik Rio Setiawan",
							Password:  "erikriosetiawan",
							Role:      "user",
							IsActive:  false,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}
					},
					func(ctx context.Context, username string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name: "it should return service.ErrCredentialNotMatch error, when comparing password is return an error",
			inputPayload: payload.Login{
				Username: "erikrios",
				Password: "erikriosetiawan",
			},
			expectedResponse: response.Login{},
			expectedError:    service.ErrCredentialNotMatch,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, username string) entity.User {
						return entity.User{
							ID:        "u-abcdef",
							Username:  "erikrios",
							Email:     "erikriosetiawan15@gmail.com",
							Name:      "Erik Rio Setiawan",
							Password:  "erikriosetiawan",
							Role:      "user",
							IsActive:  true,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}
					},
					func(ctx context.Context, username string) error {
						return nil
					},
				).Once()

				mockPwdGen.On(
					"CompareHashAndPassword",
					mock.AnythingOfType(fmt.Sprintf("%T", []byte{})),
					mock.AnythingOfType(fmt.Sprintf("%T", []byte{})),
				).Return(
					func(hashedPassword []byte, password []byte) error {
						return errors.New("password not match")
					},
				).Once()
			},
		},
		{
			name: "it should return service.ErrRepository error, when generate token return an error",
			inputPayload: payload.Login{
				Username: "erikrios",
				Password: "erikriosetiawan",
			},
			expectedResponse: response.Login{},
			expectedError:    service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, username string) entity.User {
						return entity.User{
							ID:        "u-abcdef",
							Username:  "erikrios",
							Email:     "erikriosetiawan15@gmail.com",
							Name:      "Erik Rio Setiawan",
							Password:  "erikriosetiawan",
							Role:      "user",
							IsActive:  true,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}
					},
					func(ctx context.Context, username string) error {
						return nil
					},
				).Once()

				mockPwdGen.On(
					"CompareHashAndPassword",
					mock.AnythingOfType(fmt.Sprintf("%T", []byte{})),
					mock.AnythingOfType(fmt.Sprintf("%T", []byte{})),
				).Return(
					func(hashedPassword []byte, password []byte) error {
						return nil
					},
				).Once()

				mockTokenGen.On(
					"GenerateToken",
					mock.AnythingOfType(fmt.Sprintf("%T", generator.TokenPayload{})),
				).Return(
					func(generator.TokenPayload) string {
						return ""
					},
					func(generator.TokenPayload) error {
						return errors.New("failed to generate token")
					},
				).Once()
			},
		},
		{
			name: "it should return valid response, when no error is returned",
			inputPayload: payload.Login{
				Username: "erikrios",
				Password: "erikriosetiawan",
			},
			expectedResponse: response.Login{
				Token: "generatedtoken",
				Role:  "user",
			},
			expectedError: nil,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, username string) entity.User {
						return entity.User{
							ID:        "u-abcdef",
							Username:  "erikrios",
							Email:     "erikriosetiawan15@gmail.com",
							Name:      "Erik Rio Setiawan",
							Password:  "erikriosetiawan",
							Role:      "user",
							IsActive:  true,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}
					},
					func(ctx context.Context, username string) error {
						return nil
					},
				).Once()

				mockPwdGen.On(
					"CompareHashAndPassword",
					mock.AnythingOfType(fmt.Sprintf("%T", []byte{})),
					mock.AnythingOfType(fmt.Sprintf("%T", []byte{})),
				).Return(
					func(hashedPassword []byte, password []byte) error {
						return nil
					},
				).Once()

				mockTokenGen.On(
					"GenerateToken",
					mock.AnythingOfType(fmt.Sprintf("%T", generator.TokenPayload{})),
				).Return(
					func(p generator.TokenPayload) string {
						return "generatedtoken"
					},
					func(p generator.TokenPayload) error {
						return nil
					},
				).Once()
			},
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

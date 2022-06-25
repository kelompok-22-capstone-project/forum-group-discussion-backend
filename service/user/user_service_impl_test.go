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
	mtr "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/thread/mocks"
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
	mockThreadRepository := &mtr.ThreadRepository{}
	mockIDGen := &mig.IDGenerator{}
	mockPwdGen := &mpg.PasswordGenerator{}
	mockTokenGen := &mtg.TokenGenerator{}

	var userService UserService = NewUserServiceImpl(
		mockUserRepository,
		mockThreadRepository,
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
	mockThreadRepository := &mtr.ThreadRepository{}
	mockIDGen := &mig.IDGenerator{}
	mockPwdGen := &mpg.PasswordGenerator{}
	mockTokenGen := &mtg.TokenGenerator{}

	var userService UserService = NewUserServiceImpl(
		mockUserRepository,
		mockThreadRepository,
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

func TestGetAll(t *testing.T) {
	mockUserRepository := &mur.UserRepository{}
	mockThreadRepository := &mtr.ThreadRepository{}
	mockIDGen := &mig.IDGenerator{}
	mockPwdGen := &mpg.PasswordGenerator{}
	mockTokenGen := &mtg.TokenGenerator{}

	now := time.Now()

	var userService UserService = NewUserServiceImpl(
		mockUserRepository,
		mockThreadRepository,
		mockIDGen,
		mockPwdGen,
		mockTokenGen,
	)

	testCases := []struct {
		name                string
		inputAccessorUserID string
		inputOrderBy        string
		inputStatus         string
		inputPage           uint
		inputLimit          uint
		inputKeyword        string
		expectedResponse    response.Pagination[response.User]
		expectedError       error
		mockBehaviours      func()
	}{
		{
			name:                "it should return service.ErrRepository, when user repository return an repository.ErrDatabase error",
			inputAccessorUserID: "u-ZrxmQS",
			inputOrderBy:        "ranking",
			inputStatus:         "banned",
			inputPage:           0,
			inputLimit:          0,
			inputKeyword:        "random keyword",
			expectedResponse:    response.Pagination[response.User]{},
			expectedError:       service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindAllWithStatusAndPagination",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.UserOrderBy(entity.Ranking))),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.UserStatus(entity.Banned))),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.PageInfo{})),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						orderBy entity.UserOrderBy,
						userStatus entity.UserStatus,
						pageInfo entity.PageInfo,
						keyword string,
					) entity.Pagination[entity.User] {
						return entity.Pagination[entity.User]{}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						orderBy entity.UserOrderBy,
						userStatus entity.UserStatus,
						pageInfo entity.PageInfo,
						keyword string,
					) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return valid users pagination, when no error is returned",
			inputAccessorUserID: "u-ZrxmQS",
			inputOrderBy:        "registered_date",
			inputStatus:         "active",
			inputPage:           1,
			inputLimit:          35,
			inputKeyword:        "random keyword",
			expectedResponse: response.Pagination[response.User]{
				List: []response.User{
					{
						UserID:         "u-ZrxmQS",
						Username:       "erikrios",
						Email:          "erikriosetiawan@gmail.com",
						Name:           "Erik Rio Setiawan",
						Role:           "user",
						IsActive:       true,
						RegisteredOn:   now.Format(time.RFC822),
						TotalThread:    15,
						TotalFollower:  200,
						TotalFollowing: 2,
						IsFollowed:     false,
					},
				},
				PageInfo: response.PageInfo{
					Limit:     35,
					Page:      1,
					PageTotal: 1,
					Total:     29,
				},
			},
			expectedError: nil,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindAllWithStatusAndPagination",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.UserOrderBy(entity.Ranking))),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.UserStatus(entity.Banned))),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.PageInfo{})),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						orderBy entity.UserOrderBy,
						userStatus entity.UserStatus,
						pageInfo entity.PageInfo,
						keyword string,
					) entity.Pagination[entity.User] {
						return entity.Pagination[entity.User]{
							List: []entity.User{
								{
									ID:             "u-ZrxmQS",
									Username:       "erikrios",
									Email:          "erikriosetiawan@gmail.com",
									Name:           "Erik Rio Setiawan",
									Role:           "user",
									IsActive:       true,
									TotalThread:    15,
									TotalFollower:  200,
									TotalFollowing: 2,
									IsFollowed:     false,
									CreatedAt:      now,
									UpdatedAt:      now,
								},
							},
							PageInfo: entity.PageInfo{
								Limit:     35,
								Page:      1,
								PageTotal: 1,
								Total:     29,
							},
						}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						orderBy entity.UserOrderBy,
						userStatus entity.UserStatus,
						pageInfo entity.PageInfo,
						keyword string,
					) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviours()

			gotPagination, gotError := userService.GetAll(
				context.Background(),
				testCase.inputAccessorUserID,
				testCase.inputOrderBy,
				testCase.inputStatus,
				testCase.inputPage,
				testCase.inputLimit,
				testCase.inputKeyword,
			)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotError, testCase.expectedError)
			} else {
				assert.NoError(t, testCase.expectedError)
				assert.ElementsMatch(t, testCase.expectedResponse.List, gotPagination.List)
				assert.Equal(t, testCase.expectedResponse.PageInfo, gotPagination.PageInfo)
			}
		})
	}
}

func TestGetOwn(t *testing.T) {
	mockUserRepository := &mur.UserRepository{}
	mockThreadRepository := &mtr.ThreadRepository{}
	mockIDGen := &mig.IDGenerator{}
	mockPwdGen := &mpg.PasswordGenerator{}
	mockTokenGen := &mtg.TokenGenerator{}

	now := time.Now()

	var userService UserService = NewUserServiceImpl(
		mockUserRepository,
		mockThreadRepository,
		mockIDGen,
		mockPwdGen,
		mockTokenGen,
	)

	testCases := []struct {
		name                  string
		inputAccessorUserID   string
		inputAccessorUsername string
		expectedResponse      response.User
		expectedError         error
		mockBehaviours        func()
	}{
		{
			name:                  "it should return service.ErrRepository, when user repository return an repository.ErrDatabase error",
			inputAccessorUserID:   "u-ZrxmQS",
			inputAccessorUsername: "erikrios",
			expectedResponse:      response.User{},
			expectedError:         service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsernameWithAccessor",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) entity.User {
						return entity.User{}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                  "it should return valid user, when no error is returned",
			inputAccessorUserID:   "u-ZrxmQS",
			inputAccessorUsername: "erikrios",
			expectedResponse: response.User{
				UserID:         "u-ZrxmQS",
				Username:       "erikrios",
				Email:          "erikriosetiawan@gmail.com",
				Name:           "Erik Rio Setiawan",
				Role:           "user",
				IsActive:       true,
				RegisteredOn:   now.Format(time.RFC822),
				TotalThread:    15,
				TotalFollower:  200,
				TotalFollowing: 2,
				IsFollowed:     false,
			},
			expectedError: nil,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsernameWithAccessor",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) entity.User {
						return entity.User{
							ID:             "u-ZrxmQS",
							Username:       "erikrios",
							Email:          "erikriosetiawan@gmail.com",
							Name:           "Erik Rio Setiawan",
							Role:           "user",
							IsActive:       true,
							TotalThread:    15,
							TotalFollower:  200,
							TotalFollowing: 2,
							IsFollowed:     false,
							CreatedAt:      now,
							UpdatedAt:      now,
						}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviours()

			gotUser, gotError := userService.GetOwn(
				context.Background(),
				testCase.inputAccessorUserID,
				testCase.inputAccessorUsername,
			)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotError, testCase.expectedError)
			} else {
				assert.NoError(t, testCase.expectedError)
				assert.Equal(t, testCase.expectedResponse, gotUser)
			}
		})
	}
}

func TestGetByUsername(t *testing.T) {
	mockUserRepository := &mur.UserRepository{}
	mockThreadRepository := &mtr.ThreadRepository{}
	mockIDGen := &mig.IDGenerator{}
	mockPwdGen := &mpg.PasswordGenerator{}
	mockTokenGen := &mtg.TokenGenerator{}

	now := time.Now()

	var userService UserService = NewUserServiceImpl(
		mockUserRepository,
		mockThreadRepository,
		mockIDGen,
		mockPwdGen,
		mockTokenGen,
	)

	testCases := []struct {
		name                  string
		inputAccessorUserID   string
		inputAccessorUsername string
		expectedResponse      response.User
		expectedError         error
		mockBehaviours        func()
	}{
		{
			name:                  "it should return service.ErrRepository, when user repository return an repository.ErrDatabase error",
			inputAccessorUserID:   "u-ZrxmQS",
			inputAccessorUsername: "erikrios",
			expectedResponse:      response.User{},
			expectedError:         service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsernameWithAccessor",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) entity.User {
						return entity.User{}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                  "it should return valid user, when no error is returned",
			inputAccessorUserID:   "u-ZrxmQS",
			inputAccessorUsername: "erikrios",
			expectedResponse: response.User{
				UserID:         "u-ZrxmQS",
				Username:       "erikrios",
				Email:          "erikriosetiawan@gmail.com",
				Name:           "Erik Rio Setiawan",
				Role:           "user",
				IsActive:       true,
				RegisteredOn:   now.Format(time.RFC822),
				TotalThread:    15,
				TotalFollower:  200,
				TotalFollowing: 2,
				IsFollowed:     false,
			},
			expectedError: nil,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsernameWithAccessor",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) entity.User {
						return entity.User{
							ID:             "u-ZrxmQS",
							Username:       "erikrios",
							Email:          "erikriosetiawan@gmail.com",
							Name:           "Erik Rio Setiawan",
							Role:           "user",
							IsActive:       true,
							TotalThread:    15,
							TotalFollower:  200,
							TotalFollowing: 2,
							IsFollowed:     false,
							CreatedAt:      now,
							UpdatedAt:      now,
						}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviours()

			gotUser, gotError := userService.GetByUsername(
				context.Background(),
				testCase.inputAccessorUserID,
				testCase.inputAccessorUsername,
			)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotError, testCase.expectedError)
			} else {
				assert.NoError(t, testCase.expectedError)
				assert.Equal(t, testCase.expectedResponse, gotUser)
			}
		})
	}
}

func TestChangeBannedState(t *testing.T) {
	mockUserRepository := &mur.UserRepository{}
	mockThreadRepository := &mtr.ThreadRepository{}
	mockIDGen := &mig.IDGenerator{}
	mockPwdGen := &mpg.PasswordGenerator{}
	mockTokenGen := &mtg.TokenGenerator{}

	now := time.Now()

	var userService UserService = NewUserServiceImpl(
		mockUserRepository,
		mockThreadRepository,
		mockIDGen,
		mockPwdGen,
		mockTokenGen,
	)

	testCases := []struct {
		name              string
		inputAccessorRole string
		inputUsername     string
		expectedError     error
		mockBehaviours    func()
	}{
		{
			name:              "it should return service.ErrAccessForbidden, when accessor role is not admin",
			inputAccessorRole: "user",
			inputUsername:     "naruto",
			expectedError:     service.ErrAccessForbidden,
			mockBehaviours:    func() {},
		},
		{
			name:              "it should return service.ErrRepository, when user repository return a repository.ErrDatabase error",
			inputAccessorRole: "admin",
			inputUsername:     "naruto",
			expectedError:     service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						username string,
					) entity.User {
						return entity.User{}
					},
					func(
						ctx context.Context,
						username string,
					) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:              "it should return service.ErrRepository, when banned user return a repository.ErrDatabase error",
			inputAccessorRole: "admin",
			inputUsername:     "naruto",
			expectedError:     service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						username string,
					) entity.User {
						return entity.User{
							ID:             "u-ZrxmQS",
							Username:       "erikrios",
							Email:          "erikriosetiawan@gmail.com",
							Name:           "Erik Rio Setiawan",
							Role:           "user",
							IsActive:       true,
							TotalThread:    15,
							TotalFollower:  200,
							TotalFollowing: 2,
							IsFollowed:     false,
							CreatedAt:      now,
							UpdatedAt:      now,
						}
					},
					func(
						ctx context.Context,
						username string,
					) error {
						return nil
					},
				).Once()

				mockUserRepository.On(
					"BannedUser",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:              "it should return service.ErrRepository, when unbanned user return a repository.ErrDatabase error",
			inputAccessorRole: "admin",
			inputUsername:     "naruto",
			expectedError:     service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						username string,
					) entity.User {
						return entity.User{
							ID:             "u-ZrxmQS",
							Username:       "erikrios",
							Email:          "erikriosetiawan@gmail.com",
							Name:           "Erik Rio Setiawan",
							Role:           "user",
							IsActive:       false,
							TotalThread:    15,
							TotalFollower:  200,
							TotalFollowing: 2,
							IsFollowed:     false,
							CreatedAt:      now,
							UpdatedAt:      now,
						}
					},
					func(
						ctx context.Context,
						username string,
					) error {
						return nil
					},
				).Once()

				mockUserRepository.On(
					"UnbannedUser",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:              "it should return nil error, when no error is returned from user repository",
			inputAccessorRole: "admin",
			inputUsername:     "naruto",
			expectedError:     nil,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						username string,
					) entity.User {
						return entity.User{
							ID:             "u-ZrxmQS",
							Username:       "erikrios",
							Email:          "erikriosetiawan@gmail.com",
							Name:           "Erik Rio Setiawan",
							Role:           "user",
							IsActive:       false,
							TotalThread:    15,
							TotalFollower:  200,
							TotalFollowing: 2,
							IsFollowed:     false,
							CreatedAt:      now,
							UpdatedAt:      now,
						}
					},
					func(
						ctx context.Context,
						username string,
					) error {
						return nil
					},
				).Once()

				mockUserRepository.On(
					"UnbannedUser",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userID string) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviours()

			gotError := userService.ChangeBannedState(
				context.Background(),
				testCase.inputAccessorRole,
				testCase.inputUsername,
			)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotError, testCase.expectedError)
			} else {
				assert.NoError(t, testCase.expectedError)
			}
		})
	}
}

func TestChangeFollowingState(t *testing.T) {
	mockUserRepository := &mur.UserRepository{}
	mockThreadRepository := &mtr.ThreadRepository{}
	mockIDGen := &mig.IDGenerator{}
	mockPwdGen := &mpg.PasswordGenerator{}
	mockTokenGen := &mtg.TokenGenerator{}

	now := time.Now()

	var userService UserService = NewUserServiceImpl(
		mockUserRepository,
		mockThreadRepository,
		mockIDGen,
		mockPwdGen,
		mockTokenGen,
	)

	testCases := []struct {
		name                  string
		inputAccessorUserID   string
		inputUsernameToFollow string
		expectedError         error
		mockBehaviours        func()
	}{
		{
			name:                  "it should return service.ErrRepository, when user repository return a repository.ErrDatabase error",
			inputAccessorUserID:   "u-ZrxmQS",
			inputUsernameToFollow: "naruto",
			expectedError:         service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsernameWithAccessor",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) entity.User {
						return entity.User{}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                  "it should return service.ErrRepository, when unfollow user return a repository.ErrDatabase error",
			inputAccessorUserID:   "u-ZrxmQS",
			inputUsernameToFollow: "erikrios",
			expectedError:         service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsernameWithAccessor",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) entity.User {
						return entity.User{
							ID:             "u-ZrxmQS",
							Username:       "erikrios",
							Email:          "erikriosetiawan@gmail.com",
							Name:           "Erik Rio Setiawan",
							Role:           "user",
							IsActive:       false,
							TotalThread:    15,
							TotalFollower:  200,
							TotalFollowing: 2,
							IsFollowed:     true,
							CreatedAt:      now,
							UpdatedAt:      now,
						}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) error {
						return nil
					},
				).Once()

				mockUserRepository.On(
					"UnfollowUser",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID, userID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                  "it should return service.ErrRepository, when generate id return an error",
			inputAccessorUserID:   "u-ZrxmQS",
			inputUsernameToFollow: "erikrios",
			expectedError:         service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsernameWithAccessor",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) entity.User {
						return entity.User{
							ID:             "u-ZrxmQS",
							Username:       "erikrios",
							Email:          "erikriosetiawan@gmail.com",
							Name:           "Erik Rio Setiawan",
							Role:           "user",
							IsActive:       false,
							TotalThread:    15,
							TotalFollower:  200,
							TotalFollowing: 2,
							IsFollowed:     false,
							CreatedAt:      now,
							UpdatedAt:      now,
						}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) error {
						return nil
					},
				).Once()

				mockIDGen.On(
					"GenerateUserFollowID",
				).Return(
					func() string {
						return ""
					},
					func() error {
						return errors.New("Something went wrong.")
					},
				).Once()
			},
		},
		{
			name:                  "it should return service.ErrRepository, when follow user return a repository.ErrDatabase error",
			inputAccessorUserID:   "u-ZrxmQS",
			inputUsernameToFollow: "erikrios",
			expectedError:         service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsernameWithAccessor",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) entity.User {
						return entity.User{
							ID:             "u-ZrxmQS",
							Username:       "erikrios",
							Email:          "erikriosetiawan@gmail.com",
							Name:           "Erik Rio Setiawan",
							Role:           "user",
							IsActive:       false,
							TotalThread:    15,
							TotalFollower:  200,
							TotalFollowing: 2,
							IsFollowed:     false,
							CreatedAt:      now,
							UpdatedAt:      now,
						}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) error {
						return nil
					},
				).Once()

				mockIDGen.On(
					"GenerateUserFollowID",
				).Return(
					func() string {
						return "f-abcdefg"
					},
					func() error {
						return nil
					},
				).Once()

				mockUserRepository.On(
					"FollowUser",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, id, accessorUserID, userID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                  "it should return nil error, when no error is returned",
			inputAccessorUserID:   "u-ZrxmQS",
			inputUsernameToFollow: "erikrios",
			expectedError:         nil,
			mockBehaviours: func() {
				mockUserRepository.On(
					"FindByUsernameWithAccessor",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) entity.User {
						return entity.User{
							ID:             "u-ZrxmQS",
							Username:       "erikrios",
							Email:          "erikriosetiawan@gmail.com",
							Name:           "Erik Rio Setiawan",
							Role:           "user",
							IsActive:       false,
							TotalThread:    15,
							TotalFollower:  200,
							TotalFollowing: 2,
							IsFollowed:     false,
							CreatedAt:      now,
							UpdatedAt:      now,
						}
					},
					func(
						ctx context.Context,
						accessorUserID string,
						username string,
					) error {
						return nil
					},
				).Once()

				mockIDGen.On(
					"GenerateUserFollowID",
				).Return(
					func() string {
						return "f-abcdefg"
					},
					func() error {
						return nil
					},
				).Once()

				mockUserRepository.On(
					"FollowUser",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, id, accessorUserID, userID string) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviours()

			gotError := userService.ChangeFollowingState(
				context.Background(),
				testCase.inputAccessorUserID,
				testCase.inputUsernameToFollow,
			)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotError, testCase.expectedError)
			} else {
				assert.NoError(t, testCase.expectedError)
			}
		})
	}
}

// func TestGetAllThreadByUsername(t *testing.T) {
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)
// 	err := godotenv.Load("./../../.env.example")
// 	if err != nil {
// 		panic(err)
// 	}

// 	db, err := config.NewPostgreSQLDatabase()
// 	if err != nil {
// 		panic(err)
// 	}

// 	var repo user.UserRepository = user.NewUserRepositoryImpl(db)
// 	var tRepo thread.ThreadRepository = thread.NewThreadRepositoryImpl(db)
// 	var idGen generator.IDGenerator = generator.NewNanoidIDGenerator()
// 	var pwdGen generator.PasswordGenerator = generator.NewBcryptPasswordGenerator()
// 	var tknGen generator.TokenGenerator = generator.NewJWTTokenGenerator()

// 	var service UserService = NewUserServiceImpl(repo, tRepo, idGen, pwdGen, tknGen)

// 	accessorUserID := "u-kt56R1"
// 	username := "erikrios"
// 	page := 1
// 	limit := 20

// 	if pagination, err := service.GetAllThreadByUsername(context.Background(), accessorUserID, username, uint(page), uint(limit)); err != nil {
// 		t.Logf("Error happened: %s", err)
// 	} else {
// 		t.Logf("Pagination: %+v", pagination)
// 	}
// }

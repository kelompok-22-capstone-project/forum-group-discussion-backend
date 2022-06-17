package category

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
	mcr "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/category/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	mig "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAll(t *testing.T) {
	mockRepo := &mcr.CategoryRepository{}
	mockIDGen := &mig.IDGenerator{}

	var categoryService CategoryService = NewCategoryServiceImpl(mockRepo, mockIDGen)
	now := time.Now()

	testCases := []struct {
		name               string
		expectedError      error
		expectedCategories []response.Category
		mockBehaviour      func()
	}{
		{
			name:               "it should return service.ErrRepository, when repository.ErrDatabase return an error",
			expectedError:      service.ErrRepository,
			expectedCategories: []response.Category{},
			mockBehaviour: func() {
				mockRepo.On(
					"FindAll",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
				).Return(
					func(ctx context.Context) []entity.Category {
						return []entity.Category{}
					},
					func(ctx context.Context) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:          "it should return valid categories, when repository return nil error",
			expectedError: nil,
			expectedCategories: []response.Category{
				{
					ID:          "c-abc",
					Name:        "Tech",
					Description: "This is tech category.",
					CreatedOn:   now.Format(time.RFC822),
				},
			},
			mockBehaviour: func() {
				mockRepo.On(
					"FindAll",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
				).Return(
					func(ctx context.Context) []entity.Category {
						return []entity.Category{
							{
								ID:          "c-abc",
								Name:        "Tech",
								Description: "This is tech category.",
								CreatedAt:   now,
								UpdatedAt:   now,
							},
						}
					},
					func(ctx context.Context) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			gotCategories, gotError := categoryService.GetAll(context.Background())

			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotError, testCase.expectedError)
			} else {
				assert.ElementsMatch(t, gotCategories, testCase.expectedCategories)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	mockRepo := &mcr.CategoryRepository{}
	mockIDGen := &mig.IDGenerator{}

	var categoryService CategoryService = NewCategoryServiceImpl(mockRepo, mockIDGen)

	testCases := []struct {
		name              string
		inputPayload      payload.CreateCategory
		inputTokenPayload generator.TokenPayload
		expectedError     error
		expectedID        string
		mockBehaviour     func()
	}{
		{
			name: "it should return service.ErrAccessForbidden, when role is not admin",
			inputPayload: payload.CreateCategory{
				Name:        "Music",
				Description: "THis is a music description",
			},
			inputTokenPayload: generator.TokenPayload{
				ID:       "u-abcd",
				Username: "user",
				Role:     "user",
				IsActive: true,
			},
			expectedError: service.ErrAccessForbidden,
			expectedID:    "",
			mockBehaviour: func() {},
		},
		{
			name: "it should return service.ErrInvalidPayload, when payload is invalid",
			inputPayload: payload.CreateCategory{
				Name:        "M",
				Description: "THis is a music description",
			},
			inputTokenPayload: generator.TokenPayload{
				ID:       "u-abcd",
				Username: "admin",
				Role:     "admin",
				IsActive: true,
			},
			expectedError: service.ErrInvalidPayload,
			expectedID:    "",
			mockBehaviour: func() {},
		},
		{
			name: "it should return service.ErrRepository, when idGenerator return an error",
			inputPayload: payload.CreateCategory{
				Name:        "Music",
				Description: "This is a music description",
			},
			inputTokenPayload: generator.TokenPayload{
				ID:       "u-abcd",
				Username: "admin",
				Role:     "admin",
				IsActive: true,
			},
			expectedError: service.ErrRepository,
			expectedID:    "",
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateCategoryID",
				).Return(
					func() string {
						return ""
					},
					func() error {
						return errors.New("failed to generate category id")
					},
				).Once()
			},
		},
		{
			name: "it should return service.ErrRepository, when repository.ErrDatabase return an error",
			inputPayload: payload.CreateCategory{
				Name:        "Music",
				Description: "This is a music description",
			},
			inputTokenPayload: generator.TokenPayload{
				ID:       "u-abcd",
				Username: "admin",
				Role:     "admin",
				IsActive: true,
			},
			expectedError: service.ErrRepository,
			expectedID:    "",
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateCategoryID",
				).Return(
					func() string {
						return "d-adsf"
					},
					func() error {
						return nil
					},
				).Once()

				mockRepo.On(
					"Insert",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Category{})),
				).Return(
					func(ctx context.Context, p entity.Category) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name: "it should be valid id, when repository return nil error",
			inputPayload: payload.CreateCategory{
				Name:        "Music",
				Description: "This is a music description",
			},
			inputTokenPayload: generator.TokenPayload{
				ID:       "u-abcd",
				Username: "admin",
				Role:     "admin",
				IsActive: true,
			},
			expectedError: nil,
			expectedID:    "d-adsf",
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateCategoryID",
				).Return(
					func() string {
						return "d-adsf"
					},
					func() error {
						return nil
					},
				).Once()

				mockRepo.On(
					"Insert",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Category{})),
				).Return(
					func(ctx context.Context, p entity.Category) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			id, err := categoryService.Create(context.Background(), testCase.inputTokenPayload, testCase.inputPayload)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.Equal(t, testCase.expectedID, id)
			}
		})
	}

}

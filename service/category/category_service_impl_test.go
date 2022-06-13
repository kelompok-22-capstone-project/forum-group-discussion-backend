package category

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
	mcr "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/category/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
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

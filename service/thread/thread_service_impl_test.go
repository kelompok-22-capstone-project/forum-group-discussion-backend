package thread

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
	mtr "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/thread/mocks"
	mur "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/user/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	mig "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAll(t *testing.T) {
	mockThreadRepo := &mtr.ThreadRepository{}
	mockCategoryRepo := &mcr.CategoryRepository{}
	mockUserRepo := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}
	now := time.Now()

	var threadService ThreadService = NewThreadServiceImpl(mockThreadRepo, mockCategoryRepo, mockUserRepo, mockIDGen)

	testCases := []struct {
		name                string
		inputAccessorUserID string
		inputPage           uint
		inputLimit          uint
		inputQuery          string
		expectedError       error
		expectedPagination  response.Pagination[response.ManyThread]
		mockBehaviour       func()
	}{
		{
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputAccessorUserID: "",
			inputPage:           0,
			inputLimit:          0,
			inputQuery:          "",
			expectedError:       service.ErrRepository,
			expectedPagination:  response.Pagination[response.ManyThread]{},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindAllWithQueryAndPagination",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.PageInfo{})),
				).Return(
					func(ctx context.Context,
						accessorUserID string,
						query string,
						pageInfo entity.PageInfo,
					) entity.Pagination[entity.Thread] {
						return entity.Pagination[entity.Thread]{}
					},
					func(ctx context.Context,
						accessorUserID string,
						query string,
						pageInfo entity.PageInfo,
					) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when repository.ErrDatabase return an error",
			inputAccessorUserID: "",
			inputPage:           1,
			inputLimit:          10,
			inputQuery:          "",
			expectedError:       service.ErrRepository,
			expectedPagination:  response.Pagination[response.ManyThread]{},
			mockBehaviour: func() {
				mockCategoryRepo.On(
					"FindAll",
					"FindAll",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
				).Return(
					func(ctx context.Context) []entity.Category {
						return []entity.Category{}
					},
					func(ctx context.Context) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllWithQueryAndPagination",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.PageInfo{})),
				).Return(
					func(ctx context.Context,
						accessorUserID string,
						query string,
						pageInfo entity.PageInfo,
					) entity.Pagination[entity.Thread] {
						return entity.Pagination[entity.Thread]{}
					},
					func(ctx context.Context,
						accessorUserID string,
						query string,
						pageInfo entity.PageInfo,
					) error {
						return service.ErrRepository
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputAccessorUserID: "",
			inputPage:           1,
			inputLimit:          10,
			inputQuery:          "",
			expectedError:       nil,
			expectedPagination: response.Pagination[response.ManyThread]{
				List: []response.ManyThread{
					{
						ID:              "d-Casfkj",
						Title:           "Technology",
						CategoryID:      "g-jMds",
						CategoryName:    "Tech",
						PublishedOn:     now.Format(time.RFC822),
						IsLiked:         false,
						IsFollowed:      false,
						Description:     "Technology is the result of accumulated knowledge and application of skills, methods, and processes used in industrial production and scientific research.",
						TotalViewer:     234,
						TotalLike:       243,
						TotalFollower:   674,
						TotalComment:    23,
						CreatorID:       "d-MDje",
						CreatorUsername: "erikrio",
						CreatorName:     "erik",
					},
				},
				PageInfo: response.PageInfo{
					Limit:     10,
					Page:      1,
					PageTotal: 1,
					Total:     1,
				},
			},
			mockBehaviour: func() {
				mockCategoryRepo.On(
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

				mockThreadRepo.On(
					"FindAllWithQueryAndPagination",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.PageInfo{})),
				).Return(
					func(ctx context.Context,
						accessorUserID string,
						query string,
						pageInfo entity.PageInfo,
					) entity.Pagination[entity.Thread] {
						return entity.Pagination[entity.Thread]{
							List: []entity.Thread{
								{
									ID:            "d-Casfkj",
									Title:         "Technology",
									Description:   "Technology is the result of accumulated knowledge and application of skills, methods, and processes used in industrial production and scientific research.",
									TotalViewer:   234,
									TotalLike:     243,
									TotalFollower: 674,
									TotalComment:  23,
									Creator: entity.User{
										ID:       "d-MDje",
										Username: "erikrio",
										Name:     "erik",
									},
									Category: entity.Category{
										ID:   "g-jMds",
										Name: "Tech",
									},
									IsLiked:    false,
									IsFollowed: false,
									CreatedAt:  now,
								},
							},
							PageInfo: entity.PageInfo{
								Limit:     10,
								Page:      1,
								PageTotal: 1,
								Total:     1,
							},
						}
					},
					func(ctx context.Context,
						accessorUserID string,
						query string,
						pageInfo entity.PageInfo,
					) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			pagination, err := threadService.GetAll(
				context.Background(),
				testCase.inputAccessorUserID,
				testCase.inputPage,
				testCase.inputLimit,
				testCase.inputQuery,
			)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.ElementsMatch(t, pagination.List, testCase.expectedPagination.List)
				assert.Equal(t, pagination.List, testCase.expectedPagination.PageInfo)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	mockThreadRepo := &mtr.ThreadRepository{}
	mockCategoryRepo := &mcr.CategoryRepository{}
	mockUserRepo := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}

	var threadService ThreadService = NewThreadServiceImpl(mockThreadRepo, mockCategoryRepo, mockUserRepo, mockIDGen)

	testCases := []struct {
		name                string
		inputAccessorUserID string
		inputPayload        payload.CreateThread
		expectedError       error
		expectedID          string
		mockBehaviour       func()
	}{
		{
			name:                "it should return service.ErrInvalidPayload, when payload is invalid",
			inputAccessorUserID: "K-amUdnk",
			inputPayload: payload.CreateThread{
				Title:       "T",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "d45Nks",
			},
			expectedError: service.ErrInvalidPayload,
			expectedID:    "",
			mockBehaviour: func() {},
		},
		{
			name:                "it should return service.ErrRepository, repository.ErrDatabase return an error",
			inputAccessorUserID: "K-amUdnk",
			inputPayload: payload.CreateThread{
				Title:       "Technology",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "d45Nks",
			},
			expectedError: service.ErrRepository,
			expectedID:    "",
			mockBehaviour: func() {
				mockCategoryRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, ID string) entity.Category {
						return entity.Category{}
					},
					func(ctx context.Context, ID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when Thread ID Generator return an error",
			inputAccessorUserID: "K-amUdnk",
			inputPayload: payload.CreateThread{
				Title:       "Technology",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "D-abc",
			},
			expectedError: service.ErrRepository,
			expectedID:    "",
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateThreadID",
				).Return(
					func() string {
						return ""
					},
					func() error {
						return errors.New("failed to generate thread id")
					},
				).Once()

				mockCategoryRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, ID string) entity.Category {
						return entity.Category{}
					},
					func(ctx context.Context, ID string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when repository.ErrDatabase return an error",
			inputAccessorUserID: "K-amUdnk",
			inputPayload: payload.CreateThread{
				Title:       "Technology",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "D-abc",
			},
			expectedError: service.ErrRepository,
			expectedID:    "",
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateThreadID",
				).Return(
					func() string {
						return "ahds6sk9"
					},
					func() error {
						return nil
					},
				).Once()

				mockCategoryRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, ID string) entity.Category {
						return entity.Category{}
					},
					func(ctx context.Context, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"Insert",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Thread{})),
				).Return(
					func(ctx context.Context, thread entity.Thread) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when repository.ErrDatabase return an error",
			inputAccessorUserID: "K-amUdnk",
			inputPayload: payload.CreateThread{
				Title:       "Technology",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "D-abc",
			},
			expectedError: service.ErrRepository,
			expectedID:    "",
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateThreadID",
				).Return(
					func() string {
						return "wdsh"
					},
					func() error {
						return nil
					},
				).Once()
				mockIDGen.On(
					"GenerateModeratorID",
				).Return(
					func() string {
						return ""
					},
					func() error {
						return errors.New("failed to generate moderator id")
					},
				).Once()

				mockCategoryRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, ID string) entity.Category {
						return entity.Category{}
					},
					func(ctx context.Context, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"Insert",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Thread{})),
				).Return(
					func(ctx context.Context, thread entity.Thread) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when Moderator ID Generator return an error",
			inputAccessorUserID: "K-amUdnk",
			inputPayload: payload.CreateThread{
				Title:       "Technology",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "D-abc",
			},
			expectedError: service.ErrRepository,
			expectedID:    "",
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateThreadID",
				).Return(
					func() string {
						return "wdsh"
					},
					func() error {
						return nil
					},
				).Once()
				mockIDGen.On(
					"GenerateModeratorID",
				).Return(
					func() string {
						return "dk-Dkj"
					},
					func() error {
						return nil
					},
				).Once()

				mockCategoryRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, ID string) entity.Category {
						return entity.Category{}
					},
					func(ctx context.Context, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"Insert",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Thread{})),
				).Return(
					func(ctx context.Context, thread entity.Thread) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"InsertModerator",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Moderator{})),
				).Return(
					func(ctx context.Context, moderator entity.Moderator) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should be valid id, when repository return nil error",
			inputAccessorUserID: "K-amUdnk",
			inputPayload: payload.CreateThread{
				Title:       "Technology",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "D-abc",
			},
			expectedError: nil,
			expectedID:    "",
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateThreadID",
				).Return(
					func() string {
						return "wdsh"
					},
					func() error {
						return nil
					},
				).Once()
				mockIDGen.On(
					"GenerateModeratorID",
				).Return(
					func() string {
						return "dk-Dkj"
					},
					func() error {
						return nil
					},
				).Once()

				mockCategoryRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, ID string) entity.Category {
						return entity.Category{}
					},
					func(ctx context.Context, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"Insert",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Thread{})),
				).Return(
					func(ctx context.Context, thread entity.Thread) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"InsertModerator",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Moderator{})),
				).Return(
					func(ctx context.Context, moderator entity.Moderator) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			id, err := threadService.Create(context.Background(), testCase.inputAccessorUserID, testCase.inputPayload)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.Equal(t, testCase.expectedID, id)
			}
		})
	}
}

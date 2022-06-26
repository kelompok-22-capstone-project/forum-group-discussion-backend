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
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputAccessorUserID: "",
			inputPage:           1,
			inputLimit:          10,
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
			name:                "it should return nil error, when no error is returned",
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
									Moderators: []entity.Moderator{
										{
											ID: "j-Ksmdh",
											User: entity.User{
												ID:            "d-Casfkj",
												Username:      "tomo12",
												Email:         "tomo@gmail.com",
												Name:          "tomo",
												Password:      "entah",
												Role:          "user",
												IsActive:      false,
												TotalThread:   16,
												TotalFollower: 357,
												IsFollowed:    false,
												CreatedAt:     now,
												UpdatedAt:     now,
											},
											ThreadID:  "h-Kldm",
											CreatedAt: now,
											UpdatedAt: now,
										},
									},
									CreatedAt: now,
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
			name:                "it should return service.ErrRepository, when category repository return a repository.ErrDatabase error",
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
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
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
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
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
			name:                "it should return nil error,  when no error is returned",
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
						return ""
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

func TestGetByID(t *testing.T) {
	mockThreadRepo := &mtr.ThreadRepository{}
	mockCategoryRepo := &mcr.CategoryRepository{}
	mockUserRepo := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}
	now := time.Now()

	var threadService ThreadService = NewThreadServiceImpl(mockThreadRepo, mockCategoryRepo, mockUserRepo, mockIDGen)

	testCases := []struct {
		name                string
		inputAccessorUserID string
		inputID             string
		expectedError       error
		expectedThread      response.Thread
		mockBehaviour       func()
	}{
		{
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputAccessorUserID: "S-MKbduK",
			inputID:             "t-123",
			expectedError:       service.ErrRepository,
			expectedThread:      response.Thread{},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputAccessorUserID: "S-MKbduK",
			inputID:             "t-123",
			expectedError:       service.ErrRepository,
			expectedThread:      response.Thread{},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{}
					},
					func(ctx context.Context, threadID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should be valid id, when the repository return nil error",
			inputAccessorUserID: "u-123",
			inputID:             "t-123",
			expectedError:       nil,
			expectedThread: response.Thread{
				ID:           "d-Casfkj",
				Title:        "Technology",
				CategoryID:   "g-jMds",
				CategoryName: "Tech",
				PublishedOn:  now.Format(time.RFC822),
				IsLiked:      false,
				IsFollowed:   false,
				Moderators: []response.Moderator{
					{
						ID:           "j-Ksmdh",
						UserID:       "d-Casfkj",
						Username:     "tomo12",
						Email:        "tomo@gmail.com",
						Name:         "tomo",
						Role:         "user",
						IsActive:     false,
						RegisteredOn: now.Format(time.RFC822),
					},
				},
				Description:     "Technology is the result of accumulated knowledge and application of skills, methods, and processes used in industrial production and scientific research.",
				TotalViewer:     320,
				TotalLike:       243,
				TotalFollower:   674,
				TotalComment:    23,
				CreatorID:       "d-MDje",
				CreatorUsername: "budi",
				CreatorName:     "budiman",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{
							ID:            "d-Casfkj",
							Title:         "Technology",
							Description:   "Technology is the result of accumulated knowledge and application of skills, methods, and processes used in industrial production and scientific research.",
							TotalViewer:   320,
							TotalLike:     243,
							TotalFollower: 674,
							TotalComment:  23,
							Creator: entity.User{
								ID:       "d-MDje",
								Username: "budi",
								Name:     "budiman",
							},
							Category: entity.Category{
								ID:   "g-jMds",
								Name: "Tech",
							},
							IsLiked:    false,
							IsFollowed: false,
							Moderators: []entity.Moderator{
								{
									ID: "j-Ksmdh",
									User: entity.User{
										ID:         "d-Casfkj",
										Username:   "tomo12",
										Email:      "tomo@gmail.com",
										Name:       "tomo",
										Role:       "user",
										IsActive:   false,
										IsFollowed: false,
									},
								},
							},
							CreatedAt: now,
							UpdatedAt: now,
						}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{
							{
								ID: "j-Ksmdh",
								User: entity.User{
									ID:       "d-Casfkj",
									Username: "tomo12",
									Email:    "tomo@gmail.com",
									Name:     "tomo",
									Role:     "user",
									IsActive: false,
								},
								ThreadID:  "d-Casfkj",
								CreatedAt: now,
								UpdatedAt: now,
							},
						}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"IncrementTotalViewer",
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ID string) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			rs, err := threadService.GetByID(
				context.Background(),
				testCase.inputAccessorUserID,
				testCase.inputID,
			)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.Equal(t, rs, testCase.expectedThread)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	mockThreadRepo := &mtr.ThreadRepository{}
	mockCategoryRepo := &mcr.CategoryRepository{}
	mockUserRepo := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}

	var threadService ThreadService = NewThreadServiceImpl(mockThreadRepo, mockCategoryRepo, mockUserRepo, mockIDGen)

	testCases := []struct {
		name                string
		inputAccessorUserID string
		inputID             string
		inputPayload        payload.UpdateThread
		expectedError       error
		mockBehaviour       func()
	}{
		{
			name:                "it should return service.ErrInvalidPayload, when payload is invalid",
			inputAccessorUserID: "",
			inputID:             "",
			inputPayload: payload.UpdateThread{
				Title:       "T",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "d45Nks",
			},
			expectedError: service.ErrInvalidPayload,
			mockBehaviour: func() {},
		},
		{
			name:                "it should return service.ErrRepository, when category repository return a  repository.ErrDatabase error",
			inputAccessorUserID: "",
			inputID:             "",
			inputPayload: payload.UpdateThread{
				Title:       "Technology",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "d45Nks",
			},
			expectedError: service.ErrRepository,
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
			name:                "it should return service.ErrRepository, when category repository return a repository.ErrDatabase error",
			inputAccessorUserID: "",
			inputID:             "",
			inputPayload: payload.UpdateThread{
				Title:       "Technology",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "d45Nks",
			},
			expectedError: service.ErrRepository,
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
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrAccessForbidden, when thread is not created by the creator id",
			inputAccessorUserID: "d-Casfkj",
			inputID:             "",
			inputPayload: payload.UpdateThread{
				Title:       "Technology",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "d45Nks",
			},
			expectedError: service.ErrAccessForbidden,
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
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()
			},
		},

		{
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputAccessorUserID: "",
			inputID:             "",
			inputPayload: payload.UpdateThread{
				Title:       "Technology",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "d45Nks",
			},
			expectedError: service.ErrRepository,
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
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"Update",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Thread{})),
				).Return(
					func(ctx context.Context, ID string, thread entity.Thread) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return nil error, when no error is returned",
			inputAccessorUserID: "",
			inputID:             "",
			inputPayload: payload.UpdateThread{
				Title:       "Technology",
				Description: "Technology is the result of accumulated knowledge and application of skills, methods, and processes",
				CategoryID:  "d45Nks",
			},
			expectedError: nil,
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
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"Update",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Thread{})),
				).Return(
					func(ctx context.Context, ID string, thread entity.Thread) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			err := threadService.Update(context.Background(), testCase.inputAccessorUserID, testCase.inputID, testCase.inputPayload)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	mockThreadRepo := &mtr.ThreadRepository{}
	mockCategoryRepo := &mcr.CategoryRepository{}
	mockUserRepo := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}

	var threadService ThreadService = NewThreadServiceImpl(mockThreadRepo, mockCategoryRepo, mockUserRepo, mockIDGen)

	testCases := []struct {
		name                string
		inputAccessorUserID string
		inputRole           string
		inputID             string
		expectedError       error
		mockBehaviour       func()
	}{
		{
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputAccessorUserID: "",
			inputRole:           "",
			expectedError:       service.ErrRepository,
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrAccessForbidden, when role is not admin and thread is not created by the creator id",
			inputAccessorUserID: "d-Caf",
			inputRole:           "user",
			expectedError:       service.ErrAccessForbidden,
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputAccessorUserID: "",
			inputRole:           "",
			expectedError:       service.ErrRepository,
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"Delete",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, ID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return nil error,  when no error is returned",
			inputAccessorUserID: "",
			inputRole:           "",
			expectedError:       nil,
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"Delete",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, ID string) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			err := threadService.Delete(context.Background(), testCase.inputAccessorUserID, testCase.inputRole, testCase.inputID)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetComments(t *testing.T) {
	mockThreadRepo := &mtr.ThreadRepository{}
	mockCategoryRepo := &mcr.CategoryRepository{}
	mockUserRepo := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}
	now := time.Now()

	var threadService ThreadService = NewThreadServiceImpl(mockThreadRepo, mockCategoryRepo, mockUserRepo, mockIDGen)

	testCases := []struct {
		name               string
		inputThreadID      string
		inputPage          uint
		inputLimit         uint
		expectedError      error
		expectedPagination response.Pagination[response.Comment]
		mockBehaviour      func()
	}{
		{
			name:               "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputThreadID:      "",
			inputPage:          0,
			inputLimit:         0,
			expectedError:      service.ErrRepository,
			expectedPagination: response.Pagination[response.Comment]{},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:               "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputThreadID:      "",
			inputPage:          0,
			inputLimit:         0,
			expectedError:      service.ErrRepository,
			expectedPagination: response.Pagination[response.Comment]{},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllCommentByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.PageInfo{})),
				).Return(
					func(ctx context.Context, threadID string, pageInfo entity.PageInfo) entity.Pagination[entity.Comment] {
						return entity.Pagination[entity.Comment]{}
					},
					func(ctx context.Context, threadID string, pageInfo entity.PageInfo) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:          "it should return nil error, when no error is returned",
			inputThreadID: "t-123",
			inputPage:     1,
			inputLimit:    3,
			expectedError: nil,
			expectedPagination: response.Pagination[response.Comment]{
				List: []response.Comment{
					{
						ID:          "c-123",
						UserID:      "u-123",
						Username:    "someone",
						Name:        "Jane Doe",
						Comment:     "Jane Doe commented",
						PublishedOn: now.Format(time.RFC822),
					},
				},
				PageInfo: response.PageInfo{Limit: 10, Page: 1, PageTotal: 1, Total: 1},
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllCommentByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.PageInfo{})),
				).Return(
					func(ctx context.Context, threadID string, pageInfo entity.PageInfo) entity.Pagination[entity.Comment] {
						return entity.Pagination[entity.Comment]{
							List: []entity.Comment{
								{
									ID:        "c-123",
									User:      entity.User{ID: "u-123", Username: "someone", Name: "Jane Doe"},
									Comment:   "Jane Doe commented",
									CreatedAt: now,
								},
							},
							PageInfo: entity.PageInfo{Limit: 10, Page: 1, PageTotal: 1, Total: 1},
						}
					},
					func(ctx context.Context, threadID string, pageInfo entity.PageInfo) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			pagination, err := threadService.GetComments(context.Background(), testCase.inputThreadID, testCase.inputPage, testCase.inputLimit)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.ElementsMatch(t, pagination.List, testCase.expectedPagination.List)
				assert.Equal(t, pagination.PageInfo, testCase.expectedPagination.PageInfo)
			}
		})
	}
}

func TestCreateComment(t *testing.T) {
	mockThreadRepo := &mtr.ThreadRepository{}
	mockCategoryRepo := &mcr.CategoryRepository{}
	mockUserRepo := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}

	var threadService ThreadService = NewThreadServiceImpl(mockThreadRepo, mockCategoryRepo, mockUserRepo, mockIDGen)

	testCases := []struct {
		name                string
		inputThreadID       string
		inputAccessorUserID string
		expectedError       error
		inputPayload        payload.CreateComment
		expectedID          string
		mockBehaviour       func()
	}{
		{
			name:                "it should return service.ErrInvalidPayload, when payload is invalid",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrInvalidPayload,
			inputPayload:        payload.CreateComment{},
			mockBehaviour:       func() {},
		},
		{
			name:                "it should return service.ErrRepository, repository.ErrDatabase return an error",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrRepository,
			inputPayload: payload.CreateComment{
				Comment: "nice",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when Comment ID Generator return an error",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrRepository,
			inputPayload: payload.CreateComment{
				Comment: "nice",
			},
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateCommentID",
				).Return(
					func() string {
						return ""
					},
					func() error {
						return errors.New("failed to generate comment id")
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when repository.ErrDatabase return an error",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrRepository,
			inputPayload: payload.CreateComment{
				Comment: "nice",
			},
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateCommentID",
				).Return(
					func() string {
						return "X-aBc"
					},
					func() error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"InsertComment",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Comment{})),
				).Return(
					func(ctx context.Context, comment entity.Comment) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should be valid id, when repository return nil error",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       nil,
			inputPayload: payload.CreateComment{
				Comment: "nice",
			},
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateCommentID",
				).Return(
					func() string {
						return ""
					},
					func() error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"InsertComment",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Comment{})),
				).Return(
					func(ctx context.Context, comment entity.Comment) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			id, err := threadService.CreateComment(context.Background(), testCase.inputThreadID, testCase.inputAccessorUserID, testCase.inputPayload)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.Equal(t, testCase.expectedID, id)
			}
		})
	}
}

func TestChangeFollowingState(t *testing.T) {
	mockThreadRepo := &mtr.ThreadRepository{}
	mockCategoryRepo := &mcr.CategoryRepository{}
	mockUserRepo := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}
	now := time.Now()

	var threadService ThreadService = NewThreadServiceImpl(mockThreadRepo, mockCategoryRepo, mockUserRepo, mockIDGen)

	testCases := []struct {
		name                string
		inputThreadID       string
		inputAccessorUserID string
		expectedError       error
		mockBehaviour       func()
	}{
		{
			name:                "it should return service.ErrRepository, when repository.ErrDatabase return an error",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrRepository,
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when Thread Follow ID Generator return an error",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrRepository,
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateThreadFollowID",
				).Return(
					func() string {
						return ""
					},
					func() error {
						return errors.New("failed to generate comment id")
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when Thread repository return a repository.ErrDatabase error.",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrRepository,
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateThreadFollowID",
				).Return(
					func() string {
						return "x-abc"
					},
					func() error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{
							ID:            "d-Casfkj",
							Title:         "Technology",
							Description:   "Technology is the result of accumulated knowledge and application of skills, methods, and processes used in industrial production and scientific research.",
							TotalViewer:   320,
							TotalLike:     243,
							TotalFollower: 674,
							TotalComment:  23,
							IsLiked:       false,
							IsFollowed:    true,
							CreatedAt:     now,
							UpdatedAt:     now,
						}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"DeleteFollowThread",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.ThreadFollow{})),
				).Return(
					func(ctx context.Context, threadFollow entity.ThreadFollow) error {
						return repository.ErrDatabase
					},
				).Once()

				mockThreadRepo.On(
					"InsertFollowThread",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.ThreadFollow{})),
				).Return(
					func(ctx context.Context, threadFollow entity.ThreadFollow) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when Thread repository return a repository.ErrDatabase error.",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrRepository,
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateThreadFollowID",
				).Return(
					func() string {
						return "x-abc"
					},
					func() error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"DeleteFollowThread",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.ThreadFollow{})),
				).Return(
					func(ctx context.Context, threadFollow entity.ThreadFollow) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"InsertFollowThread",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.ThreadFollow{})),
				).Return(
					func(ctx context.Context, threadFollow entity.ThreadFollow) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return nil error, when no error is returned.",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       nil,
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateThreadFollowID",
				).Return(
					func() string {
						return "x-abc"
					},
					func() error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{
							ID:            "d-Casfkj",
							Title:         "Technology",
							Description:   "Technology is the result of accumulated knowledge and application of skills, methods, and processes used in industrial production and scientific research.",
							TotalViewer:   320,
							TotalLike:     243,
							TotalFollower: 674,
							TotalComment:  23,
							IsLiked:       false,
							IsFollowed:    true,
							CreatedAt:     now,
							UpdatedAt:     now,
						}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"DeleteFollowThread",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.ThreadFollow{})),
				).Return(
					func(ctx context.Context, threadFollow entity.ThreadFollow) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"InsertFollowThread",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.ThreadFollow{})),
				).Return(
					func(ctx context.Context, threadFollow entity.ThreadFollow) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			err := threadService.ChangeFollowingState(context.Background(), testCase.inputThreadID, testCase.inputAccessorUserID)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChangeLikeState(t *testing.T) {
	mockThreadRepo := &mtr.ThreadRepository{}
	mockCategoryRepo := &mcr.CategoryRepository{}
	mockUserRepo := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}
	now := time.Now()

	var threadService ThreadService = NewThreadServiceImpl(mockThreadRepo, mockCategoryRepo, mockUserRepo, mockIDGen)

	testCases := []struct {
		name                string
		inputThreadID       string
		inputAccessorUserID string
		expectedError       error
		mockBehaviour       func()
	}{
		{
			name:                "it should return service.ErrRepository, when repository.ErrDatabase return an error",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrRepository,
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when Thread Follow ID Generator return an error",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrRepository,
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateLikeID",
				).Return(
					func() string {
						return ""
					},
					func() error {
						return errors.New("failed to generate comment id")
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when Thread repository return a repository.ErrDatabase error",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrRepository,
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateLikeID",
				).Return(
					func() string {
						return "P-sk8d"
					},
					func() error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{
							ID:            "d-Casfkj",
							Title:         "Technology",
							Description:   "Technology is the result of accumulated knowledge and application of skills, methods, and processes used in industrial production and scientific research.",
							TotalViewer:   320,
							TotalLike:     243,
							TotalFollower: 674,
							TotalComment:  23,
							IsLiked:       true,
							IsFollowed:    true,
							CreatedAt:     now,
							UpdatedAt:     now,
						}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"DeleteLike",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Like{})),
				).Return(
					func(ctx context.Context, like entity.Like) error {
						return repository.ErrDatabase
					},
				).Once()

				mockThreadRepo.On(
					"InsertLike",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Like{})),
				).Return(
					func(ctx context.Context, like entity.Like) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when Thread repository return a repository.ErrDatabase error",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       service.ErrRepository,
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateLikeID",
				).Return(
					func() string {
						return "P-sk8d"
					},
					func() error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"DeleteLike",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Like{})),
				).Return(
					func(ctx context.Context, like entity.Like) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"InsertLike",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Like{})),
				).Return(
					func(ctx context.Context, like entity.Like) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return nil error, when no error is returned.",
			inputThreadID:       "",
			inputAccessorUserID: "",
			expectedError:       nil,
			mockBehaviour: func() {
				mockIDGen.On(
					"GenerateLikeID",
				).Return(
					func() string {
						return "P-sk8d"
					},
					func() error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{
							ID:            "d-Casfkj",
							Title:         "Technology",
							Description:   "Technology is the result of accumulated knowledge and application of skills, methods, and processes used in industrial production and scientific research.",
							TotalViewer:   320,
							TotalLike:     243,
							TotalFollower: 674,
							TotalComment:  23,
							IsLiked:       true,
							IsFollowed:    true,
							CreatedAt:     now,
							UpdatedAt:     now,
						}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"DeleteLike",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Like{})),
				).Return(
					func(ctx context.Context, like entity.Like) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"InsertLike",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.Like{})),
				).Return(
					func(ctx context.Context, like entity.Like) error {
						return nil
					},
				).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			err := threadService.ChangeLikeState(context.Background(), testCase.inputThreadID, testCase.inputAccessorUserID)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddModerator(t *testing.T) {
	mockThreadRepo := &mtr.ThreadRepository{}
	mockCategoryRepo := &mcr.CategoryRepository{}
	mockUserRepo := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}

	var threadService ThreadService = NewThreadServiceImpl(mockThreadRepo, mockCategoryRepo, mockUserRepo, mockIDGen)

	testCases := []struct {
		name                string
		inputThreadID       string
		inputAccessorUserID string
		expectedError       error
		expectedID          string
		inputPayload        payload.AddRemoveModerator
		mockBehaviour       func()
	}{
		{
			name:                "it should return service.ErrInvalidPayload, when payload is invalid",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-1243",
			expectedError:       service.ErrInvalidPayload,
			inputPayload:        payload.AddRemoveModerator{},
			mockBehaviour:       func() {},
		},
		{
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-1243",
			expectedError:       service.ErrRepository,
			expectedID:          "",
			inputPayload: payload.AddRemoveModerator{
				Username: "Jody",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrAccessForbidden, when accessor user ID is not thread creator ID",
			inputThreadID:       "U-abcd",
			inputAccessorUserID: "b-123",
			expectedError:       service.ErrAccessForbidden,
			inputPayload: payload.AddRemoveModerator{
				Username: "Jody",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "b-124"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-1243",
			expectedError:       service.ErrRepository,
			inputPayload: payload.AddRemoveModerator{
				Username: "tomo12",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-1243"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{}
					},
					func(ctx context.Context, threadID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when user repository return a repository.ErrDatabase error",
			inputThreadID:       "t-xyz",
			inputAccessorUserID: "u-1243",
			expectedError:       service.ErrRepository,
			inputPayload: payload.AddRemoveModerator{
				Username: "Jody",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-1243"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockUserRepo.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userName string) entity.User {
						return entity.User{}
					},
					func(ctx context.Context, userName string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrAccessForbidden, when accessorUserID is userToAddID",
			inputThreadID:       "t-xyz",
			inputAccessorUserID: "u-1243",
			expectedError:       service.ErrAccessForbidden,
			inputPayload: payload.AddRemoveModerator{
				Username: "tomo12",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-1243"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockUserRepo.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userName string) entity.User {
						return entity.User{ID: "u-1243"}
					},
					func(ctx context.Context, userName string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrDataAlreadyExists, when mod.User.ID is userToAdded.ID",
			inputThreadID:       "t-xyz",
			inputAccessorUserID: "u-1243",
			expectedError:       service.ErrDataAlreadyExists,
			inputPayload: payload.AddRemoveModerator{
				Username: "tomo12",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-1243"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{
							{
								User: entity.User{ID: "u-1244"},
							},
						}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockUserRepo.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userName string) entity.User {
						return entity.User{ID: "u-1244"}
					},
					func(ctx context.Context, userName string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when generate ID return an error",
			inputThreadID:       "t-xyz",
			inputAccessorUserID: "u-1243",
			expectedError:       service.ErrRepository,
			inputPayload: payload.AddRemoveModerator{
				Username: "tomo12",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-1243"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{
							{
								User: entity.User{ID: "u-1244"},
							},
						}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockUserRepo.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userName string) entity.User {
						return entity.User{ID: "u-1245"}
					},
					func(ctx context.Context, userName string) error {
						return nil
					},
				).Once()

				mockIDGen.On("GenerateModeratorID").Return(
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
			name:                "it should return service.ErrRepository, when insert moderator return repository.ErrDatabase error",
			inputThreadID:       "t-xyz",
			inputAccessorUserID: "u-1243",
			expectedError:       service.ErrRepository,
			inputPayload: payload.AddRemoveModerator{
				Username: "tomo12",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-1243"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{
							{
								User: entity.User{ID: "u-1244"},
							},
						}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockUserRepo.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userName string) entity.User {
						return entity.User{ID: "u-1245"}
					},
					func(ctx context.Context, userName string) error {
						return nil
					},
				).Once()

				mockIDGen.On("GenerateModeratorID").Return(
					func() string {
						return "m-123"
					},
					func() error {
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
			name:                "it should return valid ID, when no error is returned",
			inputThreadID:       "t-xyz",
			inputAccessorUserID: "u-1243",
			expectedError:       nil,
			inputPayload: payload.AddRemoveModerator{
				Username: "tomo12",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-1243"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{
							{
								User: entity.User{ID: "u-1244"},
							},
						}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockUserRepo.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userName string) entity.User {
						return entity.User{ID: "u-1245"}
					},
					func(ctx context.Context, userName string) error {
						return nil
					},
				).Once()

				mockIDGen.On("GenerateModeratorID").Return(
					func() string {
						return "m-123"
					},
					func() error {
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

			err := threadService.AddModerator(context.Background(), testCase.inputPayload, testCase.inputThreadID, testCase.inputAccessorUserID)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveModerator(t *testing.T) {
	mockThreadRepo := &mtr.ThreadRepository{}
	mockCategoryRepo := &mcr.CategoryRepository{}
	mockUserRepo := &mur.UserRepository{}
	mockIDGen := &mig.IDGenerator{}

	var threadService ThreadService = NewThreadServiceImpl(mockThreadRepo, mockCategoryRepo, mockUserRepo, mockIDGen)

	testCases := []struct {
		name                string
		inputThreadID       string
		inputAccessorUserID string
		expectedError       error
		inputPayload        payload.AddRemoveModerator
		mockBehaviour       func()
	}{
		{
			name:                "it should return service.ErrInvalidPayload, when payload is invalid",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-abcde",
			expectedError:       service.ErrInvalidPayload,
			inputPayload:        payload.AddRemoveModerator{},
			mockBehaviour:       func() {},
		},
		{
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase error",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-abcde",
			expectedError:       service.ErrRepository,
			inputPayload: payload.AddRemoveModerator{
				Username: "tomo12",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrAccessForbidden, when accessor user ID is not thread creator ID",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-abcde",
			expectedError:       service.ErrAccessForbidden,
			inputPayload: payload.AddRemoveModerator{
				Username: "Jody",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-abcdd"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when thread repository return a repository.ErrDatabase",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-abcde",
			expectedError:       service.ErrRepository,
			inputPayload: payload.AddRemoveModerator{
				Username: "Jody",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-abcde"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{}
					},
					func(ctx context.Context, threadID string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when user repository return a repository.ErrDatabase",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-abcde",
			expectedError:       service.ErrRepository,
			inputPayload: payload.AddRemoveModerator{
				Username: "Jody",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-abcde"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockUserRepo.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userName string) entity.User {
						return entity.User{}
					},
					func(ctx context.Context, userName string) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrAccessForbidden, when accessor want to remove itself",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-abcde",
			expectedError:       service.ErrAccessForbidden,
			inputPayload: payload.AddRemoveModerator{
				Username: "Jody",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-abcde"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockUserRepo.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userName string) entity.User {
						return entity.User{ID: "u-abcde"}
					},
					func(ctx context.Context, userName string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrUsernameNotFound, when user is not moderator",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-abcde",
			expectedError:       service.ErrUsernameNotFound,
			inputPayload: payload.AddRemoveModerator{
				Username: "Jody",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-abcde"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{{
							User: entity.User{
								ID: "u-abcde",
							},
						}}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockUserRepo.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userName string) entity.User {
						return entity.User{ID: "u-abcdf"}
					},
					func(ctx context.Context, userName string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when delete moderator return repository.ErrDatabase error",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-abcde",
			expectedError:       service.ErrRepository,
			inputPayload: payload.AddRemoveModerator{
				Username: "Jody",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-abcde"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{{
							User: entity.User{
								ID: "u-abcdf",
							},
						}}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockUserRepo.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userName string) entity.User {
						return entity.User{ID: "u-abcdf"}
					},
					func(ctx context.Context, userName string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"DeleteModerator",
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
			name:                "it should return nil error, when no error is returned",
			inputThreadID:       "t-123",
			inputAccessorUserID: "u-abcde",
			expectedError:       nil,
			inputPayload: payload.AddRemoveModerator{
				Username: "Jody",
			},
			mockBehaviour: func() {
				mockThreadRepo.On(
					"FindByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, accessorUserID string, ID string) entity.Thread {
						return entity.Thread{Creator: entity.User{ID: "u-abcde"}}
					},
					func(ctx context.Context, accessorUserID string, ID string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"FindAllModeratorByThreadID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, threadID string) []entity.Moderator {
						return []entity.Moderator{{
							User: entity.User{
								ID: "u-abcdf",
							},
						}}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockUserRepo.On(
					"FindByUsername",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, userName string) entity.User {
						return entity.User{ID: "u-abcdf"}
					},
					func(ctx context.Context, userName string) error {
						return nil
					},
				).Once()

				mockThreadRepo.On(
					"DeleteModerator",
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

			err := threadService.RemoveModerator(context.Background(), testCase.inputPayload, testCase.inputThreadID, testCase.inputAccessorUserID)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

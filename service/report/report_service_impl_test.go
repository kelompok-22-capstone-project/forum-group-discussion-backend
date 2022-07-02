package report

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
	mr "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/report/mocks"
	mt "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/thread/mocks"
	mu "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/user/mocks"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	mi "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAll(t *testing.T) {
	mockReportRepo := &mr.ReportRepository{}
	mockUserRepo := &mu.UserRepository{}
	mockThreadRepo := &mt.ThreadRepository{}
	mockIDGen := &mi.IDGenerator{}

	var reportService ReportService = NewReportServiceImpl(mockReportRepo, mockUserRepo, mockThreadRepo, mockIDGen)

	now := time.Now()

	testCases := []struct {
		name              string
		inputAccessorRole string
		inputStatus       string
		inputPage         uint
		inputLimit        uint
		expectedResponse  response.Pagination[response.Report]
		expectedError     error
		mockBehaviours    func()
	}{
		{
			name:              "it should return service.ErrAccessForbidden, when accessor role is not admin",
			inputAccessorRole: "user",
			inputStatus:       "accepted",
			inputPage:         1,
			inputLimit:        10,
			expectedResponse:  response.Pagination[response.Report]{},
			expectedError:     service.ErrAccessForbidden,
			mockBehaviours:    func() {},
		},
		{
			name:              "it should return service.ErrRepository, when report repository return repository.ErrDatabase error",
			inputAccessorRole: "admin",
			inputStatus:       "accepted",
			inputPage:         0,
			inputLimit:        0,
			expectedResponse:  response.Pagination[response.Report]{},
			expectedError:     service.ErrRepository,
			mockBehaviours: func() {
				mockReportRepo.On(
					"GetReportsWithPagination",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.PageInfo{})),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.ReportStatus(entity.Accepted))),
				).Return(
					func(
						ctx context.Context,
						pageInfo entity.PageInfo,
						reportStatus entity.ReportStatus,
					) entity.Pagination[entity.UserBanned] {
						return entity.Pagination[entity.UserBanned]{}
					},
					func(
						ctx context.Context,
						pageInfo entity.PageInfo,
						reportStatus entity.ReportStatus,
					) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:              "it should return valid reports pagination, when no error is returned",
			inputAccessorRole: "admin",
			inputStatus:       "review",
			inputPage:         1,
			inputLimit:        20,
			expectedResponse: response.Pagination[response.Report]{
				List: []response.Report{
					{
						ID:                 "r-ErLN4lS",
						ModeratorID:        "m-QwROlyYS",
						ModeratorUsername:  "erikrios",
						ModeratorName:      "Erik Rio Setiawan",
						UserID:             "u-ZrxmQq",
						Username:           "naruto",
						Name:               "Naruto Uzumaki",
						Reason:             "Harrashment",
						Status:             "review",
						ThreadID:           "t-Ku7Pi",
						ThreadTitle:        "Go Programming Going Hype",
						ReportedOn:         now.Format(time.RFC822),
						Comment:            "Bocil ini sangat meresahkan",
						CommentPublishedOn: now.Format(time.RFC822),
					},
				},
				PageInfo: response.PageInfo{
					Limit:     20,
					Page:      1,
					PageTotal: 1,
					Total:     15,
				},
			},
			expectedError: nil,
			mockBehaviours: func() {
				mockReportRepo.On(
					"GetReportsWithPagination",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.PageInfo{})),
					mock.AnythingOfType(fmt.Sprintf("%T", entity.ReportStatus(entity.Accepted))),
				).Return(
					func(
						ctx context.Context,
						pageInfo entity.PageInfo,
						reportStatus entity.ReportStatus,
					) entity.Pagination[entity.UserBanned] {
						return entity.Pagination[entity.UserBanned]{
							List: []entity.UserBanned{
								{
									ID: "r-ErLN4lS",
									Moderator: entity.Moderator{
										ID: "m-QwROlyYS",
										User: entity.User{
											Username: "erikrios",
											Name:     "Erik Rio Setiawan",
										},
									},
									User: entity.User{
										ID:       "u-ZrxmQq",
										Username: "naruto",
										Name:     "Naruto Uzumaki",
									},
									Thread: entity.Thread{
										ID:    "t-Ku7Pi",
										Title: "Go Programming Going Hype",
									},
									Reason:    "Harrashment",
									Status:    "review",
									CreatedAt: now,
									Comment: entity.Comment{
										Comment:   "Bocil ini sangat meresahkan",
										CreatedAt: now,
										UpdatedAt: now,
									},
								},
							},
							PageInfo: entity.PageInfo{
								Limit:     20,
								Page:      1,
								PageTotal: 1,
								Total:     15,
							},
						}
					},
					func(
						ctx context.Context,
						pageInfo entity.PageInfo,
						reportStatus entity.ReportStatus,
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

			gotPagination, gotErr := reportService.GetAll(
				context.Background(),
				testCase.inputAccessorRole,
				testCase.inputStatus,
				testCase.inputPage,
				testCase.inputLimit,
			)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotErr, testCase.expectedError)
			} else {
				assert.NoError(t, gotErr)
				assert.ElementsMatch(t, testCase.expectedResponse.List, gotPagination.List)
				assert.Equal(t, testCase.expectedResponse.PageInfo, gotPagination.PageInfo)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	mockReportRepo := &mr.ReportRepository{}
	mockUserRepo := &mu.UserRepository{}
	mockThreadRepo := &mt.ThreadRepository{}
	mockIDGen := &mi.IDGenerator{}

	var reportService ReportService = NewReportServiceImpl(mockReportRepo, mockUserRepo, mockThreadRepo, mockIDGen)

	testCases := []struct {
		name                string
		inputAccessorUserID string
		inputPayload        payload.CreateReport
		expectedID          string
		expectedError       error
		mockBehaviours      func()
	}{
		{
			name:                "it should return service.ErrInvalidPayload, when the payload is invalid",
			inputAccessorUserID: "u-ZrxmQS",
			inputPayload: payload.CreateReport{
				Username:  "erikrios",
				CommentID: "c-Ku7Pi",
				Reason:    "",
			},
			expectedID:     "",
			expectedError:  service.ErrInvalidPayload,
			mockBehaviours: func() {},
		},
		{
			name:                "it should return service.ErrDataNotFound, when find by username return repository.ErrRecordNotFound error",
			inputAccessorUserID: "u-ZrxmQS",
			inputPayload: payload.CreateReport{
				Username:  "erikrios",
				CommentID: "c-Ku7Pi",
				Reason:    "Harrashment",
			},
			expectedID:    "",
			expectedError: service.ErrDataNotFound,
			mockBehaviours: func() {
				mockUserRepo.On(
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
			name:                "it should return service.ErrDataNotFound, when find comment by ID return repository.ErrRecordNotFound error",
			inputAccessorUserID: "u-ZrxmQS",
			inputPayload: payload.CreateReport{
				Username:  "erikrios",
				CommentID: "c-Ku7Pi",
				Reason:    "Harrashment",
			},
			expectedID:    "",
			expectedError: service.ErrDataNotFound,
			mockBehaviours: func() {
				mockUserRepo.On(
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

				mockThreadRepo.On(
					"FindCommentByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, commentID string) entity.Comment {
						return entity.Comment{}
					},
					func(ctx context.Context, commentID string) error {
						return repository.ErrRecordNotFound
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when find all moderators by thread ID return repository.ErrDatabase error",
			inputAccessorUserID: "u-ZrxmQS",
			inputPayload: payload.CreateReport{
				Username:  "erikrios",
				CommentID: "c-Ku7Pi",
				Reason:    "Harrashment",
			},
			expectedID:    "",
			expectedError: service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepo.On(
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

				mockThreadRepo.On(
					"FindCommentByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, commentID string) entity.Comment {
						return entity.Comment{}
					},
					func(ctx context.Context, commentID string) error {
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
			name:                "it should return service.ErrAccessForbidden, when accessor is not moderator",
			inputAccessorUserID: "u-ZrxmQS",
			inputPayload: payload.CreateReport{
				Username:  "erikrios",
				CommentID: "c-Ku7Pi",
				Reason:    "Harrashment",
			},
			expectedID:    "",
			expectedError: service.ErrAccessForbidden,
			mockBehaviours: func() {
				mockUserRepo.On(
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

				mockThreadRepo.On(
					"FindCommentByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, commentID string) entity.Comment {
						return entity.Comment{}
					},
					func(ctx context.Context, commentID string) error {
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
								User: entity.User{
									ID: "u-ZrxmQq",
								},
							},
						}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()
			},
		},
		{
			name:                "it should return service.ErrRepository, when generate report ID return an error",
			inputAccessorUserID: "u-ZrxmQS",
			inputPayload: payload.CreateReport{
				Username:  "erikrios",
				CommentID: "c-Ku7Pi",
				Reason:    "Harrashment",
			},
			expectedID:    "",
			expectedError: service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepo.On(
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

				mockThreadRepo.On(
					"FindCommentByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, commentID string) entity.Comment {
						return entity.Comment{}
					},
					func(ctx context.Context, commentID string) error {
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
								User: entity.User{
									ID: "u-ZrxmQS",
								},
							},
						}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockIDGen.On(
					"GenerateReportID",
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
			name:                "it should return service.ErrRepository, when insert report return repository.ErrDatabase",
			inputAccessorUserID: "u-ZrxmQS",
			inputPayload: payload.CreateReport{
				Username:  "erikrios",
				CommentID: "c-Ku7Pi",
				Reason:    "Harrashment",
			},
			expectedID:    "",
			expectedError: service.ErrRepository,
			mockBehaviours: func() {
				mockUserRepo.On(
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

				mockThreadRepo.On(
					"FindCommentByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, commentID string) entity.Comment {
						return entity.Comment{}
					},
					func(ctx context.Context, commentID string) error {
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
								User: entity.User{
									ID: "u-ZrxmQS",
								},
							},
						}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockIDGen.On(
					"GenerateReportID",
				).Return(
					func() string {
						return "r-ErLN4lS"
					},
					func() error {
						return nil
					},
				).Once()

				mockReportRepo.On(
					"Insert",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						ID,
						moderatorID,
						userID,
						threadID,
						reason string,
					) error {
						return repository.ErrDatabase
					},
				).Once()
			},
		},
		{
			name:                "it should return valid ID, when no error is returned",
			inputAccessorUserID: "u-ZrxmQS",
			inputPayload: payload.CreateReport{
				Username:  "erikrios",
				CommentID: "c-Ku7Pi",
				Reason:    "Harrashment",
			},
			expectedID:    "r-ErLN4lS",
			expectedError: nil,
			mockBehaviours: func() {
				mockUserRepo.On(
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

				mockThreadRepo.On(
					"FindCommentByID",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(ctx context.Context, commentID string) entity.Comment {
						return entity.Comment{}
					},
					func(ctx context.Context, commentID string) error {
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
								User: entity.User{
									ID: "u-ZrxmQS",
								},
							},
						}
					},
					func(ctx context.Context, threadID string) error {
						return nil
					},
				).Once()

				mockIDGen.On(
					"GenerateReportID",
				).Return(
					func() string {
						return "r-ErLN4lS"
					},
					func() error {
						return nil
					},
				).Once()

				mockReportRepo.On(
					"Insert",
					mock.AnythingOfType(fmt.Sprintf("%T", context.Background())),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
					mock.AnythingOfType(fmt.Sprintf("%T", "")),
				).Return(
					func(
						ctx context.Context,
						ID,
						moderatorID,
						userID,
						threadID,
						reason string,
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

			gotID, gotErr := reportService.Create(context.Background(), testCase.inputAccessorUserID, testCase.inputPayload)

			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotErr, testCase.expectedError)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, testCase.expectedID, gotID)
			}
		})
	}
}

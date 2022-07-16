package report

import (
	"context"
	"time"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/report"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/thread"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/user"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	"gopkg.in/validator.v2"
)

type reportServiceImpl struct {
	reportRepository report.ReportRepository
	userRepository   user.UserRepository
	threadRepository thread.ThreadRepository
	idGenerator      generator.IDGenerator
}

func NewReportServiceImpl(
	reportRepository report.ReportRepository,
	userRepository user.UserRepository,
	threadRepository thread.ThreadRepository,
	idGenerator generator.IDGenerator,
) *reportServiceImpl {
	return &reportServiceImpl{
		reportRepository: reportRepository,
		userRepository:   userRepository,
		threadRepository: threadRepository,
		idGenerator:      idGenerator,
	}
}

func (r *reportServiceImpl) GetAll(
	ctx context.Context,
	accessorRole,
	status string,
	page uint,
	limit uint,
) (rs response.Pagination[response.Report], err error) {
	if accessorRole != "admin" {
		err = service.ErrAccessForbidden
		return
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	var reportStatus entity.ReportStatus
	if status == "accepted" {
		reportStatus = entity.Accepted
	} else {
		reportStatus = entity.Review
	}

	reports, repoErr := r.reportRepository.GetReportsWithPagination(
		ctx,
		entity.PageInfo{Page: page, Limit: limit},
		reportStatus,
	)

	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	rs.List = make([]response.Report, len(reports.List))

	for i, item := range reports.List {
		report := response.Report{
			ID:                 item.ID,
			ModeratorID:        item.Moderator.ID,
			ModeratorUsername:  item.Moderator.User.Username,
			ModeratorName:      item.Moderator.User.Name,
			UserID:             item.User.ID,
			Username:           item.User.Username,
			Name:               item.User.Name,
			Reason:             item.Reason,
			Status:             item.Status,
			ThreadID:           item.Thread.ID,
			ThreadTitle:        item.Thread.Title,
			ReportedOn:         item.CreatedAt.Format(time.RFC822),
			Comment:            item.Comment.Comment,
			CommentPublishedOn: item.Comment.CreatedAt.Format(time.RFC822),
		}
		rs.List[i] = report
	}

	rs.PageInfo.Page = reports.PageInfo.Page
	rs.PageInfo.Limit = reports.PageInfo.Limit
	rs.PageInfo.PageTotal = reports.PageInfo.PageTotal
	rs.PageInfo.Total = reports.PageInfo.Total

	return
}

func (r *reportServiceImpl) Create(
	ctx context.Context,
	accessorUserID string,
	p payload.CreateReport,
) (id string, err error) {
	if validateErr := validator.Validate(p); validateErr != nil {
		err = service.ErrInvalidPayload
		return
	}

	reportedUser, repoErr := r.userRepository.FindByUsername(ctx, p.Username)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	comment, repoErr := r.threadRepository.FindCommentByID(ctx, p.CommentID)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	moderators, repoErr := r.threadRepository.FindAllModeratorByThreadID(ctx, comment.Thread.ID)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	var moderatorID string
	var isModerator bool

	for _, moderator := range moderators {
		if accessorUserID == moderator.User.ID {
			isModerator = true
			moderatorID = moderator.ID
			break
		}
	}

	if !isModerator {
		err = service.ErrAccessForbidden
		return
	}

	id, genErr := r.idGenerator.GenerateReportID()
	if genErr != nil {
		err = service.MapError(genErr)
		return
	}

	if repoErr := r.reportRepository.Insert(ctx, id, moderatorID, reportedUser.ID, comment.ID, p.Reason); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	return
}

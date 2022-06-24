package user

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/thread"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/user"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	"gopkg.in/validator.v2"
)

type userServiceImpl struct {
	userRepository    user.UserRepository
	threadRepository  thread.ThreadRepository
	idGenerator       generator.IDGenerator
	passwordGenerator generator.PasswordGenerator
	tokenGenerator    generator.TokenGenerator
}

func NewUserServiceImpl(
	userRepository user.UserRepository,
	threadRepository thread.ThreadRepository,
	idGenerator generator.IDGenerator,
	passwordGenerator generator.PasswordGenerator,
	tokenGenerator generator.TokenGenerator,
) *userServiceImpl {
	return &userServiceImpl{
		userRepository:    userRepository,
		threadRepository:  threadRepository,
		idGenerator:       idGenerator,
		passwordGenerator: passwordGenerator,
		tokenGenerator:    tokenGenerator,
	}
}

func (u *userServiceImpl) Register(
	ctx context.Context,
	p payload.Register,
) (id string, err error) {
	if validateErr := validator.Validate(p); validateErr != nil {
		err = service.ErrInvalidPayload
		return
	}

	if _, parseErr := mail.ParseAddress(p.Email); parseErr != nil {
		err = service.ErrInvalidPayload
		return
	}

	id, genErr := u.idGenerator.GenerateUserID()
	if genErr != nil {
		err = service.MapError(genErr)
		return
	}

	password, genErr := u.passwordGenerator.GenerateFromPassword([]byte(p.Password), 10)
	if genErr != nil {
		err = service.MapError(genErr)
		return
	}

	user := entity.User{
		ID:       id,
		Username: p.Username,
		Email:    p.Email,
		Name:     p.Name,
		Password: string(password),
		Role:     "user",
	}

	if repoErr := u.userRepository.Insert(ctx, user); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	return
}

func (u *userServiceImpl) Login(
	ctx context.Context,
	p payload.Login,
) (r response.Login, err error) {
	if validateErr := validator.Validate(p); validateErr != nil {
		err = service.ErrInvalidPayload
		return
	}

	user, repoErr := u.userRepository.FindByUsername(ctx, p.Username)
	if repoErr != nil {
		if errors.Is(repoErr, repository.ErrRecordNotFound) {
			err = service.ErrUsernameNotFound
			return
		}
		err = service.MapError(repoErr)
		return
	}

	if !user.IsActive {
		err = service.ErrUsernameNotFound
		return
	}

	if compareErr := u.passwordGenerator.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(p.Password),
	); compareErr != nil {
		err = service.ErrCredentialNotMatch
		return
	}

	tokenPayload := generator.TokenPayload{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	token, genErr := u.tokenGenerator.GenerateToken(tokenPayload)
	if genErr != nil {
		err = service.MapError(repoErr)
		return
	}

	r.Token = token
	r.Role = user.Role

	return
}

func (u *userServiceImpl) GetAll(
	ctx context.Context,
	accessorUserID,
	orderBy,
	status string,
	page,
	limit uint,
) (r response.Pagination[response.User], err error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 20
	}

	var userOrderBy entity.UserOrderBy
	var userStatus entity.UserStatus

	if orderBy == "ranking" {
		userOrderBy = entity.Ranking
	} else {
		userOrderBy = entity.RegisteredDate
	}

	if status == "banned" {
		userStatus = entity.Banned
	} else {
		userStatus = entity.Active
	}

	pageInfo := entity.PageInfo{
		Page:  page,
		Limit: limit,
	}

	pagination, repoErr := u.userRepository.FindAllWithStatusAndPagination(
		ctx, accessorUserID, userOrderBy, userStatus, pageInfo,
	)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	r.List = make([]response.User, len(pagination.List))

	r.PageInfo.Page = pagination.PageInfo.Page
	r.PageInfo.Limit = pagination.PageInfo.Limit
	r.PageInfo.Total = pagination.PageInfo.Total
	r.PageInfo.PageTotal = pagination.PageInfo.PageTotal

	for i, user := range pagination.List {
		user := response.User{
			UserID:         user.ID,
			Username:       user.Username,
			Email:          user.Email,
			Name:           user.Name,
			Role:           user.Role,
			IsActive:       user.IsActive,
			RegisteredOn:   user.CreatedAt.Format(time.RFC822),
			TotalThread:    uint(user.TotalThread),
			TotalFollower:  uint(user.TotalFollower),
			TotalFollowing: uint(user.TotalFollowing),
			IsFollowed:     user.IsFollowed,
		}
		r.List[i] = user
	}

	return
}

func (u *userServiceImpl) GetOwn(
	ctx context.Context,
	accessorUserID,
	accessorUsername string,
) (r response.User, err error) {
	if user, repoErr := u.userRepository.FindByUsernameWithAccessor(
		ctx, accessorUserID, accessorUsername,
	); repoErr != nil {
		err = service.MapError(repoErr)
	} else {
		r = response.User{
			UserID:         user.ID,
			Username:       user.Username,
			Email:          user.Email,
			Name:           user.Name,
			Role:           user.Role,
			IsActive:       user.IsActive,
			RegisteredOn:   user.CreatedAt.Format(time.RFC822),
			TotalThread:    uint(user.TotalThread),
			TotalFollower:  uint(user.TotalFollower),
			TotalFollowing: uint(user.TotalFollowing),
			IsFollowed:     user.IsFollowed,
		}
	}

	return
}

func (u *userServiceImpl) GetByUsername(
	ctx context.Context,
	accessorUserID,
	username string,
) (r response.User, err error) {
	if user, repoErr := u.userRepository.FindByUsernameWithAccessor(
		ctx, accessorUserID, username,
	); repoErr != nil {
		err = service.MapError(repoErr)
	} else {
		r = response.User{
			UserID:         user.ID,
			Username:       user.Username,
			Email:          user.Email,
			Name:           user.Name,
			Role:           user.Role,
			IsActive:       user.IsActive,
			RegisteredOn:   user.CreatedAt.Format(time.RFC822),
			TotalThread:    uint(user.TotalThread),
			TotalFollower:  uint(user.TotalFollower),
			TotalFollowing: uint(user.TotalFollowing),
			IsFollowed:     user.IsFollowed,
		}
	}

	return
}

func (u *userServiceImpl) ChangeBannedState(
	ctx context.Context,
	accessorRole string,
	username string,
) (err error) {
	if accessorRole != "admin" {
		err = service.ErrAccessForbidden
		return
	}

	user, repoErr := u.userRepository.FindByUsername(ctx, username)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	if user.IsActive {
		if repoErr := u.userRepository.BannedUser(
			context.Background(), user.ID,
		); repoErr != nil {
			err = service.MapError(repoErr)
			return
		}
	} else {
		if repoErr := u.userRepository.UnbannedUser(
			context.Background(), user.ID,
		); repoErr != nil {
			err = service.MapError(repoErr)
			return
		}
	}

	return
}

func (u *userServiceImpl) ChangeFollowingState(
	ctx context.Context,
	accessorUserID,
	usernameToFollow string,
) (err error) {
	user, repoErr := u.userRepository.FindByUsernameWithAccessor(ctx, accessorUserID, usernameToFollow)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	if user.IsFollowed {
		if repoErr := u.userRepository.UnfollowUser(ctx, accessorUserID, user.ID); repoErr != nil {
			err = service.MapError(repoErr)
			return
		}
	} else {
		id, genErr := u.idGenerator.GenerateUserFollowID()
		if genErr != nil {
			err = service.MapError(genErr)
			return
		}

		if repoErr := u.userRepository.FollowUser(ctx, id, accessorUserID, user.ID); repoErr != nil {
			err = service.MapError(repoErr)
			return
		}
	}

	return
}

func (u *userServiceImpl) GetAllThreadByUsername(
	ctx context.Context,
	accessorUserID,
	username string,
	page uint,
	limit uint,
) (rs response.Pagination[response.ManyThread], err error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	user, repoErr := u.userRepository.FindByUsername(ctx, username)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	pagination, repoErr := u.threadRepository.FindAllByUserIDWithPagination(
		ctx,
		accessorUserID,
		user.ID,
		entity.PageInfo{
			Limit: limit,
			Page:  page,
		},
	)

	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	rs.PageInfo.Page = pagination.PageInfo.Page
	rs.PageInfo.Limit = pagination.PageInfo.Limit
	rs.PageInfo.PageTotal = pagination.PageInfo.PageTotal
	rs.PageInfo.Total = pagination.PageInfo.Total

	rs.List = make([]response.ManyThread, len(pagination.List))

	for i, item := range pagination.List {
		thread := response.ManyThread{
			ID:              item.ID,
			Title:           item.Title,
			CategoryID:      item.Category.ID,
			CategoryName:    item.Category.Name,
			PublishedOn:     item.CreatedAt.Format(time.RFC822),
			IsLiked:         item.IsLiked,
			IsFollowed:      item.IsFollowed,
			Description:     item.Description,
			TotalViewer:     item.TotalViewer,
			TotalLike:       item.TotalLike,
			TotalFollower:   item.TotalFollower,
			TotalComment:    item.TotalComment,
			CreatorID:       item.Creator.ID,
			CreatorUsername: item.Creator.Username,
			CreatorName:     item.Creator.Name,
		}
		rs.List[i] = thread
	}

	return
}

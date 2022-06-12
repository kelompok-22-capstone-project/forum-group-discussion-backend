package thread

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
)

type ThreadRepository interface {
	Insert(ctx context.Context, thread entity.Thread) (err error)
}

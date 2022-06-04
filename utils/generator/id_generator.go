package generator

import (
	"fmt"

	"github.com/aidarkhanov/nanoid/v2"
)

type IDGenerator interface {
	GenerateUserID() (id string, err error)
	GenerateCategoryID() (id string, err error)
	GenerateThreadID() (id string, err error)
	GenerateModeratorID() (id string, err error)
	GenerateReportID() (id string, err error)
	GenerateThreadFollowID() (id string, err error)
	GenerateLikeID() (id string, err error)
	GenerateCommentID() (id string, err error)
	GenerateUserFollowID() (id string, err error)
}

type nanoidIDGenerator struct{}

func NewNanoidIDGenerator() *nanoidIDGenerator {
	return &nanoidIDGenerator{}
}

func (n *nanoidIDGenerator) GenerateUserID() (id string, err error) {
	id, err = n.generate(6)
	id = fmt.Sprintf("u-%s", id)
	return
}

func (n *nanoidIDGenerator) GenerateCategoryID() (id string, err error) {
	id, err = n.generate(3)
	id = fmt.Sprintf("c-%s", id)
	return
}

func (n *nanoidIDGenerator) GenerateThreadID() (id string, err error) {
	id, err = n.generate(7)
	id = fmt.Sprintf("t-%s", id)
	return
}

func (n *nanoidIDGenerator) GenerateModeratorID() (id string, err error) {
	id, err = n.generate(4)
	id = fmt.Sprintf("m-%s", id)
	return
}

func (n *nanoidIDGenerator) GenerateReportID() (id string, err error) {
	id, err = n.generate(7)
	id = fmt.Sprintf("r-%s", id)
	return
}

func (n *nanoidIDGenerator) GenerateThreadFollowID() (id string, err error) {
	id, err = n.generate(7)
	id = fmt.Sprintf("f-%s", id)
	return
}

func (n *nanoidIDGenerator) GenerateLikeID() (id string, err error) {
	id, err = n.generate(7)
	id = fmt.Sprintf("l-%s", id)
	return
}

func (n *nanoidIDGenerator) GenerateCommentID() (id string, err error) {
	id, err = n.generate(7)
	id = fmt.Sprintf("c-%s", id)
	return
}

func (n *nanoidIDGenerator) GenerateUserFollowID() (id string, err error) {
	id, err = n.generate(7)
	id = fmt.Sprintf("o-%s", id)
	return
}

func (n *nanoidIDGenerator) generate(size int) (id string, err error) {
	id, err = nanoid.GenerateString(nanoid.DefaultAlphabet, size)
	return
}

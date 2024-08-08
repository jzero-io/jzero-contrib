package status

import (
	"github.com/pkg/errors"
	"log"
	"testing"
)

const (
	GetUserListError = Code(28001)
)

func TestError(t *testing.T) {
	err := Error(GetUserListError)
	log.Println(err)
	status := FromError(err)
	log.Println(status.code, status.message)
}

func TestUnknownError(t *testing.T) {
	err := Error(28000)
	log.Println(err)

	status := FromError(err)
	log.Println(status.code, status.message)
}

func TestWrap(t *testing.T) {
	err := Wrap(GetUserListError, errors.New("connect to db error"))
	status := FromError(err)
	log.Println(status.Error())
}

func init() {
	Register(GetUserListError, "获取用户列表失败")
}

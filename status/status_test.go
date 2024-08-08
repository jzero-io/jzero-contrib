package status

import (
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

	err = Error(28000)
	log.Println(err)

	status = FromError(err)
	log.Println(status.code, status.message)

}

func init() {
	Register(GetUserListError, "获取用户列表失败")
}

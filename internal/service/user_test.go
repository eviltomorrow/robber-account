package service

import (
	"database/sql"
	"log"
	"testing"

	"github.com/eviltomorrow/robber-account/internal/model"
	"github.com/eviltomorrow/robber-core/pkg/mysql"
	"github.com/stretchr/testify/assert"
)

func init() {
	mysql.DSN = "root:root@tcp(127.0.0.1:3306)/account?charset=utf8mb4&parseTime=true&loc=Local"
	mysql.Build()

}

func truncateUser() {
	var _sql = `truncate table user`
	_, err := mysql.DB.Exec(_sql)
	if err != nil {
		log.Fatal(err)
	}
}

var u1 = &model.User{
	NickName: sql.NullString{String: "shepard"},
	Email:    "eviltomorrow@163.com",
	Phone:    "12345678902",
}

func TestRegisterUser(t *testing.T) {
	_assert := assert.New(t)

	truncateUser()

	uuid, err := RegisterUser(u1)
	_assert.Nil(err)
	t.Logf("uuid: %v", uuid)

	_, err = RegisterUser(u1)
	_assert.NotNil(err)

}

func TestRemoveUser(t *testing.T) {
	_assert := assert.New(t)

	truncateUser()

	uuid, err := RegisterUser(u1)
	_assert.Nil(err)
	t.Logf("uuid: %v", uuid)

	err = RemoveUser(uuid)
	_assert.Nil(err)
}

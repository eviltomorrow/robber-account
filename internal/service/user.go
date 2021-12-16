package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/robber-account/internal/model"
	"github.com/eviltomorrow/robber-core/pkg/mysql"
	"github.com/google/uuid"
)

var (
	timeout = 10 * time.Second
)

func CreateUser(user *model.User) (string, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	_, err = model.UserWithSelectOneByEmail(mysql.DB, user.Email, timeout)
	if err == nil {
		return "", fmt.Errorf("email[%s] has been created", user.Email)
	}
	if err != sql.ErrNoRows {
		return "", err
	}

	_, err = model.UserWithSelectOneByPhone(mysql.DB, user.Phone, timeout)
	if err == nil {
		return "", fmt.Errorf("phone[%s] has been created", user.Phone)
	}
	if err != sql.ErrNoRows {
		return "", err
	}

	tx, err := mysql.DB.Begin()
	if err != nil {
		return "", err
	}

	user.UUID = uid.String()
	if _, err := model.UserWithInsertOne(tx, user, timeout); err != nil {
		tx.Rollback()
		return "", err
	}

	return uid.String(), tx.Commit()
}

func RemoveUser(uuid string) error {
	_, err := model.UserWithSelectOneByUUID(mysql.DB, uuid, timeout)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no user with uuid[%s]", uuid)
	}
	if err != nil {
		return err
	}

	tx, err := mysql.DB.Begin()
	if err != nil {
		return err
	}
	if _, err := model.UserWithDeleteByUUID(tx, uuid, timeout); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

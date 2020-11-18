package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zliang90/kingRest/internal/app/db"
	"github.com/zliang90/kingRest/pkg/log"
)

type User struct {
	Base
}

func NewUser(ctx *gin.Context) *User {
	s := new(User)
	s.Base.Init(ctx)
	return s
}

func (s User) DescribeUsers() ([]db.User, error) {
	var users []db.User

	log.Infof("%s, describe users", s.LogRequestIdPrefix())
	if err := s.db.Find(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []db.User{}, nil
		}
		return nil, err
	}
	return users, nil
}

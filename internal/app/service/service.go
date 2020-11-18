package service

import (
	"fmt"
	"github.com/zliang90/kingRest/internal/app/db"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Base struct {
	Ctx *gin.Context

	// default db
	db *gorm.DB
}

func (b *Base) Init(ctx *gin.Context) {
	b.Ctx = ctx

	if b.db == nil {
		b.db = db.GetDefaultDB()
	}
}

func (b *Base) GetRequestId() string {
	if b.Ctx != nil {
		return b.Ctx.GetString("Request-Id")
	}
	return ""
}

func (b *Base) LogRequestIdPrefix() string {
	reqId := b.GetRequestId()
	if reqId == "" {
		return fmt.Sprint("service req")
	}
	return fmt.Sprintf("reqId: %s", b.GetRequestId())
}

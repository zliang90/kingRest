package db

import (
	"github.com/zliang90/kingRest/pkg/log"
)

func dataInitialize() {
	// user
	u1 := User{
		Name:     "admin",
		Password: "admin123",
	}

	if err := _db.FirstOrCreate(&u1).Error; err != nil {
		log.Error(err)
	}
}

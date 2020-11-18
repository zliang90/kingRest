package db

// User user entity
type User struct {
	BaseModel

	Name     string `gorm:"type:varchar(128);"`
	Password string `gorm:"type:varchar(128);"`
}

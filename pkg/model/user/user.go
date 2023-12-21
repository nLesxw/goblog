package user

import "GoBlog/pkg/model"

//User 用户模型
type User struct{
	model.BaseModel

	Name string `gorm:"column:name;type:varchar(255);not null;unique"`
	Email string `gorm:"column:email;eype:varchar(255);default:NULL;unique;"`
	Password string `gorm:"column:password;type:varchar(255)"`
}
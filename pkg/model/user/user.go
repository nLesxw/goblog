package user

import "GoBlog/pkg/model"

//User 用户模型
type User struct{
	model.BaseModel

	Name string `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Email string `gorm:"type:varchar(255);unique;" valid:"email"`
	Password string `gorm:"type:varchar(255)" valid:"password"`
	// gorm: "-" 设置 GORM 在读写时略过此字段
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}
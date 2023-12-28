package category

import "GoBlog/pkg/model"

// Category 文章分类
type Category struct{
	model.BaseModel

	Name string `gorm:"type:varchar(255);not null" valid:"name"`
}


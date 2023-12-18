package article

import (
	"GoBlog/pkg/model"
	"GoBlog/pkg/types"
)

//Get 通过 ID 获取文章
func Get(idstr string) (Article, error){
	var article Article
	id := types.StringToUint64(idstr)
	//First是 gorm.DB 提供的用以从结果集中获取第一条数据的查询方法，需要注意的是第二个参数可以传参整型或者字符串 ID，但是传字符串会有 SQL 注入的风险
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}
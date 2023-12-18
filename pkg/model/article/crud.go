package article

import (
	"GoBlog/pkg/logger"
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

//GetAll 获取全部文章
func GetAll()([]Article, error) {
	var articles []Article
	if err := model.DB.Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

//Create 创建文章，通过 article.ID 来判断是否成功
func (article *Article) Create() (err error) {
	if err = model.DB.Create(&article).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}
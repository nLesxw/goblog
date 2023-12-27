package policies

import (
	"GoBlog/pkg/auth"
	"GoBlog/pkg/model/article"
)

//CanModifyArticle 是否允许修改话题
func CanModifyArticle(_article article.Article) bool {
	return auth.User().ID == _article.UserID
}
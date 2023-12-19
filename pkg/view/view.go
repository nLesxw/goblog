package view

import (
	"GoBlog/pkg/logger"
	"GoBlog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

//Render 渲染视图
func Render(w io.Writer, name string, data interface{}){
	//1.设置模板相对路径
	viewDir := "views/"

	//2. 语法糖，将 articles.show 更正为 articles/show
	//n 是允许替换的次数，设置为 -1 意味着替换所有
	name = strings.Replace(name, ".", "/", -1)

	//3. 所有布局模板文件 slice
	files, err := filepath.Glob(viewDir+"layouts/*.gohtml")
	logger.LogError(err)

	//4. 在slice里新增我们的目标文件
	newFiles := append(files, viewDir+name+".gohtml")

	//5. 解析所有的模板文件
	tmpl, err := template.New(name + ".gohtml").Funcs(
		template.FuncMap{
			"RouteName2URL": route.Name2URL,
	}).ParseFiles(newFiles...)
	logger.LogError(err)

	//6. 渲染模板
	err = tmpl.ExecuteTemplate(w, "app", data)
	logger.LogError(err)
}
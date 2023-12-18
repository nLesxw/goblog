package main

import (
	"GoBlog/app/http/controllers"
	"GoBlog/bootstrap"
	"GoBlog/pkg/database"
	"GoBlog/pkg/logger"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

var router *mux.Router
var db *sql.DB

// Article 对应一条文章数据
type Article struct {
    Title,Body string
    ID int64
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //1.设置标头
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        //2.继续处理
        next.ServeHTTP(w, r)
    })
}

//ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
    Title,Body string
    URL *url.URL
    Errors map[string]string
}

func removeTrailingSlash(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //1.除首页外，移除所有路径后带"/"
        if r.URL.Path != "/" {
            r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
            
        }
        //2.将请求传递下去
        next.ServeHTTP(w, r)
    })
}



func getArticleByID(id string) (Article, error) {
    article := Article{}
    query := "select * from articles where id = ?"
    err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
    return article, err
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {

    //1.获取参数
    id := getRouteVariable("id", r)

    //2.获取对应的文章数据
    article, err := getArticleByID(id)
    
    //3.如果出现了错误
    if err != nil {
        if err == sql.ErrNoRows {
            //3.1 未找到数据
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            //3.2 数据库错误
           logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    }else {
        //4. 读取成功,显示表单
        updateURL, _ := router.Get("articles.update").URL("id", id)
        data := ArticlesFormData{
            Title: article.Title,
            Body: article.Body,
            URL: updateURL,
            Errors: nil,
        }
        tmpl, err := template.ParseFiles("views/articles/edit.gohtml")
       logger.LogError(err)

        err = tmpl.Execute(w, data)
       logger.LogError(err)
    }
}

func articlesUpdateHandler(w http.ResponseWriter, r * http.Request) {
    //1. 获取 URL 参数
    id := getRouteVariable("id", r)

    //2. 读取对应的文章
    _, err := getArticleByID(id)

    //3.检查错误
    if err != nil {
        if err == sql.ErrNoRows {
            //3.1 未找到数据
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            //3.2 数据库错误
           logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    }else {
        //4. 未出现错误

        //4.1 表单验证
        title := r.PostFormValue("title")
        body := r.PostFormValue("body")

        errors := controllers.ValidateArticleFromData(title, body)

        if len(errors) == 0 {

            //4.2 表单验证通过，更新数据
            query := "update articles set title = ?, body = ? where id = ?"
            rs, err := db.Exec(query, title, body, id)

            if err != nil {
               logger.LogError(err)
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprint(w, "500 服务器内部错误")
            }

            //更新文章成功，跳转到文章详情页
            if n, _ := rs.RowsAffected(); n > 0 {
                showURL, _ := router.Get("articles.show").URL("id", id)
                http.Redirect(w, r, showURL.String(), http.StatusFound)
            }else {
                fmt.Fprint(w, "你没有做任何的修改!")
            }
        }else {
            // 4.3 表单验证不通过，显示理由

            updateURL, _ := router.Get("articles.update").URL("id", id)
            data := ArticlesFormData{
                Title: title,
                Body: body,
                URL: updateURL,
                Errors: errors,
            }
            tmpl, err := template.ParseFiles("views/articles/edit.gohtml")
           logger.LogError(err)

            err = tmpl.Execute(w, data)
           logger.LogError(err)
        }
    }
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {

    //1. 获取 URL 参数
    id := getRouteVariable("id", r)

    //2. 读取对应文章数据
    article, err := getArticleByID(id)

    //3. 如果出现错误
    if err != nil {
        if err == sql.ErrNoRows {
            //3.1 数据未找到
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        }else {
             // 3.2 数据库错误
            logger.LogError(err)
             w.WriteHeader(http.StatusInternalServerError)
             fmt.Fprint(w, "500 服务器内部错误")
        }
    }else {
        // 4. 未出现错误，执行删除操作
        rowsAffected, err := article.Delete()

        // 4.1 发生错误
        if err != nil {
            // 应该是 SQL 报错了
           logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        } else {
            // 4.2 未发生错误
            if rowsAffected > 0 {
                // 重定向到文章列表页
                indexURL, _ := router.Get("articles.index").URL()
                http.Redirect(w, r, indexURL.String(), http.StatusFound)
            }else {
                // Edge case
                w.WriteHeader(http.StatusNotFound)
                fmt.Fprint(w, "404 文章未找到")
            }
        }
    }
}

func (a Article) Delete() (rowAffected int64, err error){
    rs, err := db.Exec("delete from articles where id = "+strconv.FormatInt(a.ID, 10))

    if err != nil {
        return 0, nil
    }

    // 删除成功
    if n, _ := rs.RowsAffected(); n > 0 {
        return 0, nil
    }

    return 0, nil
}

func getRouteVariable(parameterName string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[parameterName]
}

func main() {
    database.Initialize()
    db = database.DB

    bootstrap.SetupDB()
    router = bootstrap.SetupRoute()

    router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
    router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
    router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

    //中间件：强制内容类型为 HTML
    router.Use(forceHTMLMiddleware)

    // 通过命名路由获取 URL 示例
    
    http.ListenAndServe(":3000", removeTrailingSlash(router))
}
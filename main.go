package main

import (
	"GoBlog/bootstrap"
	"GoBlog/pkg/database"
	"GoBlog/pkg/logger"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

    router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

    //中间件：强制内容类型为 HTML
    router.Use(forceHTMLMiddleware)

    // 通过命名路由获取 URL 示例
    
    http.ListenAndServe(":3000", removeTrailingSlash(router))
}
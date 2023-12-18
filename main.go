package main

import (
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
	"unicode/utf8"

	"github.com/gorilla/mux"
)

var router *mux.Router
var db *sql.DB

// Article 对应一条文章数据
type Article struct {
    Title,Body string
    ID int64
}

func (a Article) Link() string {
    showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID, 10))
    if err != nil {
       logger.LogError(err)
        return ""
    }
    return showURL.String()
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //1.设置标头
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        //2.继续处理
        next.ServeHTTP(w, r)
    })
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
    //1. 执行查询语句，返回结果集
    rows, err := db.Query("select * from articles")
   logger.LogError(err)
    defer rows.Close()

    var articles []Article
    //2. 循环读取结果
    for rows.Next() {
        var article Article
        //2.1 扫描每一行的结果并赋值到一个 article 对象中
        err := rows.Scan(&article.ID, &article.Title, &article.Body)
       logger.LogError(err)
        //2.2 将article 追加到 articles 这个数组中
        articles = append(articles, article)
    }

    //2.3 检查遍历时是否发生错误
    err = rows.Err()
   logger.LogError(err)

    //3. 加载模板
    tmpl, err := template.ParseFiles("views/articles/index.gohtml")
   logger.LogError(err)

    //4. 渲染模板，将所有文章的数据传输进去
    err = tmpl.Execute(w, articles)
   logger.LogError(err)
}

//ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
    Title,Body string
    URL *url.URL
    Errors map[string]string
}

func validateArticleFromData(title string, body string) map[string]string {
    errors := make(map[string]string)

    //验证标题
    if title == "" {
        errors["title"] = "标题不能为空"
    }else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
        errors["title"] = "标题长度需 介于 3-40"
    }

    //验证内容
    if body == "" {
        errors["body"] = "内容不能为空"
    }else if utf8.RuneCountInString(body) < 10 {
        errors["body"] = "内容不能少于10个字节"
    }

    return errors
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
    
    title := r.PostFormValue("title")
    body := r.PostFormValue("body")

    errors := validateArticleFromData(title, body)
    
    //检查是否出错
    if len(errors) == 0 {
        lastInsertID, err := saveArticleToDB(title, body)
        if lastInsertID > 0 {
            fmt.Fprint(w, "插入成功，ID为"+strconv.FormatInt(lastInsertID, 10))
        }else {
           logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    }else {

        storeURL, _ := router.Get("articles.store").URL()
        data := ArticlesFormData{
            Title: title,
            Body: body,
            URL: storeURL,
            Errors: errors,
        }
        tmpl, err := template.ParseFiles("views/articles/create.gohtml")
        if err != nil {
            panic(err)
        }

        err = tmpl.Execute(w, data)
        if err != nil {
            panic(err)
        }
    }
}

func saveArticleToDB(title string, body string) (int64, error){

    //初始化变量
    var (
        id int64
        err error
        rs sql.Result
        stmt *sql.Stmt
    )

    //1. 获取一个 prepare 声明语句
    stmt, err = db.Prepare("insert into articles (title, body) values (?,?)")
    //例行检查错误
    if err != nil {
        return 0, err
    }

    //2.在此函数运行结束后关闭此语句，防止占用SQL连接
    defer stmt.Close()

    //3. 执行请求，传参进入绑定的内容
    rs, err = stmt.Exec(title, body)
    if err != nil {
        return 0, err
    }

    //4. 插入成功的话，会返回自增 ID
    if id, err = rs.LastInsertId(); id > 0 {
        return id, err
    }

    return 0, err
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
    
    storeURL, _ := router.Get("articles.store").URL()

    data := ArticlesFormData{
        Title: "",
        Body: "",
        URL: storeURL,
        Errors: nil,
    }

    tmpl, err := template.ParseFiles("views/articles/create.gohtml")
    if err != nil {
        panic(err)
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        panic(err)
    }
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

        errors := validateArticleFromData(title, body)

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

    router = bootstrap.SetupRoute()

    router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
    router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
    router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
    router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
    router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
    router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

    //中间件：强制内容类型为 HTML
    router.Use(forceHTMLMiddleware)

    // 通过命名路由获取 URL 示例
    
    http.ListenAndServe(":3000", removeTrailingSlash(router))
}
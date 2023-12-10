package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)


var router = mux.NewRouter()
var db *sql.DB

func initDB(){
    var err error
    config := mysql.Config{
        User: "root",
        Passwd: "123456",
        Addr: "127.0.0.1:3306",
        Net: "tcp",
        DBName: "goblogNew",
        AllowNativePasswords: true,
    }

    //准备数据库连接池
    db, err = sql.Open("mysql", config.FormatDSN())
    checkError(err)

    //设置最大的连接数
    db.SetMaxOpenConns(25)
    //设置最大的空闲连接数
    db.SetMaxIdleConns(25)
    //设置每个连接的过期时间
    db.SetConnMaxLifetime(5 * time.Minute)

    //尝试连接，失败则会报错
    err = db.Ping()
    checkError(err)
}

func checkError(err error){
    if err != nil {
        log.Fatal(err)
    }
}

func createTables(){
    createArticlesTale := `create table if not exists articles(
        id bigint(20) primary key auto_increment not null,
        title varchar(255) collate utf8mb4_unicode_ci not null,
        body longtext collate utf8mb4_unicode_ci
    );`

    _, err := db.Exec(createArticlesTale)
    checkError(err)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    
    fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    
    fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
        "<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
    
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    fmt.Fprint(w, "文章 ID："+id)
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
    fmt.Fprint(w, "访问文章列表")
}

//ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
    Title,Body string
    URL *url.URL
    Errors map[string]string
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
    
    title := r.PostFormValue("title")
    body := r.PostFormValue("body")

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

    //检查是否出错
    if len(errors) == 0 {
        fmt.Fprint(w, "验证通过!<br>")
        fmt.Fprintf(w, "title 的值为: %v <br>", title)
        fmt.Fprintf(w, "title 的长度为: %v <br>", utf8.RuneCountInString(title))
        fmt.Fprintf(w, "body 的值为: %v <br>", body)
        fmt.Fprintf(w, "body 的长度为: %v <br>", utf8.RuneCountInString(body))
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

func main() {
    initDB()
    createTables()

    router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
    router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

    router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
    router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
    router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
    router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

    // 自定义 404 页面
    router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

    //中间件：强制内容类型为 HTML
    router.Use(forceHTMLMiddleware)

    // 通过命名路由获取 URL 示例
    
    http.ListenAndServe(":3000", removeTrailingSlash(router))
}
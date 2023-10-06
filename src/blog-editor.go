package editor

import (
	"database/sql"
	"net/http"

	"github.com/mattn/go-sqlite3"
)

var (
    db *sql.DB 
    sqlite3Conn sqlite3.SQLiteConn
    
    writeBlogStatement *sql.Stmt
)

type FullBlog struct {

    BlogTitle string
    BlogDate string
    BlogAuthor string
    BlogDescription string
    BlogPathname string
    BlogTopics string
    BlogNotes string
    BlogContent string

}

func main() {


    var err error

    db, err = sql.Open("sqlite3", "file:./src/blog-database?cache=shared")
    defer db.Close()

    if err != nil {
        // TODO:
    }

    if err = db.Ping(); err != nil {
        // TODO:
    }

    writeBlogStatement, err = db.Prepare(`
        UPDATE blogs
        SET blogtitle = ?, blogdate = ?, blogauthor = ?, blogdescription = ?,
        blogpathname = ?, blogtopics = ?, blognotes = ?, blogcontent = ?
        `)

    defer writeBlogStatement.Close()

    if err != nil {
        // TODO:
    }


    http.HandleFunc("/", RouteHandler)
    err = http.ListenAndServe(":8000", nil)

}


func RouteHandler(writer http.ResponseWriter, request *http.Request) {

    requestPath := request.URL.Path

    if requestPath == "/" {
        http.ServeFile(writer, request, "src/static/templates/blogedit.html")
    }
}

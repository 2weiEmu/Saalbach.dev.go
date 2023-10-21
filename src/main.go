package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/mattn/go-sqlite3"
)

var (
    db *sql.DB 
    sqlite3Conn sqlite3.SQLiteConn
    fetchBlogStatement *sql.Stmt 
    searchBlogsStatement *sql.Stmt
)

type RawBlog struct {
    BlogTitle string;
    BlogDate string;
    BlogAuthor string;
    BlogTopics string;
    BlogNotes string;
    BlogContent string;
}

type BlogHeader struct {
    BlogTitle string
    BlogDate string
    BlogAuthor string
    BlogDescription string
    BlogPathname string
    BlogTopics string
    BlogNotes string
}

// TODO: make the log.fatals maybe just logs that go into an actual log -> we don't want the site to crash everytime someone messes around

// NOTE: obviously, optimally you would want some kind of blog manifesto where you keep the key information about the blog without neccessarily 
// loading the whole file - but that A. prob won't happen anyway, because obv you only read what you read, but I would still have to open every 
// file like this, this would be the obvious optimisation, but I am not gonna do that for now - I will keep it simple (unless if loading times 
// become too big)

func GetAllBlogHeaders(statement *sql.Stmt, titleFilter string) []BlogHeader {

    rows, err := statement.Query("%" + titleFilter + "%")
    defer rows.Close()

    if err != nil {
        fmt.Println("Prepared statement failed to execute with error:", err)
    }

    var blogHeaders []BlogHeader;

    for rows.Next() {

        var header BlogHeader;
        err = rows.Scan(
            &header.BlogTitle, &header.BlogDate, &header.BlogAuthor, 
            &header.BlogDescription, &header.BlogPathname, &header.BlogTopics, 
            &header.BlogNotes,
            )

        if err != nil {
            fmt.Println("Failed to retrieve row:", rows, "With the error:", err)
            return nil
        }

        header.BlogDate = strings.TrimSuffix(header.BlogDate, "T00:00:00Z")
        blogHeaders = append(blogHeaders, header)

    }
    

    return blogHeaders

}

func GetBlogUsingPathname(statement *sql.Stmt, pathname string) RawBlog {

    rows, err := statement.Query(pathname)

    defer rows.Close()

    if err != nil {
        // TODO:
    }

    var rawBlog RawBlog

    rows.Next();

    err = rows.Scan(
        &rawBlog.BlogTitle, &rawBlog.BlogDate, &rawBlog.BlogAuthor,
        &rawBlog.BlogTopics, &rawBlog.BlogNotes, &rawBlog.BlogContent,
        )

    if err != nil {
        // TODO:
    }
    rawBlog.BlogDate = strings.TrimSuffix(rawBlog.BlogDate, "T00:00:00Z")

    return rawBlog


}

func RouteHandler(writer http.ResponseWriter, request *http.Request) {

    requestPath := request.URL.Path

    fmt.Println("Request to: ", requestPath)
    
    // Serving the main page
    if requestPath == "/" {
        http.ServeFile(writer, request, "src/static/templates/index.html")


    // Serving static files
    } else if match, _ := regexp.MatchString("^/((css)|(js))/", requestPath); match {
        fmt.Println("Serving Static File...")
        fs := http.FileServer(http.Dir("src/static"))
        //http.StripPrefix("static/", fs)
        fs.ServeHTTP(writer, request)


    // Serving any blog
    } else if match, _ := regexp.MatchString("^/blogs/", requestPath); match {
        fmt.Println("Serving blog...")
        blogToLoad, err := os.Open("./" + requestPath)
        defer blogToLoad.Close()

        if err != nil { // if invalid blog loaded, I will for now just redirect
            http.Redirect(writer, request, "/", http.StatusNotFound)
        }

        // getting the blog article, so that it may be inserted into the template
        blogArticle := GetBlogUsingPathname(fetchBlogStatement, strings.Split(requestPath, "/")[2])

        blogTemplate, err := template.ParseFiles("src/static/templates/blog-article.html")
        if err != nil {
           log.Fatal(err)
        }

        blogTemplate.Execute(writer, blogArticle)


    // Serving any images
    } else if match, _ := regexp.MatchString("^/images/", requestPath); match {
        // TODO: if statement can be improved
        fmt.Println("Serving Static File...")
        fs := http.FileServer(http.Dir("src/static"))
        http.StripPrefix("/static", fs)
        fs.ServeHTTP(writer, request)


    // Serving the about, contact or projects page
    } else if requestPath == "/about" || requestPath == "/contact" || requestPath == "/projects" {
        http.ServeFile(writer, request, "src/static/templates" + requestPath + ".html")


    // Serving the main blogs page
    } else if requestPath == "/blog" {
        // NOTE: strategy for loading search... simply don't include the things that don't match.... yes very fun.
        // TODO: make it so that on page reload it empties the search params (perhaps requires more javascript)
        
        fmt.Println("filtering blogs")

        searchParams := request.URL.Query()
        search := ""

        if searchParams.Has("search") {
            search = searchParams["search"][0]
        }

        var filteredBlogs []BlogHeader
        filteredBlogs = GetAllBlogHeaders(searchBlogsStatement, search)

        fmt.Println("filteredBlogs for blog page:", filteredBlogs)

        blogPage, err := template.ParseFiles("src/static/templates/blog.html")

        if err != nil {
            log.Fatal(err)
        }

        err = blogPage.Execute(writer, filteredBlogs)

        if err != nil {
            log.Fatal(err)
        }


    // Serving the 404 page
    } else {
        http.ServeFile(writer, request, "src/static/templates/404.html")


    }
}


func main() {

    fmt.Println("Received Arguments:", os.Args)

    /*
     * NOTE: Recommended format for easily launching
     * ./main deploy|test [port number] [certificate location] [private key location]
     * TODO: make this use the flags package that go has available
     */

    // default values for flags
    deployType := "test"
    port := ""
    certificate := ""
    privateKey := ""

    recommendedFormat := "./main deploy|test [port number] [certificate location] [private key location]"
    osArgsLen := len(os.Args) 

    // must give deploy or test - otherwise invalid
    if osArgsLen >= 2 {
        if os.Args[1] == "deploy" || os.Args[1] == "test" {
            deployType = os.Args[1]
        } else {
            fmt.Println("Invalid deploy type given:\n", recommendedFormat)
            os.Exit(1)
        }
    } else {
        fmt.Println("Too few arguments given:\n", recommendedFormat)
        os.Exit(1)
    }

    if osArgsLen < 3 {
        fmt.Println("No port number given, defaulting to port 8000.\n", recommendedFormat)
    } else {
        _, err := strconv.Atoi(os.Args[2])

        port = os.Args[2]

        if err != nil {
            fmt.Println("Given port number is not valid:\n", recommendedFormat)
            os.Exit(1)
        }
    }

    http.HandleFunc("/", RouteHandler)
    fmt.Println("Server is almost ready.")

    var err error
    db, err = sql.Open("sqlite3", "file:./src/blog-database?cache=shared")

    defer db.Close()

    if err != nil {
        // TODO:
        log.Fatal("failed to open db with error ", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal("failed to ping db with error ", err)
    }

    fetchBlogStatement, err = db.Prepare(`SELECT blogtitle, blogdate, blogauthor, blogtopics, blognotes, blogcontent FROM blogs WHERE blogpathname = ?`)
    searchBlogsStatement, err = db.Prepare(`SELECT blogtitle, blogdate, blogauthor, blogdescription, blogpathname, blogtopics, blognotes FROM blogs WHERE blogtitle LIKE ?`)

    defer fetchBlogStatement.Close()
    defer searchBlogsStatement.Close()

    if deployType == "test" {
        err = http.ListenAndServe(":" + port, nil)
    } else {
        
        if osArgsLen < 5 { 
            fmt.Println("Either a certificate or private key path was not given.\n", recommendedFormat)
            os.Exit(1)
        }

        certificate = os.Args[3]
        privateKey = os.Args[4]

        err = http.ListenAndServeTLS(":" + port, certificate, privateKey, nil)
    }

    if err != nil {
        fmt.Println(err)
    }
}


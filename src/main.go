package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"text/template"
)

type BlogParagraph struct {
    BlogParagraph string;
}

type BlogArticle struct {
    BlogTitle string;
    BlogDate string;
    BlogContent []BlogParagraph;
}

type Blog struct {
    BlogTitle string;
    BlogDate string;
    BlogDescription string;
    BlogPathName string;
}

type BlogOverview struct {
    AllBlogs []Blog;
}


// TODO: make the log.fatals maybe just logs that go into an actual log -> we don't want the site to crash everytime someone messes around

var Blogs BlogOverview // TODO: make this not be a global variable (I mean its my website so no one really cares but meh)

// NOTE: obviously, optimally you would want some kind of blog manifesto where you keep the key information about the blog without neccessarily 
// loading the whole file - but that A. prob won't happen anyway, because obv you only read what you read, but I would still have to open every 
// file like this, this would be the obvious optimisation, but I am not gonna do that for now - I will keep it simple (unless if loading times 
// become too big)
func InitialBlogRead() (BlogOverview) {

    allBlogs, err := os.ReadDir("./blogs/")

    if err != nil {
        log.Fatal(err)
    }

    var TotalBlogs []Blog;

    for _, e := range allBlogs { 

        blog, err := os.Open("blogs/" + e.Name())
        defer blog.Close()

        if err != nil {
            log.Fatal(err)
        }

        scanner := bufio.NewScanner(blog)

        scanner.Scan()
        title := scanner.Text()

        scanner.Scan()
        date := scanner.Text()

        scanner.Scan()
        description := scanner.Text()

        scanner.Scan()
        path := scanner.Text()

        newBlog := Blog {
            BlogTitle: title,
            BlogDate: date,
            BlogDescription: description,
            BlogPathName: path,
        }

        TotalBlogs = append(TotalBlogs, newBlog)

    }

    fmt.Println(TotalBlogs)

    return BlogOverview {
        AllBlogs: TotalBlogs,

    }
}

func RouteHandler(writer http.ResponseWriter, request *http.Request) {

    requestPath := request.URL.Path

    fmt.Println("Request to: ", requestPath)

    if requestPath == "/" {
        http.ServeFile(writer, request, "src/static/templates/index.html")


    } else if match, _ := regexp.MatchString("^/((css)|(js))/", requestPath); match {

        fmt.Println("Serving Static File...")
        fs := http.FileServer(http.Dir("src/static"))
        //http.StripPrefix("static/", fs)
        fs.ServeHTTP(writer, request)
    } else if match, _ := regexp.MatchString("^/blogs/", requestPath); match {

        fmt.Println("Serving blog...")
        blogToLoad, err := os.Open("./" + requestPath)
        defer blogToLoad.Close()

        if err != nil { // if invalid blog loaded, I will for now just redirect
            http.Redirect(writer, request, "/", http.StatusNotFound)
        }

        scanner := bufio.NewScanner(blogToLoad)

        scanner.Scan()
        title := scanner.Text()
        scanner.Scan()
        date := scanner.Text()
        scanner.Scan();
        scanner.Scan();

        var blogContent []BlogParagraph

        for scanner.Scan() { 
            blogContent = append(blogContent, BlogParagraph { BlogParagraph: scanner.Text() })
        }

        blogArticle := BlogArticle {
            BlogTitle: title,
            BlogDate: date,
            BlogContent: blogContent,

        }

        blogTemplate, err := template.ParseFiles("src/static/templates/blog-article.html")
        if err != nil {
            log.Fatal(err)
        }

        blogTemplate.Execute(writer, blogArticle)
        


    } else if match, _ := regexp.MatchString("^/images/", requestPath); match {
        // TODO: if statement can be improved
        fmt.Println("Serving Static File...")
        fs := http.FileServer(http.Dir("src/static"))
        http.StripPrefix("/static", fs)
        fs.ServeHTTP(writer, request)

    } else if requestPath == "/about" || requestPath == "/contact" || requestPath == "/projects" {
        http.ServeFile(writer, request, "src/static/templates" + requestPath + ".html")
    } else if requestPath == "/blog" {
        // NOTE: strategy for loading search... simply don't include the things that don't match.... yes very fun.
        // TODO: make it so that on page reload it empties the search params (perhaps requires more javascript)
        
        fmt.Println("filtering blogs")

        searchParams := request.URL.Query()
        search := ""

        if searchParams.Has("search") {
            search = searchParams["search"][0]
        }

        regex, err := regexp.Compile(`(?i)` + search)

        if err != nil {
            log.Println("Failed to search.")
            http.Redirect(writer, request, "/", http.StatusBadRequest)
        }

        var filteredBlogs BlogOverview

        for i := 0; i < len(Blogs.AllBlogs); i++ {
            title := Blogs.AllBlogs[i].BlogTitle
            if regex.MatchString(title) {
                filteredBlogs.AllBlogs = append(filteredBlogs.AllBlogs, Blogs.AllBlogs[i])
            }
        }


        blogPage, err := template.ParseFiles("src/static/templates/blog.html")

        if err != nil {
            log.Fatal(err)
        }

        err = blogPage.Execute(writer, filteredBlogs)

        if err != nil {
            log.Fatal(err)
        }

    } else {
        http.ServeFile(writer, request, "src/static/templates/404.html")
    }
}


func main() {

    fmt.Println("Received Arguments:", os.Args)

    /*
     * NOTE: Recommended format for easily launching
     * ./main deploy|test [port number] [certificate location] [private key location]
     */

    // default values for flags
    deployType := "test"
    port := ""
    certificate := ""
    privateKey := ""

    recommendedFormat := "./main deploy|test [port number] [certificate location] [private key location]"
    osArgsLen := len(os.Args) 


    Blogs = InitialBlogRead()

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

    var err error;

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


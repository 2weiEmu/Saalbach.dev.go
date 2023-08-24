package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadHeaderBlogV2(blogFile *os.File) Blog {
    scanner := bufio.NewScanner(blogFile)
    newBlog := Blog{ }

    // Reading the version 1 blog

    scanner.Scan()
    newBlog.BlogVersion, _ = strconv.Atoi(scanner.Text())

    scanner.Scan()
    newBlog.BlogTitle = scanner.Text()

    scanner.Scan()
    newBlog.BlogDate = scanner.Text()

    scanner.Scan()
    newBlog.BlogAuthor = scanner.Text()

    scanner.Scan()
    newBlog.BlogDescription = scanner.Text()

    scanner.Scan()
    newBlog.BlogPathName = scanner.Text()

    scanner.Scan()
    newBlog.BlogTopics = scanner.Text() // TODO: have to be upgraded, but you can't search for it yet so its okay

    scanner.Scan()
    newBlog.BlogNotes = scanner.Text()

    return newBlog;
}

func ReadBodyBlogV2(blogFile *os.File) BlogArticle {
    scanner := bufio.NewScanner(blogFile)
    newBlog := BlogArticle{ }


    scanner.Scan()

    scanner.Scan()
    newBlog.BlogTitle = scanner.Text()

    scanner.Scan()
    newBlog.BlogDate = scanner.Text()

    scanner.Scan()
    newBlog.BlogAuthor = scanner.Text()

    scanner.Scan()
    newBlog.BlogDescription = scanner.Text()

    scanner.Scan()

    scanner.Scan()
    newBlog.BlogTopics = scanner.Text() // TODO: have to be upgraded, but you can't search for it yet so its okay

    scanner.Scan()
    newBlog.BlogNotes = scanner.Text()

    var text []BlogParagraph;

    for scanner.Scan() {
        text = append(text, BlogParagraph { scanner.Text() } )
    }

    newBlog.BlogContent = text;

    return newBlog;
}

func UpgradeBlogV1ToV2(blogFile *os.File) { // TODO: alright this function can wiat, I have one blog I shall do it manually

    scanner := bufio.NewScanner(blogFile)
    newBlog := Blog{ }

    // Reading the version 1 blog

    scanner.Scan()
    newBlog.BlogTitle = scanner.Text()

    scanner.Scan()
    newBlog.BlogDate = scanner.Text()

    scanner.Scan()
    newBlog.BlogDescription = scanner.Text()

    scanner.Scan()
    newBlog.BlogPathName = scanner.Text()

    var text []string;

    for scanner.Scan() {
        text = append(text, scanner.Text())
    }

    // Writing the version 2 blog

    os.Truncate("blogs/" + newBlog.BlogPathName, 0)
    
    fmt.Fprintln(blogFile, "2\n" + newBlog.BlogTitle)

    // formating the date, i.e moving every character after 
    dateSplits := strings.Split(newBlog.BlogDate, " ")

    fmt.Fprintln(blogFile, dateSplits[0])

    fmt.Fprintln(blogFile, "robert arno saalbach")
    fmt.Fprintln(blogFile, newBlog.BlogDescription)
    fmt.Fprintln(blogFile, newBlog.BlogPathName)
    fmt.Fprintln(blogFile, "-")

    for i := 1; i < len(dateSplits); i++ {
        fmt.Fprintf(blogFile, dateSplits[i] + " ")
    }

    for _, line := range text { 
        fmt.Fprintln(blogFile, line)
    }
}

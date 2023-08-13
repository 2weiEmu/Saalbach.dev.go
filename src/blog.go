package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func ReadVersion1BlogHeaderAndConvert(blogFile *os.File, scanner *bufio.Scanner) Blog {

    scanner.Scan()
    title := scanner.Text()

    scanner.Scan()
    date := scanner.Text()

    scanner.Scan()
    description := scanner.Text()

    scanner.Scan()
    path := scanner.Text()

    var lines []string

    // I will be real, I was not the most awake when I wrote all of this code.
    for scanner.Scan() {
        line := scanner.Text()
        lines = append(lines, line)
    }

    newBlog := Blog {
        BlogTitle: title,
        BlogDate: date,
        BlogDescription: description,
        BlogPathName: path,
    }

    UpgradeVersion1To2(blogFile, newBlog, lines)

    return newBlog

}

func ReadVersion1Blog() {

}

func UpgradeVersion1To2(blogFile *os.File, blogData Blog, blogText []string) {
    
    fmt.Fprintln(blogFile, "")

    os.Truncate(blogFile.Name(), 0)

    fmt.Fprintln(blogFile, "2\n" + blogData.BlogTitle + "\n" + blogData.BlogDate)
    fmt.Fprintln(blogFile, "-\n" + blogData.BlogDescription + "\n" + blogData.BlogPathName)
    fmt.Fprintln(blogFile, "-\n-")

    for _, line := range blogText {
        fmt.Fprintln(blogFile, line) // prob not the best way, since you know, individual writes to buffer, but meh
    }
   
}

func ReadVersion2BlogHeader(blogfile *os.File) Blog {

    scanner := bufio.NewScanner(blogfile)

    blogFile := Blog{}

    scanner.Scan()
    blogFile.BlogVersion, _ = strconv.Atoi(scanner.Text())

    scanner.Scan()
    blogFile.BlogTitle = scanner.Text()

    scanner.Scan()
    blogFile.BlogDate = scanner.Text()

    scanner.Scan()
    blogFile.BlogAuthor = scanner.Text()

    scanner.Scan()
    blogFile.BlogDescription = scanner.Text()

    scanner.Scan()
    blogFile.BlogPathName = scanner.Text()

    scanner.Split())
    scanner.Scan()
    blogFile.BlogTopics = scanner.Text()

    scanner.Scan()
    blogFile.BlogNotes = scanner.Text()

}

func ReadVersion2Blog() {

}

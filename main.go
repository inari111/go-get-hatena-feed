package main

import (
	"encoding/xml"
	"fmt"
    "io/ioutil"
    "net/http"
    "time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type XML struct {
	Bookmarks []struct {
		Title         string `xml:"title"`
		Link          string `xml:"link"`
		Description   string `xml:"description"`
		Date          string `xml:"date"`
		BookmarkCount int    `xml:"bookmarkcount"`
	} `xml:"item"`
}

type hotentry struct {
	title string
	link string
	description string
	bookmarkcount int
	date string
}

func main() {
	data := httpGet("http://b.hatena.ne.jp/hotentry/it.rss")

	result := XML{}
	err := xml.Unmarshal([]byte(data), &result)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	for _, bookmark := range result.Bookmarks {
		datetime, _ := time.Parse(time.RFC3339, bookmark.Date)

		fmt.Printf("%v\n", datetime.Format("2006/01/02 15:04:05"))
		fmt.Printf("%s - %dbookmark\n", bookmark.Title, bookmark.BookmarkCount)
		fmt.Printf("%v\n", bookmark.Link)
		fmt.Println()

		h := hotentry{
			title: bookmark.Title,
			link: bookmark.Link,
			description: bookmark.Description,
			bookmarkcount: bookmark.BookmarkCount,
			date: datetime.Format("2006/01/02 15:04:05"),
		}
		fmt.Println(h)
		save(h)
	}

}

func httpGet(url string) string {
	response, _ := http.Get(url)
    body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return string(body)
}

func save(h hotentry) bool {
	db, err := sql.Open("mysql", "user@tcp(localhost:3306)/test?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT test SET title=?")
	if err != nil {
		panic(err.Error())
	}

	now := time.Now()
	const layout = "2006-01-02 15:04:05"
	now.Format(layout)

	res, err :=  stmt.Exec(h.title)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(res)
	return true
}

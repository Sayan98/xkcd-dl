package main

import (
	"os"
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"
	"sync"
	"log"
)

type xkcdPost struct {
	Num int
	Month, Link, Year, News, Safe_title, Transcript, Alt, Img, Title, Day string
}

func getPost(url string) (xkcdPost) {
	post := xkcdPost{}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	    panic(err.Error())
	}

	json.Unmarshal(body, &post)

	return post
}

func saveImage(url, path string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	file, err := os.Create(path)
	if err != nil {
        log.Fatal(err)
    }

    defer file.Close()

    _, err = io.Copy(file, resp.Body)
    if err != nil {
        log.Fatal(err)
    }
}


func main() {
	var wg sync.WaitGroup

	url := "https://xkcd.com/info.0.json"
	latestPost := getPost(url)
	maxId := latestPost.Num

	os.Mkdir("images", os.ModePerm)

	f, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
	    log.Fatalf("error opening file: %v", err)
	}

	defer f.Close()
	log.SetOutput(f)

	for Id := 1; Id <= maxId; Id++ {
		if Id != 404 {
			wg.Add(1)

			go func(Id int) {
				url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", Id)
				post := getPost(url)

				src := post.Img
				name := fmt.Sprintf("#%d%s", post.Num, filepath.Ext(src))
				path := fmt.Sprintf("images/%s", name)

				fmt.Println(url)
				saveImage(src, path)

				wg.Done()
			}(Id)
		}
	}

	wg.Wait()
}

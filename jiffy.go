package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/paddycarey/gophy"
)

func main() {
	log.SetFlags(log.Lshortfile)
	if len(os.Args) < 2 {
		log.Fatal("no arguments received")
	}

	// search Giphy
	q := strings.Join(os.Args[1:], " ")
	opts := &gophy.ClientOptions{}
	giphyClient := gophy.NewClient(opts)
	gifs, _, err := giphyClient.SearchGifs(q, "", 100, 0)
	if err != nil {
		log.Fatal(err)
	}

	if len(gifs) < 1 {
		log.Fatal("no results")
	}

	r := rand.New(rand.NewSource(99))
	g := gifs[r.Intn(len(gifs))]

	out, err := os.Create(g.Id + ".gif")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	go spinner(100 * time.Millisecond)

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Get(gifs[0].Images.Original.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("\r%s downloaded.\n", out.Name())
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cespedes/planet/feed"
)

func main() {
	//f, err := feed.Parse("https://go.dev/blog/feed.atom")
	if len(os.Args) != 2 {
		log.Fatal("Wrong number of arguments")
	}
	f, err := feed.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(f)
	fmt.Printf("# %s\n", f.Title)
	for _, i := range f.Items {
		fmt.Printf("- %s %s\n", i.PublishedParsed.Format("2006-01-02"), i.Title)
	}
}

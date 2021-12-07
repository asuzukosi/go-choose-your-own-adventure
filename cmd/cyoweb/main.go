package main

import (
	"cyoa"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	file := flag.String("file", "stories.json", "the JSON file with the CYOA story")
	port := flag.Int("port", 8080, "port where the CYOA story server runs")
	flag.Parse()

	fmt.Printf("Using the story in file %s\n", *file)

	f_obj, err := os.Open(*file)
	if err != nil {
		panic(err)
	}
	story, err := cyoa.JSONStory(f_obj)
	if err != nil {
		panic(err)
	}
	handler := cyoa.NewHandler(story)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}

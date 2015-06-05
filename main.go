package main

import (
	"fmt"
	"github.com/jBugman/go-pocket/pocket"
	"os"
	"pocket/analysis"
)

const CONSUMER_KEY = "<my-key>"
const ACCESS_TOKEN = "<access-token>"

func dump(model analysis.Model, filename string) {
	fmt.Println("Dumping model to " + filename)

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for tag, counter := range model.Source {
		f.WriteString(tag + "\n")
		for _, countItem := range counter.ItemsWithThreshold(1) {
			f.WriteString(fmt.Sprintf("\t%s: %d\n", countItem.Key, countItem.Count))
		}
		f.WriteString("\n")
	}
}

func main() {
	api := pocket.Api{CONSUMER_KEY, ACCESS_TOKEN}

	// items, err := api.Retrieve(pocket.Request{Count: 1})
	items, err := api.RetrieveAllArticles()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total items: %d\n", len(items))

	corpus := analysis.ConvertItems(items)
	model := analysis.TrainModel(corpus)
	dump(model, "model-src.txt")
	model.Dump("model.txt")
}

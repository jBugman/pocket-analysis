package main

import (
	"fmt"
	"github.com/jBugman/go-pocket/pocket"
	"os"
	"pocket/analysis"
)

const CONSUMER_KEY = "13888-e9be4bfc69cef5f8917d1ca6"
const ACCESS_TOKEN = "581257dc-0915-6b6c-bbc9-b12a22"

func dump(model analysis.Model, filename string) {
	fmt.Println("Dumping model to " + filename)

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for tag, counter := range model {
		f.WriteString(tag + "\n")
		for _, countItem := range counter.Items() {
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

	corpus := analysis.ConvertItems(items[0:50])
	model := analysis.TrainModel(corpus)
	dump(model, "model.txt")
}

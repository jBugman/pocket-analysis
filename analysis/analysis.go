package analysis

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/jBugman/go-pocket/pocket"
	"os"
	"strings"
)

type Item struct {
	Title   string
	Excerpt string
	Tags    []string
}

func ConvertItem(src pocket.Item) Item {
	return Item{src.Title, src.Excerpt, src.GetTags()}
}

func ConvertItems(source []pocket.Item) (result []Item) {
	for _, item := range source {
		result = append(result, ConvertItem(item))
	}
	return
}

func (this *Item) tokens() []string {
	tokenList := []string{}
	for _, x := range append(getTockens(this.Title), getTockens(this.Excerpt)...) {
		if !trash.Contains(x) {
			tokenList = append(tokenList, x)
		}
	}
	return tokenList
}

var trash = mapset.NewSetFromSlice([]interface{}{"-", "--", "·", "—", "–", "/"})
var stopwords = mapset.NewSetFromSlice([]interface{}{
	"the", "to", "a", "i", "is", "that", "of",
	"it", "and", "or", "at", "in", "for", "we", "as",
})

func (this *Item) Tokens() []string {
	tokenList := []string{}
	for _, x := range this.tokens() {
		if !stopwords.Contains(x) {
			tokenList = append(tokenList, x)
		}
	}
	return append(tokenList, this.Bigrams()...)
}

func (this *Item) Bigrams() []string {
	return bigrams(this.tokens())
}

func getTockens(source string) (result []string) {
	for _, value := range strings.Fields(source) {
		token := strings.ToLower(value)
		token = strings.Trim(token, "()[]{}“”«»\"'`")
		token = strings.TrimRight(token, ",.!?;:")
		token = strings.TrimSuffix(token, "’s")
		result = append(result, token)
	}
	return
}

func bigrams(tokens []string) (result []string) {
	const N = 2
	for i := 0; i < len(tokens)-N; i++ {
		result = append(result, strings.Join([]string{tokens[i], tokens[i+1]}, " "))
	}
	return
}

type ModelSource map[string]Counter
type Features map[string]mapset.Set
type Model struct {
	Source   ModelSource
	Features Features
}

func (this ModelSource) Add(tag string, key string) {
	counter, exists := this[tag]
	if exists {
		counter.Add(key)
	} else {
		this[tag] = Counter{key: 1}
	}
}

func TrainModel(items []Item) (model Model) {
	model = Model{Source: ModelSource{}}
	for _, item := range items {
		for _, tag := range item.Tags {
			for _, bigram := range item.Tokens() {
				model.Source.Add(tag, bigram)
			}
		}
	}

	features := model.crossfit()
	model.Features = features
	// model.crossEliminate(features)
	return
}

func (this *Model) crossfit() Features {
	rawFeatures := Features{}
	for tag, counter := range this.Source {
		items := mapset.NewSet()
		// for _, countItem := range counter.ItemsWithThreshold(1) {
			// items.Add(countItem.Key)
		for k, _ := range counter.Weights(-0.025) {
			items.Add(k)
		}
		rawFeatures[tag] = items
	}
	return rawFeatures
}

func (this *Model) crossEliminate(rawFeatures Features) {
	this.Features = Features{}
	for tag, features := range rawFeatures {
		for otherTag, otherFeatures := range rawFeatures {
			if otherTag != tag {
				features = features.Difference(otherFeatures)
			}
		}
		this.Features[tag] = features
	}
}

func (this *Model) Predict(item Item) Weights {
	result := Counter{}
	tokens := item.Tokens()
	for _, token := range tokens {
		for tag, features := range this.Features {
			for feature := range features.Iter() {
				if token == feature {
					result.Add(tag)
				}
			}
		}
	}
	return result.Weights(0.05)
}

func (this *Model) Dump(filename string) {
	fmt.Println("Dumping model to " + filename)

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for tag, features := range this.Features {
		f.WriteString(tag + "\n")
		// sort.Sort(features)
		for x := range features.Iter() {
			f.WriteString(fmt.Sprintf("\t%s\n", x))
		}
		f.WriteString("\n")
	}
}

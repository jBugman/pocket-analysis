package analysis

import (
	"github.com/jBugman/go-pocket/pocket"
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

func (this *Item) Tokens() []string {
	return append(getTockens(this.Title), getTockens(this.Excerpt)...)
}

func (this *Item) Bigrams() []string {
	return bigrams(this.Tokens())
}

func getTockens(source string) (result []string) {
	for _, value := range strings.Fields(source) {
		token := strings.ToLower(value)
		token = strings.Trim(token, "()[]{}“”«»")
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

type Model map[string]Counter

func (this Model) Add(tag string, key string) {
	counter, exists := this[tag]
	if exists {
		counter.Add(key)
	} else {
		this[tag] = Counter{key: 1}
	}
}

func TrainModel(items []Item) (result Model) {
	result = Model{}
	for _, item := range items {
		for _, tag := range item.Tags {
			for _, bigram := range item.Bigrams() {
				result.Add(tag, bigram)
			}
		}
	}
	return
}

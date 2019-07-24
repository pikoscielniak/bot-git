package jokes

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/mattermost/mattermost-bot-sample-golang/bot/abstract"
	"github.com/mattermost/mattermost-bot-sample-golang/bot/blacklists"
	"github.com/mattermost/mattermost-bot-sample-golang/bot/limit"
	"github.com/mattermost/mattermost-bot-sample-golang/config"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type getJoke func() []string

// all getJoke functions' algorithm:
// 1. open the joke website
// 2. get the joke's div
// 3. get rid of unnecessary text and whitespace

var jokeList []string

func Fetch() string {
	limit.AddRequest(abstract.GetUserId(), "joke")
	jokeSources := jokersPl
	if checkDay() {
		jokeSources = jokersEn
	}
	var jokeFunction getJoke
	if len(jokeList) == 0 {
		jokeFunction = jokeSources[rand.Intn(len(jokeSources))]
		jokeList = jokeFunction()
	}
	joke := jokeList[rand.Intn(len(jokeList))]
	handleBlacklist(jokeFunction, joke)
	return joke
}

func checkDay() bool {
	return time.Now().Weekday().String() == config.BotCfg.EnglishDay
}

func getFunctionName(functionReturningJoke getJoke) string {
	return runtime.FuncForPC(reflect.ValueOf(functionReturningJoke).Pointer()).Name()
}

func handleBlacklist(functionReturningJoke getJoke, jokeReturned string) {
	blacklist := blacklists.BlacklistsMap[getFunctionName(functionReturningJoke)]

	if blacklist.Contains(jokeReturned) {
		functionReturningJoke()
	}

	blacklist.AddElement(jokeReturned)

	removeFromJokeList(jokeReturned)
}

func removeFromJokeList(joke string) {
	for i,v := range jokeList {
		if v == joke {
			jokeList[i] = jokeList[len(jokeList)-1]
			jokeList = jokeList[:len(jokeList)-1]
			return
		}
	}
}

func getJokesList(selectionsToFormat *goquery.Selection) []string {

	var jokes []string
	selectionsToFormat.Each(func(i int, s *goquery.Selection) {
		selectionHTML,_ := s.Html()
		jokes = append(jokes, fixFormat(selectionHTML))
	})

	return jokes
}

func fixFormat(HTMLtoFormat string) string {
	formattedString := strings.ReplaceAll(HTMLtoFormat, "<br>", "\n")
	formattedString = strings.ReplaceAll(formattedString, "<br/>", "\n")
	formattedString = strings.ReplaceAll(formattedString, "<p>", "")
	formattedString = strings.ReplaceAll(formattedString, "</p>","")

	// markdown escape
	formattedString = strings.ReplaceAll(formattedString, "-", "\\-")
	formattedString = strings.ReplaceAll(formattedString, "*", "\\*")

	formattedString = strings.TrimSpace(formattedString)

	return formattedString
}
package core

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/rexagod/newman/core/queries"
	"k8s.io/klog/v2"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

//go:embed help.json
var marshalledHelp []byte

//go:generate echo "{\"commands\":{"

//go:generate echo "\"quote\": \"Fetches a random quote from Seinfeld, the TV show.\","
func (b *Bot) Quote(*gateway.MessageCreateEvent) (string, error) {
	quotes := R.loader.Quotes
	// assumption: quote is always of the form: `[author]: [dialog]`.
	splitQuote := strings.Split(quotes[rand.Intn(len(quotes))], ":")
	return /* emphasize the author */ "**" + splitQuote[0] + "**" + ": " +
		/* italicize the dialog*/ "_" + splitQuote[1] + "_", nil
}

//go:generate echo "\"ping\": \"Is anyone there?\","
func (b *Bot) Ping(*gateway.MessageCreateEvent) (string, error) {
	quotes := R.loader.Quotes
	// assumption: quote is always of the form: `[author]: [dialog]`.
	quote := quotes[rand.Intn(len(quotes))]
	author := strings.Split(quote, ":")[0]
	author = strings.ToUpper(author) + "!"
	var modifiedAuthor string
	for _, c := range author {
		modifiedAuthor += string(c)
		if c == ' ' {
			continue
		}
		repeatCount := rand.Intn(len(author))/2 + 1
		for i := 0; i < repeatCount; i++ {
			modifiedAuthor += string(c)
		}
	}
	return "**_" + modifiedAuthor + "_**", nil
}

//go:generate echo "\"nshow\": \"Show all notepad entries.\","
func (b *Bot) Nshow(*gateway.MessageCreateEvent) (string, error) {
	s, err := showRows(queries.Q[queries.SHOW])
	if err != nil {
		return "", fmt.Errorf("failed to show rows: %w", err)
	}
	return s, nil
}

//go:generate echo "\"nsdel\": \"Show all deleted messages.\","
func (b *Bot) Nsdel(*gateway.MessageCreateEvent) (string, error) {
	s, err := showRows(queries.Q[queries.SHOWDELETEDMESSAGES])
	if err != nil {
		return "", fmt.Errorf("failed to show rows: %w", err)
	}
	return s, nil
}

//go:generate echo "\"nadd\": \"Add a notepad entry.\","
func (b *Bot) Nadd(e *gateway.MessageCreateEvent) (string, error) {
	var err error
	op := "nadd"
	content := e.Content
	prefixLength := strings.Index(e.Content, op) + len(op) + 1
	if prefixLength >= len(content) {
		return "No input provided.", fmt.Errorf("invalid input: %s", content)
	}
	message := content[prefixLength:]
	klog.Infof("Adding entry: %s", message)
	var s string
	s, err = addRow(queries.Q[queries.ADD], message)
	if err != nil {
		return "", fmt.Errorf("failed to add row: %w", err)
	}
	return s, nil
}

//go:generate echo "\"nremove\": \"Remove a notepad entry.\","
func (b *Bot) Nremove(e *gateway.MessageCreateEvent) (string, error) {
	var err error
	op := "nremove"
	response := e.Content
	id := response[strings.Index(e.Content, op)+len(op)+1:]
	_, err = strconv.Atoi(id)
	if err != nil {
		return "", fmt.Errorf("invalid id: %w", err)
	}
	query := queries.Q[queries.REMOVE]
	_, err = R.database.ExecContext(R.databaseContext, query, id)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %w", err)
	}
	return "`Entry removed!`", nil
}

//go:generate echo "\"xkcd\": \"Fetch an XKCD comic.\","
func (b *Bot) Xkcd(e *gateway.MessageCreateEvent) (string, error) {
	var err error
	type comic struct {
		Num   int    `json:"num"`
		Image string `json:"img"`
	}
	fetchComic := func(url string) (comic, error) {
		var response *http.Response
		response, err = http.Get(url)
		if err != nil {
			return comic{}, fmt.Errorf("failed to fetch xkcd: %w", err)
		}
		defer func(response *http.Response) {
			err := response.Body.Close()
			if err != nil {
				panic(err)
			}
		}(response)
		c := comic{}
		err = json.NewDecoder(response.Body).Decode(&c)
		if err != nil {
			return comic{}, fmt.Errorf("failed to decode xkcd: %w", err)
		}
		return c, nil
	}
	var c comic
	c, err = fetchComic("https://xkcd.com/info.0.json")
	if err != nil {
		return "", fmt.Errorf("failed to fetch xkcd: %w", err)
	}
	totalComics := c.Num
	c, err = fetchComic(fmt.Sprintf("https://xkcd.com/%d/info.0.json", rand.Intn(totalComics)))
	if err != nil {
		return "", fmt.Errorf("failed to fetch xkcd: %w", err)
	}
	return c.Image, nil
}

//go:generate echo "\"short\": \"Shorten any link.\","
func (b *Bot) Short(e *gateway.MessageCreateEvent) (string, error) {
	var err error
	content := e.Content
	url := content[strings.Index(e.Content, "short")+len("short")+1:]
	response, err := http.Get(fmt.Sprintf("https://api.1pt.co/addURL?long=%s", url))
	if err != nil {
		return "", err
	}
	defer func(response *http.Response) {
		err := response.Body.Close()
		if err != nil {
			panic(err)
		}
	}(response)
	var shortURL struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Short   string `json:"short"`
		Long    string `json:"long"`
	}
	err = json.NewDecoder(response.Body).Decode(&shortURL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://1pt.co/%s", shortURL.Short), nil
}

//go:generate echo "\"help\": \"Display this message.\""
func (b *Bot) Help(*gateway.MessageCreateEvent) (string, error) {
	var err error
	type help struct {
		Commands map[string]string `json:"commands"`
	}
	var h help
	err = json.Unmarshal(marshalledHelp, &h)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal help: %w", err)
	}
	commands := []string{"```yaml\n---"}
	for k := range h.Commands {
		commands = append(commands, k+": "+h.Commands[k])
	}
	commands = append(commands, "---\n```")
	return strings.Join(commands, "\n"), nil
}

//go:generate echo "}}"

package cyoa

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
)

/*
 This is a choose your own you adventure game built in go
*/

const handlerTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Choose your own adventure</title>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{ range .Paragraphs }}
    <p>
        {{.}}
    </p>
    {{end}}

    <h3>Continuation Options</h3>
<ul>
    {{range .Options}}
    <li><a href="/{{.Arc}}" rel="noopener noreferrer">{{.Text}}</a></li>
    {{end}}
</ul>
</body>
</html>`

var tmp *template.Template

func init() {
	tmp = template.Must(template.New("").Parse(handlerTemplate))
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Story map[string]Chapter

func JSONStory(r io.Reader) (Story, error) {
	decoder := json.NewDecoder(r)
	var story Story
	err := decoder.Decode(&story)

	return story, err
}

type PathFunction func(r *http.Request) string

type ChapterHandler struct {
	story  Story
	t      *template.Template
	pathFn PathFunction
}

func WithPathFunction(pathFn PathFunction) HandlerOption {
	return func(h *ChapterHandler) {
		h.pathFn = pathFn
	}
}

func defaultPathFunction(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/" || path == "" {
		path = "/intro"
	}

	path = path[1:]

	return string(path)
}

func (c ChapterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if chapter, ok := c.story[c.pathFn((r))]; ok {
		err := tmp.Execute(w, chapter)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Fprintf(w, "Invalid chapter")
	}
}

type HandlerOption func(*ChapterHandler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *ChapterHandler) {
		h.t = t
	}
}
func NewHandler(s Story, ops ...HandlerOption) http.Handler {

	ch := ChapterHandler{s, tmp, defaultPathFunction}

	for _, opt := range ops {
		opt(&ch)
	}
	return ch
}

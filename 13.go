package main
import "fmt"
import "html/template"
import "net/http"
import "os"
import "strings"
import "time"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const TIME_FORMAT = "Mon Jan 2 15:04:05 MST 2006"

var VERSION, ADDRESS string

func init() {
	s := strings.Split(os.Args[0], "/")
	VERSION = s[len(s) - 1]

	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}
}

type Message struct {
	TimeStamp, Author, Content string
}
type PageConfiguration struct {
	Version string
	Messages []Message
}

func main() {
	p :=  PageConfiguration{ Version: VERSION }

	html := LoadTemplate(VERSION + ".html")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "text/html")
			Abort(BAD_TEMPLATE, html.Execute(w, p))

		case "POST":
			w.Header().Set("Content-Type", "text/plain")
			r.ParseForm()
			m := Message {
				TimeStamp: time.Now().Format(TIME_FORMAT),
				Author: r.PostForm.Get("a"),
				Content: r.PostForm.Get("m"),
			}
			p.Messages = append(p.Messages, m)
			fmt.Fprintf(w, "%v\n%v\n%v", m.Author, m.TimeStamp, m.Content)
		}
	})

	js_file := VERSION + ".js"
	js := LoadTemplate(js_file)
	http.HandleFunc("/" + js_file, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		Abort(BAD_TEMPLATE, js.Execute(w, p))
	})

	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func Abort(n int, e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(n)
	}
}

func LoadTemplate(s string) (r *template.Template) {
	var e error
	r, e = template.ParseFiles(s)
	Abort(FILE_READ, e)
	return
}
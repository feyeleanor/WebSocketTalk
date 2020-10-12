package main
import "fmt"
import H "html/template"
import "net/http"
import "os"
import "strings"
import "time"
import T "text/template"

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
	var e error
	var html *H.Template
	var js *T.Template

	p :=  PageConfiguration{ Version: VERSION }

	html, e = H.ParseFiles(VERSION + ".html")
	Abort(FILE_READ, e)

	js_file := VERSION + ".js"
	js, e = T.ParseFiles(js_file)
	Abort(FILE_READ, e)

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
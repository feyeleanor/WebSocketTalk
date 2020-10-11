package main
import "fmt"
import "html/template"
import "net/http"
import "os"
import "strconv"
import "time"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const VERSION = 14

const TIME_FORMAT = "Mon Jan 2 15:04:05 MST 2006"

var ADDRESS string

func init() {
	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}
}

type Message struct {
	TimeStamp, Author, Content string
}

type Handler func(http.ResponseWriter, *http.Request)
type PageConfiguration struct {
	Version int
	Messages []Message
}

func main() {
	p :=  PageConfiguration{ Version: VERSION }

	html := LoadTemplate(Filename(VERSION, "html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		Abort(BAD_TEMPLATE, html.Execute(w, p))
	})

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "text/plain")
			if i := r.URL.Query()["i"]; len(i) == 0 {
				http.NotFound(w, r)
			} else {
				if i, e := ParseIndex(i[0]); e == nil && i < len(p.Messages) {
					m := p.Messages[i]
					fmt.Fprintf(w, "%v\n%v\n%v", m.Author, m.TimeStamp, m.Content)
				} else {
					http.NotFound(w, r)
				}
			}
		case "POST":
			r.ParseForm()
			m := Message {
				TimeStamp: time.Now().Format(TIME_FORMAT),
				Author: r.PostForm.Get("a"),
				Content: r.PostForm.Get("m"),
			}
			p.Messages = append(p.Messages, m)
		}
	})

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, len(p.Messages))
	})

	js_file := Filename(VERSION, "js")
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

func Filename(n int, ext string) string {
	return fmt.Sprintf("%v.%v", n, ext)
}

func ParseIndex(s string) (int, error) {
	u, e := strconv.ParseUint(s, 10, 64)
	return int(u), e
}
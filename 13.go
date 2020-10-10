package main
import "fmt"
import "html/template"
import "net/http"
import "os"
import "time"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const VERSION = 13

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
	time.Time
	Author, Content string
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
		switch r.Method {
		case "GET":
			Abort(BAD_TEMPLATE, html.Execute(w, p))

		case "POST":
			r.ParseForm()
			m := Message {
				Time: time.Now(),
				Author: r.PostForm.Get("a"),
				Content: r.PostForm.Get("m"),
			}
			p.Messages = append(p.Messages, m)
			w.Write([]byte(
				fmt.Sprintf("%v\n%v\n%v", m.Author, m.Time.Format(TIME_FORMAT), m.Content),
			))
		}
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
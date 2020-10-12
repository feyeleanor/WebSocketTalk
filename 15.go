package main
import "fmt"
import "html/template"
import "net/http"
import "os"
import "strconv"
import "strings"
import "time"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const TIME_FORMAT = "Mon Jan 2 15:04:05 MST 2006"
const PUBLIC_ID = "public"

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
type PigeonHole []Message
type PigeonHoles map[string] PigeonHole
type PageConfiguration struct {
	Version string
	Clients int
	PigeonHoles
}

func main() {
	p :=  &PageConfiguration {
		Version: VERSION,
		PigeonHoles: make(PigeonHoles),
	}

	html := LoadTemplate(VERSION + ".html")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		Abort(BAD_TEMPLATE, html.Execute(w, p))
		p.Clients += 1
	})

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		r.ParseForm()
		switch r.Method {
		case "GET":
			q := Feed("r", r)
			ph := p.PigeonHoles[q]
			if i := MessageIndex(r); i < len(ph) {
				m := ph[i]
				fmt.Fprintf(w, "%v\n%v\n%v", m.Author, m.TimeStamp, m.Content)
			} else {
				http.NotFound(w, r)
			}
		case "POST":
			q := Feed("r", r)
			p.PigeonHoles[q] = append(p.PigeonHoles[q], Message {
				TimeStamp: time.Now().Format(TIME_FORMAT),
				Author: r.PostForm.Get("a"),
				Content: r.PostForm.Get("m"),
			})
		}
	})

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		r.ParseForm()
		fmt.Fprint(w, len(p.PigeonHoles[Feed("a", r)]))
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

func ParseIndex(s string) (int, error) {
	u, e := strconv.ParseUint(s, 10, 64)
	return int(u), e
}

func MessageIndex(r *http.Request) (i int) {
	if len(r.Form["i"]) != 0 {
		if n, e := ParseIndex(r.Form["i"][0]); e == nil {
			i = n
		}
	}
	return
}

func Feed(n string, r *http.Request) (s string) {
	if id := r.Form[strings.ToLower(n)]; len(id) == 0 {
		s = PUBLIC_ID
	} else if s = id[0]; len(s) == 0 {
		s = PUBLIC_ID
	}
	return
}
package main
import "bytes"
import "errors"
import "fmt"
import "io/ioutil"
import "net/http"
import "os"
import "strconv"
import "strings"
import "text/template"
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

type WebHandler func(http.ResponseWriter, *http.Request)
type Message struct {
	TimeStamp, Author, Content string
}
type PigeonHole []Message
type PigeonHoles map[string] PigeonHole
type PageConfiguration struct {
	Clients int
	PigeonHoles
}
func (p *PageConfiguration) AddClient(f func()) {
	f()
	p.Clients += 1
}

func main() {
	p :=  &PageConfiguration { PigeonHoles: make(PigeonHoles) }

	html, e := ioutil.ReadFile(VERSION + ".html")
	Abort(FILE_READ, e)

	js, e := template.ParseFiles(VERSION + ".js")
	Abort(FILE_READ, e)

	http.HandleFunc("/", ServeContent("text/html", html))
	http.HandleFunc("/js",
		ServeContent("application/javascript", func(*http.Request) interface{} {
			var b bytes.Buffer
			p.AddClient(func() {
				Abort(BAD_TEMPLATE, js.Execute(&b, p))			
			})
			return b.String()
		}),
	)

	http.HandleFunc("/messages",
		ServeContent("text/plain", func(r *http.Request) interface{} {
			r.ParseForm()
			return len(p.PigeonHoles[ClientID("a", r)])
		}),
	)

	http.HandleFunc("/message",
		ServeContent("text/plain", func(r *http.Request) (x interface{}) {
			r.ParseForm()
			switch r.Method {
			case "GET":
				q := ClientID("r", r)
				ph := p.PigeonHoles[q]
				if i := MessageIndex(r); i < len(ph) {
					m := ph[i]
					x = fmt.Sprintf("%v\t%v\t%v", m.Author, m.TimeStamp, m.Content)
				} else {
					x = errors.New("invalid index")
				}

			case "POST":
				q := ClientID("r", r)
				p.PigeonHoles[q] = append(p.PigeonHoles[q], Message {
					TimeStamp: time.Now().Format(TIME_FORMAT),
					Author: r.PostForm.Get("a"),
					Content: r.PostForm.Get("m"),
				})
			}
			return
		}),
	)

	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func Abort(n int, e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(n)
	}
}

func ParseIndex(s string) (int, error) {
	u, e := strconv.ParseUint(s, 10, 64)
	return int(u), e
}

func ServeContent(mime_type string, v interface{}) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime_type)
		switch v := v.(type) {
		case func(*http.Request) interface{}:
			x := v(r)
			if _, ok := x.(error); ok {
				http.NotFound(w, r)
			} else {
				fmt.Fprintf(w, "%v", x)
			}
		case []byte:
			fmt.Fprint(w, string(v))
		default:
			fmt.Fprintf(w, "%v", v)
		}
	}
}

func MessageIndex(r *http.Request) (i int) {
	if len(r.Form["i"]) != 0 {
		if n, e := ParseIndex(r.Form["i"][0]); e == nil {
			i = n
		}
	}
	return
}

func ClientID(n string, r *http.Request) (s string) {
	switch id := r.Form[strings.ToLower(n)]; {
	case len(id) == 0:
		fallthrough
	case len(id[0]) == 0:
		s = PUBLIC_ID
	default:
		s = id[0]
	}
	return
}

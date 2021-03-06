package main
import "encoding/json"
import "fmt"
import "io/ioutil"
import "net/http"
import "os"
import "strconv"
import "strings"
import "time"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const ADDRESS = ":3000"
const TIME_FORMAT = "Mon Jan 2 15:04:05 MST 2006"

type WebHandler func(http.ResponseWriter, *http.Request)

type Message struct {
	TimeStamp, Author, Content string
}

type PigeonHole []Message

func main() {
	p := make([]PigeonHole, 1)

	html, e := ioutil.ReadFile(BaseName() + ".html")
	Abort(FILE_READ, e)

	js, e := ioutil.ReadFile(BaseName() + ".js")
	Abort(FILE_READ, e)

	http.HandleFunc("/", ServeContent("text/html", html))
	http.HandleFunc("/js", ServeContent("application/javascript", js))

	http.HandleFunc("/register",
		ServeContent("application/json", func(r *http.Request) string {
			id := len(p)
			p = append(p, make(PigeonHole, 0))
			b, _ := json.Marshal(id)
			return string(b)
		}),
	)

	http.HandleFunc("/messages",
		ServeContent("application/json", func(r *http.Request) string {
			r.ParseForm()
			id := Index("a", r)
			b, _ := json.Marshal(len(p[id]))
			return string(b)
		}),
	)

	http.HandleFunc("/message",
		ServeContent("application/json", func(r *http.Request) (s string) {
			r.ParseForm()
			switch r.Method {
			case "GET":
				id := Index("r", r)
				ph := p[id]
				if i := Index("i", r); i < len(ph) {
					m := ph[i]
					b, _ := json.Marshal(m)
					s = string(b)
				}

			case "POST":
				id := Index("r", r)
				p[id] = append(p[id], Message {
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

func BaseName() string {
	s := strings.Split(os.Args[0], "/")
	return s[len(s) - 1]	
}

func Index(s string, r *http.Request) (i int) {
	switch id := r.Form[s]; {
	default:
		i, _ = strconv.Atoi(id[0])
	case len(id) == 0:
		fallthrough
	case len(id[0]) == 0:
	}
	return
}

func ServeContent(mime_type string, v interface{}) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime_type)
		switch v := v.(type) {
		case func(*http.Request) string:
			fmt.Fprintf(w, "%v", v(r))
		case []byte:
			fmt.Fprint(w, string(v))
		default:
			fmt.Fprintf(w, "%v", v)
		}
	}
}
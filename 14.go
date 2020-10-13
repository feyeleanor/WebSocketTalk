package main
import "errors"
import "fmt"
import "io/ioutil"
import "net/http"
import "os"
import "strconv"
import "strings"
import "time"

const LAUNCH_FAILED = 1
const FILE_READ = 2

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

type WebHandler func(http.ResponseWriter, *http.Request)
type Message struct {
	TimeStamp, Author, Content string
}
type PageConfiguration struct {
	Messages []Message
}

func main() {
	var p PageConfiguration

	html, e := ioutil.ReadFile(VERSION + ".html")
	Abort(FILE_READ, e)

	js, e := ioutil.ReadFile(VERSION + ".js")
	Abort(FILE_READ, e)

	http.HandleFunc("/", ServeContent("text/html", html))
	http.HandleFunc("/js", ServeContent("application/javascript", js))
	http.HandleFunc("/messages",
		ServeContent("text/plain", func(*http.Request) interface{} {
			return len(p.Messages)
		}),
	)

	http.HandleFunc("/message",
		ServeContent("text/plain", func(r *http.Request) (x interface{}) {
			r.ParseForm()
			switch r.Method {
			case "GET":
				if i := r.Form["i"]; len(i) == 0 {
					x = errors.New("no index provided")
				} else {
					if i, e := ParseIndex(i[0]); e == nil && i < len(p.Messages) {
						m := p.Messages[i]
						x = fmt.Sprintf("%v\t%v\t%v", m.Author, m.TimeStamp, m.Content)
					} else {
						x = errors.New("invalid index")
					}
				}
			case "POST":
				p.Messages = append(p.Messages, Message {
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
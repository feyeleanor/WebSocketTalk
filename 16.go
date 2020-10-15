package main
import "encoding/json"
import "fmt"
import "golang.org/x/net/websocket"
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
	events := make(chan string, 4096)
	events_broadcast := 0

	html, e := ioutil.ReadFile(BaseName() + ".html")
	Abort(FILE_READ, e)

	js, e := ioutil.ReadFile(BaseName() + ".js")
	Abort(FILE_READ, e)

	http.HandleFunc("/", Monitor(events, ServeContent("text/html", html)))
	http.HandleFunc("/js", Monitor(events, ServeContent("application/javascript", js)))

	http.HandleFunc("/register",
		Monitor(events,
			ServeContent("application/json", func(r *http.Request) string {
				id := len(p)
				p = append(p, make(PigeonHole, 0))
				b, _ := json.Marshal(id)
				return string(b)
			}),
		),
	)

	http.HandleFunc("/messages",
		Monitor(events,
			ServeContent("application/json", func(r *http.Request) string {
				r.ParseForm()
				id := Index("a", r)
				b, _ := json.Marshal(len(p[id]))
				return string(b)
			}),
		),
	)

	http.HandleFunc("/message",
		Monitor(events,
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
		),
	)

	http.HandleFunc("/events",
		Monitor(events,
			ServeContent("application/json", func(*http.Request) string {
				b, _ := json.Marshal(events_broadcast)
				return string(b)
			}),
		),
	)

	var monitors []*websocket.Conn
	go func() {
		for {
			select {
			case ev := <- events:
				events_broadcast += 1
				for _, ws := range monitors {
					websocket.JSON.Send(ws, ev)
				}
			default:
				time.Sleep(time.Millisecond)
			}
		}
	}()

	http.Handle("/monitor", websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			if e := ws.Close(); e != nil {
				events <- e.Error()
			}
		}()
		monitors = append(monitors, ws)
		ioutil.ReadAll(ws)
	}))

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

func Monitor(c chan string, f WebHandler) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		c <- r.URL.String()
		f(w, r)
	}
}
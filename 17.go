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
	monitors := make([]*websocket.Conn, 1)

	events := make(chan string, 4096)
	events_broadcast := 0

	html, e := ioutil.ReadFile(BaseName() + ".html")
	Abort(FILE_READ, e)

	js, e := ioutil.ReadFile(BaseName() + ".js")
	Abort(FILE_READ, e)

	http.HandleFunc("/", Monitor(events, ServeContent("text/html", html)))
	http.HandleFunc("/js", Monitor(events, ServeContent("application/javascript", js)))

	http.Handle("/register", websocket.Handler(func(ws *websocket.Conn) {
		events <- fmt.Sprintf("/register: %v", ws.RemoteAddr())
		defer func() {
			if e := ws.Close(); e != nil {
				events <- e.Error()
			}
		}()
		id := len(p)
		p = append(p, make(PigeonHole, 0))
		monitors = append(monitors, ws)
		websocket.JSON.Send(ws, id)
		ioutil.ReadAll(ws)
	}))

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

	go func() {
		for {
			select {
			case ev := <- events:
				events_broadcast += 1
				for i, ws := range monitors {
					if ws != nil {
						websocket.JSON.Send(ws, []interface{} {
							events_broadcast, ev, len(p[0]), len(p[i]),
						})
					}
				}
			default:
				time.Sleep(time.Millisecond)
			}
		}
	}()

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
		f(w, r)
		c <- r.URL.String()
	}
}
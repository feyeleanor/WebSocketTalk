package main
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

type PigeonHole struct {
	*websocket.Conn
	Messages []Message
	Cursor int
}

type PigeonHoles []*PigeonHole
func (p PigeonHoles) Tick(d time.Duration) {
	if p[0].Cursor < len(p[0].Messages) {
		m := p[0].Messages[p[0].Cursor]
		for _, ph := range p[1:] {
			if ph.Conn != nil {
				SendJSON(ph.Conn, "broadcast", m)
			}
		}
		p[0].Cursor += 1
	}

	for _, ph := range p[1:] {
		if ph.Conn != nil && ph.Cursor < len(ph.Messages) {
			SendJSON(ph.Conn, "private", ph.Messages[ph.Cursor])
			ph.Cursor += 1
		}
	}
	time.Sleep(d)
}

func main() {
	p := PigeonHoles{ new(PigeonHole) }

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
		p = append(p, &PigeonHole { Conn: ws })
		websocket.JSON.Send(ws, id)

		t := time.NewTicker(1 * time.Second)
		go func() {
			for {
				select {
				case ev := <- events:
					events_broadcast += 1
					for i, ph := range p {
						if ph.Conn != nil {
							websocket.JSON.Send(ws, []interface{} {
								"status",
								events_broadcast,
								ev,
								len(p[0].Messages),
								len(p[i].Messages),
							})
						}
					}

				case <- t.C:
					for i, ph := range p {
						if ph.Conn != nil {
							websocket.JSON.Send(ws, []interface{} {
								"status",
								events_broadcast,
								"",
								len(p[0].Messages),
								len(p[i].Messages),
							})
						}
					}

				default:
					events <- fmt.Sprint("send messages")
					p.Tick(100 * time.Millisecond)
				}
			}
		}()

		var b struct { Recipient, Content string }
		for {
			if e := websocket.JSON.Receive(ws, &b); e == nil {
				r, _ := strconv.Atoi(b.Recipient)
				events <- fmt.Sprint("recv:", b)
				p[r].Messages = append(p[r].Messages, Message {
					TimeStamp: time.Now().Format(TIME_FORMAT),
					Author: fmt.Sprint(id),
					Content: b.Content,
				})
			} else {
				fmt.Println("socket receive error:", e)
				break
			}
		}
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
		f(w, r)
		c <- r.URL.String()
	}
}

func SendJSON(ws *websocket.Conn, v ...interface{}) {
	websocket.JSON.Send(ws, v)
}
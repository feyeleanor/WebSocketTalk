package main
import "fmt"
import "golang.org/x/net/websocket"
import "io/ioutil"
import "log"
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

type Messages []*Message

type PigeonHole struct {
	*websocket.Conn
	Messages []*Message
	Queue chan *Message
	Quit chan struct {}
	int
}

func (p *PigeonHole) Pending() (r int) {
	if p != nil {
		r = len(p.Messages) - p.int
	}
	return
}

func (p *PigeonHole) Unsent() (m Messages) {
	return p.Messages[p.int:]
}

func (p *PigeonHole) Peek() (m *Message) {
	if p != nil {
		m = p.Messages[p.int]
	}
	return
}

func (p *PigeonHole) Advance() {
	if p != nil {
		p.int += 1
	}
}

func (p *PigeonHole) Send(s string, m interface{}) {
	if p!= nil && p.Conn != nil {
		websocket.JSON.Send(p.Conn, []interface{} { s, m })
	}	
}

func NewPigeonHole(c *websocket.Conn, f func(*PigeonHole, *Message)) (p *PigeonHole) {
	p = &PigeonHole {
		Conn: c,
		Queue: make(chan *Message, 256),
		Quit: make(chan struct{}),
	}
	go func(p *PigeonHole) {
		for p != nil {
			m := <- p.Queue
			p.Messages = append(p.Messages, m)
			f(p, m)
		}
	}(p)
	return
}

type PigeonHoles []*PigeonHole

func main() {
	var p PigeonHoles

	p = PigeonHoles{
		NewPigeonHole(nil, func(_ *PigeonHole, m *Message) {
			for ; p[0].Pending() > 0; p[0].Advance() {}
			for _, ph := range p[1:] {
				if ph.Conn != nil {
					go func(ph *PigeonHole) {
						SendStatus(ph.Conn, p[0].Messages, ph.Messages)
						for _, m := range p[0].Messages {
							ph.Send("broadcast", m)
						}
					}(ph)
				}
			}
		}),
	}

	html, js := LoadApplication(BaseName())
	http.HandleFunc("/", ServeContent("text/html", html))
	http.HandleFunc("/js", ServeContent("application/javascript", js))

	http.Handle("/register", websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			if x := recover(); x != nil {
				log.Println("socket error:", x)
			} else if e := ws.Close(); e != nil {
				log.Println("socket error:", e.Error())
			}
		}()
		id := len(p)
		p = append(p, NewPigeonHole(ws, func(ph *PigeonHole, m *Message) {
			if ph.Conn != nil {
				SendStatus(ph.Conn, p[0].Messages, ph.Messages)
				for ; ph.int < len(ph.Messages); ph.Advance() {
					ph.Send("private", m)
				}
			}
		}))
		websocket.JSON.Send(ws, id)
		SendStatus(ws, p[0].Messages, p[id].Messages)

		var b struct { Recipient, Content string }
		for {
			if e := websocket.JSON.Receive(ws, &b); e == nil {
				r, _ := strconv.Atoi(b.Recipient)
				p[r].Queue <- &Message {
					TimeStamp: time.Now().Format(TIME_FORMAT),
					Author: fmt.Sprint(id),
					Content: b.Content,
				}
			} else {
				panic(e)
			}
		}
	}))

	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func Abort(n int, e error) {
	if e != nil {
		log.Println(e)
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

func LoadApplication(name string) (html, js []byte) {
	var e error

	html, e = ioutil.ReadFile(name + ".html")
	Abort(FILE_READ, e)

	js, e = ioutil.ReadFile(name + ".js")
	Abort(FILE_READ, e)
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

func SendStatus(ws *websocket.Conn, b, p Messages) {
	websocket.JSON.Send(ws, []interface{} { "status", len(b), len(p) })
}
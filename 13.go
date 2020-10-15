package main
import "encoding/json"
import "fmt"
import "io"
import "io/ioutil"
import "net/http"
import "os"
import "strings"
import "time"
import "text/template"

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const ADDRESS = ":3000"
const TIME_FORMAT = "Mon Jan 2 15:04:05 MST 2006"

type WebHandler func(http.ResponseWriter, *http.Request)

type Message struct {
	TimeStamp, Author, Content string
}

type PageConfiguration struct {
	Messages []Message
}
type Template interface {
	Execute(io.Writer, interface{}) error
}

func main() {
	var p PageConfiguration

	html, e := ioutil.ReadFile(BaseName() + ".html")
	Abort(FILE_READ, e)

	js, e := template.ParseFiles(BaseName() + ".js")
	Abort(FILE_READ, e)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, string(html))

		case "POST":
			w.Header().Set("Content-Type", "application/json")
			r.ParseForm()
			m := Message {
				TimeStamp: time.Now().Format(TIME_FORMAT),
				Author: r.PostForm.Get("a"),
				Content: r.PostForm.Get("m"),
			}
			p.Messages = append(p.Messages, m)
			b, _ := json.Marshal(m)
			fmt.Fprint(w, string(b))
		}
	})

	http.HandleFunc("/js", ServeTemplate(js, "application/javascript", Tap(p)))

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

func Tap(v interface{}) func() interface{} {
	return func() interface{} {
		return v
	}
}

func ServeTemplate(t Template, mime_type string, f func() interface{}) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime_type)
		Abort(BAD_TEMPLATE, t.Execute(w, f()))
	}
}
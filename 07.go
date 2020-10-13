package main
import "fmt"
import "io/ioutil"
import "net/http"
import "os"
import "strings"

const LAUNCH_FAILED = 1
const FILE_READ = 2

var VERSION, ADDRESS string

type WebHandler func(w http.ResponseWriter, r *http.Request)

func init() {
	s := strings.Split(os.Args[0], "/")
	VERSION = s[len(s) - 1]

	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}
}

func main() {
	html, e := ioutil.ReadFile(VERSION + ".html")
	Abort(FILE_READ, e)

	http.HandleFunc("/", ServeStatic("text/html", html))
	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func Abort(n int, e error) {
    if e != nil {
    	fmt.Println(e)
        os.Exit(n)
    }	
}

func ServeStatic(mime_type string, b []byte) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime_type)
		fmt.Fprint(w, string(b))
	}
}

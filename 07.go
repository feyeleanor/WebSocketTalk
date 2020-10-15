package main
import "fmt"
import "io/ioutil"
import "net/http"
import "os"
import "strings"

const LAUNCH_FAILED = 1
const FILE_READ = 2

const ADDRESS = ":3000"

type WebHandler func(w http.ResponseWriter, r *http.Request)

func main() {
	html, e := ioutil.ReadFile(BaseName() + ".html")
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

func BaseName() string {
	s := strings.Split(os.Args[0], "/")
	return s[len(s) - 1]	
}

func ServeStatic(mime_type string, b []byte) WebHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime_type)
		fmt.Fprint(w, string(b))
	}
}

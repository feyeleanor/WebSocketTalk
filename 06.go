package main
import "fmt"
import "io/ioutil"
import "net/http"
import "os"
import "strings"

const LAUNCH_FAILED = 1
const FILE_READ = 2

var ADDRESS, MESSAGE string

func init() {
	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}

	s := strings.Split(os.Args[0], "/")
	html, e := ioutil.ReadFile(s[len(s) - 1] + ".html")
	Abort(FILE_READ, e)
	MESSAGE = string(html)
}

func main() {
	http.HandleFunc("/", billboard)
	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func billboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, MESSAGE)
}

func Abort(n int, e error) {
    if e != nil {
    	fmt.Println(e)
        os.Exit(n)
    }	
}
package main
import "fmt"
import "io/ioutil"
import "net/http"
import "os"

const LAUNCH_FAILED = 1
const FILE_READ = 2

var ADDRESS string

func init() {
	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}
}

func main() {
	html, e := ioutil.ReadFile("07.html")
	Abort(FILE_READ, e)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, string(html))
	})
	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func Abort(n int, e error) {
    if e != nil {
    	fmt.Println(e)
        os.Exit(n)
    }	
}
package main
import . "fmt"
import "io/ioutil"
import "net/http"
import "os"

const LAUNCH_FAILED = 1
const FILE_READ = 2

var ADDRESS, MESSAGE string

func init() {
	if p := os.Getenv("PORT"); len(p) == 0 {
		ADDRESS = ":3000"
	} else {
		ADDRESS = ":" + p
	}

	html, e := ioutil.ReadFile("06.html")
	halt_on_error(FILE_READ, e)
	MESSAGE = string(html)
}

func main() {
	http.HandleFunc("/", billboard)
	halt_on_error(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func billboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, MESSAGE)
}

func halt_on_error(n int, e error) {
    if e != nil {
    	fmt.Println(e)
        os.Exit(n)
    }	
}
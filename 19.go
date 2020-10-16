package main
import "fmt"
import "golang.org/x/net/websocket"
import "log"

const SERVER = "ws://localhost:3000/register"
const ADDRESS = "http://localhost/"

func main() {
	if ws, e := websocket.Dial(SERVER, "", ADDRESS); e == nil {
		var id int
		if e := websocket.JSON.Receive(ws, &id); e == nil {
			fmt.Println("Connected as user", id)

			var event_count, pub_count, priv_count int
			for {
				var m []interface{}
				switch e := websocket.JSON.Receive(ws, &m); {
				case e != nil:
					log.Fatal(e)

				case m[0] == "broadcast":
					x := m[1].(map[string] interface{})
					fmt.Printf("*** BROADCAST from %v ***\n", x["Author"])
					fmt.Printf("\t%v\n", x["TimeStamp"])
					fmt.Printf("\t%v\n", x["Content"])

				case m[0] == "private":
					x := m[1].(map[string] interface{})

					fmt.Printf("*** PRIVATE MESSAGE from %v ***\n", x["Author"])
					fmt.Printf("\t%v\n", x["TimeStamp"])
					fmt.Printf("\t%v\n", x["Content"])

				case m[0] == "status":
					x := int(m[3].(float64))
					if x != pub_count {
						pub_count = x
						fmt.Printf("*** STATUS UPDATE ***\n\tbroadcasts: %v\n", pub_count)
					}

					y := int(m[4].(float64))
					if y != priv_count {
						priv_count = y
						fmt.Printf("*** STATUS UPDATE ***\n\tprivate messages: %v\n", priv_count)
					}
				default:
					fmt.Printf("*** UNKNOWN ***\n\t%v\n", m)
				}
				event_count += 1
			}
		} else {
			log.Fatal(e)
		}
	} else {
	    log.Fatal(e)
	}
}

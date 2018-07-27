package _go

import "time"
import "fmt"

func main() {
	var ts = time.Now()
	println("time", ts.UnixNano())
	fmt.Printf("%x\n", ts.UnixNano())
}

package main

import (
	"os"
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
	"strconv"
)

func getClean(val int, myBucket *gocb.Bucket, ch chan bool) int {

	if val == 0 {
		/*fmt.Println("Done: ", ch)
		go done()*/
		ch <- true
	}

	fmt.Println("got here")

	var keya string
	var keyc string
	str := strconv.Itoa(val)
	keya = ("ACCT::"+ str)
	keyc = ("CUST::"+ str)
	myBucket.Remove(keya, 0)
	myBucket.Remove(keyc, 0)
	return val

}

func main() {

	var seedNode string
	var i = 50
	ch := make(chan bool)
	seedNode = ("couchbase://" + os.Args[1])

	// Connect to Couchbase
	myC, _ := gocb.Connect(seedNode)
	myB, _ := myC.OpenBucket("default", "")


	for i <= 50 && i >0 {
		fmt.Println("got here")
		go getClean(i, myB, ch)
		i = i - 1
	}

	<-ch
	if <-ch == true {
		fmt.Println("Done: ", ch)
		go thend()
	}
	fmt.Println("done")
}

// Shut down.
func thend() {
	fmt.Println("Good Bye")
	os.Exit(111)
}

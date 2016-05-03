package main

import (
	"os"
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
	"strconv"
)

func main() {

	var seedNode string
	var key string
	var i int
	seedNode = ("couchbase://" + os.Args[1])

	// Connect to Couchbase
	myC, _ := gocb.Connect(seedNode)
	myB, _ := myC.OpenBucket("default", "")

	for i <= 50 && i >0 {
		fmt.Println("got here")
		str := strconv.Itoa(i)
		key = ("ACCT::"+str)
		myB.Remove(key, nil)

		i = i - 1
	}
	fmt.Println("done")
}

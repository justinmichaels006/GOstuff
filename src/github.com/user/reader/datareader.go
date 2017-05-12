package main

import (
	"fmt"
	"os"
	"strconv"
	"gopkg.in/couchbase/gocb.v1"
	//"gopkg.in/matryer/try.v1"
	//"log"
	"time"
)

func checker(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

// Shut down.
func done() {
	fmt.Println("Good Bye")
	os.Exit(111)
}

func getOne(val int, myBucket *gocb.Bucket, ch chan bool) int {

	if val == 0 {
		/*fmt.Println("Done: ", ch)
		go done()*/
		ch <- true
	}

	var err error
	var theDoc map[string]interface{}

	str := strconv.Itoa(val)
	//var rtnValue bool
	var key string

	key = ("ACCT::"+str)
	// Retrieve Document
	_, err = myBucket.Get(key, &theDoc)
	if err != nil {
		fmt.Println("Doh! ", key)
		//time.Sleep(100000)
		return getOne(val, myBucket, ch)
	}

	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	theDoc["getStamp"] = millis
	myBucket.Upsert(key, theDoc, 0)

	myBucket.Get(key, theDoc)
	fmt.Println("Got It:", key)
	return val + getOne(val-1, myBucket, ch)

}

func main() {

	var seedNode string
	// holds the arguments for Couchbase seed node
	seedNode = ("couchbase://" + os.Args[1])

	myC, _ := gocb.Connect(seedNode)
	myB, _ := myC.OpenBucket("default", "")

	var opsGroups = 50
	ch := make(chan bool)

	for opsGroups <= 50 && opsGroups >0 {
		fmt.Println("got here")
		go getOne(opsGroups, myB, ch)
		opsGroups = opsGroups - 1
	}

		//One way to do it
		/*err := try.Do(func(attempt int) (bool, error) {
			var err error
			rtnValue, err = getOne(key, myB)
			return attempt < 5, err // try 5 times
		})*/
		/*if err != nil {
			log.Fatalln("error:", err)
		}*/

	<-ch
	if <-ch == true {
		fmt.Println("Done: ", ch)
		go done()
	}

}


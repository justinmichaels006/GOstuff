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

func getOne(val int, myBucket *gocb.Bucket, ch chan bool) { //(bool, error)

	var err error
	var theDoc map[string]interface{}

	str := strconv.Itoa(val)
	//var rtnValue bool
	var key string

	key = ("ACCT::"+str)
	// Retrieve Document
	_, err = myBucket.Get(key, &theDoc)
	if err != nil {
		fmt.Println("Doh! ", str)
		getOne(val, myBucket, ch)
	}
	if err == nil {
		//itemsGet = append(itemsGet, &gocb.GetOp{Key: val, Value: &theDoc{}})
		//Add time stamp for when this occured
		//Update log file? or maybe the doc itself
		now := time.Now()
		nanos := now.UnixNano()
		millis := nanos / 1000000
		theDoc["getStamp"] = millis
		myBucket.Upsert(str, theDoc, 0)
		fmt.Println("Got It:", str)
		return
	}
	if val == 0 {
		ch <- true
		return
	}

	//return status, err

}

func main() {

	var seedNode string
	// holds the arguments for Couchbase seed node
	seedNode = ("couchbase://" + os.Args[1])

	myC, _ := gocb.Connect(seedNode)
	myB, _ := myC.OpenBucket("testload", "")

	var opsGroups = 5
	var i = 0
	ch := make(chan bool)

	for i < opsGroups {

		go getOne(opsGroups, myB, ch)

		//One way to do it
		/*err := try.Do(func(attempt int) (bool, error) {
			var err error
			rtnValue, err = getOne(key, myB)
			return attempt < 5, err // try 5 times
		})*/
		/*if err != nil {
			log.Fatalln("error:", err)
		}*/
	i++
	}

	if <-ch {
		fmt.Println("Done: ", ch)
		go done()
	}

}


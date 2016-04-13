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
	os.Exit(101)
}

func getOne(val string, myBucket *gocb.Bucket, ch chan bool) { //(bool, error)

	var err error
	//var status bool
	//var itemsGet []gocb.BulkOp
	var theDoc map[string]interface{}

	// Retrieve Document
	_, err = myBucket.Get(val, &theDoc)
	if err != nil {
		fmt.Println("Doh! ", val)
		//status = false
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
		myBucket.Upsert(val, theDoc, 0)
		fmt.Println("Got It:", val)
		//status = true
		ch <- true
		return
	}

	//return status, err

}

func main() {

	myC, _ := gocb.Connect("couchbase://10.111.94.102")
	myB, _ := myC.OpenBucket("testload", "")

	var opsGroups = 5
	var i = 0
	ch := make(chan bool)

	for i < opsGroups {

		str := strconv.Itoa(i)
		//var rtnValue bool
		var key string

		key = ("ACCT::"+str)

		go getOne(key, myB, ch)

		//One way to do it
/*		err := try.Do(func(attempt int) (bool, error) {
			var err error
			rtnValue, err = getOne(key, myB)
			return attempt < 5, err // try 5 times
		})*/
		/*if err != nil {
			log.Fatalln("error:", err)
		}*/
	i++
	}

	<-ch
	go done()
}


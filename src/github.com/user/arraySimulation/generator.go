package main

import (
"gopkg.in/couchbase/gocb.v1"
"fmt"
"os"
"encoding/json"
"strconv"
"io"
"crypto/rand"
)

func main() {

	// Configuration
	// 1 Customer -> 800 to 1200 Groups -> 40 to 80 Devices -> 70 to 200 App
	// 10,000 Customers -> 1M Groups -> 50M Devices -> 2B App
	var currentGroup = 0;
	var cusomterTotal = 2;
	var groupTotal = cusomterTotal * 1000
	var deviceTotal = groupTotal * 50
	var appTotal = deviceTotal * 100
	// Create an Array of BulkOps for Insert
	var itemGroups []gocb.BulkOp
	var itemDevice []gocb.BulkOp
	var itemApp []gocb.BulkOp
	var runningLoad = false
	var seedNode string
	// holds the arguments for Couchbase seed node
	//seedNode = ("couchbase://" + os.Args[1])
	seedNode = ("couchbase://" + "127.0.0.1")

	// Connect to Couchbase
	myC, _ := gocb.Connect(seedNode)
	myB, _ := myC.OpenBucket("streamer", "")

	// Read the Greoup file
	tmpGROUP, err := os.OpenFile("/Users/justin/Documents/Symantec/sampledata/GROUP.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docGROUP map[string]interface{}
	err = json.NewDecoder(tmpGROUP).Decode(&docGROUP)
	if err != nil {
		fmt.Println(err)
		os.Exit(52)
	}

	tmpDEVICE, err := os.OpenFile("/Users/justin/Documents/Symantec/sampledata/DEVICE.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docDEVICE map[string]interface{}
	err = json.NewDecoder(tmpDEVICE).Decode(&docDEVICE)
	if err != nil {
		fmt.Println(err)
		os.Exit(53)
	}

	tmpAPP, err := os.OpenFile("/Users/justin/Documents/Symantec/sampledata/APP.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docAPP map[string]interface{}
	err = json.NewDecoder(tmpAPP).Decode(&docAPP)
	if err != nil {
		fmt.Println(err)
		os.Exit(54)
	}

	for x := cusomterTotal; x != 0; x-- {
		// Create the Customer that will anchor the rest of the relationships
		uuid, err := newUUID()

		var appArray [appTotal]string

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		//fmt.Printf("%s\n", uuid) // Debug

		for i := 0; i < groupTotal; i++ {
			itemGroups = append(itemGroups, &gocb.InsertOp{Key: uuid + "::GROUP::" + strconv.Itoa(i), Value: &docGROUP})
			for j := 0; j < deviceTotal; j++ {
				for y := 0; y < appTotal; y++ {
					appArray[y] = uuid + "::APP::" + strconv.Itoa(y)
				}
					docDEVICE["APP_install"] = appArray
					itemDevice = append(itemDevice, &gocb.InsertOp{Key: uuid + "::DEVICE::" + strconv.Itoa(j), Value: &docDEVICE})
				for k := 0; k < appTotal; k++ {
					itemApp = append(itemApp, &gocb.InsertOp{Key: uuid + "::APP::" + strconv.Itoa(k), Value: &docAPP})
				}
				// Perform the bulk operation
				err = myB.Do(itemApp)
				if err != nil {
					fmt.Println("ERRROR PERFORMING BULK INSERT:", err)
				}
			}
			// Perform the bulk operation
			err = myB.Do(itemDevice)
			if err != nil {
				fmt.Println("ERRROR PERFORMING BULK INSERT:", err)
			}
		}
		// Perform the bulk operation
		err = myB.Do(itemGroups)
		if err != nil {
			fmt.Println("ERRROR PERFORMING BULK INSERT:", err)
		}

	finish(myB)
	}

	//if (runningLoad == false && currentGroup == 0) {
	//	FlowControl(runningLoad, uuid, cusomterTotal, docGROUP, docDEVICE, docAPP, myB)
	//}
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// Shut down.
func finish(cBucket *gocb.Bucket) {
	cBucket.Close()
	fmt.Println("Good Bye")
	os.Exit(101)
}

/*func FlowControl(controller bool, number int, uuid string,
	jsonGROUP map[string]interface{},
	jsonDEVICE map[string]interface{},
	jsonAPP map[string]interface{}, myBucket *gocb.Bucket) int {

	if number == 0 {
		fmt.Println("Done: ", number)
		go finish(myBucket)
	}

	str := strconv.Itoa(number)
	//fmt.Println("Upsert: ", str)

	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000

	jsonGROUP["upStamp"] = millis
	jsonDEVICE["upStamp"] = millis
	jsonAPP["upStamp"] = millis
	myBucket.Upsert("GROUP::"+str, jsonGROUP, 0)
	myBucket.Upsert("DEVICE::"+str, jsonDEVICE, 0)
	myBucket.Upsert("APP::"+str, jsonAPP, 0)

	return number + FlowControl(controller, uuid, number-1, jsonGROUP, jsonDEVICE, jsonAPP, myBucket)
}*/

/*func upsertOne(num int, jsonACCT map[string]interface{}, jsonCUST map[string]interface{}, myBucket *gocb.Bucket) {

	// Upsert One Document
	str := strconv.Itoa(num)
	fmt.Println("Upsert: ", str)
	myBucket.Upsert("ACCT::"+str, jsonACCT, 0)
	myBucket.Upsert("CUST::"+str, jsonCUST, 0)
}*/

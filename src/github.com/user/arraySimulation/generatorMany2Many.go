package main

import (
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
	"os"
	"encoding/json"
	"strconv"
	"math/rand"
	"github.com/nu7hatch/gouuid"
	"time"
)

func main() {

	// Configuration
	// 1 Customer -> 800 to 1200 Groups -> 40 to 80 Devices -> 70 to 200 App
	// 10,000 Customers -> 1M Groups -> 50M Devices -> 2B App
	//TODO: Create ranges during simulation
	var cusomterTotal = 2 //10000;
	var groupTotal = cusomterTotal * 2
	var deviceTotal = groupTotal * 2
	//var appTotal = 200
	var appCatalog = 20
	var MaxBatch = 1
	rand.Seed(time.Now().Unix())
	// Create an Array of BulkOps for Insert
	var itemCust []gocb.BulkOp
	var itemGroups []gocb.BulkOp
	var itemDevice []gocb.BulkOp
	var itemApp []gocb.BulkOp
	//var itemREL []gocb.BulkOp
	// Channel controls
	fCust := make(chan bool)
	var seedNode string

	// holds the arguments for Couchbase seed node
	//seedNode = ("couchbase://" + os.Args[1])
	seedNode = ("couchbase://" + "192.168.61.101")

	// Connect to Couchbase
	myC, _ := gocb.Connect(seedNode)
	myB, _ := myC.OpenBucket("testload", "")

	// Read the file for each doc type
	tmpGROUP, err := os.OpenFile("/Users/justin/Documents/Symantec/sampledata/GROUP.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
		os.Exit(50)
	}
	var docGROUP map[string]interface{}
	err = json.NewDecoder(tmpGROUP).Decode(&docGROUP)
	if err != nil {
		fmt.Println(err)
		os.Exit(51)
	}

	tmpDEVICE, err := os.OpenFile("/Users/justin/Documents/Symantec/sampledata/DEVICE.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
		os.Exit(52)
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

	/* For when the relationship is created
	tmpREL, err := os.OpenFile("/Users/justin/Documents/Symantec/sampledata/REL.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docREL map[string]interface{}
	err = json.NewDecoder(tmpAPP).Decode(&docAPP)
	if err != nil {
		fmt.Println(err)
		os.Exit(54)
	}*/

	var docCUST map[string]interface{}
	str1 := `{"TYPE": "CUSTONMER", "ID": "uuid", "NUM": "x"}`
	json.Unmarshal([]byte(str1), &docCUST)

	var docAPPFile map[string]interface{}
	str2 := `{"TYPE": "APPFILE", "NUM": "x"}`
	json.Unmarshal([]byte(str2), &docAPPFile)

	//Create Simulated App Catalog
	for y := 0; y < appCatalog; y++ {
		docAPP["TYPE"] = "APP"
		itemApp = append(itemApp, &gocb.InsertOp{Key: "APP::" + strconv.Itoa(y), Value: &docAPP})
		if len(itemApp) >= MaxBatch {
			ops := itemApp
			itemApp = nil
			myB.Do(ops)
		}
	}

	for x := cusomterTotal; x != 0; x-- {
		// Create the Customer that will anchor the rest of the relationships
		uuid, err := newUUIDMany()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		docCUST["TYPE"] = "CUSTOMER"
		docCUST["ID"] = uuid
		docCUST["NUM"] = x

		itemCust = append(itemCust, &gocb.InsertOp{Key: uuid + "::" + strconv.Itoa(x), Value: &docCUST})
		if len(itemCust) >= MaxBatch {
			ops1 := itemCust
			itemCust = nil
			myB.Do(ops1)
		}

		//Create Devices and Groups
		for i := 0; i < groupTotal; i++ {
			docGROUP["TYPE"] = "GROUP"
			docGROUP["GROUP_id"] = uuid

			itemGroups = append(itemGroups, &gocb.InsertOp{Key: uuid + "::GROUP::" + strconv.Itoa(i), Value: &docGROUP})
			if len(itemGroups) >= MaxBatch {
				ops2 := itemGroups
				itemGroups = nil
				myB.Do(ops2)
			}


			for j := 0; j < deviceTotal; j++ {
				docDEVICE["TYPE"] = "DEVICE"

				itemDevice = append(itemDevice, &gocb.InsertOp{Key: uuid + "::DEVICE::" + strconv.Itoa(j), Value: &docDEVICE})
				if len(itemDevice) >= MaxBatch {
					ops3 := itemDevice
					itemDevice = nil
					myB.Do(ops3)
				}
			// fmt.Println("Got this far DEVICE") //debug
			//fDevice <- true
			}
		// fmt.Println("Got this far GROUP") //debug
		//fGroup <- true
		}
	//fmt.Println("Got this far CUST") //debug
	//fCust <- true
	}

<-fCust

finishMany(myB)
}

// newUUID generates a random UUID according to RFC 4122
func newUUIDMany() (string, error) {
	/*This won't work on Windows
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}*/
	out, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return out.String(), err
	}
	//fmt.Println(out)
	return out.String(), err
}

// Shut down.
func finishMany(cBucket *gocb.Bucket) {
	cBucket.Close()
	fmt.Println("Good Bye")
	os.Exit(101)
}

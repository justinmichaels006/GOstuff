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
	var cusomterTotal = 100 //10000;
	var groupTotal = cusomterTotal * 1000
	var deviceTotal = 50
	var appTotal = 200
	var appCatalog = 2000
	var MaxBatch = 50
	rand.Seed(time.Now().Unix())
	// Create an Array of BulkOps for Insert
	var itemCust []gocb.BulkOp
	var itemGroups []gocb.BulkOp
	var itemDevice []gocb.BulkOp
	var itemApp []gocb.BulkOp
	var itemAppFile []gocb.BulkOp
	var seedNode string
	appArray := make([]string, appTotal)

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
		for m := 0; m < 10; m++ {
			docAPPFile["TYPE"] = "APPFILE"
			docAPPFile["NAME"] = "some.dll"
			itemAppFile = append(itemAppFile, &gocb.InsertOp{Key: "APP::" + strconv.Itoa(y) + "::" + strconv.Itoa(m), Value: &docAPPFile})
			if len(itemAppFile) >= 10 {
				ops := itemAppFile
				itemAppFile = nil
				myB.Do(ops)
			}
		}
		//fmt.Println("Got this far APP") //debug
	}

	for x := cusomterTotal; x != 0; x-- {
		// Create the Customer that will anchor the rest of the relationships
		uuid, err := newUUID()
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
				for m := 0; m < appTotal; m++ {
					var a = rand.Intn(appCatalog)
					// fmt.Println(a) //debug
					appArray[m] = "APP::" + strconv.Itoa(a) + "::" + strconv.Itoa(rand.Intn(10))
				}
				docDEVICE["APP_install"] = appArray

				itemDevice = append(itemDevice, &gocb.InsertOp{Key: uuid + "::DEVICE::" + strconv.Itoa(j), Value: &docDEVICE})
				if len(itemDevice) >= MaxBatch {
					ops3 := itemDevice
					itemDevice = nil
					myB.Do(ops3)
				}
			// fmt.Println("Got this far DEVICE") //debug
			}
		// fmt.Println("Got this far GROUP") //debug
		}
	//fmt.Println("Got this far CUST") //debug
	}
finish(myB)
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
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
func finish(cBucket *gocb.Bucket) {
	cBucket.Close()
	fmt.Println("Good Bye")
	os.Exit(101)
}

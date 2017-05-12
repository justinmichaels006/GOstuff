package main

import (
	"gopkg.in/couchbase/gocb.v1"
	"fmt"
	"os"
	"encoding/json"
	"strconv"
	"github.com/nu7hatch/gouuid"
	"time"
)

func main() {

	var userTotal = 10 //100 //10000;
	var activityTotal = userTotal * 10 //1000
	fControl := make(chan bool)
	var seedNode string
	// holds the arguments for Couchbase seed node
	//seedNode = ("couchbase://" + os.Args[1])
	seedNode = ("couchbase://" + "192.168.61.101")

	// Connect to Couchbase
	myC, _ := gocb.Connect(seedNode)
	myB, _ := myC.OpenBucket("testload", "")

	// Read the user profile file
	tmpACTIVITY, err := os.OpenFile("/Users/justin/Documents/Symantec/sampledata/GROUP.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
		os.Exit(50)
	}
	var docACTIVITY map[string]interface{}
	err = json.NewDecoder(tmpACTIVITY).Decode(&docACTIVITY)
	if err != nil {
		fmt.Println(err)
		os.Exit(51)
	}

	// Activuty document for each user
	var docUSER map[string]interface{}
	str := `{"TYPE": "USER", "ID": "uuid", "DESCRIPTION": "xxx"}`
	json.Unmarshal([]byte(str), &docUSER)

	for x := userTotal; x != 0; x-- {

		// Create the UUID that will anchor the rest of the relationships
		uuid, err := newUUID2()

		go func() {
			for {
				if err != nil {
					fmt.Println(err)
					os.Exit(101)
				}
				// TOOO Add faker code here
				docUSER["TYPE"] = "USER"
				docUSER["ID"] = "USER::" + uuid
				userID := "USER::" + uuid
				now := time.Now()
				nanos := now.UnixNano()
				millis := nanos / 1000000
				docUSER["upStamp"] = millis

				myB.Upsert(userID, docUSER, 0)
			}
		}()

		//Create Groups
		for i := 0; i < activityTotal; i++ {
			docUSER["TYPE"] = "ACTIVITY"
			docUSER["ACTIVITY_id"] = uuid + "::ACTIVITY::" + strconv.Itoa(i)
			activityID := uuid + "::ACTIVITY::" + strconv.Itoa(i)

			//(int, uuid, type, jsonDOC, myBucket, chan)
			go FlowControl22(activityID, docUSER, myB)

		}
		// Where does my final channel control
		// fControl <- true
	}

	fmt.Println("Got this far")

	if <-fControl {
		fmt.Println("Done")
		go finish2(myB)
	}

}

func FlowControl22(theID string, jsonDOC map[string]interface{}, myBucket *gocb.Bucket) {
	myBucket.Upsert(theID, jsonDOC, 0)
	return
}

// newUUID generates a random UUID according to RFC 4122
func newUUID2() (string, error) {
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
func finish2(cBucket *gocb.Bucket) {
	cBucket.Close()
	fmt.Println("Good Bye")
	os.Exit(101)
}

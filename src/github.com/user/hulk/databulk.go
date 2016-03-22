package main

import (
	"github.com/couchbaselabs/gocb"
	"fmt"
	"strconv"
	"sync"
	"os"
	//"io/ioutil"
	"math/rand"
	"encoding/json"
)

func f(jsonA map[string]interface{}, jsonCH map[string]interface{}, jsonCONT map[string]interface{}, myBucket *gocb.Bucket, wg sync.WaitGroup) {
	for i := 1; i < 200; i++ {

		i := strconv.Itoa(i)
		chID := ("CH::" + i);
		contID := ("CONT::" + i);

		//Element a3 and ch11 should be the same
		Aelement3_CHelement11 := RandIntBytes(4)
		//ch.element41=13
		CHelement41 := RandIntBytes(2)
		//ch.element57="WNET"
		CHelement57 := RandStringBytes(4)

		/*var datCH map[string]interface{}
		if err := json.Unmarshal(jsonCH, datCH);
		err != nil {
			fmt.Println(jsonCH)
			fmt.Println("marshal error ... ", err)
		}*/

		jsonCH["Ch_element11"] = Aelement3_CHelement11
		jsonCH["Ch_element41"] = CHelement41
		jsonCH["Ch_element57"] = CHelement57

		//jsonCH, _ := json.Marshal(jsonCH)

		myBucket.Upsert(chID, jsonCH, 0);
		myBucket.Upsert(contID, jsonCONT, 0);

		//For every CH create 3000 A documents

		for y := 1; y < 3; y++ {

			y := strconv.Itoa(y)
			aID := ("A::" + i + "::" + y)

			/*var datA map[string]interface{}
			if err := json.Unmarshal(jsonA, datA);
			err != nil {
				panic(err)
			}*/

			jsonA["A_element3"] = Aelement3_CHelement11
			//jsonA, _ := json.Marshal(datA)

			myBucket.Upsert(aID, jsonA, 0)
			//fmt.Println("Debug ... y" + Aelement3_CHelement11)
		};

		//fmt.Println("Debug ... i" + i)
	}

	//fmt.Println("Wait Group Done ...")
	wg.Done()
}

/*func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func RandIntBytes(n int) string {
	const numBytes = "0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = numBytes[rand.Intn(len(numBytes))]
	}
	return string(b)
}*/

func main() {

	myCluster, _ := gocb.Connect("couchbase://192.168.5.81")
	myBucket, _ := myCluster.OpenBucket("default", "")

	// read whole the file
	tmpA, err := os.OpenFile("/Users/justin/GOstuff/a.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docA map[string]interface{}
	err = json.NewDecoder(tmpA).Decode(&docA)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	tmpCH, err := os.OpenFile("/Users/justin/GOstuff/ch.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docCH map[string]interface{}
	err = json.NewDecoder(tmpCH).Decode(&docCH)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	tmpCONT, err := os.OpenFile("/Users/justin/GOstuff/cont.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	var docCONT map[string]interface{}
	err = json.NewDecoder(tmpCONT).Decode(&docCONT)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	var wg sync.WaitGroup

	fmt.Println("Data creation ...")

	wg.Add(1)
	go f(docA, docCH, docCONT, myBucket, wg)
	wg.Wait()

	os.Exit(0)
}

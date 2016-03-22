package main

import (
	"github.com/couchbaselabs/gocb"
	"fmt"
	"strconv"
	"sync"
	"encoding/json"
)

type data struct {
	segment_num string `json:"segment_num"`
	jump_time string `json:"jump_time"`
	stream_id string `json:"stream_id"`
	segment string `json:"segment"`
	period_first string `json:"period_first"`
	period_current string `json:"period_current"`
	period_lat string `json:"period_last"`
	rep_id string `json:"rep_id"`
	start_time string `json:"start_time"`
	duration_pte string `json:"duration_pts"`
	actual_copy_cnt string `json:"actual_copy_cnt"`
	is_archived string `json:"is_archived"`
	batch_list string `json:"batch_list`
}

func f(batch string, myBucket *gocb.Bucket, wg sync.WaitGroup) {
	for i := 1; i < 1801; i++ {

		i := strconv.Itoa(i)
		fmt.Println("conversion check" + batch, ":", i)

		var segmentDoc interface {}
			myBucket.Get(batch + "::" + i, &segmentDoc)
		fmt.Println(batch + "::" + i)

		err := json.Unmarshal(segmentDoc, &data)
		if err != nil {
			panic(err)
		}
		fmt.Println(data.segment_num)

	}
	wg.Done()
}


func main() {

	myCluster, _ := gocb.Connect("couchbase://192.168.106.101")
	myBucket, _ := myCluster.OpenBucket("recordings", "")

	var wg sync.WaitGroup

	fmt.Println("Get segments for recording 10 ...")

	wg.Add(1)
	go f("batch10", myBucket, wg)
	wg.Wait()
}

package main

import (
	nodecapacitylist "costkube/nodecapacity"
	typelister "costkube/typelist"
	"fmt"
	"time"
	"github.com/go-co-op/gocron"
)

func main() {
	// fmt.Println("Staring Data collection")
	// // dt := time.Now()
	// fmt.Println(dt.Format("01-02-2006 15:04:05 Monday"))
	// usageinfo := typelister.KindLister()
	// // fmt.Println(usageinfo)
	// nodecapacitylist.NodeCacityLister(usageinfo)
	// fmt.Println("Data Sent to Elatic Search")
	// // end := time.Now()
	// fmt.Println(end.Format("01-02-2006 15:04:05 Monday"))

	runcron()




}

func runcron(){
	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Seconds().Do(func() {
		dataCollection()
	})
	s.StartBlocking()
}

func dataCollection(){
	fmt.Println("Staring Data collection")
	usageinfo := typelister.KindLister()
	nodecapacitylist.NodeCacityLister(usageinfo)
	fmt.Println("Data Sent to Elatic Search")

}
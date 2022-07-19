package main

import (
	nodecapacitylist "costkube/nodecapacity"
	typelister "costkube/typelist"
	"fmt"
)

func main() {
	fmt.Println("Staring Data collection")
	usageinfo := typelister.KindLister()
	fmt.Println(usageinfo)
	nodecapacitylist.NodeCacityLister(usageinfo)
	fmt.Println("Data Sent to Elatic Search")
}

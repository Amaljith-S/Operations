package nodecapacitylist

import (
	"context"
	nodeusagemetrics "costkube/nodeinfo"
	"log"
	"strconv"
	"reflect"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	
)

// type NodeAllocatedResources struct {
// 	MemoryCapacity int64 `json:"memoryCapacity"`

// }

func NodeCacityLister(usageinfo []string) {

	config, err := clientcmd.BuildConfigFromFlags("", "/home/amaljith/.kube/config")
	if err != nil {
		panic(err.Error())
	}
	// create the kubeClient
	kubeClient, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	var nodeinfo []string

	nodes, _ := kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	for _, node := range nodes.Items {
		memoryAsKb, okay := node.Status.Capacity.Memory().AsInt64()
		cpuCoreCount, okay := node.Status.Capacity.Cpu().AsInt64()
		fmt.Println(reflect.TypeOf(cpuCoreCount))
		for key, value := range node.Status.Capacity {
			if key == "ephemeral-storage"{
				diskDetails:= value
				fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@",diskDetails)
				fmt.Printf("%s value is %v\n", key, value)
			} 
			
			}


		memoryAsGb := memoryAsKb / (1024 * 1024)

		coreTotalM := (cpuCoreCount * 1000)
		// fmt.Println("----------------------------------->>>>>",reflect.TypeOf(coreTotalM))


		if err != nil {
			log.Printf(err.Error())
		}
		if !okay {
			return
		}
		coreTotalMString := strconv.FormatInt(coreTotalM, 10)
		memoryAsGbString := strconv.FormatInt(memoryAsGb, 10)
		nodeinfo = append(nodeinfo, node.Name, coreTotalMString, memoryAsGbString)
		// fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@",nodeinfo)

	}

	nodeusagemetrics.NodeUsageMetricsCollector(nodeinfo, usageinfo)
	if err != nil {
		log.Printf(err.Error())
	}
}

package nodeusagemetrics

import (
	"context"
	elastic "costkube/elastic"
	elasticdatatype "costkube/elasticindex"
	"encoding/json"
	"fmt"
	
	"strconv"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NodeMetricsList : NodeMetricsList
type NodeMetricsList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		Metadata struct {
			Name              string    `json:"name"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Label             *struct {
			} `jason:"lables"`
		} `json:"metadata"`
		Timestamp time.Time `json:"timestamp"`
		Window    string    `json:"window"`
		Usage     struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"usage"`
	} `json:"items"`
}

func sum(array []int) int { // int to float
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}

func NodeUsageMetricsCollector(nodeCapacity []string, usageinfo []string) {

	config, err := clientcmd.BuildConfigFromFlags("", "/home/amaljith/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	var nodes NodeMetricsList
	data, err := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").DoRaw(context.TODO())
	// fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^",data)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &nodes)

	var nodeInfoUsage []string

	// fmt.Println("******************************************************",nodes)

	for _, m := range nodes.Items {
		trimmedCpu := strings.Trim(m.Usage.CPU, "n")
		trimmedMemory := strings.Trim(m.Usage.Memory, "Ki")
		intTrimmedCpu, err := strconv.Atoi(trimmedCpu)
		intTrimmedMemory, err := strconv.Atoi(trimmedMemory)
		finalTrimMemory := intTrimmedMemory / 1024
		finalTrimCpu := intTrimmedCpu / (1000 * 1000)

		if err != nil {
			return
		}
		elastic.NodeInfoInserter(m.Metadata.Name, finalTrimCpu, finalTrimMemory)
		intfinalTrimCpu := strconv.Itoa(finalTrimCpu)
		intfinalTrimMemory := strconv.Itoa(finalTrimMemory)
		nodeInfoUsage = append(nodeInfoUsage, m.Metadata.Name, intfinalTrimCpu, intfinalTrimMemory)
	}

	NodePercentageCalculation(nodeInfoUsage, nodeCapacity, usageinfo)

	return
}

func NodePercentageCalculation(nodeInfoUsage []string, nodeCapacity []string, usageinfo []string) {

	var totalNodeInfo []string
	finalintcpuCapacity := 0
	finalintMemoryCapacity := 0
	RangeValue := len(nodeInfoUsage)
	for i, j, k := 0, 1, 2; i <= RangeValue-1; i, j, k = i+3, j+3, k+3 {
		nodeName := nodeInfoUsage[i]
		nodeCpuUsage := nodeInfoUsage[j]
		nodeMemoryUsage := nodeInfoUsage[k]

		intNodeCpuUsage, err := strconv.Atoi(nodeCpuUsage)

		intNodeMemoryUsage, err := strconv.Atoi(nodeMemoryUsage)
		// fmt.Println("intNodeMemoryUsage = ", reflect.TypeOf(intNodeMemoryUsage))
		intcputotalCapacity, err := strconv.Atoi(nodeCapacity[j])
		inttotalMemoryCapacity, err := strconv.Atoi(nodeCapacity[k])
		if err != nil {
			return
		}

		if nodeName == nodeCapacity[i] {
			intcpuCapacity, err := strconv.Atoi(nodeCapacity[j])
			intMemoryCapacity, err := strconv.Atoi(nodeCapacity[k])
			// fmt.Println("#################################################################",intMemoryCapacity)
			cpuPercentage := (float64(intNodeCpuUsage) / float64(intcpuCapacity)) * 100
			memoryPercentage := (float64(intNodeMemoryUsage) / float64(intMemoryCapacity)) * 100
			fmt.Println("cpu on", nodeName, "is", int(cpuPercentage), "%")    // int need to change as float ********
			fmt.Println("mem on", nodeName, "is", int(memoryPercentage), "%") // int need to change as float*************
			stringcpuPercentage := strconv.Itoa(int(cpuPercentage))           // int need to change as float******
			stringmemoryPercentage := strconv.Itoa(int(memoryPercentage))     // int need to change as float
			totalNodeInfo = append(totalNodeInfo, nodeName, stringcpuPercentage, stringmemoryPercentage)
			elastic.NodePercentageInserter(nodeName, int(cpuPercentage), int(memoryPercentage)) // int need to change as float

			if err != nil {
				return
			}

		}
		finalintcpuCapacity = finalintcpuCapacity + intcputotalCapacity
		finalintMemoryCapacity = finalintMemoryCapacity + inttotalMemoryCapacity

	}

	TotalNodeCapacity(totalNodeInfo)
	NamespaceUsagePercentageFinder(usageinfo, finalintcpuCapacity, finalintMemoryCapacity)
	// fmt.Println("****************************************",totalNodeInfo)
}

func TotalNodeCapacity(nodecapacityforcalc []string) {
	fmt.Println("capacity of all nodes", nodecapacityforcalc)

}

func NamespaceUsagePercentageFinder(usageinfo []string, finalintcpuCapacity int, finalintMemoryCapacity int) {
	fmt.Println("from namespace clusterCpuCapacity", finalintcpuCapacity, "clusterMemoryCapacity", finalintMemoryCapacity)
	NamespceCpu := make(map[string][]int)
	NamespceMemory := make(map[string][]int)
	KindCpu := make(map[string][]int)
	KindMemory := make(map[string][]int)
	DeploymentNamecpu := make(map[string][]int)
	DeploymentNameMemory := make(map[string][]int)
	var DeploymentNamespace = make(map[string]string)

	podRangeValue := len(usageinfo) - 1
	for a, b, c, l, m, n := 0, 1, 2, 3, 4, 5; a <= (podRangeValue - 6); a, b, c, l, m, n = a+6, b+6, c+6, l+6, m+6, n+6 {
		podName := usageinfo[a]
		fmt.Sprintln(podName)
		podCpuUsage := usageinfo[b]
		podMemoryUsage := usageinfo[c]
		podNamespace := usageinfo[l]
		podDeploymentName := usageinfo[m]
		podDeploymentKind := usageinfo[n]
		intpodCpuUsage, err := strconv.Atoi(podCpuUsage)
		intpodMemoryUsage, err := strconv.Atoi(podMemoryUsage)
		NamespceCpu[podNamespace] = append(NamespceCpu[podNamespace], intpodCpuUsage)
		NamespceMemory[podNamespace] = append(NamespceMemory[podNamespace], intpodMemoryUsage)
		KindCpu[podDeploymentKind] = append(KindCpu[podDeploymentKind], intpodCpuUsage)
		KindMemory[podDeploymentKind] = append(KindMemory[podDeploymentKind], intpodMemoryUsage)
		DeploymentNamecpu[podDeploymentName] = append(DeploymentNamecpu[podDeploymentName], intpodCpuUsage)
		DeploymentNameMemory[podDeploymentName] = append(DeploymentNameMemory[podDeploymentName], intpodMemoryUsage)
		value, ok := DeploymentNamespace[podDeploymentName+"_randomstringtoavoidconflict_"+podNamespace]
		if ok {
			fmt.Sprintln(value)
		} else {
			elasticdatatype.DeploymentaNameAsDataInserter(podDeploymentName+"_randomstringtoavoidconflict_"+podNamespace, podNamespace)
		}
		DeploymentNamespace[podDeploymentName+"_randomstringtoavoidconflict_"+podNamespace] = podNamespace

		if err != nil {
			return
		}
	}
	var namespaceCpuData []string
	var namespaceMemoryData []string

	var kindCpuData []string
	var kindMemoryData []string

	var deploymentCpuData []string
	var deploymentMemoryData []string

	for namespace, val := range NamespceCpu {
		NamespceCpuPer := (float64(sum(val)) / float64(finalintcpuCapacity)) * 100
		stringNamespceCpuPer := strconv.Itoa(int(NamespceCpuPer))
		fmt.Sprintln(stringNamespceCpuPer)
		Stringnamespacecpuval := strconv.Itoa(int(float64(sum(val))))
		namespaceCpuData = append(namespaceCpuData, namespace, Stringnamespacecpuval)
		elastic.CpuValueInserter(namespace, int(float64(sum(val))))

	}

	for memnamespace, memval := range NamespceMemory {
		NamespceMemoryPer := (float64(sum(memval))) / float64(finalintMemoryCapacity) * 100
		stringNamespceMemoryPer := strconv.Itoa(int(NamespceMemoryPer))
		fmt.Sprintln(stringNamespceMemoryPer)
		Stringmemnamespace := strconv.Itoa(int(float64(sum(memval))))
		namespaceMemoryData = append(namespaceMemoryData, memnamespace, Stringmemnamespace)
		elastic.NamespaceMemoryValueInserter(memnamespace, int(float64(sum(memval))))
		elasticdatatype.NamespaceaNameAsDataInserter(memnamespace, memnamespace)

	}

	for kind, kindcpuval := range KindCpu {

		KindCpuPer := (float64(sum(kindcpuval))) / float64(finalintcpuCapacity) * 100
		stringKindCpuPer := strconv.Itoa(int(KindCpuPer))
		fmt.Sprintln(stringKindCpuPer)
		Stringkindcpuval := strconv.Itoa(int(float64(sum(kindcpuval))))
		kindCpuData = append(kindCpuData, kind, Stringkindcpuval)
		elastic.KindcpuValueInserter(kind, int(float64(sum(kindcpuval))))
		elasticdatatype.KindaNameAsDataInserter(kind, kind)

	}

	for kindCpuName, kindmemval := range KindMemory {
		KindMemoryPer := (float64(sum(kindmemval))) / float64(finalintMemoryCapacity) * 100
		stringKindMemoryPer := strconv.Itoa(int(KindMemoryPer))
		fmt.Sprintln(stringKindMemoryPer)
		Stringkindmemval := strconv.Itoa(int(float64(sum(kindmemval))))
		kindMemoryData = append(kindMemoryData, kindCpuName, Stringkindmemval)
		elastic.KindMemoryValueInserter(kindCpuName, int(float64(sum(kindmemval))))

	}

	for DeployName, Deploycpuval := range DeploymentNamecpu {
		DeploymentNamecpuPer := (float64(sum(Deploycpuval))) / float64(finalintcpuCapacity) * 100
		stringDeploymentNamecpuPer := strconv.Itoa(int(DeploymentNamecpuPer))
		fmt.Sprintln(stringDeploymentNamecpuPer)
		StringDeploycpuval := strconv.Itoa(int(float64(sum(Deploycpuval))))
		deploymentCpuData = append(deploymentCpuData, DeployName, StringDeploycpuval)
		elastic.DeploymentCpuValueInserter(DeployName, int(float64(sum(Deploycpuval))))

	}

	for DeployMemoryName, Deploymemval := range DeploymentNameMemory {

		DeploymentNameMemoryPer := (float64(sum(Deploymemval))) / float64(finalintMemoryCapacity) * 100
		stringDeploymentNameMemoryPer := strconv.Itoa(int(DeploymentNameMemoryPer))
		fmt.Sprintln(stringDeploymentNameMemoryPer)
		StringDeploymemval := strconv.Itoa(int(float64(sum(Deploymemval))))
		deploymentMemoryData = append(deploymentMemoryData, DeployMemoryName, StringDeploymemval)
		elastic.DeploymentMemoryValueInserter(DeployMemoryName, int(float64(sum(Deploymemval))))

	}

	fmt.Println("namespaceCpuData", namespaceCpuData)
	fmt.Println("namespaceMemoryData", namespaceMemoryData)

	fmt.Println("kindCpuData", kindCpuData)
	fmt.Println("kindMemoryData", kindMemoryData)

	fmt.Println("deploymentCpuData", deploymentCpuData)
	fmt.Println("deploymentMemoryData", deploymentMemoryData)

	return
}

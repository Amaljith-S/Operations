package typelister

import (
	"context"
	costkube "costkube/elastic"
	"fmt"
	"log"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

func KindLister() []string {
	config, err := clientcmd.BuildConfigFromFlags("", "/home/amaljith/.kube/config")
	if err != nil {
		panic(err.Error())
	}
	

	// create the kubeClient
	kubeClient, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	metricsClientset, err := metricsv.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := kubeClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	var usageinfo []string
	//var usageinfowithoutref []string
	

	for _, pod := range pods.Items {

		if len(pod.OwnerReferences) == 0 {
			
			podMetricsListgettest, err := metricsClientset.MetricsV1beta1().PodMetricses(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
			if err != nil {
				log.Printf(err.Error())
			}

			podWithoutOwner := pod.Name
			ownerNameWithoutOwner := "pod"
			ownerKindWithoutOwner := "pod"

			var cpuWithoutOwner int64    
			var memoryWithoutOwner int64 

			podContainers1 := podMetricsListgettest.Containers
			for _, container1 := range podContainers1 {
				cpuQuantity1 := container1.Usage.Cpu().MilliValue()
				memQuantity1High, ok1 := container1.Usage.Memory().AsInt64()
				memQuantity1 := memQuantity1High * float64(0.00838861)
				fmt.Println("###################################################################",memQuantity1)

				cpuWithoutOwner += cpuQuantity1
				memoryWithoutOwner += memQuantity1
				if !ok1 {
					fmt.Sprintln("hi")
				}
				stringCpuWithOutOwner := strconv.Itoa(int(cpuWithoutOwner))       
				stringmemoryWithOutOwner := strconv.Itoa(int(memoryWithoutOwner)) 
				usageinfo = append(usageinfo, podWithoutOwner, stringCpuWithOutOwner, stringmemoryWithOutOwner, pod.Namespace, ownerNameWithoutOwner, ownerKindWithoutOwner)

			}
			
			fmt.Println(podWithoutOwner, cpuWithoutOwner, memoryWithoutOwner, pod.Namespace, ownerNameWithoutOwner, ownerKindWithoutOwner)
			costkube.ElasticInsert(podWithoutOwner, cpuWithoutOwner, memoryWithoutOwner, pod.Namespace, ownerNameWithoutOwner, ownerKindWithoutOwner)

			// fmt.Println("pod.Namespace *****************", usageinfo)

			continue
		}
		

		var ownerName, ownerKind string

		// fmt.Println("pod.Namespace &&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&", pod.Namespace)

		switch pod.OwnerReferences[0].Kind {
		case "ReplicaSet":
			replica, repErr := kubeClient.AppsV1().ReplicaSets(pod.Namespace).Get(context.TODO(), pod.OwnerReferences[0].Name, metav1.GetOptions{})

			if repErr != nil {
				// fmt.Println("Zero replicas are available")
				panic(repErr.Error())
			}

			ownerName = replica.OwnerReferences[0].Name
			// fmt.Println(ownerName)
			ownerKind = "Deployment"
		case "DaemonSet", "StatefulSet", "Node":
			ownerName = pod.OwnerReferences[0].Name
			// fmt.Println(ownerName)
			ownerKind = pod.OwnerReferences[0].Kind

		default:
			fmt.Printf("Could not find resource manager for type %s\n", pod.OwnerReferences[0].Kind)
			continue
		}
		
		podName := pod.Name
		podMetricsListget, err := metricsClientset.MetricsV1beta1().PodMetricses(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
		}
		var cpuWithOwner int64    // change need to float64
		var memoryWithOwner int64 //change need to float64

		podContainers := podMetricsListget.Containers
		for _, container := range podContainers {
			cpuQuantity := container.Usage.Cpu().MilliValue()
			memQuantityHigh, ok := container.Usage.Memory().AsInt64() // change AsInt64 to float64

			memQuantity := memQuantityHigh / (1024 * 1024)

			cpuWithOwner += cpuQuantity
			memoryWithOwner += memQuantity
			if !ok {
				// fmt.Sprintln("hi")
			}

		}

		costkube.ElasticInsert(podName, cpuWithOwner, memoryWithOwner, pod.Namespace, ownerName, ownerKind)
		


		stringCpuWithOwner := strconv.Itoa(int(cpuWithOwner))
		// fmt.Printf("stringCpuWithOwner: %s\n", reflect.TypeOf(stringCpuWithOwner))
		stringmemoryWithOwner := strconv.Itoa(int(memoryWithOwner))
		usageinfo = append(usageinfo, podName, stringCpuWithOwner, stringmemoryWithOwner, pod.Namespace, ownerName, ownerKind)

	}

	return usageinfo
}

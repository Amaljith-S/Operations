package diskusagefinder
 
import (
	"fmt"
	"context"
	"encoding/json"
	"time"
	"reflect"
	// "strconv"
	// "encoding/binary"

	// "strings"
	elastic "costkube/elastic"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)
 
 
type DiskUsage struct {
	Node struct {
		NodeName         string `json:"nodeName"`
		SystemContainers []struct {
			Name      string    `json:"name"`
			StartTime time.Time `json:"startTime"`
			CPU       struct {
				Time                 time.Time `json:"time"`
				UsageNanoCores       int       `json:"usageNanoCores"`
				UsageCoreNanoSeconds int64     `json:"usageCoreNanoSeconds"`
			} `json:"cpu"`
			Memory struct {
				Time            time.Time `json:"time"`
				AvailableBytes  int64     `json:"availableBytes"`
				UsageBytes      int64     `json:"usageBytes"`
				WorkingSetBytes int64     `json:"workingSetBytes"`
				RssBytes        int64     `json:"rssBytes"`
				PageFaults      int       `json:"pageFaults"`
				MajorPageFaults int       `json:"majorPageFaults"`
			} `json:"memory,omitempty"`
			Memory0 struct {
				Time            time.Time `json:"time"`
				UsageBytes      int       `json:"usageBytes"`
				WorkingSetBytes int       `json:"workingSetBytes"`
				RssBytes        int       `json:"rssBytes"`
				PageFaults      int       `json:"pageFaults"`
				MajorPageFaults int       `json:"majorPageFaults"`
			} `json:"memory,omitempty"`
			Memory1 struct {
				Time            time.Time `json:"time"`
				UsageBytes      int       `json:"usageBytes"`
				WorkingSetBytes int       `json:"workingSetBytes"`
				RssBytes        int       `json:"rssBytes"`
				PageFaults      int       `json:"pageFaults"`
				MajorPageFaults int       `json:"majorPageFaults"`
			} `json:"memory,omitempty"`
		} `json:"systemContainers"`
		StartTime time.Time `json:"startTime"`
		CPU       struct {
			Time                 time.Time `json:"time"`
			UsageNanoCores       int       `json:"usageNanoCores"`
			UsageCoreNanoSeconds int64     `json:"usageCoreNanoSeconds"`
		} `json:"cpu"`
		Memory struct {
			Time            time.Time `json:"time"`
			AvailableBytes  int64     `json:"availableBytes"`
			UsageBytes      int64     `json:"usageBytes"`
			WorkingSetBytes int64     `json:"workingSetBytes"`
			RssBytes        int64     `json:"rssBytes"`
			PageFaults      int       `json:"pageFaults"`
			MajorPageFaults int       `json:"majorPageFaults"`
		} `json:"memory"`
		Network struct {
			Time       time.Time `json:"time"`
			Name       string    `json:"name"`
			RxBytes    int       `json:"rxBytes"`
			RxErrors   int       `json:"rxErrors"`
			TxBytes    int       `json:"txBytes"`
			TxErrors   int       `json:"txErrors"`
			Interfaces []struct {
				Name     string `json:"name"`
				RxBytes  int    `json:"rxBytes"`
				RxErrors int    `json:"rxErrors"`
				TxBytes  int    `json:"txBytes"`
				TxErrors int    `json:"txErrors"`
			} `json:"interfaces"`
		} `json:"network"`
		Fs struct {
			Time           time.Time `json:"time"`
			AvailableBytes int64     `json:"availableBytes"`
			CapacityBytes  int64     `json:"capacityBytes"`
			UsedBytes      int64     `json:"usedBytes"`
			InodesFree     int       `json:"inodesFree"`
			Inodes         int       `json:"inodes"`
			InodesUsed     int       `json:"inodesUsed"`
		} `json:"fs"`
		Runtime struct {
			ImageFs struct {
				Time           time.Time `json:"time"`
				AvailableBytes int64     `json:"availableBytes"`
				CapacityBytes  int64     `json:"capacityBytes"`
				UsedBytes      int64     `json:"usedBytes"`
				InodesFree     int       `json:"inodesFree"`
				Inodes         int       `json:"inodes"`
				InodesUsed     int       `json:"inodesUsed"`
			} `json:"imageFs"`
		} `json:"runtime"`
		Rlimit struct {
			Time    time.Time `json:"time"`
			Maxpid  int       `json:"maxpid"`
			Curproc int       `json:"curproc"`
		} `json:"rlimit"`
	} `json:"node"`
	Pods []struct {
		PodRef struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
			UID       string `json:"uid"`
		} `json:"podRef"`
		StartTime  time.Time `json:"startTime"`
		Containers []struct {
			Name      string    `json:"name"`
			StartTime time.Time `json:"startTime"`
			CPU       struct {
				Time                 time.Time `json:"time"`
				UsageNanoCores       int       `json:"usageNanoCores"`
				UsageCoreNanoSeconds int       `json:"usageCoreNanoSeconds"`
			} `json:"cpu"`
			Memory struct {
				Time            time.Time `json:"time"`
				UsageBytes      int       `json:"usageBytes"`
				WorkingSetBytes int       `json:"workingSetBytes"`
				RssBytes        int       `json:"rssBytes"`
				PageFaults      int       `json:"pageFaults"`
				MajorPageFaults int       `json:"majorPageFaults"`
			} `json:"memory"`
			Rootfs struct {
				Time           time.Time `json:"time"`
				AvailableBytes int64     `json:"availableBytes"`
				CapacityBytes  int64     `json:"capacityBytes"`
				UsedBytes      int       `json:"usedBytes"`
				InodesFree     int       `json:"inodesFree"`
				Inodes         int       `json:"inodes"`
				InodesUsed     int       `json:"inodesUsed"`
			} `json:"rootfs"`
			Logs struct {
				Time           time.Time `json:"time"`
				AvailableBytes int64     `json:"availableBytes"`
				CapacityBytes  int64     `json:"capacityBytes"`
				UsedBytes      int       `json:"usedBytes"`
				InodesFree     int       `json:"inodesFree"`
				Inodes         int       `json:"inodes"`
				InodesUsed     int       `json:"inodesUsed"`
			} `json:"logs"`
		} `json:"containers"`
		CPU struct {
			Time                 time.Time `json:"time"`
			UsageNanoCores       int       `json:"usageNanoCores"`
			UsageCoreNanoSeconds int       `json:"usageCoreNanoSeconds"`
		} `json:"cpu"`
		Memory struct {
			Time            time.Time `json:"time"`
			UsageBytes      int       `json:"usageBytes"`
			WorkingSetBytes int       `json:"workingSetBytes"`
			RssBytes        int       `json:"rssBytes"`
			PageFaults      int       `json:"pageFaults"`
			MajorPageFaults int       `json:"majorPageFaults"`
		} `json:"memory,omitempty"`
		Network struct {
			Time       time.Time `json:"time"`
			Name       string    `json:"name"`
			RxBytes    int       `json:"rxBytes"`
			RxErrors   int       `json:"rxErrors"`
			TxBytes    int       `json:"txBytes"`
			TxErrors   int       `json:"txErrors"`
			Interfaces []struct {
				Name     string `json:"name"`
				RxBytes  int    `json:"rxBytes"`
				RxErrors int    `json:"rxErrors"`
				TxBytes  int    `json:"txBytes"`
				TxErrors int    `json:"txErrors"`
			} `json:"interfaces"`
		} `json:"network"`
		Volume []struct {
			Time           time.Time `json:"time"`
			AvailableBytes int64     `json:"availableBytes"`
			CapacityBytes  int64     `json:"capacityBytes"`
			UsedBytes      int       `json:"usedBytes"`
			InodesFree     int       `json:"inodesFree"`
			Inodes         int       `json:"inodes"`
			InodesUsed     int       `json:"inodesUsed"`
			Name           string    `json:"name"`
		} `json:"volume,omitempty"`
		EphemeralStorage struct {
			Time           time.Time `json:"time"`
			AvailableBytes int64     `json:"availableBytes"`
			CapacityBytes  int64     `json:"capacityBytes"`
			UsedBytes      int       `json:"usedBytes"`
			InodesFree     int       `json:"inodesFree"`
			Inodes         int       `json:"inodes"`
			InodesUsed     int       `json:"inodesUsed"`
		} `json:"ephemeral-storage"`
		ProcessStats struct {
			ProcessCount int `json:"process_count"`
		} `json:"process_stats"`
		Memory0 struct {
			Time            time.Time `json:"time"`
			AvailableBytes  int       `json:"availableBytes"`
			UsageBytes      int       `json:"usageBytes"`
			WorkingSetBytes int       `json:"workingSetBytes"`
			RssBytes        int       `json:"rssBytes"`
			PageFaults      int       `json:"pageFaults"`
			MajorPageFaults int       `json:"majorPageFaults"`
		} `json:"memory,omitempty"`
		Memory1 struct {
			Time            time.Time `json:"time"`
			AvailableBytes  int       `json:"availableBytes"`
			UsageBytes      int       `json:"usageBytes"`
			WorkingSetBytes int       `json:"workingSetBytes"`
			RssBytes        int       `json:"rssBytes"`
			PageFaults      int       `json:"pageFaults"`
			MajorPageFaults int       `json:"majorPageFaults"`
		} `json:"memory,omitempty"`
		Memory2 struct {
			Time            time.Time `json:"time"`
			AvailableBytes  int       `json:"availableBytes"`
			UsageBytes      int       `json:"usageBytes"`
			WorkingSetBytes int       `json:"workingSetBytes"`
			RssBytes        int       `json:"rssBytes"`
			PageFaults      int       `json:"pageFaults"`
			MajorPageFaults int       `json:"majorPageFaults"`
		} `json:"memory,omitempty"`
		Memory3 struct {
			Time            time.Time `json:"time"`
			AvailableBytes  int       `json:"availableBytes"`
			UsageBytes      int       `json:"usageBytes"`
			WorkingSetBytes int       `json:"workingSetBytes"`
			RssBytes        int       `json:"rssBytes"`
			PageFaults      int       `json:"pageFaults"`
			MajorPageFaults int       `json:"majorPageFaults"`
		} `json:"memory,omitempty"`
	} `json:"pods"`
}

func NodeDiskdata(nodeName string) {
	var url string
	url= "api/v1/nodes/" + nodeName + "/proxy/stats/summary"
	fmt.Println("$$$$$$$$$$$$$$$$$$-- urls ", url)
	config, err := clientcmd.BuildConfigFromFlags("", "/home/amaljith/.kube/config")
	if err != nil {
		panic(err.Error())
	}
 
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	var nodes DiskUsage
	data, err := clientset.RESTClient().Get().AbsPath(url).DoRaw(context.TODO())
	fmt.Println(reflect.TypeOf(data))
 
	if err != nil {
		return
	}
 
	err = json.Unmarshal(data, &nodes)
	//  fmt.Println("data", nodes)
	// fmt.Println( nodes)
	
	// var diskUsage [] string
	var podName string

	var rootfsusage int
	var volumedata int
	for _, i:= range nodes.Pods{
		podName = i.PodRef.Name
		fmt.Println(podName)
		for _, j:= range i.Containers{
			rootfsusage = j.Rootfs.UsedBytes
			// rootfsUsageString = strconv.Itoa(rootfsusage)

			// fmt.Println("----------------------------------->>>>>",reflect.TypeOf(rootfsUsageString))
			if err != nil {
				return
			}
		// fmt.Println("data to insert------------------------------",diskUsage )
		}	
		for _, k:= range i.Volume{
			volumedata = k.UsedBytes
			// volumeUsageString = strconv.Itoa(volumedata)
			// fmt.Println("------------------------------",reflect.TypeOf(volumedata))


			if err != nil {
				return
			}
		}
		// diskUsage= append(diskUsage,podName,rootfsUsageString,volumeUsageString )

		elastic.NodeDiskUsageInserter(podName,rootfsusage,volumedata)
		// for _, i:= range diskusage {

			



		

		// fmt.Println("____________________________",diskUsage)


		// return volumedata

		// return rootfsusage
	}
	// for _,k:= range nodes.Pods{
	// 	for _, j:= range k.Volume{
	// 		volumedata:= j.UsedBytes
	// 		fmt.Println("**************************",volumedata)

	// 		if err != nil {
	// 			return
	// 		}
	// 	}
	// 	return volumedata
	// }
	// elasticinsert.NodeInfoInserter(rootfsusage, volumedata)


}

type NodeName struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			UID               string    `json:"uid"`
			ResourceVersion   string    `json:"resourceVersion"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Labels            struct {
				BetaKubernetesIoArch                             string `json:"beta.kubernetes.io/arch"`
				BetaKubernetesIoOs                               string `json:"beta.kubernetes.io/os"`
				KubernetesIoArch                                 string `json:"kubernetes.io/arch"`
				KubernetesIoHostname                             string `json:"kubernetes.io/hostname"`
				KubernetesIoOs                                   string `json:"kubernetes.io/os"`
				MinikubeK8SIoCommit                              string `json:"minikube.k8s.io/commit"`
				MinikubeK8SIoName                                string `json:"minikube.k8s.io/name"`
				MinikubeK8SIoPrimary                             string `json:"minikube.k8s.io/primary"`
				MinikubeK8SIoUpdatedAt                           string `json:"minikube.k8s.io/updated_at"`
				MinikubeK8SIoVersion                             string `json:"minikube.k8s.io/version"`
				NodeRoleKubernetesIoControlPlane                 string `json:"node-role.kubernetes.io/control-plane"`
				NodeKubernetesIoExcludeFromExternalLoadBalancers string `json:"node.kubernetes.io/exclude-from-external-load-balancers"`
			} `json:"labels"`
			Annotations struct {
				KubeadmAlphaKubernetesIoCriSocket                string `json:"kubeadm.alpha.kubernetes.io/cri-socket"`
				NodeAlphaKubernetesIoTTL                         string `json:"node.alpha.kubernetes.io/ttl"`
				VolumesKubernetesIoControllerManagedAttachDetach string `json:"volumes.kubernetes.io/controller-managed-attach-detach"`
			} `json:"annotations"`
			ManagedFields []struct {
				Manager    string    `json:"manager"`
				Operation  string    `json:"operation"`
				APIVersion string    `json:"apiVersion"`
				Time       time.Time `json:"time"`
				FieldsType string    `json:"fieldsType"`
				FieldsV1   struct {
					FMetadata struct {
						FAnnotations struct {
							NAMING_FAILED struct {
							} `json:"."`
							FVolumesKubernetesIoControllerManagedAttachDetach struct {
							} `json:"f:volumes.kubernetes.io/controller-managed-attach-detach"`
						} `json:"f:annotations"`
						FLabels struct {
							NAMING_FAILED struct {
							} `json:"."`
							FBetaKubernetesIoArch struct {
							} `json:"f:beta.kubernetes.io/arch"`
							FBetaKubernetesIoOs struct {
							} `json:"f:beta.kubernetes.io/os"`
							FKubernetesIoArch struct {
							} `json:"f:kubernetes.io/arch"`
							FKubernetesIoHostname struct {
							} `json:"f:kubernetes.io/hostname"`
							FKubernetesIoOs struct {
							} `json:"f:kubernetes.io/os"`
						} `json:"f:labels"`
					} `json:"f:metadata"`
				} `json:"fieldsV1,omitempty"`
				FieldsV10 struct {
					FMetadata struct {
						FAnnotations struct {
							FNodeAlphaKubernetesIoTTL struct {
							} `json:"f:node.alpha.kubernetes.io/ttl"`
						} `json:"f:annotations"`
					} `json:"f:metadata"`
					FSpec struct {
						FPodCIDR struct {
						} `json:"f:podCIDR"`
						FPodCIDRs struct {
							NAMING_FAILED struct {
							} `json:"."`
							V102440024 struct {
							} `json:"v:"10.244.0.0/24""`
						} `json:"f:podCIDRs"`
					} `json:"f:spec"`
				} `json:"fieldsV1,omitempty"`
				FieldsV11 struct {
					FStatus struct {
						FAllocatable struct {
							FMemory struct {
							} `json:"f:memory"`
						} `json:"f:allocatable"`
						FCapacity struct {
							FMemory struct {
							} `json:"f:memory"`
						} `json:"f:capacity"`
						FConditions struct {
							KTypeDiskPressure struct {
								FLastHeartbeatTime struct {
								} `json:"f:lastHeartbeatTime"`
								FLastTransitionTime struct {
								} `json:"f:lastTransitionTime"`
								FMessage struct {
								} `json:"f:message"`
								FReason struct {
								} `json:"f:reason"`
								FStatus struct {
								} `json:"f:status"`
							} `json:"k:{"type":"DiskPressure"}"`
							KTypeMemoryPressure struct {
								FLastHeartbeatTime struct {
								} `json:"f:lastHeartbeatTime"`
								FLastTransitionTime struct {
								} `json:"f:lastTransitionTime"`
								FMessage struct {
								} `json:"f:message"`
								FReason struct {
								} `json:"f:reason"`
								FStatus struct {
								} `json:"f:status"`
							} `json:"k:{"type":"MemoryPressure"}"`
							KTypePIDPressure struct {
								FLastHeartbeatTime struct {
								} `json:"f:lastHeartbeatTime"`
								FLastTransitionTime struct {
								} `json:"f:lastTransitionTime"`
								FMessage struct {
								} `json:"f:message"`
								FReason struct {
								} `json:"f:reason"`
								FStatus struct {
								} `json:"f:status"`
							} `json:"k:{"type":"PIDPressure"}"`
							KTypeReady struct {
								FLastHeartbeatTime struct {
								} `json:"f:lastHeartbeatTime"`
								FLastTransitionTime struct {
								} `json:"f:lastTransitionTime"`
								FMessage struct {
								} `json:"f:message"`
								FReason struct {
								} `json:"f:reason"`
								FStatus struct {
								} `json:"f:status"`
							} `json:"k:{"type":"Ready"}"`
						} `json:"f:conditions"`
						FImages struct {
						} `json:"f:images"`
						FNodeInfo struct {
							FBootID struct {
							} `json:"f:bootID"`
						} `json:"f:nodeInfo"`
					} `json:"f:status"`
				} `json:"fieldsV1,omitempty"`
				Subresource string `json:"subresource,omitempty"`
			} `json:"managedFields"`
		} `json:"metadata"`
	} `json:"items"`
}

func GetNodes(){

	var nameOfNode [] string
	config, err := clientcmd.BuildConfigFromFlags("", "/home/amaljith/.kube/config")
	if err != nil {
		panic(err.Error())
	}
 
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}


	var nodes NodeName
	nodenames, err := clientset.RESTClient().Get().AbsPath("api/v1/nodes").DoRaw(context.TODO())
	//fmt.Println("+++++++++++++++++++++++++++++",name)
	if err != nil {
		return
	}
	err = json.Unmarshal(nodenames, &nodes)
	
	for _, i:= range nodes.Items {
		
		nameOfNode = append(nameOfNode, i.Metadata.Name)

		// fmt.Println("^^^^^^^^^^^^^",nameOfNode)	
		if err != nil {
			return
		}
		
	}
	for _, j:= range nameOfNode{
		fmt.Println("node names ",j)
		NodeDiskdata(j)
	}
	
	// return nameOfNode

}
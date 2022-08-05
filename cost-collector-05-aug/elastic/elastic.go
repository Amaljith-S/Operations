package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Declare a struct for Elasticsearch fields"ioutil"
type ElasticPodMetrics struct {
	PodName        string
	PodCpu         int64
	PodMemory      int64
	Namespace      string
	DeploymentName string
	KindName       string
	
}

// A function for marshaling structs to JSON string
func jsonStruct(doc ElasticPodMetrics) string {

	// Create struct instance of the Elasticsearch fields struct object
	docStruct := &ElasticPodMetrics{
		PodName:        doc.PodName,
		PodCpu:         doc.PodCpu,
		PodMemory:      doc.PodMemory,
		Namespace:      doc.Namespace,
		DeploymentName: doc.DeploymentName,
		KindName:       doc.KindName,
	}

	// Marshal the struct to JSON and check for errors
	b, err := json.Marshal(docStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(b)
}

func ElasticInsert(PodName string, PodCpu int64, PodMemory int64, Namespace string, DeploymentName string, KindName string) {

	
	// Allow for custom formatting of log output
	log.SetFlags(0)

	// Create a context object for the API calls
	ctx := context.Background()

	elasticHost := os.Getenv("es_host")
	elasticUser := os.Getenv("es_user")
	elasticPass := os.Getenv("es_pass")
	

	// Declare an Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
		Username: elasticUser,
		Password: elasticPass,
	}
	

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)
	

	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)


	}
	
	// Have the client instance return a response
	res, err := client.Info()	 

	// checking the status code to whether the client authentication successful or not
	
	if res.StatusCode == 200{
		// fmt.Println("")
	}else{
		fmt.Println("Connection failed")
		os.Exit(1)
	}
	// Deserialize the response into a map.
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	} else {
		fmt.Sprintf("client response:", res)
	}
	
	// fmt.Println("Connection successful")

	// Declare empty array for the document strings
	var docs1 []string
	t := time.Now()
	timeNow := (t.Format("20060102150405"))
	timeInt, err := strconv.Atoi(timeNow)
	rand.Seed(time.Now().UnixNano())
	esId := timeInt + rand.Intn(100)

	// Declare documents to be indexed using struct
	doc2 := ElasticPodMetrics{}
	doc2.PodName = PodName
	doc2.PodCpu = PodCpu
	doc2.PodMemory = PodMemory
	doc2.Namespace = Namespace
	doc2.DeploymentName = DeploymentName
	doc2.KindName = KindName

	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr1 := jsonStruct(doc2)
	// fmt.Println("docS1tr1 = ", reflect.TypeOf(docS1tr1))
	// fmt.Println(docS1tr1)

	// Check whether data is passing into the elastic search
	if(docS1tr1 == ""){
		fmt.Println("No data ")
		os.Exit(1)
	}else{
		// fmt.Println("Data insertion started")
	
	// Append the doc strings to an array
	docs1 = append(docs1, docS1tr1)
	// fmt.Println(docs1, "**********************************")
	elastic, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{elasticHost},
	})
	// creating indices inside elastic search
	create, err := elastic.Indices.Create("elastic_kind_details")
	if err != nil {
		fmt.Println(err)

	}
	if create.IsError() {
		fmt.Println(err)
	}

	createNs, err := elastic.Indices.Create("elastic_namespace_details")
	if err != nil {
		fmt.Println(err)

	}
	if createNs.IsError() {
		fmt.Println(err)
	}

	createDp, err := elastic.Indices.Create("elastic_deployment_details")
	if err != nil {
		fmt.Println(err)
	}
	if createDp.IsError() {
		fmt.Println(err)
	}
	createNo, err := elastic.Indices.Create("elastic_node_details")
	if err != nil {
		fmt.Println(err)
	}
	if createNo.IsError() {
		fmt.Println(err)
	}

	// Iterate the array of string documents
	for k, bod := range docs1 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "some_index2",
			DocumentID: strconv.Itoa(esId + 2),
			Body:       strings.NewReader(bod),
			Refresh:    "false",
		}

		// fmt.Println(req)

		// Return an API response object from request
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), k+1)
		} else {

			// Deserialize the response into a map.
			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				fmt.Sprintf("\nIndexRequest() RESPONSE:")

			}
		}
	}

	}
}

// Node Usage Insert

// Node struct

type KubeNodeMetrics struct {
	NodeName   string
	NodeCpu    int
	NodeMemory int
}

// A function for marshaling structs to JSON string
func jsonStructForNode(doc4 KubeNodeMetrics) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDocStruct := &KubeNodeMetrics{
		NodeName:   doc4.NodeName,
		NodeCpu:    doc4.NodeCpu,
		NodeMemory: doc4.NodeMemory,
	}

	// Marshal the struct to JSON and check for errors
	c, err := json.Marshal(nodeDocStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(c)
}

func NodeInfoInserter(Nodename string, Nodecpu int, Nodememory int) {

	// Allow for custom formatting of log output
	log.SetFlags(0)

	// Create a context object for the API calls
	ctx := context.Background()

	elasticHost := os.Getenv("es_host")
	elasticUser := os.Getenv("es_user")
	elasticPass := os.Getenv("es_pass")

	// Declare an Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
		Username: elasticUser,
		Password: elasticPass,
	}

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)
	}

	// Have the client instance return a response
	res, err := client.Info()

	// Deserialize the response into a map.
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	} else {
		fmt.Sprintf("client response:", res)
	}

	// Declare empty array for the document strings
	var docs3 []string

	t := time.Now()
	timeNowDate := (t.Format("20060102150405"))
	timeIntnode, err := strconv.Atoi(timeNowDate)
	rand.Seed(time.Now().UnixNano())
	esIdNode := timeIntnode + rand.Intn(100)

	// Declare documents to be indexed using struct
	docs4 := KubeNodeMetrics{}
	docs4.NodeName = Nodename
	docs4.NodeCpu = Nodecpu
	docs4.NodeMemory = Nodememory
	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr2 := jsonStructForNode(docs4)

	// Append the doc strings to an array
	docs3 = append(docs3, docS1tr2)

	// Iterate the array of string documents
	for kl, NodeBod := range docs3 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "some_index2",
			DocumentID: strconv.Itoa(esIdNode + 1),
			Body:       strings.NewReader(NodeBod),
			Refresh:    "false",
		}

		// Return an API response object from request
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), kl+1)
		} else {

			// Deserialize the response into a map.
			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				fmt.Sprintf("\nIndexRequest() RESPONSE:")

			}
		}

	}
}

// Node Percentage Insert

// Node struct

type KubeNodeUsagePercentageMetrics struct {
	KubeNodeName      string
	NodeCpuPercent    int
	NodeMemoryPercent int
}

// A function for marshaling structs to JSON string
func jsonStructForNodePercentage(doc5 KubeNodeUsagePercentageMetrics) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDocStruct := &KubeNodeUsagePercentageMetrics{
		KubeNodeName:      doc5.KubeNodeName,
		NodeCpuPercent:    doc5.NodeCpuPercent,
		NodeMemoryPercent: doc5.NodeMemoryPercent,
	}

	// Marshal the struct to JSON and check for errors
	d, err := json.Marshal(nodeDocStruct)
	// fmt.Println(d, "**********************************************")
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(d)
}

func NodePercentageInserter(KubeNodeName string, NodeCpuPercent int, NodeMemoryPercent int) {

	// Allow for custom formatting of log output
	log.SetFlags(0)

	// Create a context object for the API calls
	ctx := context.Background()

	elasticHost := os.Getenv("es_host")
	elasticUser := os.Getenv("es_user")
	elasticPass := os.Getenv("es_pass")

	// Declare an Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
		Username: elasticUser,
		Password: elasticPass,
	}

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)
	}

	// Have the client instance return a response
	res, err := client.Info()

	// Deserialize the response into a map.
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	} else {
		fmt.Sprintf("client response:", res)
	}

	// Declare empty array for the document strings
	var docs4 []string

	t := time.Now()
	timeNowDate := (t.Format("20060102150405"))
	timeIntnode, err := strconv.Atoi(timeNowDate)
	rand.Seed(time.Now().UnixNano())
	esIdNode := timeIntnode + rand.Intn(100)

	// Declare documents to be indexed using struct
	doc5 := KubeNodeUsagePercentageMetrics{}
	doc5.KubeNodeName = KubeNodeName
	doc5.NodeCpuPercent = NodeCpuPercent
	doc5.NodeMemoryPercent = NodeMemoryPercent
	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr3 := jsonStructForNodePercentage(doc5)

	// Append the doc strings to an array
	docs4 = append(docs4, docS1tr3)

	// Iterate the array of string documents
	for kl, NodeBod := range docs4 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "some_index2",
			DocumentID: strconv.Itoa(esIdNode + 1),
			Body:       strings.NewReader(NodeBod),
			Refresh:    "false",
		}

		// Return an API response object from request
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), kl+1)
		} else {

			// Deserialize the response into a map.
			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				fmt.Sprintf("\nIndexRequest() RESPONSE:")

			}
		}

	}
}

/// namespace values inserter

type NamespaceCpuUsagePercentageMetrics struct {
	NamespaceName     string
	NamespaceCpuValue int
}

// A function for marshaling structs to JSON string
func jsonStructForCpuValue(doc7 NamespaceCpuUsagePercentageMetrics) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDocStruct := &NamespaceCpuUsagePercentageMetrics{
		NamespaceName:     doc7.NamespaceName,
		NamespaceCpuValue: doc7.NamespaceCpuValue,
	}

	// Marshal the struct to JSON and check for errors
	d, err := json.Marshal(nodeDocStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(d)
}

func CpuValueInserter(NamespaceName string, NamespaceCpuValue int) {

	// Allow for custom formatting of log output
	log.SetFlags(0)

	// Create a context object for the API calls
	ctx := context.Background()

	elasticHost := os.Getenv("es_host")
	elasticUser := os.Getenv("es_user")
	elasticPass := os.Getenv("es_pass")

	// Declare an Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
		Username: elasticUser,
		Password: elasticPass,
	}

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)
	}

	// Have the client instance return a response
	res, err := client.Info()

	// Deserialize the response into a map.
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	} else {
		fmt.Sprintf("client response:", res)
	}

	// Declare empty array for the document strings
	var docs7 []string

	t := time.Now()
	timeNowDate := (t.Format("20060102150405"))
	timeIntnode, err := strconv.Atoi(timeNowDate)
	rand.Seed(time.Now().UnixNano())
	esIdCPU := timeIntnode + rand.Intn(100)

	// Declare documents to be indexed using struct
	doc7 := NamespaceCpuUsagePercentageMetrics{}
	doc7.NamespaceName = NamespaceName
	doc7.NamespaceCpuValue = NamespaceCpuValue
	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr3 := jsonStructForCpuValue(doc7)

	// Append the doc strings to an array
	docs7 = append(docs7, docS1tr3)

	// Iterate the array of string documents
	for kl, NodeBod := range docs7 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "some_index2",
			DocumentID: strconv.Itoa(esIdCPU + 1),
			Body:       strings.NewReader(NodeBod),
			Refresh:    "false",
		}

		// Return an API response object from request
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), kl+1)
		} else {

			// Deserialize the response into a map.
			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				fmt.Sprintf("\nIndexRequest() RESPONSE:")

			}
		}

	}
}

/// namespace MEMORY values inserter

type NamespaceMemoryValueMetrics struct {
	MemoryNamespaceName  string
	NamespaceMemoryValue int
}

// A function for marshaling structs to JSON string
func jsonStructForNamespaceMemoryValue(doc6 NamespaceMemoryValueMetrics) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDocStruct := &NamespaceMemoryValueMetrics{
		MemoryNamespaceName:  doc6.MemoryNamespaceName,
		NamespaceMemoryValue: doc6.NamespaceMemoryValue,
	}

	// Marshal the struct to JSON and check for errors
	d, err := json.Marshal(nodeDocStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(d)
}

func NamespaceMemoryValueInserter(MemoryNamespaceName string, NamespaceMemoryValue int) {

	// Allow for custom formatting of log output
	log.SetFlags(0)

	// Create a context object for the API calls
	ctx := context.Background()

	elasticHost := os.Getenv("es_host")
	elasticUser := os.Getenv("es_user")
	elasticPass := os.Getenv("es_pass")

	// Declare an Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
		Username: elasticUser,
		Password: elasticPass,
	}
	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)
	}

	// Have the client instance return a response
	res, err := client.Info()

	// Deserialize the response into a map.
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	} else {
		fmt.Sprintf("client response:", res)
	}

	// Declare empty array for the document strings
	var docs8 []string

	t := time.Now()
	timeNowDate := (t.Format("20060102150405"))
	timeIntnode, err := strconv.Atoi(timeNowDate)
	rand.Seed(time.Now().UnixNano())
	esIdMEM := timeIntnode + rand.Intn(100)

	// Declare documents to be indexed using struct
	doc6 := NamespaceMemoryValueMetrics{}
	doc6.MemoryNamespaceName = MemoryNamespaceName
	doc6.NamespaceMemoryValue = NamespaceMemoryValue
	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr4 := jsonStructForNamespaceMemoryValue(doc6)

	// Append the doc strings to an array
	docs8 = append(docs8, docS1tr4)

	// Iterate the array of string documents
	for kl, NodeBod := range docs8 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "some_index2",
			DocumentID: strconv.Itoa(esIdMEM + 1),
			Body:       strings.NewReader(NodeBod),
			Refresh:    "false",
		}

		// Return an API response object from request
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), kl+1)
		} else {

			// Deserialize the response into a map.
			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				fmt.Sprintf("\nIndexRequest() RESPONSE:")

			}
		}

	}
}

/// Kind CPU insert

type KindcpuValueMetrics struct {
	KindCpuName  string
	KindcpuValue int
}

// A function for marshaling structs to JSON string
func jsonStructForKindcpuValue(doc7 KindcpuValueMetrics) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDocStruct := &KindcpuValueMetrics{
		KindCpuName:  doc7.KindCpuName,
		KindcpuValue: doc7.KindcpuValue,
	}

	// Marshal the struct to JSON and check for errors
	d, err := json.Marshal(nodeDocStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(d)
}

func KindcpuValueInserter(KindCpuName string, KindcpuValue int) {

	// Allow for custom formatting of log output
	log.SetFlags(0)

	// Create a context object for the API calls
	ctx := context.Background()

	elasticHost := os.Getenv("es_host")
	elasticUser := os.Getenv("es_user")
	elasticPass := os.Getenv("es_pass")

	// Declare an Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
		Username: elasticUser,
		Password: elasticPass,
	}

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)
	}

	// Have the client instance return a response
	res, err := client.Info()

	// Deserialize the response into a map.
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	} else {
		fmt.Sprintf("client response:", res)
	}

	// Declare empty array for the document strings
	var docs9 []string

	t := time.Now()
	timeNowDate := (t.Format("20060102150405"))
	timeIntnode, err := strconv.Atoi(timeNowDate)
	rand.Seed(time.Now().UnixNano())
	esIdCPU := timeIntnode + rand.Intn(100)

	// Declare documents to be ineladexed using struct
	doc7 := KindcpuValueMetrics{}
	doc7.KindCpuName = KindCpuName
	doc7.KindcpuValue = KindcpuValue
	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr5 := jsonStructForKindcpuValue(doc7)

	// Append the doc strings to an array
	docs9 = append(docs9, docS1tr5)

	// Iterate the array of string documents
	for kl, NodeBod := range docs9 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "some_index2",
			DocumentID: strconv.Itoa(esIdCPU + 1),
			Body:       strings.NewReader(NodeBod),
			Refresh:    "false",
		}

		// Return an API response object from request
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), kl+1)
		} else {

			// Deserialize the response into a map.
			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				fmt.Sprintf("\nIndexRequest() RESPONSE:")

			}
		}

	}
}

/// Kind Memory insert

type KindMemoryValueMetrics struct {
	KindMemoryName  string
	KindMemoryValue int
}

// A function for marshaling structs to JSON string
func jsonStructForKindMemoryValue(doc8 KindMemoryValueMetrics) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDocStruct := &KindMemoryValueMetrics{
		KindMemoryName:  doc8.KindMemoryName,
		KindMemoryValue: doc8.KindMemoryValue,
	}

	// Marshal the struct to JSON and check for errors
	d, err := json.Marshal(nodeDocStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(d)
}

func KindMemoryValueInserter(KindMemoryName string, KindMemoryValue int) {

	// Allow for custom formatting of log output
	log.SetFlags(0)

	// Create a context object for the API calls
	ctx := context.Background()

	elasticHost := os.Getenv("es_host")
	elasticUser := os.Getenv("es_user")
	elasticPass := os.Getenv("es_pass")

	// Declare an Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
		Username: elasticUser,
		Password: elasticPass,
	}

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)
	}

	// Have the client instance return a response
	res, err := client.Info()

	// Deserialize the response into a map.
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	} else {
		fmt.Sprintf("client response:", res)
	}

	// Declare empty array for the document strings
	var docs10 []string

	t := time.Now()
	timeNowDate := (t.Format("20060102150405"))
	timeIntnode, err := strconv.Atoi(timeNowDate)
	rand.Seed(time.Now().UnixNano())
	esIdKINDMEM := timeIntnode + rand.Intn(100)

	// Declare documents to be indexed using struct
	doc8 := KindMemoryValueMetrics{}
	doc8.KindMemoryName = KindMemoryName
	doc8.KindMemoryValue = KindMemoryValue
	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr6 := jsonStructForKindMemoryValue(doc8)

	// Append the doc strings to an array
	docs10 = append(docs10, docS1tr6)

	// Iterate the array of string documents
	for kl, NodeBod := range docs10 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "some_index2",
			DocumentID: strconv.Itoa(esIdKINDMEM + 1),
			Body:       strings.NewReader(NodeBod),
			Refresh:    "false",
		}

		// Return an API response object from request
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), kl+1)
		} else {

			// Deserialize the response into a map.
			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				fmt.Sprintf("\nIndexRequest() RESPONSE:")

			}
		}

	}
}

/// Kind Memory insert

type DeploymentCpuValueMetrics struct {
	DeploymentCpuName  string
	DeploymentCpuValue int
}

// A function for marshaling structs to JSON string
func jsonStructForDeploymentCpuValue(doc9 DeploymentCpuValueMetrics) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDocStruct := &DeploymentCpuValueMetrics{
		DeploymentCpuName:  doc9.DeploymentCpuName,
		DeploymentCpuValue: doc9.DeploymentCpuValue,
	}

	// Marshal the struct to JSON and check for errors
	d, err := json.Marshal(nodeDocStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(d)
}

func DeploymentCpuValueInserter(DeploymentCpuName string, DeploymentCpuValue int) {

	// Allow for custom formatting of log output
	log.SetFlags(0)

	// Create a context object for the API calls
	ctx := context.Background()

	elasticHost := os.Getenv("es_host")
	elasticUser := os.Getenv("es_user")
	elasticPass := os.Getenv("es_pass")

	// Declare an Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
		Username: elasticUser,
		Password: elasticPass,
	}
	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)
	}

	// Have the client instance return a response
	res, err := client.Info()

	// Deserialize the response into a map.
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	} else {
		fmt.Sprintf("client response:", res)
	}

	// Declare empty array for the document strings
	var docs11 []string

	t := time.Now()
	timeNowDate := (t.Format("20060102150405"))
	timeIntnode, err := strconv.Atoi(timeNowDate)
	rand.Seed(time.Now().UnixNano())
	esIdDCPU := timeIntnode + rand.Intn(100)

	// Declare documents to be indexed using struct
	doc9 := DeploymentCpuValueMetrics{}
	doc9.DeploymentCpuName = DeploymentCpuName
	doc9.DeploymentCpuValue = DeploymentCpuValue
	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr7 := jsonStructForDeploymentCpuValue(doc9)

	// Append the doc strings to an array
	docs11 = append(docs11, docS1tr7)

	// Iterate the array of string documents
	for kl, NodeBod := range docs11 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "some_index2",
			DocumentID: strconv.Itoa(esIdDCPU + 1),
			Body:       strings.NewReader(NodeBod),
			Refresh:    "false",
		}

		// Return an API response object from request
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), kl+1)
		} else {

			// Deserialize the response into a map.
			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				fmt.Sprintf("\nIndexRequest() RESPONSE:")

			}
		}

	}
}

/// Deployment Memory

type DeploymentMemoryValueMetrics struct {
	DeploymentMemoryName  string
	DeploymentMemoryValue int
}

// A function for marshaling structs to JSON string
func jsonStructForDeploymentMemoryValue(doc10 DeploymentMemoryValueMetrics) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDocStruct := &DeploymentMemoryValueMetrics{
		DeploymentMemoryName:  doc10.DeploymentMemoryName,
		DeploymentMemoryValue: doc10.DeploymentMemoryValue,
	}

	// Marshal the struct to JSON and check for errors
	d, err := json.Marshal(nodeDocStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(d)
}

func DeploymentMemoryValueInserter(DeploymentMemoryName string, DeploymentMemoryValue int) {

	// Allow for custom formatting of log output
	log.SetFlags(0)

	// Create a context object for the API calls
	ctx := context.Background()

	// Declare an Elasticsearch configuration
	elasticHost := os.Getenv("es_host")
	elasticUser := os.Getenv("es_user")
	elasticPass := os.Getenv("es_pass")

	// Declare an Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
		Username: elasticUser,
		Password: elasticPass,
	}

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)
	}

	// Have the client instance return a response
	res, err := client.Info()

	// Deserialize the response into a map.
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	} else {
		fmt.Sprintf("client response:", res)
	}

	// Declare empty array for the document strings
	var docs12 []string

	t := time.Now()
	timeNowDate := (t.Format("20060102150405"))
	timeIntnode, err := strconv.Atoi(timeNowDate)
	rand.Seed(time.Now().UnixNano())
	esIdDMEM := timeIntnode + rand.Intn(100)

	// Declare documents to be indexed using struct
	doc10 := DeploymentMemoryValueMetrics{}
	doc10.DeploymentMemoryName = DeploymentMemoryName
	doc10.DeploymentMemoryValue = DeploymentMemoryValue
	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr8 := jsonStructForDeploymentMemoryValue(doc10)

	// Append the doc strings to an array
	docs12 = append(docs12, docS1tr8)

	// Iterate the array of string documents
	for kl, NodeBod := range docs12 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "some_index2",
			DocumentID: strconv.Itoa(esIdDMEM + 1),
			Body:       strings.NewReader(NodeBod),
			Refresh:    "false",
		}

		// Return an API response object from request
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), kl+1)
		} else {

			// Deserialize the response into a map.
			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				fmt.Sprintf("\nIndexRequest() RESPONSE:")

			}
		}

	}
}


type NodeDiskUsageMetric struct {
	PodName  	string
	RoofFsUsage int
	VolumeUsage int
}

// Function to unmarshal the struct to a JSON string.

func jsonStructForNodeDiskUsage(doc31 NodeDiskUsageMetric) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDiskStruct :=&NodeDiskUsageMetric{
		PodName:        doc31.PodName,
		RoofFsUsage:    doc31.RoofFsUsage,
		VolumeUsage: 	doc31.VolumeUsage,
	}
	// Marshal the struct to JSON and check for errors
	d, err := json.Marshal(nodeDiskStruct)
	// fmt.Println(d, "**********************************************")
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(d)

}

func NodeDiskUsageInserter(PodName string, RoofFsUsage int, VolumeUsage int){

	// Allow for custom formatting of log output
	log.SetFlags(0)

	// Create a context object for the API calls
	ctx := context.Background()

	elasticHost := os.Getenv("es_host")
	elasticUser := os.Getenv("es_user")
	elasticPass := os.Getenv("es_pass")

	// Declare an Elasticsearch configuration
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
		Username: elasticUser,
		Password: elasticPass,
	}

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)
	}

	// Have the client instance return a response
	res, err := client.Info()

	// Deserialize the response into a map.
	if err != nil {
		log.Fatalf("client.Info() ERROR:", err)
	} else {
		fmt.Sprintf("client response:", res)
	}

	// Declare empty array for the document strings
	var docs30 []string

	t := time.Now()
	timeNowDate := (t.Format("20060102150405"))
	timeIntnode, err := strconv.Atoi(timeNowDate)
	rand.Seed(time.Now().UnixNano())
	esIdNode := timeIntnode + rand.Intn(100)

	doc31 := NodeDiskUsageMetric{}
	doc31.PodName = PodName
	doc31.RoofFsUsage = RoofFsUsage
	doc31.VolumeUsage = VolumeUsage

	// Marshal Elasticsearch document struct objects to JSON string
	
	docS1tr8 := jsonStructForNodeDiskUsage(doc31)

	// Append the doc strings to an array
	
	docs30 = append(docs30, docS1tr8)

	fmt.Println("docs30____________________________",docs30)

	// Iterate the array of string documents
	for kl, NodeBod := range docs30 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "some_index2",
			DocumentID: strconv.Itoa(esIdNode + 1),
			Body:       strings.NewReader(NodeBod),
			Refresh:    "false",
		}

		// Return an API response object from request
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), kl+1)
		} else {

			// Deserialize the response into a map.
			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				fmt.Sprintf("\nIndexRequest() RESPONSE:")

			}
		}


	}
}
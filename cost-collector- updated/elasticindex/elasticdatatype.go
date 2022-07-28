package elasticdatatype

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

/// Kind data values inserter

type KindNamedata struct {
	Kindname        string
	KindaNameAsData string
}

// A function for marshaling structs to JSON string
func jsonStructForKindaNameAsData(doc20 KindNamedata) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDocStruct := &KindNamedata{
		Kindname:        doc20.Kindname,
		KindaNameAsData: doc20.KindaNameAsData,
	}

	// Marshal the struct to JSON and check for errors
	d, err := json.Marshal(nodeDocStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(d)
}

func KindaNameAsDataInserter(Kindname string, KindaNameAsData string) {

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
	var docs21 []string

	// Declare documents to be indexed using struct
	doc20 := KindNamedata{}
	doc20.Kindname = Kindname
	doc20.KindaNameAsData = KindaNameAsData
	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr22 := jsonStructForKindaNameAsData(doc20)

	// Append the doc strings to an array
	docs21 = append(docs21, docS1tr22)

	// Iterate the array of string documents
	for kl, NodeBod := range docs21 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "elastic_kind_details",
			DocumentID: Kindname,
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

/// Namespace data Inserter

type NamespaceNamedata struct {
	Namespacename        string
	NamespaceaNameAsData string
}

// A function for marshaling structs to JSON string
func jsonStructForNamespaceaNameAsData(doc23 NamespaceNamedata) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDocStruct := &NamespaceNamedata{
		Namespacename:        doc23.Namespacename,
		NamespaceaNameAsData: doc23.NamespaceaNameAsData,
	}

	// Marshal the struct to JSON and check for errors
	d, err := json.Marshal(nodeDocStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(d)
}

func NamespaceaNameAsDataInserter(Namespacename string, NamespaceaNameAsData string) {

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
	var docs24 []string

	// Declare documents to be indexed using struct
	doc23 := NamespaceNamedata{}
	doc23.Namespacename = Namespacename
	doc23.NamespaceaNameAsData = NamespaceaNameAsData
	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr25 := jsonStructForNamespaceaNameAsData(doc23)

	// Append the doc strings to an array
	docs24 = append(docs24, docS1tr25)
	// fmt.Println("#######################################",docs24)

	// Iterate the array of string documents
	for kl, NodeBod := range docs24 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "elastic_namespace_details",
			DocumentID: Namespacename,
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

/// Deployment details Inserter

type DeploymentNamedata struct {
	Deploymentname        string
	DeploymentaNameAsData string
}

// A function for marshaling structs to JSON string
func jsonStructForDeploymentaNameAsData(doc26 DeploymentNamedata) string {

	// Create struct instance of the Elasticsearch fields struct object
	nodeDocStruct := &DeploymentNamedata{
		Deploymentname:        doc26.Deploymentname,
		DeploymentaNameAsData: doc26.DeploymentaNameAsData,
	}

	// Marshal the struct to JSON and check for errors
	d, err := json.Marshal(nodeDocStruct)
	if err != nil {
		fmt.Println("json.Marshal ERROR:", err)
		return string(err.Error())
	}
	return string(d)
}

func DeploymentaNameAsDataInserter(Deploymentname string, DeploymentaNameAsData string) {

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
	var docs27 []string

	// Declare documents to be indexed using struct
	doc26 := DeploymentNamedata{}
	doc26.Deploymentname = Deploymentname
	doc26.DeploymentaNameAsData = DeploymentaNameAsData
	// Marshal Elasticsearch document struct objects to JSON string
	docS1tr28 := jsonStructForDeploymentaNameAsData(doc26)

	// Append the doc strings to an array
	docs27 = append(docs27, docS1tr28)
	// fmt.Println("#############################################################",docs27)

	// Iterate the array of string documents
	for kl, NodeBod := range docs27 {

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "elastic_deployment_details",
			DocumentID: Deploymentname,
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

// /// Node data inserter //////////////////////////////////////////

// type NodeNamedata struct {
// 	Nodename        string
// 	NodeNameAsData 	string
// }

// // A function for marshaling structs to JSON string

// func jsonStructForNodeNameAsData(doc29 NodeNamedata) string {

// 	// Create struct instance of the Elasticsearch fields struct object


// 	node2DocStruct := &NodeNamedata{
// 		Nodename:        	doc29.Nodename,
// 		NodeNameAsData: 	doc29.NodeNameAsData,
// 	}

// // Marshal the struct to JSON and check for errors

// 	d, err := json.Marshal(node2DocStruct)
// 	if err != nil {
// 		fmt.Println("json.Marshal ERROR:", err)
// 		return string(err.Error())
// 	}

// 	return string(d)
// }



// func NodeNameAsDataInserter(Nodename string, NodeNameAsData string) {

// 	// Allow for custom formatting of log output
// 	log.SetFlags(0)

// 	// Create a context object for the API calls
// 	ctx := context.Background()
// 	elasticHost := os.Getenv("es_host")
// 	elasticUser := os.Getenv("es_user")
// 	elasticPass := os.Getenv("es_pass")

// 	// Declare an Elasticsearch configuration
// 	cfg := elasticsearch.Config{
// 		Addresses: []string{
// 			elasticHost,
// 		},
// 		Username: elasticUser,
// 		Password: elasticPass,
// 	}

// 	// Instantiate a new Elasticsearch client object instance
// 	client, err := elasticsearch.NewClient(cfg)

// 	if err != nil {
// 		fmt.Println("Elasticsearch connection error:", err)
// 	}

// 	// Have the client instance return a response
// 	res, err := client.Info()

// 	// Deserialize the response into a map.
// 	if err != nil {
// 		log.Fatalf("client.Info() ERROR:", err)
// 	} else {
// 		fmt.Sprintf("client response:", res)
// 	}

// 	// Declare empty array for the document strings
// 	var docs30 []string

// 	// Declare documents to be indexed using struct
// 	doc29 := NodeNamedata{}
// 	doc29.Nodename = Nodename
// 	doc29.NodeNameAsData = NodeNameAsData
// 	// Marshal Elasticsearch document struct objects to JSON string
// 	docS1tr31 := jsonStructForNodeNameAsData(doc29)

// 	// Append the doc strings to an array
// 	docs30 = append(docs30, docS1tr31)
// 	fmt.Println("#######################################",docs30)

// 	for kl, NodeBod := range docs30 {

// 		// Instantiate a request object
// 		req := esapi.IndexRequest{
// 			Index:      "elastic_node_details",
// 			DocumentID: Nodename,
// 			Body:       strings.NewReader(NodeBod),
// 			Refresh:    "false",
// 		}
// 		// Return an API response object from request
// 		res, err := req.Do(ctx, client)
// 		if err != nil {
// 			log.Fatalf("IndexRequest ERROR: %s", err)
// 		}
// 		defer res.Body.Close()

// 		if res.IsError() {
// 			log.Printf("%s ERROR indexing document ID=%d", res.Status(), kl+1)
// 		} else {

// 			// Deserialize the response into a map.
// 			var resMap map[string]interface{}
// 			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
// 				log.Printf("Error parsing the response body: %s", err)
// 			} else {
// 				fmt.Sprintf("\nIndexRequest() RESPONSE:")

// 			}
// 		}

// 	}
// }

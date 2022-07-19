package namespaceusage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/montanaflynn/stats"
)

type NamespaceJsonOutput struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore interface{}   `json:"max_score"`
		Hits     []interface{} `json:"hits"`
	} `json:"hits"`
	Aggregations struct {
		Langs struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"langs"`
	} `json:"aggregations"`
}

type InputCostData struct {
	CalcCPU    bool    `json:"Calc_Cpu"`
	CalcMemory bool    `json:"Calc_Memory"`
	SearchFor  string  `json:"Search_for"`
	SearchFrom string  `json:"Search_From"`
	SearchTil  string  `json:"Search_Til"`
	TimeZone   string  `json:"TimeZone"`
	CPUCost    float64 `json:"Cpu_Cost"`
	MemoryCost float64 `json:"Memory_cost"`
}

type json_inputvars struct {
	Todate   string
	FromDate string
	TermKey  string
	Term     string
	AvgField string
	TimeZone string
}

func InputCostReader() (coastInput InputCostData) {
	var data []byte
	data, _ = ioutil.ReadFile("input.json")

	_ = json.Unmarshal(data, &coastInput)
	return coastInput
}

func NamespaceFinder() (namespaces []string) {
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
	es, err := elasticsearch.NewClient(cfg)

	esquery := `{
		
			"size": 0,
			"aggs" : {
				"langs" : {
					"terms" : { "field" : "NamespaceaNameAsData.keyword",  "size" : 5000 }
				}
			}}
	
	`

	res3, err := es.Search(
		es.Search.WithIndex("elastic_namespace_details"),
		es.Search.WithContext(ctx),
		es.Search.WithBody(strings.NewReader(esquery)),
		es.Search.WithPretty(),
	)
	defer res3.Body.Close()
	fmt.Println(res3)
	if err != nil {
		log.Fatal(err)
	}
	var generic NamespaceJsonOutput
	err = json.NewDecoder(res3.Body).Decode(&generic)
	if err != nil {
		log.Fatal(err)
	}

	jsonString, err := json.Marshal(generic)
	fmt.Println(err)
	res := NamespaceJsonOutput{}
	// fmt.Println(res)
	json.Unmarshal([]byte(jsonString), &res)
	for k := range res.Aggregations.Langs.Buckets {
		namespaces = append(namespaces, res.Aggregations.Langs.Buckets[k].Key)
		fmt.Println(namespaces)

	}

	return namespaces

}

func QueryStringGernaretor(namespaceForUsageQuery string, resourceType string, userCostInput InputCostData, fromDate string, toDate string) (finalquery string, AvgFieldName string) {
	var buf bytes.Buffer
	t, err := template.ParseFiles("templates/template_new.json")
	if err != nil {
		log.Fatal(err)
	}

	if resourceType == "Cpu" {
		Todatejson := fromDate
		FromDatejson := toDate
		TermKeyjson := "NamespaceName.keyword"
		Termjson := namespaceForUsageQuery
		AvgField := "NamespaceCpuValue"
		Timezonejson := userCostInput.TimeZone

		data := json_inputvars{
			Todate:   Todatejson,
			FromDate: FromDatejson,
			TermKey:  TermKeyjson,
			Term:     Termjson,
			AvgField: AvgField,
			TimeZone: Timezonejson,
		}

		_ = t.Execute(&buf, data)
		finalquery = buf.String()
		return finalquery, AvgField
	}

	if resourceType == "Memory" {
		Todatejson := fromDate
		FromDatejson := toDate
		TermKeyjson := "MemoryNamespaceName.keyword"
		Termjson := namespaceForUsageQuery
		AvgField := "NamespaceMemoryValue"
		Timezonejson := userCostInput.TimeZone

		data := json_inputvars{
			Todate:   Todatejson,
			FromDate: FromDatejson,
			TermKey:  TermKeyjson,
			Term:     Termjson,
			AvgField: AvgField,
			TimeZone: Timezonejson,
		}

		_ = t.Execute(&buf, data)
		finalquery = buf.String()
		return finalquery, AvgField
	}

	testError := "error"

	return testError, testError

}

func NamespaceResourceUsageFinder(namespaceForUsage string, resourceType string, userCostInput InputCostData, fromDate string, toDate string, timeinHours float64) (meadianValue float64) {

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
	es, err := elasticsearch.NewClient(cfg)

	finalquery, NamespaceResourceValue := QueryStringGernaretor(namespaceForUsage, resourceType, userCostInput, fromDate, toDate)
	var mapResp map[string]interface{}

	res3, err := es.Search(
		es.Search.WithIndex("some_index2"),
		es.Search.WithContext(ctx),
		es.Search.WithBody(strings.NewReader(finalquery)),
		es.Search.WithPretty(),
	)
	fmt.Sprintln(err)
	// fmt.Println(res3)
	if err := json.NewDecoder(res3.Body).Decode(&mapResp); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)

		// If no error, then convert response to a map[string]interface
	} else {

		// Iterate the document "hits" returned by API call

		var dataForMedian []float64
		for _, hit := range mapResp["hits"].(map[string]interface{})["hits"].([]interface{}) {
			// Parse the attributes/fields of the document
			doc := hit.(map[string]interface{})

			// The "_source" data is another map interface nested inside of doc
			source := doc["_source"]

			// Get the document's _id and print it out along with _source data
			docID := doc["_id"]
			fmt.Sprintln("docID:", docID)
			fmt.Sprintln("_source:", source, "\n")
			data := source.(map[string]interface{})

			dataForMedian = append(dataForMedian, data[NamespaceResourceValue].(float64))
			// fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@",resourceType,dataForMedian)

		}
		median, _ := stats.Median(dataForMedian)
		fmt.Println("######################",median)
		if resourceType == "Cpu" {
			fmt.Println(resourceType, "Cost of", namespaceForUsage, "is : $", ((median*userCostInput.CPUCost)*timeinHours)/1000)
			return median
		}

		if resourceType == "Memory" {
			fmt.Println(resourceType, "Cost of", namespaceForUsage, "is : $", ((median*userCostInput.MemoryCost)*timeinHours)*0.001048576)
			return median
		}

	}
	return

}

func NamespaceCostFinder(fromDate string, toDate string, timeinHours float64) {
	userCostInput := InputCostReader()
	fmt.Println("\ncalculating the cost between", fromDate, "to", toDate, "\n")
	if userCostInput.CalcCPU == true && userCostInput.CalcMemory == true {
		namespaceList := NamespaceFinder()
		// fmt.Println(namespaceList)

		for l := range namespaceList {
			fmt.Println("Hello from namespaceList")
			meadianCpu := NamespaceResourceUsageFinder(namespaceList[l], "Cpu", userCostInput, fromDate, toDate, timeinHours)
			// fmt.Println(meadianCpu)
			meadianMemory := NamespaceResourceUsageFinder(namespaceList[l], "Memory", userCostInput, fromDate, toDate, timeinHours)
			totalCost := ((meadianCpu*userCostInput.CPUCost)*timeinHours)/1000 + ((meadianMemory*userCostInput.MemoryCost)*timeinHours)*0.001048576
			fmt.Println("Toatal Cost of Namespace", namespaceList[l], "is", totalCost, "\n")

		}
	}

}

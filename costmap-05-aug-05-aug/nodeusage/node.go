package nodeUsage

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
	"math"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/montanaflynn/stats"
)

type NodeJsonOutput struct {
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
	Datapoint_Count float64 `json:"Datapoint_Count"`
	Ref_Value_PodUsage float64 `json:"Ref_Value_PodUsage"`
	Mi_to_GB_Value float64 `json:"Mi_to_GB_Value"`


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

func NodeFinder() (nodes []string) {
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
					"terms" : { "field" : "NodeName.keyword",  "size" : 5000 }
				}
			}}
	
	`
	res3, err := es.Search(
		es.Search.WithIndex("some_index2"),
		es.Search.WithContext(ctx),
		es.Search.WithBody(strings.NewReader(esquery)),
		es.Search.WithPretty(),
	)
	defer res3.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var generic NodeJsonOutput
	err = json.NewDecoder(res3.Body).Decode(&generic)
	if err != nil {
		log.Fatal(err)
	}
	jsonString, err := json.Marshal(generic)
	fmt.Println(err)
	res := NodeJsonOutput{}
	json.Unmarshal([]byte(jsonString), &res)
	for k := range res.Aggregations.Langs.Buckets {
		nodes = append(nodes, res.Aggregations.Langs.Buckets[k].Key)

	// fmt.Println("############################################",nodes)

	}

	return nodes

}

func QueryStringGernaretor(nodeForUsageQuery string, resourceType string, userCostInput InputCostData, fromDate string, toDate string) (finalquery string, AvgFieldName string) {
	var buf bytes.Buffer
	t, err := template.ParseFiles("templates/template_new.json")
	if err != nil {
		log.Fatal(err)
	}
	if resourceType == "Cpu" {
		Todatejson := fromDate
		FromDatejson := toDate
		TermKeyjson := "NodeName.keyword"
		Termjson := nodeForUsageQuery
		AvgField := "NodeCpu"
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
		TermKeyjson := "NodeName.keyword"
		Termjson := nodeForUsageQuery
		AvgField := "NodeMemory"
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

func NodeResourceUsageFinder(nodeForUsage string, resourceType string, userCostInput InputCostData, fromDate string, toDate string, timeinHours float64) (meadianValue float64) {

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

	finalquery, NodeResourceValue := QueryStringGernaretor(nodeForUsage, resourceType, userCostInput, fromDate, toDate)
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

		var dataForSum []float64
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

			dataForSum = append(dataForSum, data[NodeResourceValue].(float64))
		}
		sum, _ := stats.Sum(dataForSum)

		if resourceType == "Cpu" {
			cpuPrice:= (sum*userCostInput.CPUCost)/(userCostInput.Ref_Value_PodUsage*userCostInput.Datapoint_Count)
			exactCpuPrice:=roundFloat(cpuPrice, 4)
			fmt.Println("The total cpu cost of ", nodeForUsage, "is","$", exactCpuPrice)

			// fmt.Println(resourceType, "Cost of", nodeForUsage, "is : $", ((median*userCostInput.CPUCost)*timeinHours)/1000)
			return cpuPrice
		}
		if resourceType == "Memory" {
			//fmt.Println(resourceType, "Cost of", nodeForUsage, "is : $", ((median*userCostInput.MemoryCost)*timeinHours)*0.001048576)
			memoryPrice:= (((sum*userCostInput.MemoryCost)*userCostInput.Mi_to_GB_Value))/userCostInput.Datapoint_Count
			exactMemoryPrice:=roundFloat(memoryPrice, 4)
			fmt.Println("The total memory cost of ", nodeForUsage, "is","$", exactMemoryPrice)

			return memoryPrice

		}

	}
	return

}

func roundFloat(val float64, precision uint,) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val*ratio) / ratio
}


func NodeCostFinder(fromDate string, toDate string, timeinHours float64) {
	userCostInput := InputCostReader()
	fmt.Println("\ncalculating the cost between", fromDate, "to", toDate, "\n")
	if userCostInput.CalcCPU == true && userCostInput.CalcMemory == true {
		nodeList := NodeFinder()
		// fmt.Println(nodeList)

		for l := range nodeList {
			fmt.Println("Hello from Nodelist")
			meadianCpu := NodeResourceUsageFinder(nodeList[l], "Cpu", userCostInput, fromDate, toDate, timeinHours)
			// fmt.Println(meadianCpu)
			meadianMemory := NodeResourceUsageFinder(nodeList[l], "Memory", userCostInput, fromDate, toDate, timeinHours)
			totalCost := (meadianCpu)+(meadianMemory)
			finalCost:=roundFloat(totalCost, 4)
			fmt.Println("Toatal Cost of Node", nodeList[l], "is", finalCost, "\n")

		}
	}

}
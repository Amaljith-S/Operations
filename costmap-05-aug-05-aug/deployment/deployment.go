package deploymentusagesearch

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
	Datapoint_Count float64 `json:"Datapoint_Count"`
	Ref_Value_PodUsage float64 `json:"Ref_Value_PodUsage"`
	Mi_to_GB_Value float64 `json:"Mi_to_GB_Value"`

}

type json_inputvars struct {
	Todate   string
	FromDate string
	TermKey  string
	Term     string
	TermKey2 string
	Term2    string
	AvgField string
	TimeZone string
}

func InputCostReader() (coastInput InputCostData) {
	var data []byte
	data, _ = ioutil.ReadFile("input.json")

	_ = json.Unmarshal(data, &coastInput)
	return coastInput
}

func DeploymentDetailsFinder() (DeploymentNameMemory2 map[string][]string) {
	DeploymentNameMemory := make(map[string][]string)
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
					"terms" : { "field" : "Deploymentname.keyword",  "size" : 5000 }
				}
			}}
	
	`

	res3, err := es.Search(
		es.Search.WithIndex("elastic_deployment_details"),
		es.Search.WithContext(ctx),
		es.Search.WithBody(strings.NewReader(esquery)),
		es.Search.WithPretty(),
	)
	defer res3.Body.Close()

	fmt.Sprintln(err)
	var jsonData NamespaceJsonOutput
	err = json.NewDecoder(res3.Body).Decode(&jsonData)
	if err != nil {
		log.Fatal(err)
	}

	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	res := NamespaceJsonOutput{}
	json.Unmarshal([]byte(jsonString), &res)
	for k := range res.Aggregations.Langs.Buckets {
		slicedData := strings.Split(res.Aggregations.Langs.Buckets[k].Key, "_randomstringtoavoidconflict_")
		DeploymentNameMemory[slicedData[1]] = append(DeploymentNameMemory[slicedData[1]], slicedData[0])

	}

	return DeploymentNameMemory

}

func QueryStringGernaretor(namespaceForUsageQuery string, DeploymentNameforcalc string, resourceType string, userCostInput InputCostData, fromDate string, toDate string) (finalquery string, AvgFieldName string) {

	var buf bytes.Buffer
	t, err := template.ParseFiles("templates/template_for_name_deploy_kind.json")
	if err != nil {
		log.Fatal(err)
	}

	if resourceType == "Cpu" {
		Todatejson := fromDate
		FromDatejson := toDate
		TermKeyjson := "Namespace.keyword"
		Termjson := namespaceForUsageQuery
		TermKeyjson2 := "DeploymentName.keyword"
		Termjson2 := DeploymentNameforcalc
		AvgField := "PodCpu"
		Timezonejson := userCostInput.TimeZone

		data := json_inputvars{
			Todate:   Todatejson,
			FromDate: FromDatejson,
			TermKey:  TermKeyjson,
			Term:     Termjson,
			TermKey2: TermKeyjson2,
			Term2:    Termjson2,
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
		TermKeyjson := "Namespace.keyword"
		Termjson := namespaceForUsageQuery
		TermKeyjson2 := "DeploymentName.keyword"
		Termjson2 := DeploymentNameforcalc
		AvgField := "PodMemory"
		Timezonejson := userCostInput.TimeZone

		data := json_inputvars{
			Todate:   Todatejson,
			FromDate: FromDatejson,
			TermKey:  TermKeyjson,
			Term:     Termjson,
			TermKey2: TermKeyjson2,
			Term2:    Termjson2,
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

func DeploymentResourceUsageFinder(namespace string, DeploymentForUsage []string, resourceType string, userCostInput InputCostData, fromDate string, toDate string, timeinHours float64) (depname string, meadianValue float64) {

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
	fmt.Sprintln(err)

	for _, depname := range DeploymentForUsage {

		finalquery, NamespaceResourceValue := QueryStringGernaretor(namespace, depname, resourceType, userCostInput, fromDate, toDate)

		var mapResp map[string]interface{}

		res3, err := es.Search(
			es.Search.WithIndex("some_index2"),
			es.Search.WithContext(ctx),
			es.Search.WithBody(strings.NewReader(finalquery)),
			es.Search.WithPretty(),
		)
		fmt.Sprintln(err)
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
				dataForSum = append(dataForSum, data[NamespaceResourceValue].(float64))
			}
			sum, _ := stats.Sum(dataForSum)

			if resourceType == "Cpu" {
				// fmt.Println(resourceType, "Cost of", depname, "from namespace", namespace, "is : $", ((median*userCostInput.CPUCost)*timeinHours)/1000)
				cpuPrice:= (sum*userCostInput.CPUCost)/(userCostInput.Ref_Value_PodUsage*userCostInput.Datapoint_Count)
				exactCpuPrice:=roundFloat(cpuPrice, 4)
				fmt.Println("The total cpu cost of ", depname, "is","$", exactCpuPrice)
 

				return depname, exactCpuPrice

			}

			if resourceType == "Memory" {
				// fmt.Println(resourceType, "Cost of", depname, "from namespace", namespace, "is : $", ((median*userCostInput.MemoryCost)*timeinHours)*0.001048576)

				memoryPrice:= (((sum*userCostInput.MemoryCost)*userCostInput.Mi_to_GB_Value))/userCostInput.Datapoint_Count
				exactMemoryPrice:=roundFloat(memoryPrice, 4)
				fmt.Println("The total memory cost of ", depname, "is","$", exactMemoryPrice)

				
				return depname, exactMemoryPrice

			}

		}
	}
	return

}


func roundFloat(val float64, precision uint,) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val*ratio) / ratio
}


func DeploymentUsageFinder(fromDate string, toDate string, timeinHours float64) {
	userCostInput := InputCostReader()
	fmt.Println("\ncalculating the cost between", fromDate, "to", toDate, "\n")

	if userCostInput.CalcCPU == true && userCostInput.CalcMemory == true {

		namespaceList := DeploymentDetailsFinder()

		for k, l := range namespaceList {
			depname, meadianCpu := DeploymentResourceUsageFinder(k, l, "Cpu", userCostInput, fromDate, toDate, timeinHours)
			_, meadianMemory := DeploymentResourceUsageFinder(k, l, "Memory", userCostInput, fromDate, toDate, timeinHours)
			totalCost:= (meadianCpu)+(meadianMemory)
			finalCost:=roundFloat(totalCost, 4)

			fmt.Println("Toatal Cost of Deployment", depname,  "is","$", finalCost, "\n")
		}
	}
}
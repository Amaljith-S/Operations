package diskUsage

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

	}
	// fmt.Println("############################################",res)

	return nodes

}

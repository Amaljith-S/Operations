

docker restart es01
DELETE some_index2


PUT some_index2


PUT _ingest/pipeline/indexed_at
{
  "description": "Adds indexed_at timestamp to documents",
  "processors": [
    {
      "set": {
        "field": "_source.indexed_at",
        "value": "{{_ingest.timestamp}}"
      }
    }
  ]
}


PUT some_index2/_settings
{
  "index.default_pipeline": "indexed_at"
}

curl -XPUT "http://localhost:9200/some_index2/_settings" -H 'Content-Type: application/json' -d '{ "index" : { "max_result_window" : 500000 } }'

GET some_index2/_search/?size=1000&pretty=true


helm upgrade --namespace default helm-metrics-server bitnami/metrics-server -f ./metric-bitnami.yml     --set apiService.create=true







curl -X DELETE "localhost:9200/some_index2?pretty"


curl -X PUT "localhost:9200/some_index2?pretty"

curl -X PUT "localhost:9200/_ingest/pipeline/my-pipeline-id?pretty" -H 'Content-Type: application/json' -d'
{
    "description": "Adds indexed_at timestamp to documents",
    "processors": [
      {
        "set": {
          "field": "_source.indexed_at",
          "value": "{{_ingest.timestamp}}"
        }
      }
    ]
  }
'



curl -X PUT "localhost:9200/some_index2/_settings?pretty" -H 'Content-Type: application/json' -d'
{
    "index.default_pipeline": "indexed_at"
  }
'


PUT some_index2/_settings
{
  "index.default_pipeline": "indexed_at"
}



_search examples.

seARCH BY feild
			"_source": "*",
			"query": {
			  "match": {
				"Namespace" : "php"
			  }
			}

Search with in a timeframe

  "query": {
    "bool": {
      "filter": [
      {
        "range": {
        "indexed_at": {
          "gte": "now-25m",
          "lte": "now"
        }
        }
      }
      ],
      "must": [
      {
        "match": {
        "DeploymentName" : "kubecost-prometheus-node-exporter"
        }
      }
      ]
    }
    }



/////


	// res2, err := es.Search(
	// 	es.Search.WithIndex("some_index2"),
	// 	es.Search.WithContext(ctx),
	// 	es.Search.WithBody(strings.NewReader(`{

	// 			"query": {
	// 			  "constant_score": {
	// 				"filter": {
	// 				  "match": { "Namespace.keyword": "kube-system" }
	// 				}
	// 			  }
	// 			},
	// 			"aggs": {
	// 			  "Sum of PodMemory in kube-system Namespace ": { "sum": { "field": "PodMemory" } }
	// 			}

	//   }`)),
	//   es.Search.WithPretty(),
	// )
	// fmt.Sprintln(err)
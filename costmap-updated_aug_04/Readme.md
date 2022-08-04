# Costmap

costmap is Costmap is the tool used to search the Elasticsearch data and maps the user input cost details to the resource utilization, Written in Golang


**Docker run command :** 
```
$ docker run -v /home/apton-sooraj/go/src/costmap/input.json:/app/input.json --env es_host="http://192.168.1.223:9200" --env es_user="elastic" --env es_pass="Hb4cG2UZJNoKkWmKKQtM" -it coastmap:v9 ./costmap
```

**input.json example**

```
{
    "Calc_Cpu": true,
    "Calc_Memory": true,
    "TimeZone": "Asia/Kolkata",
    "Cpu_Cost" : 0.02,
    "Memory_cost" : 0.02
  }
```
If you are running the costmap from a docker setup, you can add input.json as volume.

## Running on Kubernetes

To run the costmap on Kubernetes, create a configmap from input.json, then add env varibles and configmap and run the pod, In the followinf exmple we are using a sleep command to run the pod. 

**create a configmap**
```
$ kubectl create configmap costmap-config --from-file=input.json
```
**Create pod**

```
apiVersion: v1
kind: Pod
metadata:
  name: costmap
spec:
  containers:
    - name: costmap
      image: coastmap:v9
      volumeMounts:
      - name: costmap-config
        mountPath: /app/input.json
        subPath: input.json
      env:
      - name: es_host
        value: "http://192.168.1.223:9200"      
      - name: es_user
        value: "elastic"
      - name: es_pass
        value: "Hb4cG2UZJNoKkWmKKQtM"
      command: ["sleep", "600"]
  volumes:
    - name: costmap-config
      configMap:
        name: costmap-config
  restartPolicy: Never

```
**Access the Pod and run the costmap**

```
kubectl exec -it costmap --  ./costmap
```
**Expected Output:**

![Screenshot](https://github.com/aptonworks/operations/blob/costkube/costkube/costmap/costmap.png)

if you  Update the configmap then run the following command:
```
$ kubectl create configmap costmap-config --from-file=input.json -o yaml --dry-run | kubectl apply -f -
 ```
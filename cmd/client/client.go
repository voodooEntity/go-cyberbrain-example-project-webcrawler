package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/voodooEntity/gits/src/query"
	"github.com/voodooEntity/gits/src/transport"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:] // Get command-line arguments excluding program name
	fmt.Println("Cyberbrain example webcrawler:")
	fmt.Println("Arguments:", args)
	switch args[0] {
	case "createtarget":
		createTarget(args[1], args[2])
	case "getdata":
		getDataByType(args[1])
	case "getdatarecursive":
		getDataByTypeRecursive(args[1])
	case "getopenjobs":
		getOpenJobs()
	case "getdeps":
		findDependencies(args[1])
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Println("Supported commands are createtarget {domain} {initialPage}, getdata {entityType}, getdatar {entityType}, getopenjobs, getdeps {entityType}")
	}
}

func getOpenJobs() {
	// HTTP endpoint
	posturl := "http://127.0.0.1:1984/v1/query"

	qry := query.New().Read("State").Match("Value", "==", "Open").Match("Context", "==", "System").From(
		query.New().Read("Job").Match("Context", "==", "System"))

	b, err := json.Marshal(qry)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendToApi(posturl, b)
}

func findDependencies(entityType string) {
	// HTTP endpoint
	posturl := "http://127.0.0.1:1984/v1/query"

	qry := query.New().Read("DependencyLookup").Match("Value", "==", entityType).To(
		query.New().Read("Dependency").From(
			query.New().Read("Action"),
		),
	)

	b, err := json.Marshal(qry)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendToApi(posturl, b)
}

func getDataByType(entityType string) {
	posturl := "http://127.0.0.1:1984/v1/query"
	qry := query.New().Read(entityType)
	b, err := json.Marshal(qry)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendToApi(posturl, b)
}

func getDataByTypeRecursive(targetType string) {
	// HTTP endpoint
	posturl := "http://127.0.0.1:1984/v1/query"

	qry := query.New().Read(targetType).TraverseOut(100)

	b, err := json.Marshal(qry)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendToApi(posturl, b)
}

func createTarget(domain string, pageUrl string) {
	// HTTP endpoint
	posturl := "http://127.0.0.1:1984/v1/learn"

	data := transport.TransportEntity{
		Type:       "Domain",
		ID:         0,
		Value:      domain,
		Context:    "data",
		Properties: make(map[string]string),
		ChildRelations: []transport.TransportRelation{
			{
				Target: transport.TransportEntity{
					Type:       "Page",
					ID:         0,
					Value:      pageUrl,
					Context:    "data",
					Properties: make(map[string]string),
				},
			},
		},
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendToApi(posturl, b)
}

func getAllJobsTraversed() {
	// HTTP endpoint
	posturl := "http://127.0.0.1:1984/v1/query"

	qry := query.New().Read("Job").TraverseOut(30)

	b, err := json.Marshal(qry)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendToApi(posturl, b)
}

func sendToApi(url string, data []byte) {
	trans := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: trans}

	// Create a HTTP post request
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	response, error := client.Do(r)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}

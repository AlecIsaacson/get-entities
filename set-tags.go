package main

import (
    "context"
    "fmt"
    "flag"
    "os"
    "bufio"
    "strings"

    "github.com/machinebox/graphql"
)

func main() {
  nrAPI := flag.String("apikey", "", "New Relic admin user API Key")
	logVerbose := flag.Bool("verbose", false, "Writes verbose logs for debugging")
  hosts := flag.String("hostfile", "", "List of hosts with GUIDs to be tagged")
	flag.Parse()

  if *logVerbose {
    fmt.Println("Entity finder v1.0")
    fmt.Println("Verbose logging enabled")
  }

  graphqlClient := graphql.NewClient("https://api.newrelic.com/graphql")

  graphqlRequest := graphql.NewRequest(`
    mutation ($guid: EntityGuid!, $key: String!, $values: String!) {
      taggingAddTagsToEntity(guid: $guid, tags: {key: $key, values: [$values]}) {
        errors {
          message
        }
      }
    }
  `)

  var graphqlResponse interface{}

  hostFile, err := os.Open(*hosts)
  if err != nil {
    //log.Fatal(err)
  }
  defer hostFile.Close()

  scanner := bufio.NewScanner(hostFile)
  for scanner.Scan() {
    hostInfo := strings.Split(scanner.Text(),",")
    guid := hostInfo[1]
    market := hostInfo[0][0:2]
    storeId := hostInfo[0][2:7]
    role := hostInfo[0][7:10]
    instance := hostInfo[0][10:12]

    fmt.Println("Hostname:", hostInfo[0])
    fmt.Println("GUID: ", guid)
    fmt.Println("Elements:", market, storeId, role, instance)

    graphqlRequest.Header.Set("API-Key",*nrAPI)
    graphqlRequest.Var("guid", guid)
    graphqlRequest.Var("key", "storeId")
    graphqlRequest.Var("values", storeId)

    if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
        panic(err)
    }
  }
}

//This is a one-off that splits hostnames that are in a very specific format so they can be
//fed back into the New Relic platform as tags.  You could use it as a pattern for your own app, but
//unless you're the team I wrote this for, it won't be much use as-is.
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
  //Define flags and defaults
  nrAPI := flag.String("apikey", "", "New Relic admin user API Key")
	logVerbose := flag.Bool("verbose", false, "Writes verbose logs for debugging")
  hosts := flag.String("hostfile", "", "List of hosts with GUIDs to be tagged")
	flag.Parse()

  if *logVerbose {
    fmt.Println("Entity set-tags v1.0")
    fmt.Println("Verbose logging enabled")
  }

  graphqlClient := graphql.NewClient("https://api.newrelic.com/graphql")

  //The request to set a tag looks like this
  graphqlRequest := graphql.NewRequest(`
    mutation ($guid: EntityGuid!, $key: String!, $values: String!) {
      taggingAddTagsToEntity(guid: $guid, tags: {key: $key, values: [$values]}) {
        errors {
          message
        }
      }
    }
  `)

  //Define an interface to hold the response.
  var graphqlResponse interface{}

  //Get the list of entities we're going to tag
  hostFile, err := os.Open(*hosts)
  if err != nil {
    //log.Fatal(err)
  }
  defer hostFile.Close()

  //For each line in the file, split the entity name into components.
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

    //Set the GraphQL request that will create the tags we need.
    graphqlRequest.Header.Set("API-Key",*nrAPI)
    graphqlRequest.Var("guid", guid)
    graphqlRequest.Var("key", "storeId")
    graphqlRequest.Var("values", storeId)

    //Execute the request.
    if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
        panic(err)
    }
  }
}

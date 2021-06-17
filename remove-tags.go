//This app will remove the speified tag from a list of New Relic entities.
//It expects a file that contains a list of entity names and GUIDs to act on.
//A good source for that file is the get-entities app in this repo.
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
  removeTag := flag.String("removetag", "", "Tag to be removed from entities")
	flag.Parse()

  if *logVerbose {
    fmt.Println("Entity tag remover v1.0")
    fmt.Println("Verbose logging enabled")
  }

  graphqlClient := graphql.NewClient("https://api.newrelic.com/graphql")

  graphqlRequest := graphql.NewRequest(`
    mutation ($guid: EntityGuid!, $key: String!) {
      taggingDeleteTagFromEntity(guid: $guid, tagKeys: [$key]) {
        errors {
          message
        }
      }
    }
  `)

  var graphqlResponse interface{}

  hostFile, err := os.Open(*hosts)
  if err != nil {
    panic(err)
  }
  defer hostFile.Close()

  scanner := bufio.NewScanner(hostFile)
  for scanner.Scan() {
    hostInfo := strings.Split(scanner.Text(),",")
    guid := hostInfo[1]

    fmt.Println("Hostname:", hostInfo[0])
    fmt.Println("GUID: ", guid)

    graphqlRequest.Header.Set("API-Key",*nrAPI)
    graphqlRequest.Var("guid", guid)
    graphqlRequest.Var("key", *removeTag)

    if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
        panic(err)
    }
  }
}

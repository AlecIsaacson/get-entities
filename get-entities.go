package main

import (
    "context"
    "fmt"
    "flag"

    "github.com/machinebox/graphql"
)

type nrEntityStruct struct {
	//Data struct {
		Actor struct {
			EntitySearch struct {
				Results struct {
          NextCursor string `json:"nextCursor"`
					Entities []struct {
						AccountID  int    `json:"accountId"`
						EntityType string `json:"entityType"`
						GUID       string `json:"guid"`
						Name       string `json:"name"`
					} `json:"entities"`
				} `json:"results"`
			} `json:"entitySearch"`
		} `json:"actor"`
	//} `json:"data"`
	// Extensions struct {
	// 	NrOnly struct {
	// 		Docs           string `json:"_docs"`
	// 		DeepTrace      string `json:"deepTrace"`
	// 		HTTPRequestLog []struct {
	// 			Body string `json:"body"`
	// 			Curl string `json:"curl"`
	// 		} `json:"httpRequestLog"`
	// 	} `json:"nrOnly"`
	// } `json:"extensions"`
}

func main() {
  nrAPI := flag.String("apikey", "", "New Relic admin user API Key")
  nrQuery := flag.String("nrql","name like '%'","A valid NRQL query")
	logVerbose := flag.Bool("verbose", false, "Writes verbose logs for debugging")
	flag.Parse()

  if *logVerbose {
    fmt.Println("Entity finder v1.0")
    fmt.Println("Verbose logging enabled")
  }

  graphqlClient := graphql.NewClient("https://api.newrelic.com/graphql")

  graphqlRequest := graphql.NewRequest(`
    query($query: String!)
    {
      actor {
        entitySearch(query: $query) {
          results {
          nextCursor
            entities {
              name
              entityType
              guid
              accountId
            }
          }
        }
      }
    }
  `)

  graphqlCursorRequest := graphql.NewRequest(`
    query($query: String!, $nextCursor: String!)
    {
      actor {
        entitySearch(query: $query) {
          results (cursor: $nextCursor){
          nextCursor
            entities {
              name
              entityType
              guid
              accountId
            }
          }
        }
      }
    }
  `)

  //nrQuery := "domain = 'APM' and accountId =" + *nrAccountID
  graphqlRequest.Var("query", *nrQuery)
  graphqlRequest.Header.Set("API-Key",*nrAPI)

  var graphqlResponse nrEntityStruct
  if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
      panic(err)
  }

  for _,entity := range graphqlResponse.Actor.EntitySearch.Results.Entities {
    fmt.Printf("%s,%s\n", entity.Name, entity.GUID)
  }

  nextCursor := graphqlResponse.Actor.EntitySearch.Results.NextCursor

  for {
    graphqlResponse = nrEntityStruct{}
    if len(nextCursor) > 0 {
      graphqlCursorRequest.Var("query", nrQuery)
      graphqlCursorRequest.Var("nextCursor", nextCursor)
      graphqlCursorRequest.Header.Set("API-Key",*nrAPI)

      if err := graphqlClient.Run(context.Background(), graphqlCursorRequest, &graphqlResponse); err != nil {
          panic(err)
      }

      for _,entity := range graphqlResponse.Actor.EntitySearch.Results.Entities {
        fmt.Printf("%s,%s\n", entity.Name, entity.GUID)
      }
      nextCursor = graphqlResponse.Actor.EntitySearch.Results.NextCursor
    } else {
      break
    }
  }
}

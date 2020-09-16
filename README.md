# get-entities

This app queries the New Relic GraphQL endpoint and dumps the name and GUID of entities.  By default, it'll dump all entities that your API key has permissions to see.

The output is formatted as *entityName,entityGUID* - one entity per line.  You'll probably want to redirect the output to a file for further use.

This app takes two mandatory and one optional command line switches:

 `-apikey : A New Relic API key that can query the GraphQL API.  `  
 `-query : A valid query string  `  
 `-verbose : An optional flag that enables verbose logging for debugging.`
 
[See the New Relic Entity GraphQL tutorial for details on query format.](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-entities-api-tutorial)

As an example:

This will dump all entities that I have permissions to see:  
`./get-entities -apikey *myAPIKey*`  

This will dump only application entities:  
 `./get-entities -apikey *myAPIKey* -query "type = 'APPLICATION'"`  

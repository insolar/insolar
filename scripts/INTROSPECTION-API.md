# How To Introspect Insolard Nodes

## How-To use API

build insolard with introspection api and run launchnet:

    BUILD_TAGS='-tags "introspection"' ./scripts/insolard/launchnet.sh -g

check is introspector API works on heavy node:

    curl -X POST -d '{}' http://127.0.0.1:55501/getMessagesStat

or:

    http -j POST http://127.0.0.1:55501/getMessagesStat

check gRPC endpoint with [grpcurl](https://github.com/fullstorydev/grpcurl):

    grpcurl -plaintext 127.0.0.1:55501 introproto.Publisher.GetMessagesStat

### fetch swagger json

    http http://127.0.0.1:55501/swagger.json

## How to use message publishing filter

check API ports after launchnet's start:

    ./scripts/glogs 'started introspection server on'

lock message type `TypeSetOutgoingRequest` on Virtual nodes:

    http POST http://127.0.0.1:55502/setMessagesFilter Name=TypeSetOutgoingRequest Enable:=true
    
check is type `TypeSetOutgoingRequest` locked and how much is filtered out:

    http POST http://127.0.0.1:55502/getMessagesFilters

filter only enabled filters with jq:

    http POST http://127.0.0.1:55502/getMessagesFilters | jq '.Filters | to_entries | .[] | select(.value.Enable==true)'

## How to develop of new APIs

1. Add types and methods to Publisher service
    * modify `./instrumentation/introspector/introproto/publisher.proto`
    * regenerate files `make generate-introspector-proto`
2. Add and implement new middleware
    * add code and test files to `./instrumentation/introspector/pubsubwrap/`
    * add middleware to `./instrumentation/introspector/pubsubwrap/service.go`
4. Create and initialize middleware here: `server/internal/pubsub_introspect.go`

# How To Introspect Insolard Nodes

## How-To use API

build insolard with introspection api and run launchnet:

    BUILD_TAGS='-tags "introspection"' ./scripts/insolard/launchnet.sh -g

check is introspector API works on heavy node:

    curl -X POST -d '{}' http://127.0.0.1:55501/getMessagesFilters
    
or:

    http -j POST http://127.0.0.1:55501/getMessagesFilters

check gRPC endpoint with [grpcurl](https://github.com/fullstorydev/grpcurl):

    grpcurl -plaintext 127.0.0.1:55501 introspector.Publisher.GetMessagesFilters

### fetch swagger json

    http http://127.0.0.1:55501/swagger.json

## How to use message publishing filter

check API ports after launchnet's start:

    ./scripts/glogs 'started introspection server on'

set lock message publishing by type on node for launchnet:

    

## How to develop of new APIs

_TODO_

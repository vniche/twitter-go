# twitter-go #

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/vniche/twitter-go?tab=doc)

twitter-go is a Go client library for accessing the [Twitter API v2](https://developer.twitter.com/en/docs/twitter-api).

## Usage ##

First we have to go to [Twitter Developer](https://developer.twitter.com/en) platform and register to get hands on a bearer token to use on our requests to the API, since all endpoints require one.

### Installing ####

```go
import "github.com/vniche/twitter-go"
```

Construct a new Twitter client, then use the various services on the client to
access different parts of the Twitter API. For example:

```go
// Fetching an user by ID

import (
    twitter "github.com/vniche/twitter-go"
)

client := twitter.WithBearerToken("mybearertoken123", nil)

// custom parameters for query
params := make(map[string][]string)
params["user.fields"] = []string{
    "public_metrics",
}

// fetch user by id
user, err := client.LookupUserByID(userID, params)
if err != nil {
    log.Panicf("unable to fetch user by id: %+v", err)
}

fmt.Printf("user: %+v\n", user)
```

## Contributing ##

I would like to cover the entire Twitter API and contributions are of course always welcome. The
calling pattern is pretty well established, so adding new methods is relatively
straightforward.

TODO: Contribution doc.

## License ##

This library is distributed under the MIT license found in the [LICENSE](./LICENSE)
file.

package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Mention struct {
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Username string `json:"username"`
	ID       string `json:"id"`
}

type Hashtag struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	Tag   string `json:"tag"`
}

type Entities struct {
	Mentions []Mention `json:"mentions,omitempty"`
	Hashtags []Hashtag `json:"hashtags,omitempty"`
}

type Includes struct {
	Users []User `json:"users"`
}
type TweetMeta struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

// Tweet represents a Twitter tweet
type Tweet struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	Text      string    `json:"text"`
	Entities  Entities  `json:"entities"`
	CreatedAt time.Time `json:"created_at"`
}

type SearchRecentTweetsResponse struct {
	Tweets   []Tweet   `json:"data"`
	Includes Includes  `json:"includes"`
	Meta     TweetMeta `json:"meta"`
}

// SearchRecentTweets returns Tweets from the last seven days that match a search query
//
// Usage example:
// params := make(map[string][]string)
// params["user.fields"] = []string{
// 	"username",
// 	"url",
// 	"location",
// 	"public_metrics",
// }
// params["query"] = []string{
// 	"has:hashtags (\"#buildinginpublic\" OR \"#buildinpublic\") -is:reply -is:retweet",
// }
// params["tweet.fields"] = []string{
// 	"entities",
// }
// params["expansions"] = []string{
// 	"author_id",
// 	"entities.mentions.username",
// }
//
// var response *twitter.SearchRecentTweetsResponse
// response, err = client.SearchRecentTweets(ctx, params)
// if err != nil {
// 	log.Panicf("unable to listen to tweets stream: %+v", err)
// }
func (client *Client) SearchRecentTweets(ctx context.Context, parameters map[string][]string) (*SearchRecentTweetsResponse, error) {
	queryParams, err := ParseURLParameters(parameters)
	if err != nil {
		return nil, err
	}

	req, err := client.buildRequest("GET", fmt.Sprintf("/tweets/search/recent%s", queryParams), nil)
	if err != nil {
		return nil, err
	}

	var response *http.Response
	response, err = client.do(ctx, req)
	if err != nil {
		return nil, err
	}

	var parsedResponse *SearchRecentTweetsResponse
	err = json.NewDecoder(response.Body).Decode(&parsedResponse)
	return parsedResponse, err
}

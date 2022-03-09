package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Rule struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Tag   string `json:"tag,omitempty"`
}

type RuleSummary struct {
	Created    int `json:"created"`
	NotCreated int `json:"not_created"`
	Valid      int `json:"valid"`
	Invalid    int `json:"invalid"`
	Deleted    int `json:"deleted"`
	NotDeleted int `json:"not_deleted"`
}

type RuleMeta struct {
	Sent    string      `json:"sent"`
	Summary RuleSummary `json:"summary,omitempty"`
}

type DeleteRules struct {
	IDs []string `json:"ids"`
}

type AddOrDeleteRulesRequest struct {
	Add    []Rule       `json:"add"`
	Delete *DeleteRules `json:"delete,omitempty"`
}

type AddOrDeleteRulesResponse struct {
	Rules  []Rule         `json:"data"`
	Meta   RuleMeta       `json:"meta"`
	Errors []GenericError `json:"errors,omitempty"`
}

// CreateRule tries to create a rule and
//
// Usage example:
// var result *twitter.AddOrDeleteRulesResponse
// result, err = client.AddOrDeleteRules(ctx, &twitter.AddOrDeleteRulesRequest{
// 	Add: []twitter.Rule{
// 		{
// 			Value: "has:hashtags (\"#buildinginpublic\" OR \"#buildinpublic\") -is:reply -is:retweet",
// 			Tag:   "non-reply non-retweet with hashtags buildinginpublic or buildinpublic",
// 		},
// 	},
// 	Delete: &twitter.DeleteRules{
// 		IDs: []string{
// 			"123123123123",
// 		},
// 	},
// }, false)
// if err != nil {
// 	log.Panicf("unable to create rule: %+v", err)
// }
func (client *Client) AddOrDeleteRules(ctx context.Context, payload *AddOrDeleteRulesRequest, dryRun bool) (*AddOrDeleteRulesResponse, error) {
	var queryParams string
	if dryRun {
		queryParams = "?dry_run=true"
	}

	req, err := client.buildRequest("POST", fmt.Sprintf("/tweets/search/stream/rules%s", queryParams), payload)
	if err != nil {
		return nil, err
	}

	var response *http.Response
	response, err = client.do(ctx, req)
	if err != nil {
		return nil, err
	}

	var parsedResponse *AddOrDeleteRulesResponse
	if err = json.NewDecoder(response.Body).Decode(&parsedResponse); err != nil {
		return nil, err
	}

	if len(parsedResponse.Errors) > 0 {
		return nil, fmt.Errorf("%+v", parsedResponse.Errors)
	}

	return parsedResponse, err
}

type GetRulesResponse struct {
	Rules []Rule   `json:"data"`
	Meta  RuleMeta `json:"meta"`
}

// GetRules Return either a single rule, or a list of rules that have been added to the stream.
//
// Usage example:
// getRulesParams := make(map[string][]string)
// var rules []twitter.Rule
// rules, err = client.GetRules(ctx, getRulesParams)
// if err != nil {
// 	log.Panicf("unable to get rules: %+v", err)
// }
func (client *Client) GetRules(ctx context.Context, parameters map[string][]string) ([]Rule, error) {
	queryParams, err := ParseURLParameters(parameters)
	if err != nil {
		return nil, err
	}

	for key := range parameters {
		if key != "ids" {
			return nil, fmt.Errorf("query parameter key invalid: %s", key)
		}
	}

	req, err := client.buildRequest("GET", fmt.Sprintf("/tweets/search/stream/rules%s", queryParams), nil)
	if err != nil {
		return nil, err
	}

	var response *http.Response
	response, err = client.do(ctx, req)
	if err != nil {
		return nil, err
	}

	var parsedResponse *GetRulesResponse
	err = json.NewDecoder(response.Body).Decode(&parsedResponse)
	return parsedResponse.Rules, err
}

type SearchStreamResponse struct {
	Tweet    Tweet          `json:"data"`
	Includes Includes       `json:"includes"`
	Errors   []GenericError `json:"errors,omitempty"`
}

// SearchStream streams Tweets in real-time that match the rules that you added to the stream
//
// Usage example:
// params := make(map[string][]string)
// 	params["user.fields"] = []string{
// 		"username",
// 		"url",
// 		"location",
// 		"public_metrics",
// 	}
// 	params["tweet.fields"] = []string{
// 		"entities",
// 	}
// 	params["expansions"] = []string{
// 		"author_id",
// 		"entities.mentions.username",
// 	}
//
// 	var channel *twitter.Channel
// 	channel, err = client.SearchStream(ctx, params)
// 	if err != nil {
// 		log.Panicf("unable to listen to tweets stream: %+v", err)
// 	}
//
// 	for {
// 		select {
// 		case <-quit:
// 			cancel()
// 		case <-ctx.Done():
// 			fmt.Printf("close stream from channel consumer\n")
// 			channel.Close()
// 			return
// 		case message, ok := <-channel.Receive():
// 			if message == nil || !ok {
// 				return
// 			}
//
// 			if len(message.Tweet.Entities.Hashtags) <= 10 {
// 				fmt.Printf("received tweet\n")
// 				fmt.Printf("%+v\n", message.Tweet)
// 			}
// 		}
// 	}
func (client *Client) SearchStream(ctx context.Context, parameters map[string][]string) (*Channel, error) {
	queryParams, err := ParseURLParameters(parameters)
	if err != nil {
		return nil, err
	}

	req, err := client.buildRequest("GET", fmt.Sprintf("/tweets/search/stream%s", queryParams), nil)
	if err != nil {
		return nil, err
	}

	var response *http.Response
	response, err = client.do(ctx, req)
	if err != nil {
		return nil, err
	}

	channel := NewChannel()

	go func(ctx context.Context) {
		dec := json.NewDecoder(response.Body)
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("close stream from channel producer\n")
				channel.Close()
				return
			default:
				var streamResponse *SearchStreamResponse
				err = dec.Decode(&streamResponse)
				if err != nil {
					if err == io.EOF {
						break
					}
					fmt.Printf("unable to decode stream message: %+v\n", err)
					return
				}

				channel.channel <- streamResponse
			}
		}
	}(ctx)

	return channel, err
}

package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// User represents a Twitter user
type User struct {
	Name          string `json:"name"`
	Username      string `json:"username"`
	ID            string `json:"id"`
	PublicMetrics struct {
		FollowersCount int `json:"followers_count"`
		FollowingCount int `json:"following_count"`
		TweetCount     int `json:"tweet_count"`
		ListedCount    int `json:"listed_count"`
	} `json:"public_metrics,omitempty"`
	URL string `json:"url"`
}

type LookupUserByIDResponse struct {
	Data *User `json:"data"`
}

var Expansions = struct {
	PinnedTweetID string
}{
	"pinned_tweet_id",
}

var TweetFields = struct {
	Attachments        string
	AuthorID           string
	ContextAnnotations string
	ConversationID     string
	CreatedAt          string
	Entities           string
	Geo                string
	ID                 string
	InReplyToUserID    string
	Lang               string
	NonPublicMetrics   string
	PublicMetrics      string
	OrganicMetrics     string
	PromotedMetrics    string
	PossiblySensitive  string
	ReferencedTweets   string
	ReplySettings      string
	Source             string
	Text               string
	Withheld           string
}{
	"attachments",
	"author_id",
	"context_annotations",
	"conversation_id",
	"created_at",
	"entities",
	"geo",
	"id",
	"in_reply_to_user_id",
	"lang",
	"non_public_metrics",
	"public_metrics",
	"organic_metrics",
	"promoted_metrics",
	"possibly_sensitive",
	"referenced_tweets",
	"reply_settings",
	"source",
	"text",
	"withheld",
}

var UserFields = struct {
	CreatedAt       string
	Description     string
	Entities        string
	ID              string
	Location        string
	Name            string
	PinnedTweetID   string
	ProfileImageURL string
	Protected       string
	PublicMetrics   string
	URL             string
	Username        string
	Verified        string
	Withheld        string
}{
	"created_at",
	"description",
	"entities",
	"id",
	"location",
	"name",
	"pinned_tweet_id",
	"profile_image_url",
	"protected",
	"public_metrics",
	"url",
	"username",
	"verified",
	"withheld",
}

var UserGenericQueryParameters = struct {
	Expansions  string
	TweetFields string
	UserFields  string
}{
	"expansions",
	"tweet.fields",
	"user.fields",
}

var userGenericParametersMap map[string]interface{}

func init() {
	userGenericParametersMap = make(map[string]interface{})
	userGenericParametersMap[UserGenericQueryParameters.Expansions] = Expansions
	userGenericParametersMap[UserGenericQueryParameters.TweetFields] = TweetFields
	userGenericParametersMap[UserGenericQueryParameters.UserFields] = UserFields
}

// UserByID returns an unique user by it's ID
func (client *Client) LookupUserByID(ctx context.Context, userID string, parameters map[string][]string) (*User, error) {
	queryParams, err := ParseURLParameters(parameters)
	if err != nil {
		return nil, err
	}

	for key, value := range parameters {
		if !IsOneOfEnum(key, UserGenericQueryParameters) {
			return nil, fmt.Errorf("query parameter key invalid: %s", key)
		}

		for _, current := range value {
			if !IsOneOfEnum(current, userGenericParametersMap[key]) {
				return nil, fmt.Errorf("query parameter key value invalid: %s=%s", key, current)
			}
		}

		queryParams += key + "=" + strings.Join(value, ",")
	}

	req, err := client.buildRequest("GET", fmt.Sprintf("/users/%s%s", userID, queryParams), nil)
	if err != nil {
		return nil, err
	}

	var response *http.Response
	response, err = client.do(ctx, req)
	if err != nil {
		return nil, err
	}

	var parsedResponse *LookupUserByIDResponse
	err = json.NewDecoder(response.Body).Decode(&parsedResponse)
	return parsedResponse.Data, err
}

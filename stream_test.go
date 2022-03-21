package twitter

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessInvalidMessageSearchStream(t *testing.T) {
	ctx := context.Background()

	channel := NewChannel()

	body := ioutil.NopCloser(strings.NewReader(`{"non_compliant":"content"}`))

	err := processSearchStream(ctx, channel, body)
	assert.NotNil(t, err)
}

func TestProcessValidMessageSearchStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	channel := NewChannel()

	body := ioutil.NopCloser(strings.NewReader(`{"data":{"id": "1067094924124872705","text": "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN"}}`))

	go func(ctx context.Context) {
		err := processSearchStream(ctx, channel, body)
		assert.Nil(t, err)
	}(ctx)

	response, ok := <-channel.Receive()
	assert.True(t, ok)
	assert.NotNil(t, response)
	cancel()
}

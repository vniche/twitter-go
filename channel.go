package twitter

import (
	"github.com/google/uuid"
)

var managedChannels = make(map[string]*Channel)

type Channel struct {
	ID      string `json:"id"`
	channel chan *SearchStreamResponse
}

// New is a constructor for new channel
func NewChannel() *Channel {
	id := uuid.New().String()
	channel := &Channel{
		ID:      id,
		channel: make(chan *SearchStreamResponse),
	}

	managedChannels[id] = channel

	return channel
}

// IsClosed returns if the channel is closed or not
func (channel *Channel) IsClosed() bool {
	return channel.channel == nil
}

// Close tries to close a managed channel
func (channel *Channel) Close() {
	if channel.IsClosed() {
		// channel already closed
		return
	}

	// closes managed channel
	close(channel.channel)

	// deletes channel from map
	delete(managedChannels, channel.ID)
}

func (channel *Channel) Receive() <-chan *SearchStreamResponse {
	return channel.channel
}

// ShutdownStream tries to close every managed channel and delete it from managed channels slice
func ShutdownStream() {
	for _, channel := range managedChannels {
		channel.Close()
	}
}

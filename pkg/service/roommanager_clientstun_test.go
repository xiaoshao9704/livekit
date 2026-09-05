package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/livekit/protocol/livekit"

	"github.com/livekit/livekit-server/pkg/config"
)

// hearth patch 4: RTCConfig.ClientSTUNServers overrides what is advertised to clients.
// TURN is left disabled so the function never touches the participant or the auth handler.
func TestIceServersForParticipantClientSTUNServers(t *testing.T) {
	newManager := func(clientSTUN *[]string, stunServers []string) *RoomManager {
		conf := &config.Config{}
		conf.RTC.STUNServers = stunServers
		conf.RTC.ClientSTUNServers = clientSTUN
		return &RoomManager{config: conf}
	}

	stunURLs := func(servers []*livekit.ICEServer) []string {
		var urls []string
		for _, s := range servers {
			for _, u := range s.Urls {
				if strings.HasPrefix(u, "stun:") {
					urls = append(urls, u)
				}
			}
		}
		return urls
	}

	t.Run("nil falls back to defaults", func(t *testing.T) {
		got := newManager(nil, nil).iceServersForParticipant("apikey", nil, false)
		require.NotEmpty(t, stunURLs(got), "upstream behaviour: default STUN servers advertised")
	})

	t.Run("empty slice advertises no stun", func(t *testing.T) {
		empty := []string{}
		got := newManager(&empty, []string{"stun.example.com:3478"}).iceServersForParticipant("apikey", nil, false)
		require.Empty(t, stunURLs(got))
	})

	t.Run("list is advertised verbatim", func(t *testing.T) {
		list := []string{"stun1.example.com:3478", "stun2.example.com:19302"}
		got := newManager(&list, nil).iceServersForParticipant("apikey", nil, false)
		require.Equal(t, []string{"stun:stun1.example.com:3478", "stun:stun2.example.com:19302"}, stunURLs(got))
	})
}

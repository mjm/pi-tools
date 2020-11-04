package trips

import (
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
)

type fakeMessagesClient struct {
	messagespb.MessagesServiceClient
}

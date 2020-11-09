package trips

import (
	"context"

	"google.golang.org/grpc"

	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
)

type fakeMessagesClient struct{}

func (fakeMessagesClient) SendTripBeganMessage(context.Context, *messagespb.SendTripBeganMessageRequest, ...grpc.CallOption) (*messagespb.SendTripBeganMessageResponse, error) {
	return &messagespb.SendTripBeganMessageResponse{}, nil
}

func (fakeMessagesClient) SendTripCompletedMessage(context.Context, *messagespb.SendTripCompletedMessageRequest, ...grpc.CallOption) (*messagespb.SendTripCompletedMessageResponse, error) {
	return &messagespb.SendTripCompletedMessageResponse{}, nil
}

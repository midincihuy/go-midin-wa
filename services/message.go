package services

import (
	"context"
	"fmt"
	"midincihuy/go-midin-wa/domains/message"
	domainMessage "midincihuy/go-midin-wa/domains/message"
	"midincihuy/go-midin-wa/pkg/whatsapp"
	"midincihuy/go-midin-wa/validations"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/proto/waSyncAction"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"time"
)

type serviceMessage struct {
	WaCli *whatsmeow.Client
}

func NewMessageService(waCli *whatsmeow.Client) domainMessage.IMessageService {
	return &serviceMessage{
		WaCli: waCli,
	}
}

func (service serviceMessage) ReactMessage(ctx context.Context, request message.ReactionRequest) (response message.GenericResponse, err error) {
	if err = validations.ValidateReactMessage(ctx, request); err != nil {
		return response, err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	msg := &waE2E.Message{
		ReactionMessage: &waE2E.ReactionMessage{
			Key: &waCommon.MessageKey{
				FromMe:    proto.Bool(true),
				ID:        proto.String(request.MessageID),
				RemoteJID: proto.String(dataWaRecipient.String()),
			},
			Text:              proto.String(request.Emoji),
			SenderTimestampMS: proto.Int64(time.Now().UnixMilli()),
		},
	}
	ts, err := service.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Reaction sent to %s (server timestamp: %s)", request.Phone, ts.Timestamp)
	return response, nil
}

func (service serviceMessage) RevokeMessage(ctx context.Context, request domainMessage.RevokeRequest) (response domainMessage.GenericResponse, err error) {
	if err = validations.ValidateRevokeMessage(ctx, request); err != nil {
		return response, err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	ts, err := service.WaCli.SendMessage(context.Background(), dataWaRecipient, service.WaCli.BuildRevoke(dataWaRecipient, types.EmptyJID, request.MessageID))
	if err != nil {
		return response, err
	}

	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Revoke success %s (server timestamp: %s)", request.Phone, ts.Timestamp)
	return response, nil
}

func (service serviceMessage) DeleteMessage(ctx context.Context, request domainMessage.DeleteRequest) (err error) {
	if err = validations.ValidateDeleteMessage(ctx, request); err != nil {
		return err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return err
	}

	isFromMe := "1"
	if len(request.MessageID) > 22 {
		isFromMe = "0"
	}

	patchInfo := appstate.PatchInfo{
		Timestamp: time.Now(),
		Type:      appstate.WAPatchRegularHigh,
		Mutations: []appstate.MutationInfo{{
			Index: []string{appstate.IndexDeleteMessageForMe, dataWaRecipient.String(), request.MessageID, isFromMe, service.WaCli.Store.ID.String()},
			Value: &waSyncAction.SyncActionValue{
				DeleteMessageForMeAction: &waSyncAction.DeleteMessageForMeAction{
					DeleteMedia:      proto.Bool(true),
					MessageTimestamp: proto.Int64(time.Now().UnixMilli()),
				},
			},
		}},
	}

	if err = service.WaCli.SendAppState(patchInfo); err != nil {
		return err
	}
	return nil
}

func (service serviceMessage) UpdateMessage(ctx context.Context, request domainMessage.UpdateMessageRequest) (response domainMessage.GenericResponse, err error) {
	if err = validations.ValidateUpdateMessage(ctx, request); err != nil {
		return response, err
	}

	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	msg := &waE2E.Message{Conversation: proto.String(request.Message)}
	ts, err := service.WaCli.SendMessage(context.Background(), dataWaRecipient, service.WaCli.BuildEdit(dataWaRecipient, request.MessageID, msg))
	if err != nil {
		return response, err
	}

	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Update message success %s (server timestamp: %s)", request.Phone, ts.Timestamp)
	return response, nil
}

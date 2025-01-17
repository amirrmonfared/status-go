package protocol

import (
	"context"
	"crypto/ecdsa"
	"errors"

	"go.uber.org/zap"

	"github.com/golang/protobuf/proto"

	"github.com/status-im/status-go/protocol/common"
	"github.com/status-im/status-go/protocol/protobuf"
	"github.com/status-im/status-go/protocol/requests"
)

var ErrInvalidEditOrDeleteAuthor = errors.New("sender is not the author of the message")
var ErrInvalidDeleteTypeAuthor = errors.New("message type cannot be deleted")
var ErrInvalidEditContentType = errors.New("only text or emoji messages can be replaced")
var ErrInvalidDeletePermission = errors.New("don't have enough permission to delete")

func (m *Messenger) EditMessage(ctx context.Context, request *requests.EditMessage) (*MessengerResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}
	message, err := m.persistence.MessageByID(request.ID.String())
	if err != nil {
		return nil, err
	}

	if message.From != common.PubkeyToHex(&m.identity.PublicKey) {
		return nil, ErrInvalidEditOrDeleteAuthor
	}

	if message.ContentType != protobuf.ChatMessage_TEXT_PLAIN && message.ContentType != protobuf.ChatMessage_EMOJI {
		return nil, ErrInvalidEditContentType
	}

	// A valid added chat is required.
	chat, ok := m.allChats.Load(message.ChatId)
	if !ok {
		return nil, errors.New("Chat not found")
	}

	clock, _ := chat.NextClockAndTimestamp(m.getTimesource())

	editMessage := &EditMessage{}

	editMessage.Text = request.Text
	editMessage.ContentType = request.ContentType
	editMessage.ChatId = message.ChatId
	editMessage.MessageId = request.ID.String()
	editMessage.Clock = clock

	err = m.applyEditMessage(&editMessage.EditMessage, message)
	if err != nil {
		return nil, err
	}

	encodedMessage, err := m.encodeChatEntity(chat, editMessage)
	if err != nil {
		return nil, err
	}

	rawMessage := common.RawMessage{
		LocalChatID:          chat.ID,
		Payload:              encodedMessage,
		MessageType:          protobuf.ApplicationMetadataMessage_EDIT_MESSAGE,
		SkipGroupMessageWrap: true,
		ResendAutomatically:  true,
	}
	_, err = m.dispatchMessage(ctx, rawMessage)
	if err != nil {
		return nil, err
	}

	if chat.LastMessage != nil && chat.LastMessage.ID == message.ID {
		chat.LastMessage = message
		err := m.saveChat(chat)
		if err != nil {
			return nil, err
		}
	}

	response := &MessengerResponse{}
	response.AddMessage(message)
	response.AddChat(chat)

	return response, nil
}

func (m *Messenger) CanDeleteMessageForEveryone(communityID string, publicKey *ecdsa.PublicKey) bool {
	if communityID != "" {
		community, err := m.communitiesManager.GetByIDString(communityID)
		if err != nil {
			m.logger.Error("failed to find community", zap.String("communityID", communityID), zap.Error(err))
			return false
		}
		return community.CanDeleteMessageForEveryone(publicKey)
	}
	return false
}

func (m *Messenger) DeleteMessageAndSend(ctx context.Context, messageID string) (*MessengerResponse, error) {
	message, err := m.persistence.MessageByID(messageID)
	if err != nil {
		return nil, err
	}

	// A valid added chat is required.
	chat, ok := m.allChats.Load(message.ChatId)
	if !ok {
		return nil, errors.New("Chat not found")
	}

	var canDeleteMessageForEveryone = false
	if message.From != common.PubkeyToHex(&m.identity.PublicKey) {
		if message.MessageType == protobuf.MessageType_COMMUNITY_CHAT {
			communityID := chat.CommunityID
			canDeleteMessageForEveryone = m.CanDeleteMessageForEveryone(communityID, &m.identity.PublicKey)
			if !canDeleteMessageForEveryone {
				return nil, ErrInvalidDeletePermission
			}
		}
		if !canDeleteMessageForEveryone {
			return nil, ErrInvalidEditOrDeleteAuthor
		}
	}

	// Only certain types of messages can be deleted
	if message.ContentType != protobuf.ChatMessage_TEXT_PLAIN &&
		message.ContentType != protobuf.ChatMessage_STICKER &&
		message.ContentType != protobuf.ChatMessage_EMOJI &&
		message.ContentType != protobuf.ChatMessage_IMAGE &&
		message.ContentType != protobuf.ChatMessage_AUDIO {
		return nil, ErrInvalidDeleteTypeAuthor
	}

	clock, _ := chat.NextClockAndTimestamp(m.getTimesource())

	deleteMessage := &DeleteMessage{}

	deleteMessage.ChatId = message.ChatId
	deleteMessage.MessageId = messageID
	deleteMessage.Clock = clock

	encodedMessage, err := m.encodeChatEntity(chat, deleteMessage)

	if err != nil {
		return nil, err
	}

	rawMessage := common.RawMessage{
		LocalChatID:          chat.ID,
		Payload:              encodedMessage,
		MessageType:          protobuf.ApplicationMetadataMessage_DELETE_MESSAGE,
		SkipGroupMessageWrap: true,
		ResendAutomatically:  true,
	}
	_, err = m.dispatchMessage(ctx, rawMessage)
	if err != nil {
		return nil, err
	}

	message.Deleted = true
	err = m.persistence.SaveMessages([]*common.Message{message})
	if err != nil {
		return nil, err
	}

	if chat.LastMessage != nil && chat.LastMessage.ID == message.ID {
		if err := m.updateLastMessage(chat); err != nil {
			return nil, err
		}
	}

	response := &MessengerResponse{}
	response.AddMessage(message)
	response.AddRemovedMessage(&RemovedMessage{MessageID: messageID, ChatID: chat.ID})
	response.AddChat(chat)

	return response, nil
}

func (m *Messenger) DeleteMessageForMeAndSync(ctx context.Context, chatID string, messageID string) (*MessengerResponse, error) {
	message, err := m.persistence.MessageByID(messageID)
	if err != nil {
		return nil, err
	}

	// A valid added chat is required.
	chat, ok := m.allChats.Load(chatID)
	if !ok {
		return nil, errors.New("Chat not found")
	}

	// Only certain types of messages can be deleted
	if message.ContentType != protobuf.ChatMessage_TEXT_PLAIN &&
		message.ContentType != protobuf.ChatMessage_STICKER &&
		message.ContentType != protobuf.ChatMessage_EMOJI &&
		message.ContentType != protobuf.ChatMessage_IMAGE &&
		message.ContentType != protobuf.ChatMessage_AUDIO {
		return nil, ErrInvalidDeleteTypeAuthor
	}

	message.DeletedForMe = true
	err = m.persistence.SaveMessages([]*common.Message{message})
	if err != nil {
		return nil, err
	}

	if chat.LastMessage != nil && chat.LastMessage.ID == message.ID {
		if err := m.updateLastMessage(chat); err != nil {
			return nil, err
		}
	}

	response := &MessengerResponse{}
	response.AddMessage(message)
	response.AddChat(chat)

	if m.hasPairedDevices() {
		clock, _ := chat.NextClockAndTimestamp(m.getTimesource())

		deletedForMeMessage := &DeleteForMeMessage{}

		deletedForMeMessage.MessageId = messageID
		deletedForMeMessage.Clock = clock

		encodedMessage, err := proto.Marshal(deletedForMeMessage.GetProtobuf())

		if err != nil {
			return response, err
		}

		rawMessage := common.RawMessage{
			LocalChatID:          chat.ID,
			Payload:              encodedMessage,
			MessageType:          protobuf.ApplicationMetadataMessage_SYNC_DELETE_FOR_ME_MESSAGE,
			SkipGroupMessageWrap: true,
			ResendAutomatically:  true,
		}
		_, err = m.dispatchMessage(ctx, rawMessage)
		if err != nil {
			return response, err
		}
	}

	return response, nil
}

func (m *Messenger) applyEditMessage(editMessage *protobuf.EditMessage, message *common.Message) error {
	if err := ValidateText(editMessage.Text); err != nil {
		return err
	}
	message.Text = editMessage.Text
	message.EditedAt = editMessage.Clock
	if editMessage.ContentType != protobuf.ChatMessage_UNKNOWN_CONTENT_TYPE {
		message.ContentType = editMessage.ContentType
	}

	// Save original message as edit so we can retrieve history
	if message.EditedAt == 0 {
		originalEdit := EditMessage{}
		originalEdit.Clock = message.Clock
		originalEdit.LocalChatID = message.LocalChatID
		originalEdit.MessageId = message.ID
		originalEdit.Text = message.Text
		originalEdit.ContentType = message.ContentType
		originalEdit.From = message.From
		err := m.persistence.SaveEdit(originalEdit)
		if err != nil {
			return err
		}
	}

	err := message.PrepareContent(common.PubkeyToHex(&m.identity.PublicKey))
	if err != nil {
		return err
	}

	return m.persistence.SaveMessages([]*common.Message{message})
}

func (m *Messenger) applyDeleteMessage(messageDeletes []*DeleteMessage, message *common.Message) error {
	if messageDeletes[0].From != message.From {
		return ErrInvalidEditOrDeleteAuthor
	}

	message.Deleted = true

	err := message.PrepareContent(common.PubkeyToHex(&m.identity.PublicKey))
	if err != nil {
		return err
	}

	err = m.persistence.SaveMessages([]*common.Message{message})
	if err != nil {
		return err
	}

	return nil
}

func (m *Messenger) applyDeleteForMeMessage(messageDeletes []*DeleteForMeMessage, message *common.Message) error {
	message.DeletedForMe = true

	err := message.PrepareContent(common.PubkeyToHex(&m.identity.PublicKey))
	if err != nil {
		return err
	}

	err = m.persistence.SaveMessages([]*common.Message{message})
	if err != nil {
		return err
	}

	return nil
}

// Code generated by gentrait; DO NOT EDIT.
package comp

import "github.com/appcrash/media/server/event"

// Message Trait Enum
const (
	MrCloneable                    = uint64(1) << 0
	MrChannelable                  = uint64(1) << 1
	MrPreComposer                  = uint64(1) << 2
	MrPostComposer                 = uint64(1) << 3
	MrInitializingNode             = uint64(1) << 4
	MrUnInitializingNode           = uint64(1) << 5
	UserMessageTraitEnumShiftBegin = 6
)

// Message Type Enum
const (
	MtRawByte = iota
	MtLinkPointRequest
	MtChannelLinkRequest
	MtUserMessageBegin
)

type RawByteConvertable interface {
	AsRawByteMessage() *RawByteMessage
}

type LinkPointRequestConvertable interface {
	AsLinkPointRequestMessage() *LinkPointRequestMessage
}

type ChannelLinkRequestConvertable interface {
	AsChannelLinkRequestMessage() *ChannelLinkRequestMessage
}

// --------Message Implementation Begin--------
func (m *RawByteMessage) Type() MessageType {
	return MtRawByte
}

func (m *RawByteMessage) AsEvent() *event.Event {
	return event.NewEvent(MtRawByte, m)
}

func (m *LinkPointRequestMessage) Type() MessageType {
	return MtLinkPointRequest
}

func (m *LinkPointRequestMessage) AsEvent() *event.Event {
	return event.NewEvent(MtLinkPointRequest, m)
}

func (m *ChannelLinkRequestMessage) Type() MessageType {
	return MtChannelLinkRequest
}

func (m *ChannelLinkRequestMessage) AsEvent() *event.Event {
	return event.NewEvent(MtChannelLinkRequest, m)
}

// --------Message Implementation End--------

func initMessageTraits() {
	AddMessageTrait(
		MT[RawByteMessage](MetaType[RawByteConvertable]()),
		MT[LinkPointRequestMessage](MetaType[LinkPointRequestConvertable]()),
		MT[ChannelLinkRequestMessage](MetaType[ChannelLinkRequestConvertable]()),
	)
}

func initMessageConversion() {
}

func InitMessage() {
	initMessageTraits()
	initMessageConversion()
}

// Copyright (c) 2023 ACOAUTO Team.
// All rights reserved.
//
// Detailed license information can be found in the LICENSE file.
//
// File: message.go Vehicle SOA protocal package.
//
// Author: Cheng.siyuan <chengsiyuan@acoinfo.com>

package protocol

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"

	"github.com/go-sylixos/go-vsoa/utils"
)

// Internal use in vsoa package
const (
	magic       byte = 0x9
	version     byte = 0x2
	magicNumber byte = magic | (version << 4)
	/* VSOA package length limit */
	HdrLength        = 10
	MaxMessageLength = 262144
	MaxDataLength    = (MaxMessageLength - HdrLength)
	/* VSOA quick channel length limit */
	MaxQMessageLength = 65507
)

var bufferPool = utils.NewLimitedPool(512, 4096)

func MagicNumber() byte {
	return magicNumber
}

var (
	ErrMessageTooLong = errors.New("message is too long")
	ErrMessageUnPad   = errors.New("raw message is not pad")
)

const (
	ServerError = "__vsoa_error__"
)

type QuickChannelFlag bool

const (
	ChannelQuick  QuickChannelFlag = true
	ChannelNormal QuickChannelFlag = false
)

type Message struct {
	*Header
	URL   []byte          // We call it URL but it more likely to be the server PATH for rpcx users
	Param json.RawMessage // It's []byte but have Marshal & UnMarshal method.
	Data  []byte
}

func smn(h *Header)

func NewMessage() *Message {
	header := Header([HdrLength]byte{})
	smn(&header)

	return &Message{
		Header: &header,
	}
}

// Header is the first part of Message and has fixed size.
// Format:
type Header [HdrLength]byte

func cmn(h Header) (ret bool)

func (h Header) checkMagicNumber() bool {
	return cmn(h)
}

func vs(h Header) (ret byte)

func (h Header) Version() byte {
	return vs(h)
}

func msgt(h Header) (ret byte)

// MessageType returns the message type.
func (h Header) MessageType() MessageType {
	return MessageType(msgt(h))
}

// MessageType returns the message type in text.
func (h Header) MessageTypeText() string {
	return TypeText(h.MessageType())
}

func smt(h *Header, mt MessageType)

// SetMessageType sets message type.
func (h *Header) SetMessageType(mt MessageType) {
	smt(h, mt)
}

func ir(h Header) (ret bool)

// IsRPC returns whether the message is RPC message.
func (h Header) IsRPC() bool {
	return ir(h)
}

func in(h Header) (ret bool)

// IsPingEcho returns whether the message is ping echo message.
func (h Header) IsNoop() bool {
	return in(h)
}

func ipe(h Header) (ret bool)

// IsPingEcho returns whether the message is ping echo message.
func (h Header) IsPingEcho() bool {
	return ipe(h)
}

func isi(h Header) (ret bool)

// IsServInfo returns whether the message is service info message.
func (h Header) IsServInfo() bool {
	return isi(h)
}

func iss(h Header) (ret bool)

// IsSubscribe returns whether the message is subscribe message.
func (h Header) IsSubscribe() bool {
	return iss(h)
}

func ius(h Header) (ret bool)

// IsUnSubscribe returns whether the message is unsubscribe info message.
func (h Header) IsUnSubscribe() bool {
	return ius(h)
}

func id(h Header) (ret bool)
func ip(h Header) (ret bool)

func (h Header) IsOneway() bool {
	if id(h) || ip(h) {
		return true
	}
	return false
}

// SetPingEcho sets the type flag to ping echo fast.
func spe(h *Header)

// SetPingEcho sets the type flag to ping echo fast.
func (h *Header) SetPingEcho() {
	spe(h)
}

func iR(h Header) (ret bool)

// IsReply returns whether the message is reply message.
func (h Header) IsReply() bool {
	return iR(h)
}

func sRt(h *Header)
func sRf(h *Header)

// SetReply sets the reply flag.
func (h *Header) SetReply(r bool) {
	if r {
		sRt(h)
	} else {
		sRf(h)
	}
}

func ivt(h Header) (ret bool)

// IsValidTunid checks if the Header is a valid Tunid.
//
// It returns a boolean indicating whether the Header is a valid Tunid.
func (h Header) IsValidTunid() bool {
	return ivt(h)
}

func svt(h *Header)

// SetValidTunid sets the valid tunid of the Header struct.
//
// It takes no parameters.
// There is no return type.
func (h *Header) SetValidTunid() {
	svt(h)
}

func mrm(h Header) (ret byte)

// MessageRpcMethod returns the rpc method.
// If it's not a RPC message in VSOA then it return 0xee
func (h Header) MessageRpcMethod() RpcMessageType {
	if h.IsRPC() {
		return RpcMessageType(mrm(h))
	} else {
		return NoneRpc
	}
}

// MessageRpcMethod returns the rpc method in text.
func (h Header) MessageRpcMethodText() string {
	return RpcMethodText(h.MessageRpcMethod())
}

func smrmg(h *Header)
func smrms(h *Header)

// MessageRpcMethod returns the rpc method.
// If it's not a RPC message in VSOA then it return 0xee
func (h *Header) SetMessageRpcMethod(t RpcMessageType) {
	if t == RpcMethodGet {
		smrmg(h)
	} else {
		smrms(h)
	}
}

func pl(h Header) (ret byte)

// Internel use padLen return the pad length for 4 byte pad for whole massage.
func (h Header) padLen() byte {
	return pl(h)
}

func spl(h *Header, pl byte)

// Internel use to pad the massage and fill the pad length flag automatic
func (h *Header) setPadLen(pl byte) {
	// clear PadLen flag
	spl(h, pl)
}

func st(h Header) (ret byte)

// StatusType returns the message status type.
func (h Header) StatusType() StatusType {
	return StatusType(st(h))
}

// StatusType returns the message status type in text.
func (h Header) StatusTypeText() string {
	return StatusText(h.StatusType())
}

func sst(h *Header, mt StatusType)

// SetStatusType sets message status type.
func (h *Header) SetStatusType(mt StatusType) {
	sst(h, mt)
}

func sn(h Header) (ret uint32)

// SeqNo returns sequence number of messages.
func (h Header) SeqNo() uint32 {
	return sn(h)
}

func ssn(h *Header, seq uint32)

// SetSeqNo sets  sequence number.
func (h *Header) SetSeqNo(seq uint32) {
	ssn(h, seq)
}

func tid(h Header) (ret uint16)

// TunID returns Tunnel id number for VSOA client.
func (h Header) TunID() uint16 {
	return tid(h)
}

func stid(h *Header, ti uint16)

// SetTunID set VSOA client Tunnel id number.
// It should be tunnel port number.
func (h *Header) SetTunId(ti uint16) {
	stid(h, ti)
}

// Clone clones from an message.
func (m Message) Clone() *Message {
	header := *m.Header
	c := NewMessage()
	c.Header = &header
	c.URL = m.URL
	c.Param = m.Param
	c.Data = m.Data
	return c
}

// CloneHeader clones header from an message.
func (m Message) CloneHeader() *Message {
	header := *m.Header
	c := NewMessage()
	c.Header = &header
	return c
}

// Encode encodes messages.
// Message can check it too long or not. Client/Server need to give the info of this Message is Quick or not.
func (m Message) Encode(q QuickChannelFlag) ([]byte, error) {
	rawMessage, err := m.encodeSlicePointer(q)
	return *rawMessage, err
}

// EncodeSlicePointer encodes messages as a byte slice pointer we can use pool to improve.
// Stream not using this.
func (m Message) encodeSlicePointer(q QuickChannelFlag) (*[]byte, error) {
	uL := len(m.URL)
	pL := len(m.Param)
	dL := len(m.Data)
	padL := 0

	// HdrL + urlL + paramL + dataL + url + param + data
	totalL := HdrLength + 10 + uL + pL + dL

	if totalL&3 != 0 {
		padL = 4 - (totalL & 3)
		totalL += padL
	}

	// We need to check Message len by type before send it,
	if !q {
		if totalL > MaxMessageLength {
			return bufferPool.Get(1), ErrMessageTooLong
		}
	} else {
		if totalL > MaxQMessageLength {
			return bufferPool.Get(1), ErrMessageTooLong
		}
	}

	urlStart := HdrLength + 10
	paramStart := urlStart + uL
	dataStart := paramStart + pL

	// Set pad length in header
	m.setPadLen(byte(padL))

	rawMessage := bufferPool.Get(totalL)
	copy(*rawMessage, m.Header[:])

	binary.BigEndian.PutUint16((*rawMessage)[10:12], uint16(uL))
	binary.BigEndian.PutUint32((*rawMessage)[12:16], uint32(pL))
	binary.BigEndian.PutUint32((*rawMessage)[16:20], uint32(dL))

	copy((*rawMessage)[urlStart:urlStart+uL], m.URL)
	copy((*rawMessage)[paramStart:paramStart+pL], m.Param)
	copy((*rawMessage)[dataStart:dataStart+dL], m.Data)

	return rawMessage, nil
}

// PutData puts the byte slice into pool.
func PutData(data *[]byte) {
	bufferPool.Put(data)
}

// WriteTo writes message to writers.
func (m Message) WriteTo(w io.Writer) (int64, error) {
	uL := len(m.URL)
	pL := len(m.Param)
	dL := len(m.Data)
	padL := 0

	// HdrL + urlL + paramL + dataL + url + param + data
	totalL := HdrLength + 10 + uL + pL + dL

	if totalL&3 != 0 {
		padL = 4 - (totalL & 3)
		totalL += padL
	}
	m.setPadLen(byte(padL))

	nn, err := w.Write(m.Header[:])
	n := int64(nn)
	if err != nil {
		return n, err
	}

	err = binary.Write(w, binary.BigEndian, uint16(uL))
	if err != nil {
		return n, err
	}

	err = binary.Write(w, binary.BigEndian, uint32(pL))
	if err != nil {
		return n, err
	}

	err = binary.Write(w, binary.BigEndian, uint32(dL))
	if err != nil {
		return n, err
	}

	_, err = w.Write(m.URL)
	if err != nil {
		return n, err
	}

	_, err = w.Write(m.Param)
	if err != nil {
		return n, err
	}

	_, err = w.Write(m.Data)
	if err != nil {
		return n, err
	}

	return int64(nn), err
}

// Read reads a message from r.
func Read(r io.Reader) (*Message, error) {
	msg := NewMessage()
	err := msg.Decode(r)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Decode decodes a message from reader.
func (m *Message) Decode(r io.Reader) error {
	defer func() {
		if err := recover(); err != nil {
			var errStack = make([]byte, 1024)
			n := runtime.Stack(errStack, true)
			log.Printf("panic in message decode: %v, stack: %s", err, errStack[:n])
		}
	}()

	// parse header
	_, err := io.ReadFull(r, m.Header[:1])
	if err != nil {
		return err
	}
	if !m.Header.checkMagicNumber() {
		// should never goes here
		return fmt.Errorf("wrong magic number: %v", m.Header[0])
	}

	_, err = io.ReadFull(r, m.Header[1:])
	if err != nil {
		return err
	}

	padL := int(m.padLen())

	// urlL
	urlLenData := make([]byte, 2)

	_, err = io.ReadFull(r, urlLenData)
	if err != nil {
		return err
	}
	uL := int(binary.BigEndian.Uint16(urlLenData))

	// paramL
	paramLenData := make([]byte, 4)
	_, err = io.ReadFull(r, paramLenData)
	if err != nil {
		return err
	}
	pL := int(binary.BigEndian.Uint32(paramLenData))

	// dataL
	dataLenData := make([]byte, 4)
	_, err = io.ReadFull(r, dataLenData)
	if err != nil {
		return err
	}
	dL := int(binary.BigEndian.Uint32(dataLenData))

	// HdrL + urlL + paramL + dataL + url + param + data
	totalL := HdrLength + 10 + uL + pL + dL + padL

	if totalL&3 != 0 {
		// this cause remain unread buffered date in the Reader cause unexpact behavior for next message.
		// TODO: avoid this behavior
		return ErrMessageUnPad
	}

	if totalL > MaxMessageLength {
		// this cause remain unread buffered date in the Reader cause unexpact behavior for next message.
		// TODO: avoid this behavior
		return ErrMessageTooLong
	}

	if cap(m.URL) >= uL { // reuse URL
		m.URL = m.URL[:uL]
	} else {
		m.URL = make([]byte, uL)
	}
	_, err = io.ReadFull(r, m.URL)
	if err != nil {
		return err
	}

	if cap(m.Param) >= pL { // reuse Param
		m.Param = m.Param[:pL]
	} else {
		m.Param = make([]byte, pL)
	}
	_, err = io.ReadFull(r, m.Param)
	if err != nil {
		return err
	}

	if cap(m.Data) >= dL { // reuse Data
		m.Data = m.Data[:dL]
	} else {
		m.Data = make([]byte, dL)
	}
	_, err = io.ReadFull(r, m.Data)
	if err != nil {
		return err
	}

	var zeros []byte
	if cap(zeros) >= padL {
		zeros = zeros[:padL]
	} else {
		zeros = make([]byte, padL)
	}

	_, err = io.ReadFull(r, zeros)
	if err != nil {
		return err
	}
	return err
}

// Reset clean data of this message but keep allocated data
func (m *Message) Reset() {
	resetHeader(m.Header)
	m.URL = m.URL[:0]
	m.Param = m.Param[:0]
	m.Data = m.Data[:0]
}

var (
	zeroHeaderArray Header
	zeroHeader      = zeroHeaderArray[1:]
)

func resetHeader(h *Header) {
	copy(h[1:], zeroHeader)
}

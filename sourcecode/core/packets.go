package core

//This file (will) contain all necessary scripts to decode/encode packets

import (
	"encoding/binary"
	"fmt"
	"gioms/utils"
	"math"
)

const (
	SEGMENT_BITS uint32 = 0x7F
	CONTINUE_BIT byte   = 0x80
)

type MinecraftPacket struct {
	Fields  map[string]interface{}
	packed  []byte
	RawInfo struct {
		State    string
		PacketID int
		Version  int16
	}
	index int
}

func NewClientboundMinecraftPacket() ClientboundMinecraftPacket {
	return ClientboundMinecraftPacket{
		&MinecraftPacket{
			Fields: make(map[string]interface{}),
			RawInfo: struct {
				State    string
				PacketID int
				Version  int16
			}{
				State:    "",
				PacketID: 0,
				Version:  utils.DEFAULT_PACKET_VERSION,
			},
		},
	}
}

func NewServerboundMinecraftPacket() ServerboundMinecraftPacket {
	return ServerboundMinecraftPacket{
		&MinecraftPacket{
			Fields: make(map[string]interface{}),
			RawInfo: struct {
				State    string
				PacketID int
				Version  int16
			}{
				State:    "",
				PacketID: 0,
				Version:  utils.DEFAULT_PACKET_VERSION,
			},
		},
	}
}

func (packet *MinecraftPacket) RawInfoFromPacked(packed []byte) {
	packet.packed = packed
	var _ = packet.UnpackVarInt()
	var id = packet.UnpackVarInt()
	packet.RawInfo.PacketID = id
}

func ClientboundPacketWithFields(state string, id int) (ClientboundMinecraftPacket, error) {
	//Protocol version 762 has a 1:1 mapping minecraft:gomc
	prot, err := GetOrLoadProtocol(762)
	if err != nil {
		CoreError(err)
		return ClientboundMinecraftPacket{}, err
	}
	return ClientboundMinecraftPacket{
		&MinecraftPacket{
			Fields: defaultVlaues(prot.ServerboundPackets[state][id].Fields),
			RawInfo: struct {
				State    string
				PacketID int
				Version  int16
			}{
				State:    state,
				PacketID: id,
				Version:  utils.DEFAULT_PACKET_VERSION,
			},
		},
	}, nil
}

func ServerboundPacketWithFields(state string, id int) (ServerboundMinecraftPacket, error) {
	//Protocol version 762 has a 1:1 mapping minecraft:gomc
	prot, err := GetOrLoadProtocol(762)
	if err != nil {
		CoreError(err)
		return ServerboundMinecraftPacket{}, err
	}
	return ServerboundMinecraftPacket{
		&MinecraftPacket{
			Fields: defaultVlaues(prot.ServerboundPackets[state][id].Fields),
			RawInfo: struct {
				State    string
				PacketID int
				Version  int16
			}{
				State:    state,
				PacketID: id,
				Version:  utils.DEFAULT_PACKET_VERSION,
			},
		},
	}, nil
}

// Functions to pack a signle type
func (packet *MinecraftPacket) PackByte(b byte) byte {
	packet.packed = append(packet.packed, b)
	return b
}

func (packet *MinecraftPacket) PackBoolean(b bool) []byte {
	if b {
		packet.packed = append(packet.packed, 0x01)
		return []byte{0x01}
	} else {
		packet.packed = append(packet.packed, 0x00)
		return []byte{0x00}
	}
}

func (packet *MinecraftPacket) PackShort(s int16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(s))
	packet.packed = append(packet.packed, b...)
	return b
}

func (packet *MinecraftPacket) PackUShort(us uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, us)
	packet.packed = append(packet.packed, b...)
	return b
}

func (packet *MinecraftPacket) PackInt(n int32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(n))
	packet.packed = append(packet.packed, b...)
	return b
}

func (packet *MinecraftPacket) PackLong(n int64) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint64(b, uint64(n))
	packet.packed = append(packet.packed, b...)
	return b
}

func (packet *MinecraftPacket) PackFloat(f float32) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, math.Float32bits(f))
	packet.packed = append(packet.packed, b...)
	return b
}

func (packet *MinecraftPacket) PackDouble(f float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(f))
	packet.packed = append(packet.packed, b...)
	return b
}

func (packet *MinecraftPacket) PackString(s string) []byte {
	l := packet.PackVarInt(32767*4 + 3)
	bin := []byte(s)
	packet.packed = append(packet.packed, bin...)
	return append(l, bin...)
}

func (packet *MinecraftPacket) PackVarInt(value int) []byte {
	packed := make([]byte, 5)
	for {
		if (value & int(^SEGMENT_BITS)) == 0 {
			packed = append(packed, byte(value))
			packet.packed = append(packet.packed, packed...)
			return packed
		}

		packed = append(packed, byte((value&int(SEGMENT_BITS)))|CONTINUE_BIT)

		value >>= 7
	}
}

func (packet *MinecraftPacket) PackVarLong(value int64) []byte {
	packed := make([]byte, 10)
	for {
		if (value & ^int64(SEGMENT_BITS)) == 0 {
			packed = append(packed, byte(value))
			packet.packed = append(packet.packed, packed...)
			return packed
		}

		packed = append(packed, byte((value&int64(SEGMENT_BITS)))|CONTINUE_BIT)

		value >>= 7
	}
}

// Functions to unpack a signle type
func (packet *MinecraftPacket) UnpackByte() byte {
	packet.index++
	return packet.packed[packet.index-1]
}

func (packet *MinecraftPacket) UnpackBoolean() bool {
	packet.index++
	if packet.packed[packet.index-1] != 0 {
		return true
	} else {
		return false
	}
}

// Updated functions to this point \/
func (packet *MinecraftPacket) UnpackShort() int16 {
	data := binary.LittleEndian.Uint16(packet.GetBytes(2))
	packet.MoveIndex(2)
	return int16(data)
}

func (packet *MinecraftPacket) UnpackUShort() uint16 {
	defer packet.MoveIndex(2)
	return binary.LittleEndian.Uint16(packet.GetBytes(2))
}

func (packet *MinecraftPacket) UnpackInt() int32 {
	data := binary.LittleEndian.Uint32(packet.GetBytes(4))
	packet.MoveIndex(4)
	return int32(data)
}

func (packet *MinecraftPacket) UnpackLong() int64 {
	data := binary.BigEndian.Uint64(packet.GetBytes(8))
	packet.MoveIndex(8)
	return int64(data)
}

func (packet *MinecraftPacket) UnpackFloat() float32 {
	bits := binary.LittleEndian.Uint32(packet.GetBytes(4))
	packet.MoveIndex(4)
	return math.Float32frombits(bits)
}

func (packet *MinecraftPacket) UnpackDouble() float64 {
	bits := binary.LittleEndian.Uint64(packet.GetBytes(8))
	packet.MoveIndex(8)
	return math.Float64frombits(bits)
}

func (packet *MinecraftPacket) UnpackVarInt() int {
	numRead := 0
	result := 0
	var read byte

	for {
		read = packet.packed[packet.index+int(numRead)]
		value := int(read & 0b01111111)
		result |= (value << (7 * numRead))

		numRead++
		if numRead > 5 {
			panic("VarInt is too big")
		}

		if (read & 0b10000000) == 0 {
			break
		}
	}
	packet.index += int(numRead)
	return result
}

func (packet *MinecraftPacket) UnpackString() string {
	var len = packet.UnpackVarInt()
	defer packet.MoveIndex(len)
	return string(packet.GetBytes(len))
}

func (packet *MinecraftPacket) MoveIndex(amount int) {
	packet.index += amount
}

func (packet *MinecraftPacket) GetBytes(amount int) []byte {
	return packet.packed[packet.index : packet.index+amount]
}

func (packet *MinecraftPacket) GetRemainingBytes() []byte {
	return packet.packed[packet.index:]
}

func (packet *MinecraftPacket) SetPacked(packed []byte) {
	packet.packed = packed
}

func defaultVlaues(types map[string]string) map[string]interface{} {
	var result = make(map[string]interface{})
	for key, value := range types {
		switch value {
		case "Boolean":
			result[key] = true
		case "Byte":
			result[key] = byte(0)
		case "UnsignedByte":
			result[key] = byte(0)
		case "Short":
			result[key] = int16(0)
		case "UnsignedShort":
			result[key] = uint16(0)
		case "Int":
			result[key] = int32(0)
		case "Long":
			result[key] = int64(0)
		case "Float":
			result[key] = float32(0)
		case "Double":
			result[key] = float64(0)
		case "String":
			result[key] = ""
		case "VarInt":
			result[key] = int(0)
		//case "VarLong":
		//	result[key] = true
		default:
			fmt.Printf("Could't find default value for type %v\n", value)
		}
	}
	return result
}

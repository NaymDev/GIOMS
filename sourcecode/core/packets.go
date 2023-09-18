package core
//This file (will) contain all necessary scripts to decode/encode packets

import (
	"encoding/binary"
	"bytes"
	"math"
	"fmt"
)

const (
	SEGMENT_BITS uint32 = 0x7F
	CONTINUE_BIT byte = 0x80
)


type MinecraftPacket struct {
	Fields map[string]interface{}
	packed []byte
	RawInfo struct{
		State string
		PacketID uint8
		Version int16
	}
	index uint8
}

func PacketWithFields(state String, id uint8, len uint, serverBound bool) (MinecraftPacket, error) {
	//Protocol version -1 is the GIOMC version which has a 1:1 mapping.
	prot, err := GetOrLoadProtocol(-1)
	if err != nil {
		CoreError(err)
		return (nil, err)
	}
	if serverBound {
		return MinecraftPacket {
			Fields: prot.ServerboundPackets[state][id].Fields,
			RawInfo: struct {
				State: state,
				PacketID: id,
				Version: -1
			}
		}
	} else {
		return MinecraftPacket {
			Fields: prot.ClientboundPackets[state][id].Fields,
			RawInfo: struct {
				State: state,
				PacketID: id,
				Version: -1
			}
		}
	}
}


func (packet *MinecraftPacket) Pack() []byte {
	protocol, err := GetOrLoadProtocol(packet.RawInfo.Version)
	if err != nil {
		CoreError(err)
		return nil
	}
	//fmt.Println(protocol.ServerboundPackets)
	packet_type := protocol.ServerboundPackets[packet.RawInfo.State][packet.RawInfo.PacketID]
	for field_name, field_type := range packet_type.Fields {
		fmt.Println("-------------------")
		fmt.Printf("%s:%s\n",field_name,field_type)
		
		val := packet.Fields[field_name]
		switch field_type {
		case "Boolean":
			packet.PackBoolean(val.(bool))
		case "Byte":
			packet.PackByte(val.(byte))
		case "UnsignedByte":
			packet.PackByte(val.(byte))
		case "Short":
			packet.PackShort(val.(int16))
		case "UnsignedShort":
			packet.PackUShort(val.(uint16))
		case "Int":
			packet.PackInt(val.(int32))
		case "Long":
			packet.PackLong(val.(int64))
		case "Float":
			packet.PackFloat(val.(float32))
		case "Double":
			packet.PackDouble(val.(float64))
		case "String":
			packet.PackString(val.(string))
		case "VarInt":
			packet.PackVarInt(val.(int))
		case "VarLong":
			packet.PackVarLong(val.(int64))
		default:
			fmt.Printf("Could't pack %v\n", val)
		}
	}
	
	return packet.packed
}

func (packet *MinecraftPacket) Unpack(packed []byte) {
	packet.packed = packed
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
	l := packet.PackVarInt(len(s))
	bin := []byte(s)
	packet.packed = append(packet.packed, bin...)
	return append(l, bin...)
}

func (packet *MinecraftPacket) PackVarInt(value int) []byte {
	value := uint32(inputValue)

	buffer := new(bytes.Buffer)

	for {
		temp := (byte)(value & 0b01111111)
		value >>= 7

		if value != 0 {
			temp |= 0b10000000
		}

		buffer.WriteByte(temp)

		if value == 0 {
			break
		}
	}

	packet.packed = append(packet.packed, buffer.Bytes()...)
	return buffer.Bytes()
}

func (packet *MinecraftPacket) PackVarLong(value int64) []byte {
	value := uint64(inputValue)

	buffer := new(bytes.Buffer)

	for {
		temp := (byte)(value & 0b01111111)
		value >>= 7

		if value != 0 {
			temp |= 0b10000000
		}

		buffer.WriteByte(temp)
		n++

		if value == 0 {
			break
		}
	}

	packet.packed = append(packet.packed, buffer.Bytes())
	return buffer.Bytes()
}


// Functions to unpack a signle type
func(packet *MinecraftPacket) UnpackByte() byte {
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

func (packet *MinecraftPacket) UnpackShort() int16 {
	data := binary.LittleEndian.Uint16(packet.packed[packet.index:packet.index+2])
	packet.index += 2
	return int16(data)
}

func (packet *MinecraftPacket) UnpackUShort() uint16 {
	packet.index += 2
	return binary.LittleEndian.Uint16(packet.index-2:packet.index)
}

func (packet *MinecraftPacket) UnpackInt() int32 {
	packet.index += 4
	data := binary.LittleEndian.Uint32(packet.packed[packet.index-4:packet.index])
	return int32(data)
}

func (packet *MinecraftPacket) UnpackLong() int64 {
	padket.index+=8
	data := binary.BigEndian.Uint64(packet.packed[packet.index-8:packet.index])
	return int64(data)
}



func (packet *MinecraftPacket) UnpackFloat() float32 {
	bits := binary.LittleEndian.Uint32(packet.packed[packet.index:packet.index+4])
	packet.index += 4
    return math.Float32frombits(bits)
}

func (packet *MinecraftPacket) UnpackDouble() float64 {
	bits := binary.LittleEndian.Uint64(packet.packed[packet.index:packet.index+8])
	packet.index += 8
    return math.Float64frombits(bits)
}


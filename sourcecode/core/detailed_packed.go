package core

import "fmt"

type ClientboundMinecraftPacket struct {
	*MinecraftPacket
}
type ServerboundMinecraftPacket struct {
	*MinecraftPacket
}

func (packet *ServerboundMinecraftPacket) Pack() []byte {
	protocol, err := GetOrLoadProtocol(packet.RawInfo.Version)
	if err != nil {
		CoreError(err)
		return nil
	}
	//fmt.Println(protocol.ServerboundPackets)
	packet_type := protocol.ServerboundPackets[packet.RawInfo.State][packet.RawInfo.PacketID]
	for field_name, field_type := range packet_type.Fields {
		fmt.Println("-------------------")
		fmt.Printf("%s:%s\n", field_name, field_type)

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

func (packet *ServerboundMinecraftPacket) Unpack(State string) error {
	packet.index = 0
	//Read Raw data
	packet.UnpackVarInt()
	packet.RawInfo.PacketID = packet.UnpackVarInt()

	//Load protocol
	protocol, err := GetOrLoadProtocol(packet.RawInfo.Version)
	if err != nil {
		return err
	}

	packet.Fields = make(map[string]interface{})
	packet_type := protocol.ServerboundPackets[State][packet.RawInfo.PacketID]
	for field_name, field_type := range packet_type.Fields {
		fmt.Println("-------------------")
		fmt.Printf("%s:%s\n", field_name, field_type)

		switch field_type {
		case "Boolean":
			packet.Fields[field_name] = packet.UnpackBoolean()
		case "Byte":
			packet.Fields[field_name] = packet.UnpackByte()
		case "UnsignedByte":
			packet.Fields[field_name] = packet.UnpackByte()
		case "Short":
			packet.Fields[field_name] = packet.UnpackShort()
		case "UnsignedShort":
			packet.Fields[field_name] = packet.UnpackUShort()
		case "Int":
			packet.Fields[field_name] = packet.UnpackInt()
		case "Long":
			packet.Fields[field_name] = packet.UnpackLong()
		case "Float":
			packet.Fields[field_name] = packet.UnpackFloat()
		case "Double":
			packet.Fields[field_name] = packet.UnpackDouble()
		case "String":
			packet.Fields[field_name] = packet.UnpackString()
		case "VarInt":
			packet.Fields[field_name] = packet.UnpackVarInt()
		//case "VarLong":
		//	packet.Fields[field_name] = packet.UnpackVar()
		default:
			fmt.Printf("Could't unpack %v\n", field_type)
		}
	}

	return nil
}

func (packet *ClientboundMinecraftPacket) Pack() []byte {
	protocol, err := GetOrLoadProtocol(packet.RawInfo.Version)
	if err != nil {
		CoreError(err)
		return nil
	}
	//fmt.Println(protocol.ServerboundPackets)
	packet_type := protocol.ClientboundPackets[packet.RawInfo.State][packet.RawInfo.PacketID]

	packet.PackVarInt(packet.RawInfo.PacketID)
	for field_name, field_type := range packet_type.Fields {
		//fmt.Println("-------------------")
		//fmt.Printf("%s:%s\n", field_name, field_type)

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

	packet.packed = append(packet.PackVarInt(len(packet.packed)), packet.packed...)
	return packet.packed
}

func (packet *ClientboundMinecraftPacket) Unpack(State string) error {
	packet.index = 0
	//Read Raw data
	packet.UnpackVarInt()
	packet.RawInfo.PacketID = packet.UnpackVarInt()

	//Load protocol
	protocol, err := GetOrLoadProtocol(packet.RawInfo.Version)
	if err != nil {
		return err
	}

	packet.Fields = make(map[string]interface{})
	packet_type := protocol.ClientboundPackets[State][packet.RawInfo.PacketID]
	for field_name, field_type := range packet_type.Fields {
		fmt.Println("-------------------")
		fmt.Printf("%s:%s\n", field_name, field_type)

		switch field_type {
		case "Boolean":
			packet.Fields[field_name] = packet.UnpackBoolean()
		case "Byte":
			packet.Fields[field_name] = packet.UnpackByte()
		case "UnsignedByte":
			packet.Fields[field_name] = packet.UnpackByte()
		case "Short":
			packet.Fields[field_name] = packet.UnpackShort()
		case "UnsignedShort":
			packet.Fields[field_name] = packet.UnpackUShort()
		case "Int":
			packet.Fields[field_name] = packet.UnpackInt()
		case "Long":
			packet.Fields[field_name] = packet.UnpackLong()
		case "Float":
			packet.Fields[field_name] = packet.UnpackFloat()
		case "Double":
			packet.Fields[field_name] = packet.UnpackDouble()
		case "String":
			packet.Fields[field_name] = packet.UnpackString()
		case "VarInt":
			packet.Fields[field_name] = packet.UnpackVarInt()
		//case "VarLong":
		//	packet.Fields[field_name] = packet.UnpackVar()
		default:
			fmt.Printf("Could't unpack %v\n", field_type)
		}
	}

	return nil
}

package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type data_type int8

const (
	Boolean data_type = iota
	Byte
	UnsignedByte
	Short
	UnsignedShort
	Int
	Long
	Float
	Double
	String
	Chat
)

var Protocols = make(map[int16]Protocol)

type Protocol struct {
	Version            int16           `json:"version"`
	Mapping            map[int16]int16 `json:"mapping"`
	ServerboundPackets map[string]map[int]struct {
		State  string
		MCID   int
		GOMCID int
		Fields map[string]string
	} `json:"serverbound_packets"`
	ClientboundPackets map[string]map[int]struct {
		State  string
		MCID   int
		GOMCID int
		Fields map[string]string
	} `json:"clientbound_packets"`
}

func GetOrLoadProtocol(version int16) (Protocol, error) {
	val, ok := Protocols[version]
	if ok {
		return val, nil
	} else {
		// Open our jsonFile
		jsonFile, err := os.Open(fmt.Sprintf("protocols/Protocol_V%v.json", version))
		// if we os.Open returns an error then handle it
		if err != nil {
			return Protocol{}, err
		}
		CoreInfo(fmt.Sprintf("Successfully Opened \"protocols/Protocol_V%v.json\"", version))
		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var protocol Protocol
		err = json.Unmarshal(byteValue, &protocol)
		if err != nil {
			CoreError(err)
		}

		Protocols[version] = protocol
		return Protocols[version], nil
	}
}

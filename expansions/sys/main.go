package main

import (
	"encoding/binary"
	"log"
	"math"
	"strings"
)

func Handler(argumentBytes []byte) []byte {

	// If the argument is less than
	if len(argumentBytes) < 2 {
		return []byte{0, 0, 0, 0}
	}

	switch argumentBytes[0] {
	case SysIFieldQuery:
		return handleFieldQuery(argumentBytes)
	case SysIDeviceQuery:

		moduleId := string(argumentBytes[1:])
		moduleId = strings.ReplaceAll(moduleId, "\x00", "")

		if val, ok := attributes["expansions"].(map[string]uint32)[moduleId]; ok {
			var bus_location []byte = make([]byte, 4)
			binary.LittleEndian.PutUint32(bus_location[:], val)

			return bus_location
		} else {
			return []byte{0, 0, 0, 0}
		}

	default:
		return []byte{0, 0, 0, 0}
	}

}

func SetAttribute(key string, value interface{}) {
	attributes[key] = value
}

func GetAttribute(key string) interface{} {
	return attributes[key]
}

func init() {
	attributes = make(map[string]interface{})

	// Set default attributes

	attributes["display_width"] = 640
	attributes["display_height"] = 480
	attributes["name"] = "sys32"
	attributes["memory"] = math.Pow(2, 16)

	attributes["expansions"] = make(map[string]uint32, 0)
	attributes["expansions"].(map[string]uint32)["goputer.sys"] = 0

	log.Println("Default goputer.sys loaded.")
}

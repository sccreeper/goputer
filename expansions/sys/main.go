package main

import (
	"encoding/binary"
	"log"
	"math"
	"strings"
)

func Handler(argument_bytes []byte) []byte {

	// If the argument is less than
	if len(argument_bytes) < 2 {
		return []byte{0, 0, 0, 0}
	}

	switch argument_bytes[0] {
	case SysIFieldQuery:
		return handle_field_query(argument_bytes)
	case SysIDeviceQuery:

		module_id := string(argument_bytes[1:])
		module_id = strings.ReplaceAll(module_id, "\x00", "")

		if val, ok := attributes["expansions"].(map[string]uint32)[module_id]; ok {
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

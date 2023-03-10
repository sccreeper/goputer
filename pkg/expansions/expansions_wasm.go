//go:build js && wasm

package expansions

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	SysLocation uint32 = 0x00
)

//Sys module constants

// Instructions
const (
	SysIFieldQuery  uint8 = 0x00
	SysIDeviceQuery uint8 = 0x01
)

// Instruction args
const (
	FieldName          uint8 = 0x00
	FieldDisplayWidth  uint8 = 0x01
	FieldDisplayHeight uint8 = 0x02
	FieldMemory        uint8 = 0x03
)

var attributes map[string]interface{}

func init() {

	attributes = make(map[string]interface{}, 0)

	attributes["name"] = "js32"
	attributes["display_width"] = 640
	attributes["display_height"] = 480
	attributes["memory"] = math.Pow(2, 16)

	attributes["expansions"] = make(map[string]uint32, 0)
	attributes["expansions"].(map[string]uint32)["goputer.sys"] = 0

	fmt.Println("Loaded goputer.sys WASM")

}

// Doesn't need to do anything for WASM just needs to be present.
func LoadExpansions() {}

func Interaction(location uint32, data []byte) []byte {

	fmt.Printf("Interaction on location %d\n", location)

	switch location {
	case SysLocation:
		return handle_sys(data)
	default:
		return []byte{0, 0, 0, 0}
	}

}

// Does a module exist at this location on the bus?
func ModuleExists(location uint32) bool {

	if location == SysLocation {
		return true
	} else {
		return false
	}

}

func GetAttribute(id string, name string) interface{} { return attributes[name] }

func SetAttribute(id string, name string, value interface{}) { attributes[name] = value }

func handle_sys(data []byte) []byte {

	switch data[0] {
	case SysIFieldQuery:
		switch data[1] {
		case FieldName:
			return []byte(attributes["name"].(string))
		case FieldDisplayWidth:
			b := make([]byte, 4)
			binary.LittleEndian.PutUint32(b, attributes["display_width"].(uint32))
			return b
		case FieldDisplayHeight:
			b := make([]byte, 4)
			binary.LittleEndian.PutUint32(b, attributes["display_height"].(uint32))
			return b
		case FieldMemory:
			b := make([]byte, 4)
			binary.LittleEndian.PutUint32(b, attributes["memory"].(uint32))
			return b
		default:
			return []byte{0, 0, 0, 0}
		}
	default:
		return []byte{0, 0, 0, 0}

	}

}

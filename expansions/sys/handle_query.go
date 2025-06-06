package main

import "encoding/binary"

func handleFieldQuery(argumentBytes []byte) []byte {

	switch argumentBytes[1] {
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

}

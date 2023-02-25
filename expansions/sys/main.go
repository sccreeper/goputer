package main

import "log"

func Handler(argument_bytes []byte) []byte {

	// If the argument is less than
	if len(argument_bytes) < 2 {
		return []byte{0, 0, 0, 0}
	}

	switch argument_bytes[0] {
	case SysIFieldQuery:
		return handle_field_query(argument_bytes)
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

	log.Println("Default goputer.sys loaded.")
}

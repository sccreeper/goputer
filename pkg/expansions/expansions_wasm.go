//go:build js && wasm

package expansions

// Doesn't need to do anything for WASM just needs to be present.
func LoadExpansions() {}

func Interaction(location uint32, data []byte) []byte { return nil }

// Does a module exist at this location on the bus?
func ModuleExists(location uint32) bool { return false }

package tests

import (
	"embed"
	_ "embed"
	"sccreeper/goputer/pkg/compiler"
	"testing"
)

//go:embed test_files/compile
var compileTestFiles embed.FS

var definitionTests map[string]string = map[string]string{
	"invalid_name" : "#def test\" 1234",
	"invalid_value_1" : "#def test 45.abc",
	"invalid_value_2" : "#def test \"hello world ",
	"invalid_value_3" : "#def test test:\"1234\"",
}

func TestDefinitionNamesShouldFail(t *testing.T)  {
	
	for k, v := range definitionTests {
		
		t.Run(k, func(t *testing.T) {

			failed := false

			codeString := "#label start\n"
			codeString += v

			p := compiler.Parser{
				CodeString:   codeString,
				FileName:     "main.gpasm",
				Verbose:      false,
				Imported:     false,
				ErrorHandler: func(errorType compiler.ErrorMessage, errorText string) { 
					 failed = true
				},
				FileReader:   func(path string) ([]byte, error) { return []byte(codeString), nil },
			}

			p.Parse()

			if failed != true {
				t.Fail()
			}

		})

	}

}
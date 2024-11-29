package tests

import (
	"testing"
	"myapp/utils"
)

// Function GetLimitAndOffset is tested here
type StrToIntTests struct { 
	input string 
	expected int
}

// test cases
var addStrToIntTests = []StrToIntTests{
	StrToIntTests{input: "1", expected: 1},
	StrToIntTests{input: "121", expected: 121},
	StrToIntTests{input: "-121", expected: -121},
}

// running tests for TestStrToInt
func TestStrToInt(t *testing.T) {
	for _, test := range addStrToIntTests{
		if output := utils.StrToInt(test.input); output != test.expected {
			t.Errorf("Expected %q but got %q", output, test.expected)
		}
	}
}

// Function GetLimitAndOffset is tested here
type GetLimitAndOffsetTests struct { 
	page string
	limit int
	offset int
}

// test cases
var getLimitAndOffsetTests = []GetLimitAndOffsetTests{
	GetLimitAndOffsetTests{page: "1", limit: 20, offset: 0},
	GetLimitAndOffsetTests{page: "2", limit: 20, offset: 20},
	GetLimitAndOffsetTests{page: "3", limit: 20, offset: 40},
	GetLimitAndOffsetTests{page: "0", limit: 20, offset: 0},
	GetLimitAndOffsetTests{page: "-1", limit: 20, offset: 0},
}

// running tests for TestGetLimitAndOffset
func TestGetLimitAndOffset(t *testing.T) {
	for _, test := range getLimitAndOffsetTests{
		if limit, offset := utils.GetLimitAndOffset(test.page); limit != test.limit || offset != test.offset {
			t.Errorf("Expected limit: %d and offset: %d, but got limit: %d and offset: %d", limit, offset, test.limit, test.offset)
		}
	}
}

// Function ApiResponder
type ApiResponderTests struct { 
	input utils.ApiResponderType
	expected map[string]interface{}
}

// test cases
var apiResponderTests = []ApiResponderTests{
	ApiResponderTests{input: utils.ApiResponderType{}, expected: map[string]interface{}{
		"data":"Successfully processed request","error":"null"},
	},
}

// running tests for TestApiResponder
func TestApiResponder(t *testing.T) {
	for _, test := range apiResponderTests{
		if output := utils.ApiResponder(test.input); output["data"] != test.expected["data"] || output["error"] != test.expected["error"] {
			t.Errorf("Expected %v but got %v", output, test.expected)
		}
	}
}

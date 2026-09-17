package main

import (
	"testing"
)

func Test_Main(t *testing.T) {

}

// BEGIN: TABLE DRIVEN TEST EXAMPLE
func Test_TableDriven(t *testing.T) {
	// Defining the columns of the table
	var tests = []struct {
		name     string // FIELD 1: TestName
		input    int    // FIELD 2: Input Value to be testesd
		expected string // FIELD 3: Desired Result when function is called
	}{
		// the table itself
		{"9 should be OK", 9, "OK"},
		{"3 should be OK", 3, "OK"},
		{"1 is not  OK", 1, "NOK"},
		{"0 should be OK", 0, "OK"},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultGot := FunctionToBeTested(tt.input)
			if resultGot != tt.expected {
				t.Errorf("tt.input[%d] resultGot %s, Expected %s", tt.input, resultGot, tt.expected)
			}
		})
	}
}

// END: TABLE DRIVEN TEST EXAMPLE

// BEGIN: FUNCTION TO BE TESTED - SHOULD EXIST ON YOUR CODE PACKET
func FunctionToBeTested(input int) string {
	var dataout string = "NONE"
	switch input {
	case 9:
		dataout = "OK"
	case 3:
		dataout = "OK"
	case 0:
		dataout = "OK"
	case 1:
		dataout = "NOK"
	default:
	}
	return dataout
}

//    END: FUNCTION TO BE TESTED - SHOULD EXIST ON YOUR CODE PACKET

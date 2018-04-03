package main

import (
	"encoding/json"
	"fmt"
)

func Example_attributeString() {
	tests := []string{
		`"text"`,
		"6",
		"9.3",
		"[9-12]",
		"{a: b}",
	}
	for _, test := range tests {
		s, err := attributeString(json.RawMessage(test))
		fmt.Println(s, err)
	}
	// Output:
	// text <nil>
	// 6 <nil>
	// 9.3 <nil>
	//  Value neither scalar nor text
	//  Value neither scalar nor text
}

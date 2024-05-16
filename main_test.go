package main

import (
    "testing"
    "reflect"
)

func TestParseResponse(t *testing.T) {
    type TestCase struct {
        input string
        expected []string
    }

    var cases []TestCase = []TestCase{
        {
            input: "There are 0 of a max of 20 players online: \n",
            expected: []string{},
        },
        {
            input: "There are 1 of a max of 20 players online: Steve\n",
            expected: []string{"Steve"},
        }, {
            input: "There are 3 of a max of 20 players online: Steve, Alex, Herobrine\n",
            expected: []string{"Steve", "Alex", "Herobrine"},
        },
    }

    for _, tc := range(cases) {
        result := parseResponse(tc.input)
        if !reflect.DeepEqual(result, tc.expected) {
            t.Fatalf("Expected %s, got %s.\n\nInput: %s", tc.expected, result, tc.input)
        }
    }
}

package main

import "testing"

func TestDecodeBencodeToken(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"i5e", 5},
		{"i-10e", -10},
		{"i0e", 0},
		{"10:hello12345", "hello12345"},
		{"0:", ""},
	}

	for _, test := range tests {
		input := test.input
		index := 0
		result, err := decodeBencodeToken(&input, &index)
		if err != nil {
			t.Errorf("Error decoding %s: %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("Expected %v, got %v for input %s", test.expected, result, test.input)
		}
	}
}

// TestDecodeInt tests the decodeInt function and checks the value of the index after decoding
// an integer from a bencoded string.
func TestDecodeInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		index    int
	}{
		{"i5e", 5, 2},
		{"i-10e", -10, 4},
		{"i0e", 0, 2},
	}

	for _, test := range tests {
		input := test.input
		index := 0
		result, err := decodeInt(&input, &index)
		if err != nil {
			t.Errorf("Error decoding %s: %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("Expected %v, got %v for input %s", test.expected, result, test.input)
		}
		if index != test.index {
			t.Errorf("Expected index %d, got %d for input %s", test.index, index, test.input)
		}
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		index    int
	}{
		{"4:spam", "spam", 5},
		{"0:", "", 1},
		{"10:hello12345", "hello12345", 12},
	}

	for _, test := range tests {
		input := test.input
		index := 0
		result, err := decodeString(&input, &index)
		if err != nil {
			t.Errorf("Error decoding %s: %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("Expected %v, got %v for input %s", test.expected, result, test.input)
		}
		if index != test.index {
			t.Errorf("Expected index %d, got %d for input %s", test.index, index, test.input)
		}
	}
}

func TestDecodeList(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
		index    int
	}{
		{"l4:spam4:eggse", []interface{}{"spam", "eggs"}, 13},
		{"l4:spami7ee", []interface{}{"spam", 7}, 10},
	}

	for _, test := range tests {
		input := test.input
		index := 0
		result, err := decodeList(&input, &index)
		if err != nil {
			t.Errorf("Error decoding %s: %v", test.input, err)
			continue
		}
		if len(result) != len(test.expected) {
			t.Errorf("Expected length %d, got %d for input %s", len(test.expected), len(result), test.input)
			continue
		}
		for i, v := range result {
			if v != test.expected[i] {
				t.Errorf("Expected %v, got %v for input %s at index %d", test.expected[i], v, test.input, i)
			}
		}
		if index != test.index {
			t.Errorf("Expected index %d, got %d for input %s", test.index, index, test.input)
		}
	}
}

func TestDecodeDict(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
		index    int
	}{
		{"d4:spami7ee", map[string]interface{}{"spam": 7}, 10},
		{"d3:foo3:bare", map[string]interface{}{"foo": "bar"}, 11},
	}

	for _, test := range tests {
		input := test.input
		index := 0
		result, err := decodeDict(&input, &index)
		if err != nil {
			t.Errorf("Error decoding %s: %v", test.input, err)
			continue
		}
		if len(result) != len(test.expected) {
			t.Errorf("Expected length %d, got %d for input %s", len(test.expected), len(result), test.input)
			continue
		}
		for k, v := range result {
			if v != test.expected[k] {
				t.Errorf("Expected %v, got %v for key %s in input %s", test.expected[k], v, k, test.input)
			}
		}
		if index != test.index {
			t.Errorf("Expected index %d, got %d for input %s", test.index, index, test.input)
		}
	}
}

func TestExtractTorrentInfoFromFile(t *testing.T) {
	tests := []struct {
		filePath string
		expected map[string]interface{}
	}{
		{"../resources/alice.torrent", map[string]interface{}{"info": map[string]interface{}{"name": "alice.txt", "length": 163783, "piece length": 16384}, "creation date": 1452468725091, "encoding": "UTF-8"}},
	}

	for _, test := range tests {
		result, err := extractTorrentInfoFromFile(test.filePath)
		if err != nil {
			t.Errorf("Error extracting torrent info from %s: %v", test.filePath, err)
			continue
		}
		creationDate, ok := result["creation date"].(int)
		if !ok {
			t.Errorf("Expected 'creation date' to be an int, got %T for file %s", result["creation date"], test.filePath)
			continue
		}
		if creationDate != test.expected["creation date"] {
			t.Errorf("Expected creation date %d, got %d for file %s", test.expected["creation date"], creationDate, test.filePath)
		}
		encoding, ok := result["encoding"].(string)
		if !ok {
			t.Errorf("Expected 'encoding' to be a string, got %T for file %s", result["encoding"], test.filePath)
			continue
		}
		if encoding != test.expected["encoding"] {
			t.Errorf("Expected encoding %s, got %s for file %s", test.expected["encoding"], encoding, test.filePath)
		}

		resultInfo, ok := result["info"].(map[string]interface{})
		if !ok {
			t.Errorf("Expected 'info' to be a map, got %T for file %s", result["info"], test.filePath)
			continue
		}
		for k, v := range test.expected["info"].(map[string]interface{}) {
			if resultInfo[k] != v {
				t.Errorf("Expected %v for key '%s', got %v for file %s", v, k, resultInfo[k], test.filePath)
			}
		}
	}
}

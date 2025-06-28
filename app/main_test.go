package main

import "testing"
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


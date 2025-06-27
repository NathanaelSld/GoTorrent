package main

import (
	"fmt"
	"os"
	"strconv"
)

func decodeBencodeToken(bencodedString *string, index *int) (interface{}, error) {
	if len((*bencodedString)) == 0 {
		return "", fmt.Errorf("empty bencoded string")
	}
	if (*bencodedString)[*index] == 'i' {
		parsedInt, err := decodeInt(bencodedString, index)
		if err != nil {
			return "", fmt.Errorf("invalid bencoded token: %s at %d", (*bencodedString), *index)
		}
		return parsedInt, nil
	}
	if '0' <= (*bencodedString)[*index] && (*bencodedString)[*index] <= '9' {
		parsedString, err := decodeString(bencodedString, index)
		if err != nil {
			return "", fmt.Errorf("invalid bencoded token: %s at %d", (*bencodedString), *index)
		}
		return parsedString, nil
	}

	if (*bencodedString)[*index] == 'l' {
		parsedList, err := decodeList(bencodedString, index)
		if err != nil {
			return nil, fmt.Errorf("invalid bencoded list: %s at %d", (*bencodedString), *index)
		}
		return parsedList, nil
	}
	if (*bencodedString)[*index] == 'd' {
		parsedDict, err := decodeDict(bencodedString, index)
		if err != nil {
			return nil, fmt.Errorf("invalid bencoded dictionary: %s at %d", (*bencodedString), *index)
		}
		return parsedDict, nil
	}
	return "", fmt.Errorf("invalid bencoded string: %s", (*bencodedString))
}

func decodeInt(bencodedString *string, index *int) (int, error) {
	*index++
	beginIndex := *index
	for (*bencodedString)[*index] != 'e' {
		*index++
	}
	endIndex := *index
	parsedInt, err := strconv.Atoi((*bencodedString)[beginIndex:endIndex])
	if err != nil {
		return 0, err
	}
	return parsedInt, nil
}

func decodeString(bencodedString *string, index *int) (string, error) {
	beginIndex := *index
	for (*bencodedString)[*index] != ':' && *index < len((*bencodedString)) {
		*index++
	}
	endIndex := *index
	strLength, err := strconv.Atoi((*bencodedString)[beginIndex:endIndex])
	if err != nil {
		return "", err
	}
	*index++
	parsedString := (*bencodedString)[*index : *index+strLength]
	*index += strLength - 1
	return parsedString, nil
}

func decodeList(bencodedString *string, index *int) ([]interface{}, error) {
	*index++
	parsedList := []interface{}{}
	for (*bencodedString)[*index] != 'e' {
		item, err := decodeBencodeToken(bencodedString, index)
		if err != nil {
			return nil, err
		}
		parsedList = append(parsedList, item)
		*index++
	}
	return parsedList, nil
}

func decodeDict(bencodedString *string, index *int) (map[string]interface{}, error) {
	*index++
	dict := map[string]interface{}{}
	for (*bencodedString)[*index] != 'e' {
		key, err := decodeString(bencodedString, index)
		if err != nil {
			return nil, err
		}
		*index++
		value, err := decodeBencodeToken(bencodedString, index)
		if err != nil {
			return nil, err
		}
		dict[key] = value
		*index++
	}
	return dict, nil
}

// function to pretty print object parsed from bencode
func prettyPrint(obj interface{}, indent string) {
	switch v := obj.(type) {
	case int:
		fmt.Printf("%s%d\n", indent, v)
	case string:
		fmt.Printf("%s%s\n", indent, v)
	case []interface{}:
		fmt.Println(indent + "[")
		for _, item := range v {
			prettyPrint(item, indent+"  ")
		}
		fmt.Println(indent + "]")
	case map[string]interface{}:
		fmt.Println(indent + "{")
		for key, value := range v {
			fmt.Printf("%s%s: ", indent+"  ", key)
			prettyPrint(value, indent+"  ")
		}
		fmt.Println(indent + "}")
	default:
		fmt.Printf("%s<unknown type>\n", indent)
	}
}

func extractTorrentInfoFromFile(filePath string) (map[string]interface{}, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	bencodedString := string(fileContent)
	index := 0
	parsedObject, err := decodeDict(&bencodedString, &index)
	if err != nil {
		return nil, err
	}
	return parsedObject, nil
}

func main() {
	path := "./resources/alice.torrent"
	//open file
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	bencodedString := string(fileContent)
	index := 0
	parsedObject, err := decodeBencodeToken(&bencodedString, &index)
	if err != nil {
		fmt.Printf("Error decoding bencoded string: %v\n", err)
		return
	}
	fmt.Println("Parsed Object:")
	prettyPrint(parsedObject, "")
	fmt.Println("Index after parsing:", index)
	if index != len(bencodedString) {
		fmt.Printf("Warning: Index %d does not match the length of the bencoded string %d\n", index, len(bencodedString))
	} else {
		fmt.Println("Index matches the length of the bencoded string.")
	}
}

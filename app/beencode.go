package main

import (
	"fmt"
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

func decodeTorrentFileString(bencodedString *string) (TorrentFile, error) {
	index := 0
	parsedTorrentFile, err := decodeDict(bencodedString, &index)
	if err != nil {
		return TorrentFile{}, err
	}
	torrent, err := torrentMapObjectToTorrenFileObject(parsedTorrentFile)
	if err != nil {
		return TorrentFile{}, err
	}
	return torrent, nil
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

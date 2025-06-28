package main

import (
	"fmt"
	"os"
)

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

func torrentMapObjectToTorrenFileObject(object map[string]interface{}) (TorrentFile, error) {
	torrent := TorrentFile{}
	infoMap, ok := object["info"].(map[string]interface{})
	if !ok {
		return TorrentFile{}, fmt.Errorf("info is not a map")
	}

	pieceLength, ok := infoMap["piece length"].(int)
	if !ok {
		return TorrentFile{}, fmt.Errorf("piece length is not an int")
	}

	pieces, ok := infoMap["pieces"].(string)
	if !ok {
		return TorrentFile{}, fmt.Errorf("pieces is not a string")
	}
	name, ok := infoMap["name"].(string)
	if !ok {
		return TorrentFile{}, fmt.Errorf("name is not a string")
	}
	var length *int
	var files *[]DictionnaryFile
	if lengthValue, ok := infoMap["length"].(int); ok {
		length = &lengthValue
	} else if filesInfoMap, ok := infoMap["files"].([]interface{}); ok {
		var fileList []DictionnaryFile
		for _, file := range filesInfoMap {
			fileMap, ok := file.(map[string]interface{})
			if !ok {
				return TorrentFile{}, fmt.Errorf("file is not a map")
			}
			lengthValue, ok := fileMap["length"].(int)
			if !ok {
				return TorrentFile{}, fmt.Errorf("file length is not an int")
			}
			path, ok := fileMap["path"].([]interface{})
			if !ok {
				return TorrentFile{}, fmt.Errorf("file path is not a list")
			}
			pathStrings := make([]string, len(path))
			for i, p := range path {
				pathString, ok := p.(string)
				if !ok {
					return TorrentFile{}, fmt.Errorf("file path element is not a string")
				}
				pathStrings[i] = pathString
			}
			fileList = append(fileList, DictionnaryFile{
				length: lengthValue,
				path:   pathStrings,
			})
		}
		files = &fileList
	}

	info := Info{
		pieceLength: pieceLength,
		pieces:      pieces,
		name:        name,
		length:      length,
		files:       files,
	}
	torrent.info = &info

	announce, ok := object["announce"].(string)
	if !ok {
		return TorrentFile{}, fmt.Errorf("announce is not a string")
	}
	torrent.announce = announce

	creationDate, ok := object["creation date"].(int)
	if !ok {
		torrent.creationDate = nil
	} else {
		torrent.creationDate = &creationDate
	}
	comment, ok := object["comment"].(string)
	if !ok {
		torrent.comment = nil
	} else {
		torrent.comment = &comment
	}
	createdBy, ok := object["created by"].(string)
	if !ok {
		torrent.createdBy = nil
	} else {
		torrent.createdBy = &createdBy
	}
	announcementList, ok := object["announce-list"].([]string)
	if !ok {
		torrent.announcementList = nil
	} else {
		torrent.announcementList = &announcementList
	}

	return torrent, nil
}

type TorrentFile struct {
	info     *Info
	announce string
	//Optionnal
	announcementList *[]string
	creationDate     *int
	comment          *string
	createdBy        *string
}

type Info struct {
	pieceLength int
	pieces      string
	name        string             //Name of the file or of the directory to store filed if mutli file
	length      *int               //Single file
	files       *[]DictionnaryFile //Multi file
}

type DictionnaryFile struct {
	length int
	path   []string
}

func main() {
	path := "../resources/debian.iso.torrent"
	//open file
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	bencodedString := string(fileContent)
	torrentFile, err := decodeTorrentFileString(&bencodedString)
	if err != nil {
		fmt.Printf("Error decoding torrent file: %v\n", err)
		return
	}
	fmt.Print(torrentFile)
	fmt.Print(torrentFile.info)
	fmt.Print(torrentFile.info.pieceLength)
	fmt.Print(torrentFile.info.length)
	
}

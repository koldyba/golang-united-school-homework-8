package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	flagId        = "id"
	flagOperation = "operation"
	flagItem      = "item"
	flagFileName  = "fileName"
)

type Usr struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Users []Usr
type Arguments map[string]string

func readWriteJson(fName string) (fileJson []byte, f *os.File, err error) {
	f, err = os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		e := fmt.Errorf("can't open %s. err caught: %w", fName, err)
		return nil, nil, e
	}

	fileJson, err = io.ReadAll(f)
	if err != nil {
		e := fmt.Errorf("can't open %s. err caught: %w", fName, err)
		return nil, nil, e
	}
	return
}

func isFlagPassed(args Arguments, name string) bool {
	if v, ok := args[name]; ok && v != "" {
		return true
	}
	return false
}

func flagErr(n string) error {
	return fmt.Errorf("-%s flag has to be specified", n)
}

func Perform(args Arguments, writer io.Writer) error {
	if !isFlagPassed(args, flagOperation) {
		return flagErr(flagOperation)
	}
	if !isFlagPassed(args, flagFileName) {
		return flagErr(flagFileName)
	}
	switch args[flagOperation] {
	case "add":
		if !isFlagPassed(args, flagItem) {
			return flagErr(flagItem)
		}
		return add(args[flagItem], args[flagFileName], writer)
	case "list":
		return list(args[flagFileName], writer)
	case "findById":
		if !isFlagPassed(args, flagId) {
			return flagErr(flagId)
		}
		return findById(args[flagId], args[flagFileName], writer)
	case "remove":
		if !isFlagPassed(args, flagId) {
			return flagErr(flagId)
		}
		return remove(args[flagId], args[flagFileName], writer)
	default:
		return fmt.Errorf("Operation %s not allowed!", args[flagOperation])
	}
}

func add(item, fName string, writer io.Writer) error {
	var newUsr Usr
	err := json.Unmarshal([]byte(item), &newUsr)
	if err != nil {
		return err
	}

	fileJson, f, err := readWriteJson(fName)
	if err != nil {
		return err
	}
	defer f.Close()

	var users Users
	if len(fileJson) != 0 {
		err = json.Unmarshal(fileJson, &users)
		if err != nil {
			return err
		}
	}
	for _, u := range users {
		if u.Id == newUsr.Id {
			msg := fmt.Sprintf("Item with id %s already exists", u.Id)
			writer.Write([]byte(msg))
			return nil
		}
	}

	users = append(users, newUsr)
	fileJson, err = json.Marshal(users)
	if err != nil {
		return err
	}

	f.Truncate(0)
	_, err = f.WriteAt(fileJson, 0)
	if err != nil {
		return err
	}

	return nil
}

func list(fName string, writer io.Writer) error {
	fileJson, f, err := readWriteJson(fName)
	if err != nil {
		return err
	}
	defer f.Close()
	writer.Write(fileJson)
	return nil
}

func findById(id, fName string, writer io.Writer) error {
	fileJson, f, err := readWriteJson(fName)
	if err != nil {
		return err
	}
	defer f.Close()

	var users Users
	err = json.Unmarshal(fileJson, &users)
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.Id == id {
			usr, err := json.Marshal(u)
			if err != nil {
				return err
			}
			writer.Write(usr)
		}
	}
	if err != nil {
		return err
	}
	writer.Write([]byte(""))
	return nil
}

func remove(id, fName string, writer io.Writer) error {
	fileJson, f, err := readWriteJson(fName)
	if err != nil {
		return err
	}
	defer f.Close()

	var users Users
	if len(fileJson) != 0 {
		err = json.Unmarshal(fileJson, &users)
		if err != nil {
			return err
		}
	}
	for k, u := range users {
		if u.Id == id {
			users[k] = users[(len(users) - 1)]
			break
		}
		if k == len(users)-1 {
			msg := fmt.Sprintf("Item with id %s not found", id)
			writer.Write([]byte(msg))
			return nil
		}
	}
	users = users[:(len(users) - 1)]
	fileJson, err = json.Marshal(users)
	if err != nil {
		return err
	}
	f.Truncate(0)
	_, err = f.WriteAt(fileJson, 0)
	if err != nil {
		return err
	}
	return nil
}

func parseArgs() Arguments {
	idPtr := flag.String(flagId, "", "pass id of user")
	operationPtr := flag.String(flagOperation, "", "pass operation")
	itemPtr := flag.String(flagItem, "", "pass item")
	fileNamePtr := flag.String(flagFileName, "", "pass filename")
	flag.Parse()

	args := make(Arguments)
	args[flagId] = *idPtr
	args[flagOperation] = *operationPtr
	args[flagItem] = *itemPtr
	args[flagFileName] = *fileNamePtr

	return args
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

type Arguments map[string]string

const filePermissionForProd = 0644

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {

	operation := args["operation"]
	fileName := args["fileName"]
	item := args["item"]
	id := args["id"]

	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}

	if operation != "add" && operation != "list" && operation != "findById" && operation != "remove" {
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	if operation == "list" {
		list(fileName, writer)
	} else if operation == "add" {
		if item == "" {
			return errors.New("-item flag has to be specified")
		} else {
			jsonString := string(unsafeRead(fileName))
			usersFromFile := getUsersFromString(jsonString)
			userFromArgs := getUserFromString(item)

			for _, v := range usersFromFile {
				if v.Id == userFromArgs.Id {
					err := fmt.Sprintf("Item with id %s already exists", userFromArgs.Id)
					writer.Write([]byte(err))
					return nil
				}
			}

			usersFromFile = append(usersFromFile, userFromArgs)

			sort.Slice(usersFromFile, func(i, j int) bool {
				return usersFromFile[i].Id < usersFromFile[j].Id
			})

			res, _ := json.Marshal(usersFromFile)

			writer.Write(res)
			fileWrite(fileName, res)
		}
	} else {
		if id == "" {
			return errors.New("-id flag has to be specified")
		} else {
			if operation == "findById" {
				jsonString := string(unsafeRead(fileName))
				usersFromFile := getUsersFromString(jsonString)

				for _, v := range usersFromFile {
					if v.Id == id {
						res, _ := json.Marshal(v)

						writer.Write(res)
						return nil
					}
				}

				writer.Write([]byte(""))
			} else if operation == "remove" {
				jsonString := string(unsafeRead(fileName))
				usersFromFile := getUsersFromString(jsonString)
				newUsers := make([]User, 0)
				for _, v := range usersFromFile {
					if v.Id != id {
						newUsers = append(newUsers, v)
					}
				}

				if len(usersFromFile) == len(newUsers) {
					writer.Write([]byte(fmt.Sprintf("Item with id %s not found", id)))
				} else {
					res, _ := json.Marshal(newUsers)

					//writer.Write(res)
					fileWrite(fileName, res)
					return nil
				}
			}
		}
	}

	return nil
}

func list(fileName string, writer io.Writer) {
	writer.Write(unsafeRead(fileName))
}

func getUsersFromString(jsonString string) []User {
	users := []User{}

	json.Unmarshal([]byte(jsonString), &users)

	return users
}

func getUserFromString(jsonString string) User {
	user := User{}

	json.Unmarshal([]byte(jsonString), &user)

	return user
}

func unsafeRead(fileName string) []byte {
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePermissionForProd)

	bytes, _ := ioutil.ReadAll(file)

	file.Close()

	return bytes
}

func fileWrite(fileName string, data []byte) {
	os.Remove(fileName)

	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePermissionForProd)

	file.Write(data)

	file.Close()
}

func main() {

	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {
	args := make(Arguments)
	operation := flag.String("operation", "", "help message for flag operation")
	item := flag.String("item", "", "help message for flag item")
	fileName := flag.String("fileName", "", "help message for flag fileName")
	id := flag.String("id", "", "help message for flag id")

	flag.Parse()

	args["operation"] = *operation
	args["item"] = *item
	args["fileName"] = *fileName
	args["id"] = *id

	return args
}

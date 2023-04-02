package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/iancoleman/orderedmap"
)

// https://stackoverflow.com/a/15323988/7543474
// Licensed under CC BY-SA 4.0
func string_in_slice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func modify_cmd_output(cmd string, cmd_output string) string {
	if cmd == "Byte Order" {
		if cmd_output == "1234" {
			return "Little Endian"
		} else if cmd_output == "4321" {
			return "Big Endian"
		}
	} else if cmd == "CPU MHz" || cmd == "CPU max MHz" || cmd == "CPU min MHz" {
		mhz, err := strconv.Atoi(cmd_output)
		if err != nil {
			log.Println("Error converting CPU MHz to int: ", err)
		}
		return strconv.Itoa(mhz / 1000000)
	} else if cmd == "L1d cache" || cmd == "L1i cache" || cmd == "L2 cache" || cmd == "L3 cache" {
		cache_size, err := strconv.Atoi(cmd_output)
		if err != nil {
			log.Println("Error converting cache size to int: ", err)
		}
		return strconv.Itoa(cache_size/1024) + "K"
	}
	return cmd_output
}

func decode_json_file(filename string) *orderedmap.OrderedMap {
	json_file, err := os.Open("cmds.json")
	if err != nil {
		log.Println("Error opening JSON file: ", err)
		return nil
	}
	defer json_file.Close()

	decoder := json.NewDecoder(json_file)
	commands := orderedmap.New()
	decoder.UseNumber()
	err = decoder.Decode(&commands)
	if err != nil {
		log.Println("Error decoding JSON file: ", err)
		return nil
	}
	return commands
}

func main() {
	commands := decode_json_file("cmds.json")
	if commands == nil {
		return
	}

	for _, cmd_name := range commands.Keys() {
		cmd, _ := commands.Get(cmd_name)
		output, err := []byte{}, error(nil)

		if strings.Contains(cmd.(string), "|") {
			output, err = exec.Command("sh", "-c", cmd.(string)).Output()
		} else {
			sub_commands := strings.Split(cmd.(string), " ")
			output, err = exec.Command(sub_commands[0], sub_commands[1:]...).Output()
		}
		if err != nil {
			log.Println("Error executing command ", cmd, " : ", err)
		}

		cmd_name = strings.TrimSpace(cmd_name)
		output_str := strings.TrimSpace(strings.TrimRight(string(output), "\n"))
		output_str = modify_cmd_output(cmd_name, output_str)
		fmt.Printf("%s: %s\n", cmd_name, output_str)
	}
}

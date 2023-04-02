package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/iancoleman/orderedmap"
)

const Cmds = `{
    "Architecture": "uname -m",
    "Byte Order": "sysctl -n hw.byteorder",
    "CPU(s)": "sysctl -n hw.ncpu",
    "On-line CPU(s)": "sysctl -n hw.activecpu",
    "Thread(s) per core": "echo \"scale=0; $(sysctl -n machdep.cpu.thread_count) / $(sysctl -n machdep.cpu.core_count)\" | bc",
    "Core(s) per socket": "echo \"scale=0; $(sysctl -n machdep.cpu.core_count) / $(sysctl -n hw.packages)\" | bc",
    "Socket(s)": "sysctl -n hw.packages",
    "Vendor ID": "sysctl -n machdep.cpu.vendor",
    "CPU family": "sysctl -n machdep.cpu.family",
    "CPU Model": "sysctl -n machdep.cpu.model",
    "Model name": "sysctl -n machdep.cpu.brand_string",
    "Stepping": "sysctl -n machdep.cpu.stepping",
    "CPU MHz": "sysctl -n hw.cpufrequency",
    "CPU max MHz": "sysctl -n hw.cpufrequency_max",
    "CPU min MHz": "sysctl -n hw.cpufrequency_min",
    "Hyper-Threading Technology": "system_profiler SPHardwareDataType | grep \"Hyper-Threading Technology\" | cut -d: -f2-",
    "L1d cache": "sysctl -n hw.l1dcachesize",
    "L1i cache": "sysctl -n hw.l1icachesize",
    "L2 cache": "sysctl -n hw.l2cachesize",
    "L3 cache": "sysctl -n hw.l3cachesize",
    "Flags": "sysctl -n machdep.cpu.features"
}
`

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

func decode_json_file() *orderedmap.OrderedMap {
	decoder := json.NewDecoder(strings.NewReader(Cmds))
	commands := orderedmap.New()
	decoder.UseNumber()
	err := decoder.Decode(&commands)
	if err != nil {
		log.Println("Error decoding JSON file: ", err)
		return nil
	}
	return commands
}

func main() {
	commands := decode_json_file()
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

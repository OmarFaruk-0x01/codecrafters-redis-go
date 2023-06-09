package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	CommandsArgMap = map[string]int{
		"PING": 0,
		"ECHO": 1,
		"SET":  2,
		"GET":  1,
	}
	CommandOptSlice = map[string]bool{
		"PX": true,
	}
	errExtraArgumentError   = errors.New("extra arguments not supported")
	errInvalidArgumentError = errors.New("arguments not supported")
	errInsufficientArgError = errors.New("insufficient arguments")
	errInvalidCmdError      = errors.New("invalid command error")
)

type Command struct {
	cmd       string
	arguments []string
	options   map[string]string
}

func parser(data string) ([]string, error) {
	cmds := make([]string, 0)
	p := regexp.MustCompile("\\$[0-9]+\r\n([^\r\n]+)\r\n")
	matches := p.FindAllStringSubmatch(data, -1)
	for _, m := range matches {
		cmds = append(cmds, m[1])
	}
	fmt.Printf("%v\n", cmds)
	return cmds, nil //errors.New("Unable to parse command.")
}

func commandParser(commands []string) (*Command, error) {
	if len(commands) == 0 {
		return nil, errInvalidCmdError
	}

	command := &Command{
		options: make(map[string]string),
	}
	command.cmd = commands[0]

	switch strings.ToUpper(command.cmd) {
	case "PING":
		if len(commands) > 1 {
			return nil, errExtraArgumentError
		}
		return command, nil
	case "ECHO":
		if len(commands) > 2 {
			return nil, errExtraArgumentError
		}
		command.arguments = append(command.arguments, commands[1:]...)
		fmt.Println(command)
		return command, nil
	case "SET":
		if len(commands) > 5 {
			return nil, errExtraArgumentError
		}
		if len(commands) < 3 {
			return nil, errInsufficientArgError
		}
		command.arguments = append(command.arguments, commands[1:3]...)

		if len(commands) > 3 && len(commands) < 5 {
			return nil, errInsufficientArgError
		} else if len(commands) > 3 && len(commands) == 5 {
			key := strings.ToUpper(commands[3])
			if _, ok := CommandOptSlice[key]; ok {
				command.options[key] = commands[4]
			}
		}

		return command, nil
	case "GET":
		if len(commands) > 2 {
			return nil, errExtraArgumentError
		}
		if len(commands) < 2 {
			return nil, errInsufficientArgError
		}
		command.arguments = append(command.arguments, commands[1:]...)
		return command, nil
	default:
		return nil, errInvalidCmdError
	}

}

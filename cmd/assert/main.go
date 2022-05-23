package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jojomi/assert"
	"github.com/jojomi/assert/exit"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type HandlerMap map[string]func([]string) (error, exit.Code)

func main() {
	args := os.Args[1:]

	handlerMap := HandlerMap{
		"command-exists":               assert.CommandExists,
		"day-of-month":                 assert.DayOfMonth,
		"dir-exists":                   assert.DirExists,
		"file-exists":                  assert.FileExists,
		"mounted":                      assert.Mounted,
		"non-empty-dir":                assert.NonEmptyDir,
		"ssh-reachable":                assert.SSHReachable,
		"ssh-reachable-noninteractive": assert.SSHReachableNonInteractive,
		"time-after":                   assert.TimeAfter,
		"time-before":                  assert.TimeBefore,
		"weekday":                      assert.Weekday,
	}

	if len(args) == 0 {
		exitWith(fmt.Errorf("no command given, existing handlers: %s", getHandlerList(handlerMap)), exit.CodeErrorFinal)
	}
	cmd := args[0]
	f, ok := handlerMap[cmd]
	if !ok {
		exitWith(fmt.Errorf("command not found: %s, existing handlers: %s", cmd, getHandlerList(handlerMap)), exit.CodeErrorFinal)
	}

	exitWith(f(args[1:]))
}

func getHandlerList(handlerMap HandlerMap) string {
	return strings.Join(getMapKeys(handlerMap), ", ")
}

func exitWith(err error, exitCode exit.Code) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
	}
	os.Exit(int(exitCode))
}

func getMapKeys(input map[string]func([]string) (error, exit.Code)) []string {
	var result = make([]string, 0)
	for k := range input {
		result = append(result, k)
	}
	return result
}

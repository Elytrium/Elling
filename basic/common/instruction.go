package common

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
)

type Instructions map[string]interface{}

var r *regexp.Regexp

func init() {
	r = regexp.MustCompile(".*/(.*)\\.yml")
}

func ReadInstructions(p string, t reflect.Type) Instructions {
	var instructions Instructions
	instructions = make(Instructions)

	instructionFiles, err := filepath.Glob(p + "/*.yml")

	if err != nil {
		log.Err(err)
	}

	for _, file := range instructionFiles {
		instruction := reflect.New(t).Interface()
		content, _ := os.ReadFile(file)
		_ = yaml.Unmarshal(content, instruction)
		instructions[r.FindStringSubmatch(file)[0]] = instruction
	}

	return instructions
}

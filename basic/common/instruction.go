package common

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
)

type Instructions map[string]interface{}

func ReadInstructions(p string, t reflect.Type) Instructions {
	var instructions Instructions
	instructions = make(Instructions)

	instructionFiles, err := filepath.Glob(p + "/*.yml")
	nameLength := len(p) + 1

	if err != nil {
		log.Error().Err(err).Send()
	}

	for _, file := range instructionFiles {
		instruction := reflect.New(t).Interface()
		content, _ := os.ReadFile(file)
		_ = yaml.Unmarshal(content, instruction)
		instructions[file[nameLength:len(file)-4]] = instruction
	}

	return instructions
}

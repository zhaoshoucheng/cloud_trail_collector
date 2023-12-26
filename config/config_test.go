package config

import (
	"fmt"
	"testing"
)

func TestNewConfig(t *testing.T) {
	file := "../config.toml"
	configTest, err := NewConfig(file)
	if err != nil {
		panic(err)
	}
	for _, input := range configTest.Inputs {
		fmt.Println(input)
	}
	for _, output := range configTest.Outputs {
		fmt.Println(output)
	}
	fmt.Println("done!")
}

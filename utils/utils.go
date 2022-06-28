package utils

import (
	"log"
	"strings"
)


func PortSet(port string) string {
  return ":" + port
}

func Error_check(err interface{}) {
	if err != nil {
		log.Panic(err)
	}
}

func Splitter(original string, split_string string) string {
	return strings.Split(original, split_string)[1]
}

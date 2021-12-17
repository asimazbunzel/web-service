package io

import (
	"fmt"
	"github.com/TwiN/go-color"
)

func LogInfo(reference, data string) {
	fmt.Println(color.Ize(color.Green, "["+reference+"] --INFO-- "+data))
}

func LogError(reference, data string) {
	fmt.Println(color.Ize(color.Red, "["+reference+"] --ERROR-- "+data))
}

func LogDebug(reference, data string) {
	fmt.Println(color.Ize(color.Yellow, "["+reference+"] --DEBUG-- "+data))
}

package file

import (
	"fmt"
	"log"
	"os"
)

func GetFilesFromSrc() {
	dir, err := os.Open("test-site/src/sections")
	if err != nil {
		fmt.Println("Something went wrong: " + err.Error())
		os.Exit(1)
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Something went wrong: " + err.Error())
		os.Exit(1)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func ReadFile() {
	content, err := os.ReadFile("test-site/src/test-file.md")
	if err != nil {
		fmt.Println("Something went wrong: " + err.Error())
		os.Exit(1)
	}
	fmt.Println(string(content))
}

func WriteFile() {
	testData := "#Hello\n\nthis is markdown"
	if err := os.WriteFile("test-site/src/test-write.md", []byte(testData), 0666); err != nil {
		log.Fatal("Could not write file.")
	}
}

package main

import (
	"fmt"
	"os"
)

func checkFileStatus(source string, target string) error {

	if _, err := os.Stat(source); os.IsNotExist(err) {
		return fmt.Errorf("সোর্স ফাইল '%s' খুঁজে পাওয়া যায়নি", source)
	}
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("লক্ষ্য ফাইল '%s' অলরেডি বিদ্যমান", target)
	}

	return nil
}

func main() {
	sourceFile := "1.txt"
	targetFile := "2.txt"

	if err := checkFileStatus(sourceFile, targetFile); err != nil {
		fmt.Println("এরর:", err)
		return
	}
	if err := os.Link(sourceFile, targetFile); err != nil {
		fmt.Printf("হার্ড লিঙ্ক তৈরি করতে সমস্যা হয়েছে: %v\n", err)
		return
	}

	fmt.Println("সফলভাবে হার্ড লিঙ্ক তৈরি হয়েছে!")
}

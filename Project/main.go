package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	sourceDir := "Configuration"
	backupDir := "Backup_Configuration"

	// ১. ব্যাকআপ ফোল্ডার তৈরি করা (যদি না থাকে)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		fmt.Printf("ফোল্ডার তৈরি করা সম্ভব হয়নি: %v\n", err)
		return
	}

	// ২. সোর্স ডিরেক্টরি স্ক্যান করা
	err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// ডিরেক্টরি হলে কিছু করব না, শুধু ফাইল হলে কাজ করব
		if !d.IsDir() {
			// ব্যাকআপ পাথের নাম তৈরি করা
			relPath, _ := filepath.Rel(sourceDir, path)
			targetPath := filepath.Join(backupDir, relPath)

			// ৩. হার্ড লিঙ্ক তৈরি করা
			// যদি আগে থেকে লিঙ্ক না থাকে, তবে তৈরি হবে
			err := os.Link(path, targetPath)
			if err != nil {
				if os.IsExist(err) {
					fmt.Printf("ফাইলটি আগেই ব্যাকআপ করা আছে: %s\n", relPath)
				} else {
					fmt.Printf("লিঙ্ক করতে সমস্যা: %s -> %v\n", relPath, err)
				}
			} else {
				fmt.Printf("সফল ব্যাকআপ: %s\n", relPath)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("ব্যাকআপ প্রক্রিয়ায় সমস্যা: %v\n", err)
	} else {
		fmt.Println("\nসব ফাইল সফলভাবে লিঙ্ক করা হয়েছে!")
	}
}

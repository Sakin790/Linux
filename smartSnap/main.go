package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// ১. লগ ফাইল সেটআপ
	logFile, err := os.OpenFile("backup.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("লগ ফাইল খুলতে ব্যর্থ:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("সিস্টেম শুরু হয়েছে। ব্যাকআপ হবে প্রতি ১ মিনিট পর পর।")

	// ২. Ticker ব্যবহার করে for range লুপ
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// সরাসরি লুপ শুরু
	for range ticker.C {
		runBackup()
	}
}

func runBackup() {
	sourceDir := "Configuration"
	backupDir := "Backup_Configuration"

	_ = os.MkdirAll(backupDir, 0755)

	err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		relPath, _ := filepath.Rel(sourceDir, path)
		targetPath := filepath.Join(backupDir, relPath)

		// হার্ড লিঙ্ক তৈরির চেষ্টা
		err = os.Link(path, targetPath)
		if err != nil && !os.IsExist(err) {
			log.Printf("এরর: %s লিঙ্ক করতে সমস্যা: %v", relPath, err)
		}
		return nil
	})

	if err != nil {
		log.Printf("ব্যাকআপ প্রক্রিয়ায় এরর: %v", err)
	} else {
		log.Println("ব্যাকআপ সম্পন্ন হয়েছে।")
	}
}

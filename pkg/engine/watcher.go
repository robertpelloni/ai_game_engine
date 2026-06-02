package engine

import (
	"log"
	"os"
	"time"
)

func WatchSchema(filepath string, callback func([]byte)) {
	var lastModTime time.Time

	for {
		info, err := os.Stat(filepath)
		if err != nil {
			log.Printf("Error stating file %s: %v", filepath, err)
			time.Sleep(1 * time.Second)
			continue
		}

		if info.ModTime().After(lastModTime) {
			lastModTime = info.ModTime()
			data, err := os.ReadFile(filepath)
			if err != nil {
				log.Printf("Error reading file %s: %v", filepath, err)
			} else {
				log.Printf("File %s changed, triggering reload", filepath)
				callback(data)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

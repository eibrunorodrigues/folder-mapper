package mapper

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Worker(fatalErrors *chan string, routineId int) {
	folderArg := os.Args[1]
	for {
		var fileList []string

		if err := filepath.Walk(folderArg, getFolderFiles(&fileList)); err != nil {
			*fatalErrors <- err.Error()
		}

		log.Println("Everything looks okay... Running routineId:" + strconv.Itoa(routineId))

		for _, loopFile := range fileList {
			// processFile(loopFile)
			// fmt.Println(path.Dir(loopFile))
		}

		time.Sleep(10)
	}
}

func processFile(fileName string) {
	lockFileName := fileName + ".lock"
	if _, err := os.Stat(lockFileName); os.IsNotExist(err) {
		lockWasGenerated := false

		defer func(lockWasGenerated *bool) {
			if !*lockWasGenerated {
				if err := os.Remove(lockFileName); err != nil {
					log.Println("Wasn't able to remove lock file from " + fileName)
				}
			}
		}(&lockWasGenerated)

		if _, err := os.OpenFile(lockFileName, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
			return
		}

		lockWasGenerated = true

	}
}

func getFolderFiles(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.Contains(info.Name(), ".json") {
			*files = append(*files, path)
		}
		return nil
	}
}

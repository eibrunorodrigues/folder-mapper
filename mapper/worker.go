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

		log.Println("Everything is ok... Running routineId " + strconv.Itoa(routineId))

		for _, loopFile := range fileList {
			processFile(loopFile)
		}

		time.Sleep(10)
	}
}

func processFile(fileName string) {

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

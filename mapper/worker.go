package mapper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	models "github.com/eibrunorodrigues/folder-mapper/mapper/models"
)

// Worker is exported when
func Worker(fatalErrors *chan string, routineID int) {
	folderArg := os.Args[1]
	// for {
	var fileList []models.FileStructure

	if err := filepath.Walk(folderArg, getFolderFiles(&fileList)); err != nil {
		*fatalErrors <- err.Error()
	}

	// log.Println("Everything looks okay... Running routineID:" + strconv.Itoa(routineID))

	for _, loopFile := range fileList {
		processFile(loopFile)
		// fmt.Println(path.Dir(loopFile))
	}

	time.Sleep(10)
	// }
}

func processFile(file models.FileStructure) {
	if _, err := os.Stat(file.LockFile); os.IsNotExist(err) {
		lockWasGenerated := false

		defer func(lockWasGenerated *bool) {
			if *lockWasGenerated {
				if err := os.Remove(file.LockFile); err != nil {
					log.Println("Wasn't able to remove lock file from " + file.Name)
				} else {
					log.Println(file.LockFile + " removed.")
				}

			}
		}(&lockWasGenerated)

		generateFolders(file, &lockWasGenerated)

		body, err := ioutil.ReadFile(file.Path)

		if err != nil {
			log.Println(err.Error())
			return
		}

		log.Println("Reading file " + file.Name)
		var fileInfos []models.FileInformation
		if err := UnmarshalBody([]byte(body), &fileInfos); err != nil {
			log.Println(err.Error())
		}

		log.Println(fileInfos)
	}
}

func UnmarshalBody(b []byte, fileInfos *[]models.FileInformation) error {
	var stuff []map[string]interface{}

	err := json.Unmarshal(b, &stuff)
	if err != nil {
		return err
	}

	for _, valueList := range stuff {
		fileInfoModel := models.FileInformation{}
		for key, value := range valueList {
			switch key {
			case "operation":
				fileInfoModel.Operation = value.(string)
			case "timestamp":
				fileInfoModel.Timestamp = value.(string)
			default:
				fileInfoModel.Data = value
			}
		}
		*fileInfos = append(*fileInfos, fileInfoModel)
	}
	return nil
}

func generateFolders(file models.FileStructure, lockWasGenerated *bool) {
	if _, err := os.OpenFile(file.LockFile, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
		log.Println("Error generating lockFile for " + file.Name)
		return
	}
	*lockWasGenerated = true

	if _, err := os.Stat(file.SentFolder); os.IsNotExist(err) {
		os.Mkdir(file.SentFolder, 0755)
	}
}

func fillFileInfo(fileData os.FileInfo, path string) models.FileStructure {
	fileDirectory, _ := filepath.Abs(filepath.Dir(path))

	return models.FileStructure{
		Name:             fileData.Name(),
		Folder:           fileDirectory,
		Path:             path,
		LockFile:         path + ".lock",
		SentFolder:       fileDirectory + "/sent",
		ErrorsFolder:     fileDirectory + "/errors",
		QueueDestination: fileDirectory[strings.LastIndex(fileDirectory, "/")+1:],
	}
}

func getFolderFiles(files *[]models.FileStructure) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() && strings.Contains(info.Name(), ".json") && !strings.Contains(info.Name(), ".lock") {
			*files = append(*files, fillFileInfo(info, path))
		}
		return nil
	}
}

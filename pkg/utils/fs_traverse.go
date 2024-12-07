package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aswin-kevin/nuclei-grpc/pkg/logger"
)

const nucleiTemplatesFolderName = "nuclei-templates"
const yamlExtension = ".yaml"

var allNucleiTemplatesMetaData = make(map[string][]string)
var currentOS = "linux"
var windowsOS = "windows"

func LoadAllNucleiTemplatesMetadata() {
	currentOS = runtime.GOOS
	homeDir, _ := os.UserHomeDir()
	basePath := filepath.Join(homeDir, nucleiTemplatesFolderName)

	err := getAllPossibleFolders(basePath, allNucleiTemplatesMetaData)

	if err != nil {
		logger.GlobalLogger.Error().Err(err).Msg("Error getting all possible folders")
	}

	for folder := range allNucleiTemplatesMetaData {
		allFiles := getAllPossibleFiles(basePath, folder)
		allNucleiTemplatesMetaData[folder] = allFiles
	}

	logger.GlobalLogger.Info().Msg("Loaded all nuclei templates metadata")
}

func getAllPossibleFolders(basePath string, fileSystemData map[string][]string) error {
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != basePath {
			relativePath := path[len(basePath)+1:]
			if len(relativePath) > 0 && relativePath[0] == '.' {
				return filepath.SkipDir
			}
			fileSystemData[relativePath] = make([]string, 0)
		}
		return nil
	})

	if err != nil {
		logger.GlobalLogger.Error().Err(err).Msg("Error walking the path " + basePath)
		return err
	}
	return nil
}

func getAllPossibleFiles(basePath string, folder string) []string {
	var allFiles = make([]string, 0)
	err := filepath.Walk(filepath.Join(basePath, folder), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filepath.Ext(path) == ".yaml" {
				filename := filepath.Base(path)
				nameWithoutExt := filename[:len(filename)-len(filepath.Ext(filename))]
				allFiles = append(allFiles, nameWithoutExt)
			}
		}
		return nil
	})

	if err != nil {
		logger.GlobalLogger.Error().Err(err).Msg("Error walking the path " + basePath)
	}
	return allFiles
}

func GetTemplateIdsFromTemplateData(templates []string) []string {
	allTemplateIds := make([]string, 0)

	for _, template := range templates {
		extention := filepath.Ext(template)
		if extention == yamlExtension {
			filename := filepath.Base(template)
			nameWithoutExt := filename[:len(filename)-5]
			allTemplateIds = append(allTemplateIds, nameWithoutExt)
		} else if extention == "" {
			if template[len(template)-1] == '/' {
				template = template[:len(template)-1]
			}
			if currentOS == windowsOS {
				template = strings.ReplaceAll(template, "/", "\\")
			}
			if _, ok := allNucleiTemplatesMetaData[template]; ok {
				allTemplateIds = append(allTemplateIds, allNucleiTemplatesMetaData[template]...)
			}
		}
	}

	return allTemplateIds
}

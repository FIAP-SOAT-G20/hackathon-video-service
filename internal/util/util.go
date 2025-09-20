package util

import (
	"os"
	"regexp"
)

func RemoveAllSpaces(s string) string {
	return regexp.MustCompile(`\s`).ReplaceAllString(s, "")
}

func ReadGoldenFile(name string) (string, error) {
	content, err := os.ReadFile("../../util/testdata/" + name + ".golden")
	if err != nil {
		return "", err
	}
	return RemoveAllSpaces(string(content)), nil
}

func ReadGoldenFiles(bussinessContext string, fileNames ...string) (map[string]string, error) {
	responses := make(map[string]string)
	for _, fileName := range fileNames {
		file, err := ReadGoldenFile(bussinessContext + "/" + fileName)
		if err != nil {
			return responses, err
		}
		responses[fileName] = string(file)

	}
	return responses, nil
}

func ReadFixtureFiles(bussinessContext string, fileNames ...string) (map[string]string, error) {
	requests := make(map[string]string)
	for _, fileName := range fileNames {
		file, err := os.ReadFile("../../util/testdata/" + bussinessContext + "/" + fileName + ".fixture")
		if err != nil {
			return requests, err
		}
		requests[fileName] = string(file)
	}
	return requests, nil
}

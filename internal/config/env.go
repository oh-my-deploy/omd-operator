package config

import (
	"github.com/joho/godotenv"
	"log"
	"path"
	"runtime"
)

func LoadEnv(relativePath string) {
	var projectDir string = getProjectDir(relativePath)
	log.Println("loading.env file")
	err := godotenv.Load(projectDir + "/.env")
	if err != nil {
		log.Panic("Error loading.env file")
	}
}

func getProjectDir(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	var projectDir string = path.Join(path.Dir(filename), relativePath)
	return projectDir
}

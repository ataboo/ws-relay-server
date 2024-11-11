package common

import (
	"io/fs"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

func GetLocalDir(folder string) (string, error) {
	rootPath := os.Getenv("GO_ATA_JWT_ROOT")
	if rootPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		rootPath = path.Join(homeDir, ".ws-relay-server")
	}

	fullPath := path.Join(rootPath, folder)

	return fullPath, nil
}

func GetAndMakeLocalDir(folder string) (string, error) {
	fullPath, err := GetLocalDir(folder)
	if err != nil {
		return "", err
	}

	os.MkdirAll(fullPath, fs.ModePerm)

	return fullPath, nil
}

func LoadDotEnv() error {
	root, err := GetLocalDir("/")
	if err != nil {
		return err
	}

	return godotenv.Load(path.Join(root, ".env"))
}

func LoadLogLevelFromEnv(fallBack log.Level) log.Level {
	logEnvStr := os.Getenv("LOG_LEVEL")
	if logEnvStr == "" {
		return fallBack
	}

	lvl, err := log.ParseLevel(logEnvStr)
	if err != nil {
		return fallBack
	}

	return lvl
}

func MustGetEnvVar(varName string) string {
	val := os.Getenv(varName)
	if val == "" {
		log.Fatalf("failed to get env var '%s'", varName)
	}

	return val
}

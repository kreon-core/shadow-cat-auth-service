package loaders

import "github.com/joho/godotenv"

func LoadEnv() {
	_ = godotenv.Load()
}

func LoadEnvFile(filePath string) error {
	return godotenv.Load(filePath)
}

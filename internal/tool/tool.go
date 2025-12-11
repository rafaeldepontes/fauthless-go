package tool

import "os"

func ChecksEnvFile(s *string) {
	_, err := os.Stat(*s)
	if err != nil {
		*s = ".env.example"
	}
}
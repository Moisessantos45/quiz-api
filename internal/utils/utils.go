package utils

import "fmt"

func BuildDriveAudioURL(fileID string) string {
	return fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", fileID)
}

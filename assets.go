package main

import (
	"github.com/google/uuid"
	"os"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(videoID uuid.UUID, mediaType string) string {
	// 1. derive extension from mediaType
	ext := getExtensionFromMediaType(mediaType)
	filename := videoID.String() + ext
	return "/assets/" + filename

	// 2. build "/assets/<videoID>.<ext>"
}

func getExtensionFromMediaType(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}

	ext := parts[1]
	if ext == "jpeg" {
		return ".jpg"
	}
	return "." + ext
}

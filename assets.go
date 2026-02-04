package main

import (
	"os"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(name string, mediaType string) string {
	ext := getExtensionFromMediaType(mediaType)
	filename := name + ext
	return "/assets/" + filename
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

// Map URL path -> disk path
func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	// assetPath is "/assets/<file>", we want "./assets/<file>"
	// so strip the leading "/assets" and prepend "./assets"
	const prefix = "/assets"
	if strings.HasPrefix(assetPath, prefix) {
		return cfg.assetsRoot + assetPath[len(prefix):]
	}
	return cfg.assetsRoot + assetPath
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return "http://localhost:" + cfg.port + assetPath
}

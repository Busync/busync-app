package main

type rgbColor struct {
	red   uint8
	green uint8
	blue  uint8
}

func addTrailingSlashIfNotExistsOnGivenPath(path string) string {
	lastCharOfPath := path[len(path)-1:]
	if lastCharOfPath != "/" {
		return path + "/"
	}

	return path
}

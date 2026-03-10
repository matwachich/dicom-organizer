package main

import (
	"path/filepath"
	"strings"
)

var forbiddenChars = strings.NewReplacer(">", "_", "<", "_", ":", "-", "\"", "'", "/", "_", "\\", "_", "|", "_", "?", "", "*", "")

func sanitizePath(s string) string {
	if abs, _ := filepath.Abs(s); abs != "" {
		s = abs
	}

	sp := strings.Split(s, string(filepath.Separator))

	i := 0
	if strings.HasSuffix(strings.TrimSpace(sp[0]), ":") {
		sp[0] = strings.TrimSpace(sp[0]) + string(filepath.Separator)
		i = 1
	}

	for ; i < len(sp); i++ {
		sp[i] = forbiddenChars.Replace(strings.TrimSpace(sp[i]))
	}

	return filepath.Clean(filepath.Join(sp...))
}

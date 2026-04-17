package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed cv-banned.csv
var csvData []byte

// messages stores simple integeration bindings for i18n translations
var messages = map[string]map[string]string{
	"en": {
		"usage":         "Usage: %s <directory-path>\n",
		"dir_not_exist": "Error: Target directory does not exist: %s\n",
		"csv_error":     "Error: Failed to parse embedded CSV: %v\n",
		"walk_warn":     "Warning: Access denied to %s: %v\n",
		"detected":      "[DETECTED] %s\n",
		"slug":          " ├─ matched slug: %s\n",
		"source":        " └─ blacklist source: %s\n",
		"walk_error":    "Error during filesystem scanning: %v\n",
		"found":         "\nScan completed. Found %d banned plugin(s).\n",
		"not_found":     "Scan completed successfully. No banned plugins found.\n",
	},
	"pl": {
		"usage":         "Użycie: %s <ścieżka-do-katalogu>\n",
		"dir_not_exist": "Błąd: Podany katalog docelowy nie istnieje: %s\n",
		"csv_error":     "Błąd podczas parsowania wbudowanego pliku CSV: %v\n",
		"walk_warn":     "Ostrzeżenie - brak dostępu do %s: %v\n",
		"detected":      "[WYKRYTO] %s\n",
		"slug":          " ├─ Slug (dopasowanie): %s\n",
		"source":        " └─ Źródło blacklisty: %s\n",
		"walk_error":    "Błąd z systemem plików podczas skanowania: %v\n",
		"found":         "\nSkanowanie zakończone. Znaleziono %d wystąpień zbanowanych pluginów.\n",
		"not_found":     "Skanowanie zakończone pomyślnie. Nie znaleziono zbanowanych pluginów w podanym folderze.\n",
	},
}

// getLang determines the current lang by checking environmental variables
func getLang() string {
	lang := os.Getenv("LANG")
	if strings.HasPrefix(lang, "pl") {
		return "pl"
	}
	return "en"
}

// msg fetches a translated string based on the current locale
func msg(key string) string {
	lang := getLang()
	if val, ok := messages[lang][key]; ok {
		return val
	}
	if val, ok := messages["en"][key]; ok {
		return val
	}
	return key
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, msg("usage"), filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	targetDir := os.Args[1]

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, msg("dir_not_exist"), targetDir)
		os.Exit(1)
	}

	reader := csv.NewReader(bytes.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, msg("csv_error"), err)
		os.Exit(1)
	}

	bannedPlugins := make(map[string]string)
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) >= 4 {
			slug := record[0]
			source := record[3]
			bannedPlugins[slug] = source
		}
	}

	foundCount := 0

	err = filepath.WalkDir(targetDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, msg("walk_warn"), path, err)
			return nil
		}

		name := d.Name()
		matchedSlug, isBanned := matchSlug(name, d.IsDir(), bannedPlugins)

		if isBanned {
			source := bannedPlugins[matchedSlug]
			fmt.Printf(msg("detected"), path)
			fmt.Printf(msg("slug"), matchedSlug)
			fmt.Printf(msg("source"), source)
			foundCount++

			if d.IsDir() {
				// Don't dive deeply into the directory to prevent duplication
				return filepath.SkipDir
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, msg("walk_error"), err)
		os.Exit(1)
	}

	if foundCount > 0 {
		fmt.Printf(msg("found"), foundCount)
		os.Exit(1)
	} else {
		fmt.Print(msg("not_found"))
	}
}

func matchSlug(name string, isDir bool, banned map[string]string) (string, bool) {
	// 1. Exact match
	if _, ok := banned[name]; ok {
		return name, true
	}

	// 2. Ext trimming for zip/tar/php payload patterns
	if !isDir {
		if strings.HasSuffix(name, ".php") {
			slug := strings.TrimSuffix(name, ".php")
			if _, ok := banned[slug]; ok {
				return slug, true
			}
		}
		if strings.HasSuffix(name, ".zip") {
			slug := strings.TrimSuffix(name, ".zip")
			if _, ok := banned[slug]; ok {
				return slug, true
			}
		}
		if strings.HasSuffix(name, ".tar.gz") {
			slug := strings.TrimSuffix(name, ".tar.gz")
			if _, ok := banned[slug]; ok {
				return slug, true
			}
		}
	}

	return "", false
}

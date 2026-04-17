package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed cv-banned.csv
var csvData []byte

var (
	bwFlag bool
)

// ANSI colors
var (
	Reset  = "\033[0m"
	Red    = "\033[31;1m" // bold red
	Green  = "\033[32;1m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Symbols
var (
	SymDetect = "🚨"
	SymBranch = " ├─"
	SymEnd    = " └─"
	SymCheck  = "✅"
	SymWarn   = "⚠️"
	SymCross  = "❌"
)

// setBW disables colors and UTF8 characters if --bw flag is used
func setBW() {
	SymDetect = "[!]"
	SymBranch = " |-"
	SymEnd = " `-"
	SymCheck = "[OK]"
	SymWarn = "[WARN]"
	SymCross = "[X]"

	Reset = ""
	Red = ""
	Green = ""
	Yellow = ""
	Cyan = ""
	Bold = ""
}

var messages = map[string]map[string]string{
	"en": {
		"use":           "wpvul <directory-path>",
		"desc":          "wpvul is a lightning-fast zero-dependency WordPress plugin scanner.",
		"dir_not_exist": "%s Error: Target directory does not exist: %s\n",
		"csv_error":     "%s Error: Failed to parse embedded CSV: %v\n",
		"walk_warn":     "%s Warning: Access denied to %s: %v\n",
		"detected":      "%s %s[DETECTED]%s %s%s%s\n",
		"slug":          "%s matched slug: %s%s%s\n",
		"source":        "%s blacklist source: %s%s%s\n",
		"walk_error":    "%s Error during filesystem scanning: %v\n",
		"found":         "\n%s Scan completed. Found %s%d%s banned plugin(s).\n",
		"not_found":     "\n%s Scan completed successfully. No banned plugins found.\n",
	},
	"pl": {
		"use":           "wpvul <ścieżka-do-katalogu>",
		"desc":          "wpvul to ekstremalnie szybki skaner podatnych pluginów WordPressa.",
		"dir_not_exist": "%s Błąd: Podany katalog docelowy nie istnieje: %s\n",
		"csv_error":     "%s Błąd podczas parsowania wbudowanego pliku CSV: %v\n",
		"walk_warn":     "%s Ostrzeżenie - brak dostępu do %s: %v\n",
		"detected":      "%s %s[WYKRYTO]%s %s%s%s\n",
		"slug":          "%s dopasowanie (slug): %s%s%s\n",
		"source":        "%s źródło blacklisty: %s%s%s\n",
		"walk_error":    "%s Błąd z systemem plików podczas skanowania: %v\n",
		"found":         "\n%s Skanowanie zakończone. Znaleziono %s%d%s wystąpień zbanowanych pluginów.\n",
		"not_found":     "\n%s Skanowanie zakończone pomyślnie. Nie znaleziono zbanowanych pluginów.\n",
	},
}

func getLang() string {
	lang := os.Getenv("LANG")
	if strings.HasPrefix(lang, "pl") {
		return "pl"
	}
	return "en"
}

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

const Version = "1.0.1"

func Execute() {
	var rootCmd = &cobra.Command{
		Use:     msg("use"),
		Short:   msg("desc"),
		Version: Version,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if bwFlag {
				setBW()
			}
			runScan(args[0])
		},
	}

	// Flag definition
	rootCmd.Flags().BoolVar(&bwFlag, "bw", false, "Disable colors and UTF-8 symbols (ASCII mode)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}

func runScan(targetDir string) {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, msg("dir_not_exist"), SymCross, targetDir)
		os.Exit(1)
	}

	reader := csv.NewReader(bytes.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, msg("csv_error"), SymCross, err)
		os.Exit(1)
	}

	bannedPlugins := make(map[string]string)
	for i, record := range records {
		if i == 0 {
			continue // Pomiń nagłówek
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
			fmt.Fprintf(os.Stderr, msg("walk_warn"), SymWarn, path, err)
			return nil
		}

		name := d.Name()
		absDir, _ := filepath.Abs(filepath.Dir(path))
		parentName := filepath.Base(absDir)
		isRoot := path == targetDir || path == "."

		matchedSlug, isBanned := matchSlug(name, d.IsDir(), bannedPlugins)

		if isBanned {
			source := bannedPlugins[matchedSlug]
			fmt.Printf(msg("detected"), SymDetect, Red, Reset, Bold, path, Reset)
			fmt.Printf(msg("slug"), SymBranch, Yellow, matchedSlug, Reset)
			fmt.Printf(msg("source"), SymEnd, Cyan, source, Reset)
			foundCount++

			if d.IsDir() && !isRoot {
				return filepath.SkipDir
			}
		} else {
			// Submodules false-positives prevention logic
			if d.IsDir() && !isRoot && (parentName == "plugins" || parentName == "mu-plugins") {
				return filepath.SkipDir
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, msg("walk_error"), SymCross, err)
		os.Exit(1)
	}

	if foundCount > 0 {
		fmt.Printf(msg("found"), SymDetect, Red, foundCount, Reset)
		os.Exit(1)
	} else {
		fmt.Printf(msg("not_found"), SymCheck)
	}
}

func matchSlug(name string, isDir bool, banned map[string]string) (string, bool) {
	if _, ok := banned[name]; ok {
		return name, true
	}

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

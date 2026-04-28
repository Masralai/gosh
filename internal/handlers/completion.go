package handlers

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	GhostDim   = "\033[2m"
	GhostReset = "\033[0m"
)

var (
	globalSuggestion string
	suggestionMu    sync.RWMutex
)

func SetSuggestion(s string) {
	suggestionMu.Lock()
	defer suggestionMu.Unlock()
	globalSuggestion = s
}

func GetSuggestion() string {
	suggestionMu.RLock()
	defer suggestionMu.RUnlock()
	return globalSuggestion
}

func ClearSuggestion() {
	suggestionMu.Lock()
	defer suggestionMu.Unlock()
	globalSuggestion = ""
}

func GetCompletionsForLine(input string) []string {
	if input == "" {
		return nil
	}

	fields := strings.Fields(input)
	if len(fields) == 0 {
		return nil
	}

	lastField := fields[len(fields)-1]
	isFirstWord := len(fields) == 1

	if isFirstWord {
		matches := GetHistoryManager().Search(lastField, 10)
		if len(matches) > 0 {
			SetSuggestion(matches[0])
			return matches[:1]
		}
	}

	prevField := ""
	if len(fields) >= 2 {
		prevField = fields[len(fields)-2]
	}

	if isFileCommand(prevField) {
		paths := getPathCompletions(lastField)
		if len(paths) > 0 {
			SetSuggestion(paths[0])
			return paths
		}
	}

	if strings.HasPrefix(lastField, "!") {
		var bangs []string
		for _, cmd := range All() {
			bangs = append(bangs, "!"+cmd.Name)
		}
		ClearSuggestion()
		return filterPrefix(bangs, lastField)
	}

	commands := []string{
		"echo", "cd", "pwd", "exit", "ls", "mkdir", "rm", "touch", "mv", "cp",
		"cat", "info", "ps", "ut", "sys", "mu", "du", "kill", "grep", "head",
		"tail", "ping", "zip", "unzip", "gosh",
	}
	filtered := filterPrefix(commands, lastField)
	if len(filtered) > 0 {
		SetSuggestion(filtered[0])
		return filtered
	}

	ClearSuggestion()
	return nil
}

func isFileCommand(cmd string) bool {
	return cmd == "ls" || cmd == "cd" || cmd == "cat" || cmd == "cp" || cmd == "mv" ||
		cmd == "rm" || cmd == "touch" || cmd == "mkdir" || cmd == "info" ||
		cmd == "du" || cmd == "zip" || cmd == "unzip" || cmd == "head" || cmd == "tail"
}

func getPathCompletions(prefix string) []string {
	dir := "."
	base := prefix

	if prefix != "" {
		if strings.Contains(prefix, "/") {
			dir = filepath.Dir(prefix)
			base = filepath.Base(prefix)
		} else if filepath.IsAbs(prefix) {
			dir = filepath.Dir(prefix)
			base = filepath.Base(prefix)
		}
	}

	if dir == "." {
		dir = ""
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var results []string
	for _, e := range entries {
		name := e.Name()
		if base == "" || strings.HasPrefix(name, base) {
			if e.IsDir() {
				results = append(results, name+"/")
			} else {
				results = append(results, name)
			}
		}
	}
	return results
}

func filterPrefix(list []string, prefix string) []string {
	var result []string
	for _, s := range list {
		if strings.HasPrefix(s, prefix) {
			result = append(result, s)
		}
	}
	return result
}

func GetGhostSuffix(input string) string {
	if input == "" {
		return ""
	}

	fields := strings.Fields(input)
	if len(fields) == 0 {
		return ""
	}

	lastField := fields[len(fields)-1]
	suggestion := GetSuggestion()

	if suggestion == "" {
		return ""
	}

	if !strings.HasPrefix(suggestion, lastField) {
		return ""
	}

	return suggestion[len(lastField):]
}
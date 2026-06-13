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

func ClearSuggestion() {
	suggestionMu.Lock()
	defer suggestionMu.Unlock()
	globalSuggestion = ""
}

func isFileCommand(cmd string) bool {
	return cmd == "ls" || cmd == "cd" || cmd == "cat" || cmd == "cp" || cmd == "mv" ||
		cmd == "rm" || cmd == "touch" || cmd == "mkdir" || cmd == "info" ||
		cmd == "du" || cmd == "zip" || cmd == "unzip" || cmd == "head" || cmd == "tail"
}

func getPathCompletions(prefix string, dirsOnly bool) []string {
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
			if dirsOnly && !e.IsDir() {
				continue
			}
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

func GetSuggestionForLine(input string) string {
	if input == "" {
		ClearSuggestion()
		return ""
	}

	fields := strings.Fields(input)
	if len(fields) == 0 {
		ClearSuggestion()
		return ""
	}

	lastField := fields[len(fields)-1]
	isFirstWord := len(fields) == 1
	inputEndsWithSpace := strings.HasSuffix(input, " ")

	// Special handling for file commands (cd, ls, etc.) with space or partial
	if len(fields) >= 1 && (isFileCommand(fields[0]) || fields[0] == "cd") {
		cmd := fields[0]
		partial := ""

		if inputEndsWithSpace {
			partial = "" // User typed "cd ", no partial yet
		} else if len(fields) >= 2 {
			partial = fields[len(fields)-1] // User typed "cd D", partial is "D"
		}

		if cmd == "cd" && (partial != "" || inputEndsWithSpace) {
			// Priority 1: Current directory subdirectories only
			paths := getPathCompletions(partial, true) // dirsOnly=true
			if len(paths) > 0 {
				fullSuggestion := cmd + " " + paths[0]
				SetSuggestion(fullSuggestion)
				if inputEndsWithSpace {
					return paths[0]
				}
				return paths[0][len(partial):]
			}

			// Priority 2: Recent directories (filtered by partial, show relative paths)
			recentDirs := GetHistoryManager().GetRecentDirs(partial, 5)
			if len(recentDirs) > 0 {
				// Convert to relative path for display
				cwd, _ := os.Getwd()
				displayPath := recentDirs[0]
				// Strip cwd prefix to show relative path
				if strings.HasPrefix(displayPath, cwd+"/") {
					displayPath = displayPath[len(cwd)+1:]
				} else if displayPath == cwd {
					displayPath = "."
				}
				fullSuggestion := cmd + " " + displayPath
				SetSuggestion(fullSuggestion)
				if inputEndsWithSpace {
					return displayPath // Show full relative path after "cd "
				}
				return displayPath[len(partial):] // Show remainder after partial
			}
		} else if isFileCommand(cmd) && (partial != "" || inputEndsWithSpace) {
			// For other file commands (ls, cat, etc.), show files too
			paths := getPathCompletions(partial, false) // dirsOnly=false
			if len(paths) > 0 {
				fullSuggestion := cmd + " " + paths[0]
				SetSuggestion(fullSuggestion)
				if inputEndsWithSpace {
					return paths[0]
				}
				return paths[0][len(partial):]
			}
		}
	}

	// Priority 1: Check known commands first (for first word only)
	if isFirstWord {
		commands := []string{
			"echo", "cd", "pwd", "exit", "ls", "mkdir", "rm", "touch", "mv", "cp",
			"cat", "info", "ps", "ut", "sys", "mu", "du", "kill", "grep", "head",
			"tail", "ping", "zip", "unzip", "gosh", "sudo", "su", "apt", "apt-get",
		}
		matches := filterPrefix(commands, lastField)
		if len(matches) > 0 {
			SetSuggestion(matches[0])
			return matches[0][len(lastField):]
		}
	}

	// Priority 2: Search history with full input line (works for both first word and after space)
	historyMatches := GetHistoryManager().Search(input, 1)
	if len(historyMatches) > 0 {
		SetSuggestion(historyMatches[0])
		return historyMatches[0][len(input):]
	}

	prevField := ""
	if len(fields) >= 2 {
		prevField = fields[len(fields)-2]
	}

	if isFileCommand(prevField) {
		paths := getPathCompletions(lastField, false)
		if len(paths) > 0 {
			SetSuggestion(paths[0])
			return paths[0][len(lastField):]
		}
	}

	if strings.HasPrefix(lastField, "!") {
		var bangs []string
		for _, cmd := range All() {
			bangs = append(bangs, "!"+cmd.Name)
		}
		matches := filterPrefix(bangs, lastField)
		if len(matches) > 0 {
			SetSuggestion(matches[0])
			return matches[0][len(lastField):]
		}
	}

	ClearSuggestion()
	return ""
}

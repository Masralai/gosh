package handlers

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

func getHomeDir() string {
	home := os.Getenv("HOME")
	if home != "" {
		return home
	}
	usr, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return usr
}

func loadBashHistoryFile(h *HistoryManager, path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, " ") {
			h.Add(line)
		}
	}
}

type HistoryManager struct {
	mu       sync.RWMutex
	items    []string
	seen     map[string]bool
	maxSize  int
}

func NewHistoryManager(maxSize int) *HistoryManager {
	if maxSize <= 0 {
		maxSize = 1000
	}
	h := &HistoryManager{
		items: make([]string, 0, maxSize),
		seen:  make(map[string]bool),
		maxSize: maxSize,
	}
	h.loadBashHistory()
	return h
}

func (h *HistoryManager) loadBashHistory() {
	home := getHomeDir()
	if home == "" {
		return
	}
	loadBashHistoryFile(h, home+"/.bash_history")
}

func (h *HistoryManager) Add(cmd string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return
	}

	if h.seen[cmd] {
		for i, item := range h.items {
			if item == cmd {
				copy(h.items[i:], h.items[i+1:])
				h.items[len(h.items)-1] = cmd
				h.items = h.items[:len(h.items)-1]
				break
			}
		}
	} else {
		h.seen[cmd] = true
		h.items = append(h.items, cmd)
	}

	if len(h.items) > h.maxSize {
		removed := h.items[0]
		delete(h.seen, removed)
		h.items = h.items[1:]
	}
}

func (h *HistoryManager) Search(prefix string, limit int) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if prefix == "" {
		if limit > len(h.items) {
			limit = len(h.items)
		}
		if limit <= 0 {
			return nil
		}
		start := len(h.items) - limit
		if start < 0 {
			start = 0
		}
		return h.items[start:]
	}

	var results []string
	prefixLower := strings.ToLower(prefix)

	for i := len(h.items) - 1; i >= 0 && len(results) < limit; i-- {
		item := h.items[i]
		if strings.HasPrefix(strings.ToLower(item), prefixLower) {
			results = append(results, item)
		} else if strings.Contains(item, " ") && strings.HasPrefix(strings.Fields(item)[0], prefix) {
			results = append(results, item)
		}
	}

	return results
}

func (h *HistoryManager) SearchPaths(prefix string, limit int) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var results []string
	prefixLower := strings.ToLower(prefix)

	for i := len(h.items) - 1; i >= 0 && len(results) < limit; i-- {
		item := h.items[i]
		fields := strings.Fields(item)

		for _, field := range fields {
			if strings.HasPrefix(field, "/") || strings.HasPrefix(field, ".") {
				if prefix == "" || strings.HasPrefix(strings.ToLower(field), prefixLower) {
					if !stringInSlice(field, results) {
						results = append(results, field)
						if len(results) >= limit {
							return results
						}
					}
				}
			}
		}
	}

	return results
}

func stringInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

var (
	globalHistory *HistoryManager
	historyOnce  sync.Once
)

func GetHistoryManager() *HistoryManager {
	historyOnce.Do(func() {
		globalHistory = NewHistoryManager(1000)
	})
	return globalHistory
}

func SearchHistoryForCompletion(input string, limit int) []string {
	h := GetHistoryManager()
	if input == "" {
		return h.Search("", limit)
	}
	return h.Search(input, limit)
}

func SearchHistoryPathsForCompletion(input string, limit int) []string {
	h := GetHistoryManager()
	return h.SearchPaths(input, limit)
}
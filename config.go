package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sync"
)

// Config holds user preferences
type Config struct {
	ZoomSensitivity  int      `json:"zoomSensitivity"`
	Theme            string   `json:"theme"`
	ZoomLevel        float64  `json:"zoomLevel"`
	FontFamily       string   `json:"fontFamily"`
	FontSize         int      `json:"fontSize"`
	Language         string   `json:"language"`
	ShowLineNumbers  bool     `json:"showLineNumbers"`
	ShowTOC          bool     `json:"showTOC"`
	WindowWidth      int      `json:"windowWidth"`
	WindowHeight     int      `json:"windowHeight"`
	WindowX          int      `json:"windowX"`
	WindowY          int      `json:"windowY"`
	LastOpenedFile   string   `json:"lastOpenedFile"`
	RecentFiles      []string `json:"recentFiles"`
	TranslateBackend string   `json:"translateBackend"`
	DebugMode        bool     `json:"debugMode"`
}

var (
	currentConfig Config
	configMu      sync.RWMutex
)

func configDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(usr.HomeDir, ".md-viewer")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func defaultConfig() Config {
	return Config{
		ZoomSensitivity: 5,
		Theme:            "auto",
		ZoomLevel:        1.0,
		FontFamily:       "-apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif",
		FontSize:         16,
		Language:        "zhTW",
		DebugMode:        false,
	}
}

// LoadConfig loads config from ~/.md-viewer/config.json.
// Creates with defaults if file does not exist.
func LoadConfig() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			configMu.Lock()
			currentConfig = defaultConfig()
			configMu.Unlock()
			return saveConfig() // create default file
		}
		return err
	}
	configMu.Lock()
	defer configMu.Unlock()
	return json.Unmarshal(data, &currentConfig)
}

// saveConfig writes currentConfig to the config file.
func saveConfig() error {
	configMu.RLock()
	data, err := json.MarshalIndent(currentConfig, "", "  ")
	configMu.RUnlock()
	
	if err != nil {
		return err
	}
	path, err := configPath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// SetZoomSensitivity updates the zoom sensitivity and persists it.
func SetZoomSensitivity(level int) error {
	configMu.Lock()
	currentConfig.ZoomSensitivity = level
	configMu.Unlock()
	return saveConfig()
}

// SetTheme updates the theme and persists it.
func SetTheme(theme string) error {
	configMu.Lock()
	currentConfig.Theme = theme
	configMu.Unlock()
	return saveConfig()
}

// SetZoomLevel updates the zoom level and persists it.
func SetZoomLevel(level float64) error {
	configMu.Lock()
	currentConfig.ZoomLevel = level
	configMu.Unlock()
	return saveConfig()
}

// SetFont updates the font and persists it.
func SetFont(family string, size int) error {
	configMu.Lock()
	currentConfig.FontFamily = family
	currentConfig.FontSize = size
	configMu.Unlock()
	return saveConfig()
}

// SetLanguage updates the language and persists it.
func SetLanguage(lang string) error {
	configMu.Lock()
	currentConfig.Language = lang
	configMu.Unlock()
	return saveConfig()
}

// SetLineNumbers updates the show line numbers preference and persists it.
func SetLineNumbers(show bool) error {
	configMu.Lock()
	currentConfig.ShowLineNumbers = show
	configMu.Unlock()
	return saveConfig()
}

// SetTOC updates the show TOC preference and persists it.
func SetTOC(show bool) error {
	configMu.Lock()
	currentConfig.ShowTOC = show
	configMu.Unlock()
	return saveConfig()
}

// SetWindowSize updates the window size and persists it.
func SetWindowSize(width, height int) error {
	configMu.Lock()
	currentConfig.WindowWidth = width
	currentConfig.WindowHeight = height
	configMu.Unlock()
	return saveConfig()
}

// SetWindowPosition updates the window position and persists it.
func SetWindowPosition(x, y int) error {
	configMu.Lock()
	currentConfig.WindowX = x
	currentConfig.WindowY = y
	configMu.Unlock()
	return saveConfig()
}

// SetDebugMode updates the debug mode and persists it.
func SetDebugMode(enabled bool) error {
	configMu.Lock()
	currentConfig.DebugMode = enabled
	configMu.Unlock()
	return saveConfig()
}

// SetLastOpenedFile updates the last opened file path and persists it.
func SetLastOpenedFile(path string) error {
	configMu.Lock()
	currentConfig.LastOpenedFile = path
	// Also add to recent files
	AddRecentFileInternal(path)
	configMu.Unlock()
	return saveConfig()
}

// AddRecentFile adds a file to the recent files list (max 10)
func AddRecentFile(path string) {
	configMu.Lock()
	defer configMu.Unlock()
	AddRecentFileInternal(path)
}

func AddRecentFileInternal(path string) {
	if path == "" {
		return
	}
	// Remove if already exists
	for i, f := range currentConfig.RecentFiles {
		if f == path {
			currentConfig.RecentFiles = append(currentConfig.RecentFiles[:i], currentConfig.RecentFiles[i+1:]...)
			break
		}
	}
	// Add to front
	currentConfig.RecentFiles = append([]string{path}, currentConfig.RecentFiles...)
	// Keep only 10
	if len(currentConfig.RecentFiles) > 10 {
		currentConfig.RecentFiles = currentConfig.RecentFiles[:10]
	}
}

// GetRecentFiles returns the list of recent files
func GetRecentFiles() []string {
	configMu.RLock()
	defer configMu.RUnlock()
	return currentConfig.RecentFiles
}

// RemoveRecentFile removes a file from the recent files list
func RemoveRecentFile(path string) error {
	configMu.Lock()
	newList := []string{}
	for _, f := range currentConfig.RecentFiles {
		if f != path {
			newList = append(newList, f)
		}
	}
	currentConfig.RecentFiles = newList
	configMu.Unlock()
	return saveConfig()
}

// GetConfig returns a copy of the current config.
func GetConfig() Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return currentConfig
}

// ConfigToJS returns a JS snippet that sets window.mdConfig.
func ConfigToJS() string {
	configMu.RLock()
	defer configMu.RUnlock()
	backend := os.Getenv("TRANSLATE_BACKEND")
	if backend == "" {
		backend = "google"
	}
	return fmt.Sprintf(
		`window.mdConfig = {zoomSensitivity: %d, theme: %q, zoomLevel: %f, fontFamily: %q, fontSize: %d, language: %q, showLineNumbers: %t, showTOC: %t, translateBackend: %q, debugMode: %t};`,
		currentConfig.ZoomSensitivity,
		currentConfig.Theme,
		currentConfig.ZoomLevel,
		currentConfig.FontFamily,
		currentConfig.FontSize,
		currentConfig.Language,
		currentConfig.ShowLineNumbers,
		currentConfig.ShowTOC,
		backend,
		currentConfig.DebugMode,
	)
}
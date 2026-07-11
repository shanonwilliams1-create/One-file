// Package detect provides operating system detection and config file discovery.
// It supports desktop (Linux, macOS, Windows) and mobile (iOS, Android) platforms.
package detect

import "runtime"

// OSInfo holds information about the detected operating system.
type OSInfo struct {
	// GOOS is the raw runtime.GOOS value (linux, darwin, windows, ios, android)
	GOOS string
	// PrettyName is a human-readable OS name
	PrettyName string
	// NativeBinary is the recommended binary name for this platform
	NativeBinary string
	// IsMobile is true for iOS and Android
	IsMobile bool
	// ConfigDirs lists common configuration directory paths
	ConfigDirs []string
	// ConfigFiles lists common individual config file paths
	ConfigFiles []string
}

// DetectOS detects the current operating system and returns OSInfo
// with platform-specific config paths, including mobile platforms.
func DetectOS() OSInfo {
	info := OSInfo{GOOS: runtime.GOOS}

	switch runtime.GOOS {
	case "linux":
		info.PrettyName = "Linux"
		info.NativeBinary = "omniconfig"
		info.ConfigDirs = []string{
			"/etc",
			"/etc/opt",
			"~/.config",
			"~/.local/share",
		}
		info.ConfigFiles = []string{
			"~/.bashrc",
			"~/.profile",
			"~/.bash_profile",
			"~/.gitconfig",
			"/etc/environment",
		}

	case "darwin":
		info.PrettyName = "macOS"
		info.NativeBinary = "omniconfig"
		info.ConfigDirs = []string{
			"~/Library/Preferences",
			"~/.config",
			"/etc",
			"/opt/homebrew/etc",
		}
		info.ConfigFiles = []string{
			"~/.bash_profile",
			"~/.zshrc",
			"~/.zshenv",
			"~/.profile",
			"~/.gitconfig",
		}

	case "windows":
		info.PrettyName = "Windows"
		info.NativeBinary = "omniconfig.exe"
		info.ConfigDirs = []string{
			`%APPDATA%`,
			`%USERPROFILE%`,
			`%PROGRAMDATA%`,
			`%LOCALAPPDATA%`,
		}
		info.ConfigFiles = []string{
			`%USERPROFILE%\.gitconfig`,
			`%WINDIR%\System32\drivers\etc\hosts`,
		}

	case "ios":
		info.PrettyName = "iOS"
		info.NativeBinary = "omniconfig"
		info.IsMobile = true
		info.ConfigDirs = []string{
			"<AppContainer>/Documents",
			"<AppContainer>/Library/Preferences",
			"<AppContainer>/Library/Application Support",
			"<SharedContainer>/Library/Preferences",
		}
		info.ConfigFiles = []string{
			"<AppContainer>/Library/Preferences/com.omnicofig.plist",
		}

	case "android":
		info.PrettyName = "Android"
		info.NativeBinary = "omniconfig"
		info.IsMobile = true
		info.ConfigDirs = []string{
			"/data/data/com.omnicofig/files",
			"/data/data/com.omnicofig/shared_prefs",
			"/sdcard/Android/data/com.omnicofig/files",
			"<AppSpecificStorage>",
		}
		info.ConfigFiles = []string{
			"/data/data/com.omnicofig/shared_prefs/omniconfig.xml",
			"/data/data/com.omnicofig/files/.env",
		}

	default:
		info.PrettyName = runtime.GOOS
		info.NativeBinary = "omniconfig"
		info.ConfigDirs = []string{
			"~/.config",
		}
		info.ConfigFiles = []string{
			"~/.config/omniconfig.conf",
		}
	}

	return info
}

// IsMobileOS returns true if the current OS is iOS or Android.
func IsMobileOS() bool {
	return runtime.GOOS == "ios" || runtime.GOOS == "android"
}
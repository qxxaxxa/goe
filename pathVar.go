package main

import (
	"runtime"
)

func cachePath() string {
	sysType := runtime.GOOS
	base := ""
	if sysType == "windows" {
		base = "D:\\cache"
		// windows系统
	}
	if sysType == "darwin" {
		// LINUX系统
		base = "/Users/xiexingan/Downloads/cache"
	}
	return base
}

func imagePath() string {
	sysType := runtime.GOOS
	base := ""
	if sysType == "windows" {
		base = "D:\\Downloads\\file.txt"
	}
	if sysType == "darwin" {
		base = "/Users/xiexingan/Downloads/file.txt"
	}
	return base
}
func savedPath() string {
	sysType := runtime.GOOS
	base := ""
	if sysType == "windows" {
		base = "D:\\Downloads\\saved.txt"
	}
	if sysType == "darwin" {
		base = "/Users/xiexingan/Downloads/saved.txt"
	}
	return base
}

func hostprefix() string {
	return "qugxjjz"
}
func hentaihost() string {
	return "pgbzwylrvito"
}

func linkPath() string {
	sysType := runtime.GOOS
	base := ""
	if sysType == "windows" {
		base = "D:\\Downloads\\links.txt"
	}
	if sysType == "darwin" {
		base = "/Users/xiexingan/Downloads/links.txt"
	}
	return base
}

func galleryPath() string {
	sysType := runtime.GOOS
	base := ""
	if sysType == "windows" {
		base = "C:/Users/Happy/Desktop/gallery.txt"
	}
	if sysType == "darwin" {
		base = "/Users/xiexingan/anime/gallery.txt"
	}
	return base
}

func logPath() string {
	sysType := runtime.GOOS
	base := ""
	if sysType == "windows" {
		base = "D:\\Downloads\\log_out"
	}
	if sysType == "darwin" {
		base = "/Users/xiexingan/Downloads/log_out"
	}
	return base
}
func sha1filePath() string {
	sysType := runtime.GOOS
	base := ""
	if sysType == "windows" {
		base = "C:\\cache\\checksha1.txt"
	}
	if sysType == "darwin" {
		base = "/Users/xiexingan/Downloads/cache/checksha1.txt"
	}
	return base

}

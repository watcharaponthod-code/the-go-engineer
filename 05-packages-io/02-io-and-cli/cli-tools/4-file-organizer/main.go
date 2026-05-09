// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 05: Packages and I/O
// Title: File Organizer
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Build a production-ready CLI tool that interacts with the filesystem.
//   - Implement the "dry-run" pattern for safe file operations.
//
// WHY THIS MATTERS:
//   - This project combines flag parsing, filesystem traversal, and directory
//     management into a single useful tool.
//
// RUN:
//   go run ./05-packages-io/02-io-and-cli/cli-tools/4-file-organizer
//
// KEY TAKEAWAY:
//   - Always implement a dry-run flag for tools that modify user data.
// ============================================================================

// Commercial use is prohibited without permission.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Stage 05: CLI Tools - File Organizer (Exercise)
//
//   - Using the `flag` package for command-line options
//   - Reading directory contents with os.ReadDir
//   - Organizing files by extension into subdirectories
//   - The --dry-run pattern for safe testing
//
// ENGINEERING DEPTH:
//   The `flag` package uses a registration pattern - you call `flag.String()`
//   to register each flag, which returns a *string pointer. The actual parsing
//   doesn't happen until you call `args := os.Args[1:]; if len(args) > 0 && args[0] == os.Args[0] { args = args[1:] }; flag.CommandLine.Parse(args)`. This two-step approach allows
//   you to define flags anywhere in your code (even in init() functions across
//   packages) before a single centralized Parse() call processes them all.
//
// USAGE: go run ./05-packages-io/02-io-and-cli/cli-tools/4-file-organizer --dir=./my-messy-folder
// FLAGS: --dir (required), --dry-run (optional, default false)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == os.Args[0] {
		args = args[1:]
	}

	fs := flag.NewFlagSet("file-organizer", flag.ExitOnError)
	dir := fs.String("dir", "", "Directory to organize (required)")
	dryRun := fs.Bool("dry-run", false, "Preview changes without moving files")

	fs.Parse(args)

	if *dir == "" {
		fmt.Println("=== File Organizer ===")
		fmt.Println()
		fmt.Println("Organizes files in a directory by extension.")
		fmt.Println()
		fs.Usage()
		fmt.Println()
		fmt.Println("Running demo mode...")
		fmt.Println()

		tmpDir, _ := os.MkdirTemp("", "organizer-demo-*")
		defer os.RemoveAll(tmpDir)

		sampleFiles := []string{
			"report.pdf", "photo.jpg", "notes.txt", "main.go",
			"styles.css", "readme.md", "data.json", "script.py",
		}
		for _, name := range sampleFiles {
			os.WriteFile(filepath.Join(tmpDir, name), []byte("sample"), 0644)
		}

		*dir = tmpDir
		*dryRun = true
	}

	fmt.Printf("DIR: Organizing: %s\n", *dir)
	if *dryRun {
		fmt.Println("PREVIEW: DRY RUN - no files will actually be moved")
	}
	fmt.Println()

	// Read directory contents
	entries, err := os.ReadDir(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		os.Exit(1)
	}

	moved := 0
	for _, entry := range entries {
		// Skip directories - we only organize files.
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		ext := strings.TrimPrefix(filepath.Ext(filename), ".") // "go", "txt", etc.

		// Skip files with no extension
		if ext == "" {
			fmt.Printf("  [-] Skipping (no extension): %s\n", filename)
			continue
		}

		// Target subdirectory named after the extension
		targetDir := filepath.Join(*dir, ext)
		targetPath := filepath.Join(targetDir, filename)
		sourcePath := filepath.Join(*dir, filename)

		if *dryRun {
			fmt.Printf("  [PREVIEW] %s -> %s/\n", filename, ext)
		} else {
			// Create the target subdirectory if it doesn't exist.
			// os.MkdirAll creates all parent directories as needed.
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "  [ERROR] Creating dir %s: %v\n", targetDir, err)
				continue
			}

			// Move the file by renaming it to the new path.
			if err := os.Rename(sourcePath, targetPath); err != nil {
				fmt.Fprintf(os.Stderr, "  [ERROR] Moving %s: %v\n", filename, err)
				continue
			}
			fmt.Printf("  [OK] %s -> %s/\n", filename, ext)
		}
		moved++
	}

	fmt.Println()
	action := "would be moved"
	if !*dryRun {
		action = "moved"
	}
	fmt.Printf("Done! %d files %s.\n", moved, action)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: EN.1 -> 05-packages-io/02-io-and-cli/encoding/1-marshalling")
	fmt.Println("Current: CL.4 (file-organizer)")
	fmt.Println("Previous: CL.3 (subcommands)")
	fmt.Println("---------------------------------------------------")
}

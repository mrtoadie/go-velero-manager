package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
)

const (
	ColReset  = "\033[0m"
	ColRed    = "\033[31m"
	ColGreen  = "\033[32m"
	ColYellow = "\033[33m"
	ColBlue   = "\033[34m"
)

type Backup struct {
	Name        string
	Status      string
	Errors      int
	Warnings    int
	Created     string
	CreatedTime time.Time
	Expires     string
}

func executeCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func fetchBackupsDetailed() ([]Backup, error) {
	output, err := executeCommand("velero", "backup", "get")
	if err != nil {
		fmt.Printf("Warning: %s\n", err)
		// --all-namespaces exists nicht, also einfach weitermachen
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var backups []Backup

	// Regex: NAME STATUS ERRORS WARNINGS CREATED EXPIRES
	// Date format: 2026-08-09 21:54:03 +0200 CEST
	re := regexp.MustCompile(`^(\S+)\s+(Completed|Failed|Partial|Pending)\s+(\d+)\s+(\d+)\s+(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}\s+[+-]\d{4}\s+\w+)\s+(\d+d|\d+h|\d+m)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Ignore headings, empty lines and <none>
		if line == "" ||
			strings.HasPrefix(line, "NAME") ||
			strings.HasPrefix(line, "<none>") ||
			strings.HasPrefix(line, "SELECTOR") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		errors := 0
		warnings := 0
		fmt.Sscanf(matches[3], "%d", &errors)
		fmt.Sscanf(matches[4], "%d", &warnings)

		// Parse timestamps
		createdTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", matches[5])
		if err != nil {
			fmt.Printf("Could not parse date for %s: %s\n", matches[1], matches[5])
			continue
		}

		backups = append(backups, Backup{
			Name:        matches[1],
			Status:      matches[2],
			Errors:      errors,
			Warnings:    warnings,
			Created:     matches[5],
			CreatedTime: createdTime,
			Expires:     matches[6],
		})
	}

	// Sort by creation date (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedTime.After(backups[j].CreatedTime)
	})

	return backups, nil
}

func formatBackupLabel(b Backup) string {
	dateStr := b.CreatedTime.Format("2006-01-02 15:04")

	statusEmoji := "✓"
	if b.Status == "Failed" {
		statusEmoji = "✗"
	} else if b.Status == "Partial" {
		statusEmoji = "⚠️"
	}

	errorIndicator := ""
	if b.Errors > 0 {
		errorIndicator = fmt.Sprintf(" ERR:%d", b.Errors)
	}

	warningIndicator := ""
	if b.Warnings > 0 {
		warningIndicator = fmt.Sprintf(" WRN:%d", b.Warnings)
	}

	label := fmt.Sprintf("%s | %-35s | %-8s | %s | %s%s%s",
		statusEmoji, b.Name, b.Status, dateStr, errorIndicator, warningIndicator)
	return label
}

func main() {
	fmt.Printf("%s=== VELERO-MANAGER ===%s\n", ColBlue, ColReset)
	fmt.Println("\nLoading Backups from Velero...")

	allBackups, err := fetchBackupsDetailed()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if len(allBackups) == 0 {
		fmt.Println("No backups found!")
		fmt.Println("Note: 'velero backup get' runs in the current namespace.")
		os.Exit(0)
	}

	opts := make([]huh.Option[int], len(allBackups))
	for i, b := range allBackups {
		opts[i] = huh.NewOption(formatBackupLabel(b), i)
	}

	var selectedIndexes []int
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[int]().
				Title(fmt.Sprintf("Select backups to delete (%d available)", len(allBackups))).
				Description("Space bar to mark, Enter to confirm").
				Options(opts...).
				Limit(0).
				Value(&selectedIndexes),
		).WithTheme(huh.ThemeBase16()),
	).Run(); err != nil {
		fmt.Println(ColGreen,"\n Bye...")
		os.Exit(0)
	}

	if len(selectedIndexes) == 0 {
		fmt.Println("No backups selected")
		os.Exit(0)
	}

	// Collect selected backups
	var selectedNames []string
	for _, idx := range selectedIndexes {
		if idx >= 0 && idx < len(allBackups) {
			selectedNames = append(selectedNames, allBackups[idx].Name)
		}
	}

	// Confirmation
	backupList := strings.Join(selectedNames, "\n  • ")
	var confirmed bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("%d Really delete backup(s)?", len(selectedNames))).
				Description(fmt.Sprintf("Selected:\n • %s This action is irrevocable!", backupList)).
				Value(&confirmed),
		),
	).Run(); err != nil {
		os.Exit(0)
	}

	if !confirmed {
		fmt.Println("Deletion process aborted")
		os.Exit(0)
	}

	// Batch deletion
	fmt.Printf("\nDelete %d backup(s)...\n\n", len(selectedNames))
	failed := 0
	success := 0

	for _, backupName := range selectedNames {
		fmt.Printf("  [%d/%d] %s... ", success+failed+1, len(selectedNames), backupName)

		if _, err := executeCommand("velero", "backup", "delete", backupName, "--confirm"); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			failed++
		} else {
			fmt.Printf("OK\n")
			success++
		}
	}

	// Summary
	fmt.Printf("\n─────────────────────────────\n")
	if failed == 0 {
		fmt.Printf("All %d backups successfully deleted!\n", success)
	} else {
		fmt.Printf("%d/%d succeeded, %d failed\n", success, len(selectedNames), failed)
		os.Exit(1)
	}
}

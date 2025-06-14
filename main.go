package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/stanley2058/hackmd-edit/lib"
)

func main() {
	ctx := lib.BootstrapApp()

	err := run(&ctx)
	if err != nil {
		lib.LogFatal(fmt.Sprintf("An error occurred: %v", err))
	}

	fmt.Printf("This session used %d tokens.\nBye!", ctx.TokenUsage)
}

func run(ctx *lib.Context) error {
	content, err := lib.GetNoteContent(ctx)
	if err != nil {
		return fmt.Errorf("failed to get note: %w", err)
	}
	ctx.LastStoredContent = content

	tmpFileName, err := lib.CreateTmpFileWithContent(content)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFileName)

	watcher, done, err := lib.CreateFileWatcher(ctx, tmpFileName)
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}
	defer watcher.Close()

	lib.OpenEditor(ctx, tmpFileName)
	slog.Debug("Editor closed.")

	close(*done)

	currentContent, _ := os.ReadFile(tmpFileName)
	if string(currentContent) != ctx.LastStoredContent {
		slog.Debug("Content changed, firing final update.")
		err := lib.DoUpdate(ctx, tmpFileName)
		if err != nil {
			return fmt.Errorf("failed to update note: %w", err)
		}
	}

	return nil
}

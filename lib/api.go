package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func ConstructRequestUrl(ctx *Context) string {
	url := ctx.BaseUrl
	if ctx.TeamPath != "" {
		url += "teams/" + ctx.TeamPath + "/"
	}
	return url + "notes/" + ctx.Note
}

func GetNoteContent(ctx *Context) (string, error) {
	requestUrl := ConstructRequestUrl(ctx)
	slog.Debug("Requesting note", "url", requestUrl)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create GET request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ctx.ApiToken)
	resp, err := ctx.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("server returned an error status: %s", resp.Status)
	}

	ctx.TokenUsage++

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var note Note
	err = json.Unmarshal(body, &note)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return note.Content, nil
}

func CreateTmpFileWithContent(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "hackmd-edit-*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write to temporary file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temporary file after writing: %w", err)
	}
	slog.Debug("Content saved to temporary file", "path", tmpFile.Name())
	return tmpFile.Name(), nil
}

func CreateFileWatcher(ctx *Context, filePath string) (*fsnotify.Watcher, *chan struct{}, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	done := make(chan struct{})

	go watchForSaves(watcher, ctx, filePath, done)

	if err = watcher.Add(filePath); err != nil {
		return nil, nil, fmt.Errorf("failed to add file to watcher: %w", err)
	}
	return watcher, &done, nil
}

func OpenEditor(ctx *Context, filePath string) error {
	slog.Debug("Opening file in editor", "file", filePath, "editor", ctx.Editor)
	cmd := exec.Command(ctx.Editor, filePath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		slog.Debug("Editor command finished with error", "error", err)
	}

	return nil
}

func watchForSaves(watcher *fsnotify.Watcher, ctx *Context, filePath string, done <-chan struct{}) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				slog.Debug("File saved, preparing to POST update...")
				err := DoUpdate(ctx, filePath)
				if err != nil {
					slog.Error("Failed to update note", "error", err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			slog.Error("Watcher error", "error", err)
		case <-done:
			slog.Debug("Stopping file watcher.")
			return
		}
	}
}

func DoUpdate(ctx *Context, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	req, err := http.NewRequest("PATCH", ConstructRequestUrl(ctx), bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to create PATCH request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", "Bearer "+ctx.ApiToken)
	resp, err := ctx.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned an error status: %s", resp.Status)
	}

	ctx.TokenUsage++
	ctx.LastStoredContent = string(content)

	return nil
}

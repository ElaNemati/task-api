package logger

import (
	"log/slog"
	"os"
)

func openOrCreateFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

func Init() (func(), error) {
	if err := os.MkdirAll("./logs", 0755); err != nil {
		return nil, err
	}

	logFiles := map[slog.Level]string{
		slog.LevelDebug: "./logs/debug.log",
		slog.LevelInfo:  "./logs/info.log",
		slog.LevelWarn:  "./logs/warn.log",
		slog.LevelError: "./logs/error.log",
	}

	var openedFiles []*os.File

	closeAll := func() {
		for _, f := range openedFiles {
			f.Close()
		}
	}

	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	fileHandlers := make(map[slog.Level]slog.Handler, len(logFiles))

	for level, path := range logFiles {
		file, err := openOrCreateFile(path)
		if err != nil {
			closeAll()
			return nil, err
		}
		openedFiles = append(openedFiles, file)
		fileHandlers[level] = slog.NewJSONHandler(file, opts)
	}

	otherFile, err := openOrCreateFile("./logs/other.log")
	if err != nil {
		closeAll()
		return nil, err
	}
	openedFiles = append(openedFiles, otherFile)

	router := &levelRouter{
		fileHandlers: fileHandlers,
		fallback:     slog.NewJSONHandler(otherFile, opts),
	}

	slog.SetDefault(slog.New(router))

	return closeAll, nil
}

package logger

import (
	"context"
	"log/slog"
)

type levelRouter struct {
	fileHandlers map[slog.Level]slog.Handler
	fallback     slog.Handler
}

func (r *levelRouter) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (r *levelRouter) Handle(ctx context.Context, record slog.Record) error {
	if handler, exists := r.fileHandlers[record.Level]; exists {
		return handler.Handle(ctx, record.Clone())
	}
	return r.fallback.Handle(ctx, record.Clone())
}

func (r *levelRouter) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make(map[slog.Level]slog.Handler, len(r.fileHandlers))
	for level, handler := range r.fileHandlers {
		newHandlers[level] = handler.WithAttrs(attrs)
	}
	return &levelRouter{
		fileHandlers: newHandlers,
		fallback:     r.fallback.WithAttrs(attrs),
	}
}

func (r *levelRouter) WithGroup(name string) slog.Handler {
	newHandlers := make(map[slog.Level]slog.Handler, len(r.fileHandlers))
	for level, handler := range r.fileHandlers {
		newHandlers[level] = handler.WithGroup(name)
	}
	return &levelRouter{
		fileHandlers: newHandlers,
		fallback:     r.fallback.WithGroup(name),
	}
}

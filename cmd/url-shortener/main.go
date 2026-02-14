package main

import (
	"github.com/famineBurgund/famiURL/internal/config"
)

func main() {
	cfg := config.MustLoad()
	// TODO: init config: cleanenv

	// TODO: init logger: log/slog

	// TODO: init storage: sqlite

	// TODO: init router: chi, render

	// TODO: run server
}

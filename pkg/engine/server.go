package engine

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServerWithGracefulShutdown(app *echo.Echo, port string) {
	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Run the server in a goroutine so that it doesn't block
	go func() {
		if err := app.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for an interrupt signal
	<-quit
	app.Logger.Info("Shutting down server...")

	// The context is used to inform the server it has 30 seconds to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatal(err)
	}

	app.Logger.Info("Server exited gracefully")
}

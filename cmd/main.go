package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/orchestra-mcp/sdk-go/plugin"
	"github.com/orchestra-mcp/plugin-ai-vision/internal"
)

func main() {
	builder := plugin.New("ai.vision").
		Version("0.1.0").
		Description("Vision tools for image analysis, text extraction, and screen understanding").
		Author("Orchestra").
		Binary("ai-vision")

	tp := &internal.ToolsPlugin{}
	tp.RegisterTools(builder)

	p := builder.BuildWithTools()
	p.ParseFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	if err := p.Run(ctx); err != nil {
		log.Fatalf("ai.vision: %v", err)
	}
}

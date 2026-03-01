package aivision

import (
	"github.com/orchestra-mcp/plugin-ai-vision/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// Register adds all vision tools to the builder.
func Register(builder *plugin.PluginBuilder) {
	tp := &internal.ToolsPlugin{}
	tp.RegisterTools(builder)
}

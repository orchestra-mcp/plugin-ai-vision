package internal

import (
	"github.com/orchestra-mcp/sdk-go/plugin"
	"github.com/orchestra-mcp/plugin-ai-vision/internal/tools"
)

// ToolsPlugin registers all ai.vision tools with the plugin builder.
type ToolsPlugin struct{}

// RegisterTools registers all 6 vision tools on the given plugin builder.
func (tp *ToolsPlugin) RegisterTools(builder *plugin.PluginBuilder) {
	builder.RegisterTool("analyze_image",
		"Analyze an image and describe its contents using Claude's vision",
		tools.AnalyzeImageSchema(), tools.AnalyzeImage())

	builder.RegisterTool("extract_text",
		"Extract all visible text from an image (OCR)",
		tools.ExtractTextSchema(), tools.ExtractText())

	builder.RegisterTool("find_elements",
		"Find and list UI elements or specific element types in an image",
		tools.FindElementsSchema(), tools.FindElements())

	builder.RegisterTool("describe_screen",
		"Describe a screen or interface in detail including app, visible content, and interactive elements",
		tools.DescribeScreenSchema(), tools.DescribeScreen())

	builder.RegisterTool("compare_images",
		"Compare two images and identify similarities and differences",
		tools.CompareImagesSchema(), tools.CompareImages())

	builder.RegisterTool("extract_data",
		"Extract structured data from an image as JSON",
		tools.ExtractDataSchema(), tools.ExtractData())
}

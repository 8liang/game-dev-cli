package cmd

import (
	"embed"
)

// ExampleProtoFS 嵌入了 proto 样例源文件。import 依赖由 proto setup --install 自动下载。
//
//go:embed example_proto_static/error.proto
//go:embed example_proto_static/messages.proto
//go:embed example_proto_static/svc/agent.svc.proto
var ExampleProtoFS embed.FS

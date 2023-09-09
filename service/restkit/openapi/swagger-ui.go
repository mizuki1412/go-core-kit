package openapi

import "embed"

//go:embed swagger-ui
var UiAssets embed.FS

//go:embed knife-ui
var KUiAssets embed.FS

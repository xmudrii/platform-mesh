package config

import "github.com/spf13/pflag"

type KcpConfig struct {
	ApiExportEndpointSliceName string
}

type OperatorConfig struct {
	Kcp       KcpConfig
	Namespace string
}

func NewOperatorConfig() OperatorConfig {
	return OperatorConfig{
		Kcp: KcpConfig{
			ApiExportEndpointSliceName: "backup.platform-mesh.io",
		},
		Namespace: "platform-mesh",
	}
}

func (c *OperatorConfig) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Kcp.ApiExportEndpointSliceName, "kcp-api-export-endpoint-slice-name", c.Kcp.ApiExportEndpointSliceName, "Set APIExportEndpointSlice name")
	fs.StringVar(&c.Namespace, "namespace", c.Namespace, "Namespace in which the operator manages resources")
}

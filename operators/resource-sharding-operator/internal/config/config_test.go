package config

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOperatorConfig_Defaults(t *testing.T) {
	cfg := NewOperatorConfig()

	assert.False(t, cfg.Kcp.Enabled, "KCP should be disabled by default")
	assert.Equal(t, "resource-sharding", cfg.Kcp.ApiExportEndpointSliceName)
	assert.True(t, cfg.WebhookEnabled, "webhook should be enabled by default")
}

func TestAddFlags_BindsAllFlags(t *testing.T) {
	cfg := NewOperatorConfig()
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	cfg.AddFlags(flags)

	require.NoError(t, flags.Parse([]string{
		"--kcp-enabled=true",
		"--kcp-api-export-endpoint-slice-name=custom-slice",
		"--webhook-enabled=false",
	}))

	assert.True(t, cfg.Kcp.Enabled)
	assert.Equal(t, "custom-slice", cfg.Kcp.ApiExportEndpointSliceName)
	assert.False(t, cfg.WebhookEnabled)
}

func TestAddFlags_DefaultValuesPreservedWhenNotSet(t *testing.T) {
	cfg := NewOperatorConfig()
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	cfg.AddFlags(flags)
	require.NoError(t, flags.Parse([]string{}))

	assert.False(t, cfg.Kcp.Enabled)
	assert.Equal(t, "resource-sharding", cfg.Kcp.ApiExportEndpointSliceName)
	assert.True(t, cfg.WebhookEnabled)
}

func TestAddFlags_FlagNames(t *testing.T) {
	cfg := NewOperatorConfig()
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	cfg.AddFlags(flags)

	assert.NotNil(t, flags.Lookup("kcp-enabled"), "kcp-enabled flag should be registered")
	assert.NotNil(t, flags.Lookup("kcp-api-export-endpoint-slice-name"), "kcp-api-export-endpoint-slice-name flag should be registered")
	assert.NotNil(t, flags.Lookup("webhook-enabled"), "webhook-enabled flag should be registered")
}

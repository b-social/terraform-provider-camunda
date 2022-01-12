package camunda

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"os"
	"terraform-provider-camunda/camunda/client"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var stderr = os.Stderr

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     *client.Client
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"endpoint": {
				Type:     types.StringType,
				Required: true,
			},
			"username": {
				Type:     types.StringType,
				Optional: true,
			},
			"password": {
				Type:      types.StringType,
				Optional:  true,
				Sensitive: true,
			},
			"insecure_skip_verify": {
				Type:     types.BoolType,
				Optional: true,
			},
			"tls_certificate": {
				Type:     types.StringType,
				Optional: true,
			},
			"tls_key": {
				Type:     types.StringType,
				Optional: true,
			},
			"tls_ca": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	Endpoint           types.String `tfsdk:"endpoint"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure_skip_verify"`
	TlsCertificate     types.String `tfsdk:"tls_certificate"`
	TlsKey             types.String `tfsdk:"tls_key"`
	TlsCA                 types.String `tfsdk:"tls_ca"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide a user to the provider
	var username string
	if config.Username.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as username",
		)
		return
	}

	if config.Username.Null {
		username = os.Getenv("CAMUNDA_USERNAME")
	} else {
		username = config.Username.Value
	}

	// User must provide a password to the provider
	var password string
	if config.Password.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as password",
		)
		return
	}

	if config.Password.Null {
		password = os.Getenv("CAMUNDA_PASSWORD")
	} else {
		password = config.Password.Value
	}

	// User must specify a endpoint
	var endpoint string
	if config.Endpoint.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as endpoint",
		)
		return
	}

	if config.Endpoint.Null {
		endpoint = os.Getenv("CAMUNDA_ENDPOINT")
	} else {
		endpoint = config.Endpoint.Value
	}

	if endpoint == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find Endpoint",
			"Endpoint cannot be an empty string",
		)
		return
	}

	// Create a new Camunda client and set it to the provider client
	c := client.NewClient(
		client.ClientOptions{
			EndpointUrl: endpoint,
			ApiUser:     username,
			ApiPassword: password,
			Timeout:     time.Second * 30,
		})

	tlsConfig := &tls.Config{
	}

	if !config.InsecureSkipVerify.Null && !config.InsecureSkipVerify.Unknown {
		tlsConfig.InsecureSkipVerify = config.InsecureSkipVerify.Value
	}

	if !config.TlsCertificate.Unknown && !config.TlsCertificate.Null && !config.TlsKey.Unknown && !config.TlsKey.Null {
		pair, err := tls.X509KeyPair([]byte(config.TlsCertificate.Value), []byte(config.TlsKey.Value))
		if err != nil {
			// Error vs warning - empty value must stop execution
			resp.Diagnostics.AddError(
				"Unable to parse certificates",
				err.Error(),
			)
			return
		}
		tlsConfig.Certificates = []tls.Certificate{pair}
	}

	if !config.TlsCA.Unknown && !config.TlsCA.Null {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM([]byte(config.TlsCA.Value))
		tlsConfig.ClientCAs = certPool
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		Proxy:           http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c.SetCustomTransport(transport)

	p.client = c
	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"camunda_deployment": resourceDeploymentType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}

// Package dqlprovider is a Terraform provider whose functions compile DQL queries.
package dqlprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.ProviderWithFunctions = (*dqlProvider)(nil)

type dqlProvider struct {
	version string
}

// New answers the provider constructor for one release version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &dqlProvider{version: version}
	}
}

func (dql *dqlProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "dql"
	response.Version = dql.version
}

func (dql *dqlProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Compiles DQL queries into each backend's native query.",
	}
}

func (dql *dqlProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

func (dql *dqlProvider) Resources(context.Context) []func() resource.Resource { return nil }

func (dql *dqlProvider) DataSources(context.Context) []func() datasource.DataSource { return nil }

func (dql *dqlProvider) Functions(context.Context) []func() function.Function {
	return []func() function.Function{
		newOpenSearchFunction,
		newPromQLFunction,
	}
}

package dqlprovider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/litebitus/dql-compiler/opensearchbackend"
)

var _ function.Function = openSearchFunction{}

type openSearchFunction struct{}

func newOpenSearchFunction() function.Function { return openSearchFunction{} }

type openSearchConfig struct {
	TimeField  types.String `tfsdk:"time_field"`
	GroupLimit types.Int64  `tfsdk:"group_limit"`
}

func (openSearchFunction) Metadata(_ context.Context, _ function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "opensearch"
}

func (openSearchFunction) Definition(_ context.Context, _ function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Summary:             "Compiles a DQL query into an OpenSearch search request body.",
		MarkdownDescription: "Compiles a DQL query, written in the stated grammar version or the latest, into an OpenSearch search request body, as a JSON string.",
		Parameters: []function.Parameter{
			queryParameter,
			function.ObjectParameter{
				Name:                "config",
				MarkdownDescription: "`time_field`: the field a time function buckets by. `group_limit`: how many groups each grouping field answers.",
				AttributeTypes: map[string]attr.Type{
					"group_limit": types.Int64Type,
					"time_field":  types.StringType,
				},
			},
		},
		VariadicParameter: grammarVersionParameter,
		Return:            function.StringReturn{},
	}
}

const openSearchVersionArgument = 2

func (openSearchFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var query string
	var config openSearchConfig
	var grammarVersions []string
	if response.Error = request.Arguments.Get(ctx, &query, &config, &grammarVersions); response.Error != nil {
		return
	}
	parsed, funcError := parsedQuery(query, grammarVersions, openSearchVersionArgument)
	if funcError != nil {
		response.Error = funcError
		return
	}
	body, err := opensearchbackend.CompileQuery(parsed, opensearchbackend.BackendConfig{
		TimeField:  config.TimeField.ValueString(),
		GroupLimit: int(config.GroupLimit.ValueInt64()),
	})
	if err != nil {
		response.Error = function.NewArgumentFuncError(queryArgument, err.Error())
		return
	}
	var encoded strings.Builder
	encoder := json.NewEncoder(&encoded)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(body); err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}
	response.Error = response.Result.Set(ctx, strings.TrimSuffix(encoded.String(), "\n"))
}

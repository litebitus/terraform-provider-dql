package dqlprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/litebitus/dql-compiler/promqlbackend"
)

var _ function.Function = promQLFunction{}

type promQLFunction struct{}

func newPromQLFunction() function.Function { return promQLFunction{} }

func (promQLFunction) Metadata(_ context.Context, _ function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "promql"
}

func (promQLFunction) Definition(_ context.Context, _ function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Summary:             "Compiles a DQL query into a PromQL expression.",
		MarkdownDescription: "Compiles a DQL query, written in the stated grammar version or the latest, into a PromQL expression.",
		Parameters:          []function.Parameter{queryParameter},
		VariadicParameter:   grammarVersionParameter,
		Return:              function.StringReturn{},
	}
}

const promQLVersionArgument = 1

func (promQLFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var query string
	var grammarVersions []string
	if response.Error = request.Arguments.Get(ctx, &query, &grammarVersions); response.Error != nil {
		return
	}
	parsed, funcError := parsedQuery(query, grammarVersions, promQLVersionArgument)
	if funcError != nil {
		response.Error = funcError
		return
	}
	expression, err := promqlbackend.CompileQuery(parsed)
	if err != nil {
		response.Error = function.NewArgumentFuncError(queryArgument, err.Error())
		return
	}
	response.Error = response.Result.Set(ctx, expression)
}

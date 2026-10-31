package dqlprovider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/litebitus/dql-compiler/dqlsyntax"
)

const queryArgument = 0

var (
	queryParameter = function.StringParameter{
		Name:                "query",
		MarkdownDescription: "The DQL query.",
	}
	grammarVersionParameter = function.StringParameter{
		Name:                "grammar_version",
		MarkdownDescription: "The grammar version the query is written in, `<major>.<minor>`; omitted, the latest the provider supports.",
	}
)

// parsedQuery reads a query in the stated grammar version, or the latest, blaming the argument that is wrong.
func parsedQuery(query string, grammarVersions []string, versionArgument int64) (dqlsyntax.ParsedQuery, *function.FuncError) {
	supported := dqlsyntax.SupportedGrammarVersions()
	grammarVersion := supported[len(supported)-1]
	switch len(grammarVersions) {
	case 0:
	case 1:
		grammarVersion = grammarVersions[0]
	default:
		return dqlsyntax.ParsedQuery{}, function.NewArgumentFuncError(versionArgument+1,
			fmt.Sprintf("a query is written in one grammar version, not %d", len(grammarVersions)))
	}
	if _, err := dqlsyntax.ParseQuery("", grammarVersion); err != nil {
		return dqlsyntax.ParsedQuery{}, function.NewArgumentFuncError(versionArgument, err.Error())
	}
	parsed, err := dqlsyntax.ParseQuery(query, grammarVersion)
	if err != nil {
		return dqlsyntax.ParsedQuery{}, function.NewArgumentFuncError(queryArgument, err.Error())
	}
	return parsed, nil
}

package dqlprovider

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var mockOpenSearchConfig = types.ObjectValueMust(
	map[string]attr.Type{"group_limit": types.Int64Type, "time_field": types.StringType},
	map[string]attr.Value{"group_limit": types.Int64Value(50), "time_field": types.StringValue("mock_time")},
)

func run(t *testing.T, called function.Function, arguments ...attr.Value) (string, *function.FuncError) {
	t.Helper()
	response := &function.RunResponse{Result: function.NewResultData(types.StringUnknown())}
	called.Run(context.Background(), function.RunRequest{Arguments: function.NewArgumentsData(arguments)}, response)
	if response.Error != nil {
		return "", response.Error
	}
	return response.Result.Value().(types.String).ValueString(), nil
}

// grammarVersions is the variadic argument as Terraform passes it: a tuple of the versions given.
func grammarVersions(versions ...string) types.Tuple {
	elementTypes := make([]attr.Type, len(versions))
	values := make([]attr.Value, len(versions))
	for index, version := range versions {
		elementTypes[index], values[index] = types.StringType, types.StringValue(version)
	}
	return types.TupleValueMust(elementTypes, values)
}

func TestPromQLCompilesAQueryInTheLatestOrTheStatedVersion(t *testing.T) {
	for name, versions := range map[string]types.Tuple{
		"latest": grammarVersions(),
		"stated": grammarVersions("1.0"),
	} {
		got, err := run(t, promQLFunction{}, types.StringValue("__name__:mock_up | rate(5m) | sum() by job"), versions)
		if want := "sum by (job) (rate(mock_up[5m]))"; err != nil || got != want {
			t.Errorf("%s: got %q (%v), want %q", name, got, err, want)
		}
	}
}

func TestOpenSearchCompilesAQueryToItsRequestBody(t *testing.T) {
	got, err := run(t, openSearchFunction{},
		types.StringValue("app:mock | count() by app | where(value > 5)"), mockOpenSearchConfig, grammarVersions())
	want := `{"aggs":{"by_app":{"aggs":{"where":{"bucket_selector":{"buckets_path":{"value":"_count"},"script":"params.value > 5"}}},"terms":{"field":"app","size":50}}},"query":{"match":{"app":"mock"}},"size":0}`
	if err != nil || got != want {
		t.Errorf("got %s (%v),\nwant %s", got, err, want)
	}
}

func TestARefusalBlamesTheArgumentThatIsWrong(t *testing.T) {
	for name, testCase := range map[string]struct {
		query    string
		versions types.Tuple
		argument int64
		says     string
	}{
		"unsupported version": {"app:mock", grammarVersions("2.0"), promQLVersionArgument, "grammar version 2.0 is not supported; supported: 1.0"},
		"two versions":        {"app:mock", grammarVersions("1.0", "1.0"), promQLVersionArgument + 1, "a query is written in one grammar version, not 2"},
		"malformed query":     {"(app:mock", grammarVersions(), queryArgument, "a parenthesis is not closed"},
		"uncompilable query":  {"mock", grammarVersions(), queryArgument, "a term with no field does not compile to PromQL: a series carries labels, not text"},
	} {
		_, err := run(t, promQLFunction{}, types.StringValue(testCase.query), testCase.versions)
		if err == nil || err.Text != testCase.says || err.FunctionArgument == nil || *err.FunctionArgument != testCase.argument {
			t.Errorf("%s: err = %+v, want %q on argument %d", name, err, testCase.says, testCase.argument)
		}
	}
}

func TestOpenSearchBlamesAnUnsupportedVersionOnItsOwnArgument(t *testing.T) {
	_, err := run(t, openSearchFunction{}, types.StringValue("app:mock"), mockOpenSearchConfig, grammarVersions("2.0"))
	if err == nil || err.FunctionArgument == nil || *err.FunctionArgument != openSearchVersionArgument {
		t.Errorf("err = %+v, want it on argument %d", err, openSearchVersionArgument)
	}
}

var protocolProviders = map[string]func() (tfprotov6.ProviderServer, error){
	"dql": providerserver.NewProtocol6WithError(New("test")()),
}

func TestTheFunctionsAnswerInAPlan(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protocolProviders,
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_8_0)},
		Steps: []resource.TestStep{{
			Config: `
output "promql" {
  value = provider::dql::promql("__name__:mock_up and job:mock-api")
}

output "promql_stated_version" {
  value = provider::dql::promql("__name__:mock_up and job:mock-api", "1.0")
}

output "opensearch" {
  value = jsondecode(provider::dql::opensearch("app:mock", { time_field = "mock_time", group_limit = 50 }))
}
`,
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue("promql", knownvalue.StringExact(`mock_up{job="mock-api"}`)),
				statecheck.ExpectKnownOutputValue("promql_stated_version", knownvalue.StringExact(`mock_up{job="mock-api"}`)),
				statecheck.ExpectKnownOutputValue("opensearch", knownvalue.ObjectExact(map[string]knownvalue.Check{
					"query": knownvalue.ObjectExact(map[string]knownvalue.Check{
						"match": knownvalue.ObjectExact(map[string]knownvalue.Check{"app": knownvalue.StringExact("mock")}),
					}),
				})),
			},
		}},
	})
}

func TestARefusalFailsThePlan(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protocolProviders,
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_8_0)},
		Steps: []resource.TestStep{{
			Config: `
output "promql" {
  value = provider::dql::promql("(app:mock")
}
`,
			ExpectError: regexp.MustCompile(`a parenthesis is not closed`),
		}},
	})
}

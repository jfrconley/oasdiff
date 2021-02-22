package diff

import "github.com/getkin/kin-openapi/openapi3"

// Diff describes changes between two OAS specs including a summary of these changes
type Diff struct {
	SpecDiff *SpecDiff `json:"spec,omitempty"`
	Summary  *Summary  `json:"summary,omitempty"`
}

/*
Get calculates the diff between two OAS specs.
Prefix is an optional path prefix that exists in s1 paths but not in s2.
If filter isn't empty, the diff will only include paths that match this regex.
*/
func Get(s1, s2 *openapi3.Swagger, prefix string, filter string) Diff {
	diff := getDiff(s1, s2, prefix)
	diff.filterByRegex(filter)

	return Diff{
		SpecDiff: diff,
		Summary:  diff.getSummary(),
	}
}

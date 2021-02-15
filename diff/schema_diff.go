package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type SchemaDiff struct {
	SchemaAdded                     bool       `json:"schemaAdded,omitempty"`
	SchemaDeleted                   bool       `json:"schemaDelete,omitempty"`
	ValueAdded                      bool       `json:"valueAdded,omitempty"`
	ValueDeleted                    bool       `json:"valueDeleted,omitempty"`
	OneOfDiff                       bool       `json:"oneOfDiff,omitempty"`
	AnyOfDiff                       bool       `json:"anyOfDiff,omitempty"`
	AllOfDiff                       bool       `json:"allOfDiff,omitempty"`
	NotDiff                         bool       `json:"notDiff,omitempty"`
	TypeDiff                        *ValueDiff `json:"typeDiff,omitempty"`
	TitleDiff                       *ValueDiff `json:"titleDiff,omitempty"`
	FormatDiff                      *ValueDiff `json:"formatDiff,omitempty"`
	DescriptionDiff                 *ValueDiff `json:"descriptionDiff,omitempty"`
	EnumDiff                        bool       `json:"enumDiff,omitempty"`
	AdditionalPropertiesAllowedDiff *ValueDiff `json:"additionalPropertiesAllowedDiff,omitempty"`
	UniqueItemsDiff                 *ValueDiff `json:"uniqueItemsDiff,omitempty"`
	ExclusiveMinDiff                *ValueDiff `json:"exclusiveMinDiff,omitempty"`
	ExclusiveMaxDiff                *ValueDiff `json:"exclusiveMaxDiff,omitempty"`
	NullableDiff                    *ValueDiff `json:"nullableDiff,omitempty"`
	ReadOnlyDiff                    *ValueDiff `json:"readOnlyDiffDiff,omitempty"`
	WriteOnlyDiff                   *ValueDiff `json:"writeOnlyDiffDiff,omitempty"`
	AllowEmptyValueDiff             *ValueDiff `json:"allowEmptyValueDiff,omitempty"`
	DeprecatedDiff                  *ValueDiff `json:"deprecatedDiff,omitempty"`
	MinDiff                         *ValueDiff `json:"minDiff,omitempty"`
	MaxDiff                         *ValueDiff `json:"maxDiff,omitempty"`
	MultipleOf                      *ValueDiff `json:"multipleOfDiff,omitempty"`
	PropertiesDiff                  bool       `json:"propertiesDiff,omitempty"`
}

func (schemaDiff SchemaDiff) empty() bool {
	return schemaDiff == SchemaDiff{}
}

func diffSchema(schema1 *openapi3.SchemaRef, schema2 *openapi3.SchemaRef) SchemaDiff {

	value1, value2, status := getSchemaValues(schema1, schema2)

	if status != schemaStatusOK {
		return getSchemaDiff(status)
	}

	result := SchemaDiff{}

	// ExtensionProps
	result.OneOfDiff = getDiffSchemas(value1.OneOf, value2.OneOf)
	result.AnyOfDiff = getDiffSchemas(value1.AnyOf, value2.AnyOf)
	result.AllOfDiff = getDiffSchemas(value1.AllOf, value2.AllOf)
	result.NotDiff = !diffSchema(value1.Not, value2.Not).empty()
	result.TypeDiff = getValueDiff(value1.Type, value2.Type)
	result.TitleDiff = getValueDiff(value1.Title, value2.Title)
	result.FormatDiff = getValueDiff(value1.Format, value2.Format)
	result.DescriptionDiff = getValueDiff(value1.Description, value2.Description)
	result.EnumDiff = getEnumDiff(value1.Enum, value2.Enum)
	// Default
	// Example
	// ExternalDocs
	result.AdditionalPropertiesAllowedDiff = getBoolRefDiff(value1.AdditionalPropertiesAllowed, value2.AdditionalPropertiesAllowed)
	result.UniqueItemsDiff = getValueDiff(value1.UniqueItems, value2.UniqueItems)
	result.ExclusiveMinDiff = getValueDiff(value1.ExclusiveMin, value2.ExclusiveMin)
	result.ExclusiveMaxDiff = getValueDiff(value1.ExclusiveMax, value2.ExclusiveMax)
	result.NullableDiff = getValueDiff(value1.Nullable, value2.Nullable)
	result.ReadOnlyDiff = getValueDiff(value1.ReadOnly, value2.ReadOnly)
	result.WriteOnlyDiff = getValueDiff(value1.WriteOnly, value2.WriteOnly)
	result.AllowEmptyValueDiff = getValueDiff(value1.AllowEmptyValue, value2.AllowEmptyValue)
	// XML
	result.DeprecatedDiff = getValueDiff(value1.Deprecated, value2.Deprecated)
	result.MinDiff = getFloat64RefDiff(value1.Min, value2.Min)
	result.MaxDiff = getFloat64RefDiff(value1.Max, value2.Max)
	result.MultipleOf = getFloat64RefDiff(value1.MultipleOf, value2.MultipleOf)
	// MultipleOf
	// MinLength
	// MaxLength
	// Pattern
	// compiledPattern
	// MinItems
	// MaxItems
	// Items
	// Required
	result.PropertiesDiff = getDiffSchemaMap(value1.Properties, value2.Properties)
	// MinProps
	// MaxProps
	// AdditionalProperties
	// Discriminator

	return result
}

type schemaStatus int

const (
	schemaStatusOK schemaStatus = iota
	schemaStatusNoSchemas
	schemaStatusSchemaAdded
	schemaStatusSchemaDeleted
	schemaStatusNoValues
	schemaStatusValueAdded
	schemaStatusValueDeleted
)

func getSchemaValues(schema1 *openapi3.SchemaRef, schema2 *openapi3.SchemaRef) (*openapi3.Schema, *openapi3.Schema, schemaStatus) {

	if schema1 == nil && schema2 == nil {
		return nil, nil, schemaStatusNoSchemas
	}

	if schema1 == nil && schema2 != nil {
		return nil, nil, schemaStatusSchemaAdded
	}

	if schema1 != nil && schema2 == nil {
		return nil, nil, schemaStatusSchemaDeleted
	}

	value1 := schema1.Value
	value2 := schema2.Value

	if value1 == nil && value2 == nil {
		return nil, nil, schemaStatusNoValues
	}

	if value1 == nil && value2 != nil {
		return nil, nil, schemaStatusValueAdded
	}

	if value1 != nil && value2 == nil {
		return nil, nil, schemaStatusValueDeleted
	}

	return value1, value2, schemaStatusOK
}

func getSchemaDiff(status schemaStatus) SchemaDiff {
	switch status {
	case schemaStatusSchemaAdded:
		return SchemaDiff{SchemaAdded: true}
	case schemaStatusSchemaDeleted:
		return SchemaDiff{SchemaDeleted: true}
	case schemaStatusValueAdded:
		return SchemaDiff{ValueAdded: true}
	case schemaStatusValueDeleted:
		return SchemaDiff{ValueDeleted: true}
	}

	// all other cases -> empty diff
	return SchemaDiff{}
}

func getDiffSchemaMap(schemas1 openapi3.Schemas, schemas2 openapi3.Schemas) bool {

	return !schemaMapContained(schemas1, schemas2) || !schemaMapContained(schemas2, schemas1)
}

func schemaMapContained(schemas1 openapi3.Schemas, schemas2 openapi3.Schemas) bool {
	for schemaName1, schemaRef1 := range schemas1 {
		schemaRef2, ok := schemas2[schemaName1]
		if !ok {
			return false
		}

		if diff := diffSchema(schemaRef1, schemaRef2); !diff.empty() {
			return false
		}
	}

	return true
}

func getDiffSchemas(schemaRefs1 openapi3.SchemaRefs, schemaRefs2 openapi3.SchemaRefs) bool {

	return !schemaRefsContained(schemaRefs1, schemaRefs2) || !schemaRefsContained(schemaRefs2, schemaRefs1)
}

func schemaRefsContained(schemaRefs1 openapi3.SchemaRefs, schemaRefs2 openapi3.SchemaRefs) bool {
	for _, schemaRef1 := range schemaRefs1 {
		if schemaRef1 != nil && schemaRef1.Value != nil {
			if !findSchema(schemaRef1, schemaRefs2) {
				return false
			}
		}
	}
	return true
}

func findSchema(schemaRef1 *openapi3.SchemaRef, schemaRefs2 openapi3.SchemaRefs) bool {
	// TODO: optimize with a map
	for _, schemaRef2 := range schemaRefs2 {
		if schemaRef2 == nil || schemaRef2.Value == nil {
			continue
		}

		if diff := diffSchema(schemaRef1, schemaRef2); diff.empty() {
			return true
		}
	}

	return false
}

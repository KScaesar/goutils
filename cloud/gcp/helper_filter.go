package gcp

import (
	"fmt"
	"strings"
)

func filterStringForInsertOpsType(resourceID string) string {
	var b filterStringBuilder
	return b.
		str("operationType", "=", "insert").
		and().
		str("targetLink", "=", resourceID).
		build()
}

type filterStringBuilder struct {
	builder *strings.Builder
}

// comparison operator possible value:
// =, !=, >, <, <=, >=
func (b filterStringBuilder) number(key string, comp string, value interface{}) filterStringBuilder {
	if b.builder == nil {
		b.builder = &strings.Builder{}
	}

	// gcp filter string rule:
	// The comparison operator must be either =, !=, >, <
	switch {
	case comp == "<=":
		b.builder.WriteString("( ")
		b.number(key, "=", value).or().number(key, "<", value)
		b.builder.WriteString(") ")

	case comp == ">=":
		b.builder.WriteString("( ")
		b.number(key, "=", value).or().number(key, ">", value)
		b.builder.WriteString(") ")

	default:
		pair := fmt.Sprintf(`%v %v %v `, key, comp, value)
		b.builder.WriteString(pair)
	}
	return b
}

func (b filterStringBuilder) bool(key string, value bool) filterStringBuilder {
	if b.builder == nil {
		b.builder = &strings.Builder{}
	}

	if value {
		pair := fmt.Sprintf(`%v = true `, key)
		b.builder.WriteString(pair)
	} else {
		pair := fmt.Sprintf(`%v = false `, key)
		b.builder.WriteString(pair)
	}
	return b
}

// comparison operator possible vale:
// =, !=
func (b filterStringBuilder) str(key string, comp string, value string) filterStringBuilder {
	if b.builder == nil {
		b.builder = &strings.Builder{}
	}

	pair := fmt.Sprintf(`%v %v "%v" `, key, comp, value)
	b.builder.WriteString(pair)
	return b
}

func (b filterStringBuilder) or() filterStringBuilder {
	if b.builder == nil {
		b.builder = &strings.Builder{}
	}

	// 用小寫會查不到
	b.builder.WriteString("OR ")
	return b
}

func (b filterStringBuilder) and() filterStringBuilder {
	if b.builder == nil {
		b.builder = &strings.Builder{}
	}

	// 用小寫會查不到
	b.builder.WriteString("AND ")
	return b
}

func (b filterStringBuilder) build() string {
	if b.builder == nil {
		b.builder = &strings.Builder{}
	}

	return b.builder.String()
}

package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_filterStringBuilder_number_equal(t *testing.T) {
	var b filterStringBuilder
	filterString := b.number("size", "=", 13).build()
	assert.Equal(t, "size = 13 ", filterString)
}

func Test_filterStringBuilder_number_gte(t *testing.T) {
	var b filterStringBuilder
	filterString := b.number("size", ">=", 13).build()
	assert.Equal(t, "( size = 13 OR size > 13 ) ", filterString)
}

func Test_filterStringBuilder_number_lte(t *testing.T) {
	var b filterStringBuilder
	filterString := b.number("size", "<=", 13).build()
	assert.Equal(t, "( size = 13 OR size < 13 ) ", filterString)
}

func Test_filterStringBuilder_bool(t *testing.T) {
	var b filterStringBuilder
	filterString := b.bool("isDefault", true).build()
	assert.Equal(t, "isDefault = true ", filterString)
}

func Test_filterStringBuilder_str(t *testing.T) {
	var b filterStringBuilder
	filterString := b.str("id", "=", "5456663258773368149").build()
	assert.Equal(t, `id = "5456663258773368149" `, filterString)
}

func Test_filterStringBuilder_or(t *testing.T) {
	var b filterStringBuilder
	filterString := b.
		str("destRange", "!=", "0.0.0.0/0").
		or().
		bool("isDefault", false).
		build()
	assert.Equal(t, `destRange != "0.0.0.0/0" OR isDefault = false `, filterString)
}

func Test_filterStringBuilder_and(t *testing.T) {
	resourceID := "https://www.googleapis.com/compute/v1/projects/goTest/global/networks/aaaaaaaaa"
	filterString := filterStringForInsertOpsType(resourceID)
	assert.Equal(
		t,
		`operationType = "insert" AND targetLink = "https://www.googleapis.com/compute/v1/projects/goTest/global/networks/aaaaaaaaa" `,
		filterString,
	)
}

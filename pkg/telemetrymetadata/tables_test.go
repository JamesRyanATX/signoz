package telemetrymetadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	dbName = ""
	Init("custom_metadata_db")
	assert.Equal(t, "custom_metadata_db", dbName)
}

func TestInitEmpty(t *testing.T) {
	dbName = ""
	Init("")
	assert.Equal(t, "", dbName)
}

func TestDBNameAfterInit(t *testing.T) {
	dbName = ""
	Init("my_custom_metadata")
	assert.Equal(t, "my_custom_metadata", DBName())
}

func TestDBNameDefault(t *testing.T) {
	dbName = ""
	assert.Equal(t, "signoz_metadata", DBName())
}

func TestInitOverwrite(t *testing.T) {
	dbName = ""
	Init("first_value")
	assert.Equal(t, "first_value", DBName())
	Init("second_value")
	assert.Equal(t, "second_value", DBName())
}

func TestInitEmptyDoesNotOverwrite(t *testing.T) {
	dbName = ""
	Init("initial_value")
	assert.Equal(t, "initial_value", DBName())
	Init("")
	assert.Equal(t, "initial_value", DBName())
}

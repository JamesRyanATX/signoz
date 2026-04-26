package telemetrytraces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	dbName = ""
	t.Cleanup(func() { dbName = "" })

	Init("custom_traces_db")

	assert.Equal(t, "custom_traces_db", dbName)
}

func TestInitEmpty(t *testing.T) {
	dbName = ""
	t.Cleanup(func() { dbName = "" })

	Init("")

	assert.Equal(t, "", dbName)
}

func TestDBNameAfterInit(t *testing.T) {
	dbName = ""
	t.Cleanup(func() { dbName = "" })

	Init("my_custom_traces")

	assert.Equal(t, "my_custom_traces", DBName())
}

func TestDBNameDefault(t *testing.T) {
	dbName = ""
	t.Cleanup(func() { dbName = "" })

	assert.Equal(t, "signoz_traces", DBName())
}

func TestInitOverwrite(t *testing.T) {
	dbName = ""
	t.Cleanup(func() { dbName = "" })

	Init("first_value")
	assert.Equal(t, "first_value", DBName())

	Init("second_value")
	assert.Equal(t, "second_value", DBName())
}

func TestInitEmptyDoesNotOverwrite(t *testing.T) {
	dbName = ""
	t.Cleanup(func() { dbName = "" })

	Init("initial_value")
	assert.Equal(t, "initial_value", DBName())

	Init("")
	// Empty string should not overwrite
	assert.Equal(t, "initial_value", DBName())
}

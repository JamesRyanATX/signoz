package telemetrymeter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	dbName = ""
	Init("custom_meter_db")
	assert.Equal(t, "custom_meter_db", dbName)
}

func TestInitEmpty(t *testing.T) {
	dbName = ""
	Init("")
	assert.Equal(t, "", dbName)
}

func TestDBNameAfterInit(t *testing.T) {
	dbName = ""
	Init("my_custom_meter")
	assert.Equal(t, "my_custom_meter", DBName())
}

func TestDBNameDefault(t *testing.T) {
	dbName = ""
	assert.Equal(t, "signoz_meter", DBName())
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

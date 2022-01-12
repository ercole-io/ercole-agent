package oracle

import (
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

const testListPDB string = `TAGGR01 ||| READ WRITE
TDPRO01     |||          READ WRITE`
const testListPDB1 string = `TAGGR01 ||| READ WRITE  ||| TDPRO01|||          READ WRITE`
const testListPDB2 string = `TAGGR01 ||| READ WRITE

TDPRO01     |||          READ WRITE`

func TestListPDB(t *testing.T) {
	expected := []model.OracleDatabasePluggableDatabase{
		{
			Name:        "TAGGR01",
			Status:      "READ WRITE",
			Tablespaces: nil,
			Schemas:     nil,
			Services:    []model.OracleDatabaseService{},
			OtherInfo:   nil,
		},
		{
			Name:        "TDPRO01",
			Status:      "READ WRITE",
			Tablespaces: nil,
			Schemas:     nil,
			Services:    []model.OracleDatabaseService{},
			OtherInfo:   nil,
		},
	}

	actual, err := ListPDB([]byte(testListPDB))

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestListPDB_Error1(t *testing.T) {
	actual, err := ListPDB([]byte(testListPDB1))

	assert.Nil(t, actual)
	assert.Error(t, err)
}

func TestListPDB_Error2(t *testing.T) {
	actual, err := ListPDB([]byte(testListPDB2))

	assert.Nil(t, actual)
	assert.Error(t, err)
}

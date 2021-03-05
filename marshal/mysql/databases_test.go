package mysql

import (
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

var testDbData = `mysql: [Warning] Using a password on the command line interface can be insecure.
"mysql";"utf8mb4";"utf8mb4_0900_ai_ci";"NO"
"classicmodels";"latin1";"latin1_swedish_ci";"NO"
`

func TestDatabases(t *testing.T) {
	cmdOutput := []byte(testDbData)

	actual := Databases(cmdOutput)

	expected := []model.MySQLDatabase{
		{
			Name:      "mysql",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_0900_ai_ci",
			Encrypted: false,
		},
		{
			Name:      "classicmodels",
			Charset:   "latin1",
			Collation: "latin1_swedish_ci",
			Encrypted: false,
		},
	}

	assert.Equal(t, expected, actual)
}

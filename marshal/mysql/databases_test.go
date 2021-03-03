package mysql

import (
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
)

var testDbData = `mysql: [Warning] Using a password on the command line interface can be insecure.
"erclinmysql:3306";"8.0.23";"COMMUNITY";"mysql";"utf8mb4";"utf8mb4_0900_ai_ci";"NO"
"erclinmysql:3306";"8.0.23";"COMMUNITY";"classicmodels";"latin1";"latin1_swedish_ci";"NO"
`

func TestDatabases(t *testing.T) {
	cmdOutput := []byte(testDbData)

	actual := Databases(cmdOutput)

	expected := []model.MySQLDatabase{
		{
			Instance:  "erclinmysql:3306",
			Version:   "8.0.23",
			Edition:   "COMMUNITY",
			Name:      "mysql",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_0900_ai_ci",
			Encrypted: false,
		},
		{
			Instance:  "erclinmysql:3306",
			Version:   "8.0.23",
			Edition:   "COMMUNITY",
			Name:      "classicmodels",
			Charset:   "latin1",
			Collation: "latin1_swedish_ci",
			Encrypted: false,
		},
	}

	assert.Equal(t, expected, actual)
}

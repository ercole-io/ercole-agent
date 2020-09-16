// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package oracle

import (
	"testing"

	"github.com/ercole-io/ercole/model"
	"github.com/stretchr/testify/assert"
)

const testDatabaseData1 string = `ERC18													|||ERC18			 |||		  1|||ERC18	      |||OPEN	     |||18.0.0.0.0 Enterprise Edition	    |||Linux x86 64-bit 							    |||ARCHIVELOG  |||AL32UTF8																					      |||AL16UTF16																						 |||8192																								    |||2																								       |||1.172 	|||.391 	 |||0.		  |||1.172	   |||	       3|||														   5||| 																								   129|||																	  12660.45|||					7.7|||																		      0|||1					  |||N|||N`
const testDatabaseData2 string = `ERC18													|||ERC18			 |||		  1|||ERC18	      |||OPEN	     |||18.0.0.0.0 Enterprise Edition	    |||Linux x86 64-bit 							    |||NOARCHIVELOG  |||AL32UTF8																					      |||AL16UTF16																						 |||8192																								    |||2																								       |||1.172 	|||.391 	 |||0.		  |||1.172	   |||	       3|||														   5||| 																								   129|||																	  12660.45|||					7.7|||																		      0|||1					  |||N|||N`
const testDatabaseData3 string = `ERC18													|||ERC18			 |||		  1|||ERC18	      |||OPEN	     |||18.0.0.0.0 Enterprise Edition	    |||Linux x86 64-bit 							    |||PIPPO|||AL32UTF8																					      |||AL16UTF16																						 |||8192																								    |||2																								       |||1.172 	|||.391 	 |||0.		  |||1.172	   |||	       3|||														   5||| 																								   129|||																	  12660.45|||					7.7|||																		      0|||1					  |||N|||N`
const testDatabaseData4 string = `ERC18													|||ERC18			 |||		  1|||ERC18	      |||OPEN	     |||18.0.0.0.0 Enterprise Edition	    |||Linux x86 64-bit 							    |||NOARCHIVELOG|||AL32UTF8																					      |||AL16UTF16																						 |||8192																								    |||2																								       |||1.172 	|||.391 	 |||0.		  |||1.172	   |||	       3|||														   5||| 																								   129|||																	  12660.45|||					7.7|||																		      N/A|||1					  |||N|||N`

func TestDatabase_WithArchivelog(t *testing.T) {
	cmdOutput := []byte(testDatabaseData1)

	actual := Database(cmdOutput)

	elapsed := (float64)(12660.45)
	dbTime := (float64)(7.7)
	dailyCPUUsage := (float64)(0)
	work := (float64)(1)

	expected := model.OracleDatabase{InstanceNumber: 1,
		InstanceName: "ERC18",

		Name:          "ERC18",
		UniqueName:    "ERC18",
		Status:        "OPEN",
		IsCDB:         false,
		Version:       "18.0.0.0.0 Enterprise Edition",
		Platform:      "Linux x86 64-bit",
		Archivelog:    true,
		Charset:       "AL32UTF8",
		NCharset:      "AL16UTF16",
		BlockSize:     8192,
		CPUCount:      2,
		SGATarget:     1.172,
		PGATarget:     0.391,
		MemoryTarget:  0,
		SGAMaxSize:    1.172,
		SegmentsSize:  3,
		DatafileSize:  5,
		Allocable:     129,
		Elapsed:       &elapsed,
		DBTime:        &dbTime,
		DailyCPUUsage: &dailyCPUUsage,
		Work:          &work,
		ASM:           false,
		Dataguard:     false,
	}

	assert.Equal(t, expected, actual)
}

func TestDatabase_WithoutArchivelog(t *testing.T) {
	cmdOutput := []byte(testDatabaseData2)

	actual := Database(cmdOutput)

	elapsed := (float64)(12660.45)
	dbTime := (float64)(7.7)
	dailyCPUUsage := (float64)(0)
	work := (float64)(1)

	expected := model.OracleDatabase{InstanceNumber: 1,
		InstanceName: "ERC18",

		Name:          "ERC18",
		UniqueName:    "ERC18",
		Status:        "OPEN",
		IsCDB:         false,
		Version:       "18.0.0.0.0 Enterprise Edition",
		Platform:      "Linux x86 64-bit",
		Archivelog:    false,
		Charset:       "AL32UTF8",
		NCharset:      "AL16UTF16",
		BlockSize:     8192,
		CPUCount:      2,
		SGATarget:     1.172,
		PGATarget:     0.391,
		MemoryTarget:  0,
		SGAMaxSize:    1.172,
		SegmentsSize:  3,
		DatafileSize:  5,
		Allocable:     129,
		Elapsed:       &elapsed,
		DBTime:        &dbTime,
		DailyCPUUsage: &dailyCPUUsage,
		Work:          &work,
		ASM:           false,
		Dataguard:     false,
	}

	assert.Equal(t, expected, actual)
}

func TestDatabase_WrongArchivelog(t *testing.T) {
	cmdOutput := []byte(testDatabaseData3)

	assert.Panics(t, func() {
		Database(cmdOutput)
	})
}

func TestDatabase_WithoutDailyCPUUsage(t *testing.T) {
	cmdOutput := []byte(testDatabaseData4)

	actual := Database(cmdOutput)

	elapsed := (float64)(12660.45)
	dbTime := (float64)(7.7)
	work := (float64)(1)

	expected := model.OracleDatabase{InstanceNumber: 1,
		InstanceName: "ERC18",

		Name:          "ERC18",
		UniqueName:    "ERC18",
		Status:        "OPEN",
		IsCDB:         false,
		Version:       "18.0.0.0.0 Enterprise Edition",
		Platform:      "Linux x86 64-bit",
		Archivelog:    false,
		Charset:       "AL32UTF8",
		NCharset:      "AL16UTF16",
		BlockSize:     8192,
		CPUCount:      2,
		SGATarget:     1.172,
		PGATarget:     0.391,
		MemoryTarget:  0,
		SGAMaxSize:    1.172,
		SegmentsSize:  3,
		DatafileSize:  5,
		Allocable:     129,
		Elapsed:       &elapsed,
		DBTime:        &dbTime,
		DailyCPUUsage: &work,
		Work:          &work,
		ASM:           false,
		Dataguard:     false,
	}

	assert.Equal(t, expected, actual)
}

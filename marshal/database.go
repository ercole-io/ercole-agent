// Copyright (c) 2019 Sorint.lab S.p.A.
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

package marshal

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ercole-io/ercole-agent/model"
)

// Database returns information about database extracted
// from the db fetcher command output.
func Database(cmdOutput []byte) model.Database {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(cmdOutput)))

	if err != nil {
		log.Fatal(err)
	}

	var db model.Database

	doc.Find("tr").Each(func(r int, row *goquery.Selection) { // there should be only one

		sel := row.Find("td")

		for i := range sel.Nodes {
			single := sel.Eq(i)
			value := cleanTr(single.Text())
			if i == 0 {
				db.Name = value
			}
			if i == 1 {
				db.UniqueName = value
			}
			if i == 2 {
				db.InstanceNumber = value
			}
			if i == 3 {
				db.Status = value
			}
			if i == 4 {
				db.Version = value
			}
			if i == 5 {
				db.Platform = value
			}
			if i == 6 {
				db.Archivelog = value
			}
			if i == 7 {
				db.Charset = value
			}
			if i == 8 {
				db.NCharset = value
			}
			if i == 9 {
				db.BlockSize = value
			}
			if i == 10 {
				db.CPUCount = value
			}
			if i == 11 {
				db.SGATarget = value
			}
			if i == 12 {
				db.PGATarget = value
			}
			if i == 13 {
				db.MemoryTarget = value
			}
			if i == 14 {
				db.SGAMaxSize = value
			}
			if i == 15 {
				db.SegmentsSize = value
			}
			if i == 16 {
				db.Used = value
			}
			if i == 17 {
				db.Allocated = value
			}
			if i == 18 {
				db.Elapsed = value
			}
			if i == 19 {
				db.DBTime = value
			}
			if i == 20 {
				db.Work = value
			}
			if i == 21 {
				db.ASM = parseBool(value)
			}
			if i == 22 {
				db.Dataguard = parseBool(value)
			}
		}

	})

	return db
}

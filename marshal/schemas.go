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

// Schemas returns information about database tablespaces extracted
// from the tablespaces fetcher command output.
func Schemas(cmdOutput []byte) []model.Schema {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(cmdOutput)))

	if err != nil {
		log.Fatal(err)
	}

	var scs []model.Schema

	doc.Find("tr").Each(func(r int, row *goquery.Selection) {

		var sc model.Schema

		sel := row.Find("td")

		for i := range sel.Nodes {
			single := sel.Eq(i)
			value := cleanTr(single.Text())
			if i == 2 {
				sc.Database = value
			}
			if i == 3 {
				sc.User = value
			}
			if i == 4 {
				sc.Total = parseInt(value)
			}
			if i == 5 {
				sc.Tables = parseInt(value)
			}
			if i == 6 {
				sc.Indexes = parseInt(value)
			}
			if i == 7 {
				sc.LOB = parseInt(value)
			}
		}

		scs = append(scs, sc)

	})

	return scs
}

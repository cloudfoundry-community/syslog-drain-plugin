package syslog

import (
	"encoding/csv"
	"io"
)

func WriteCSV(w io.Writer, sds []*Drain) error {
	var records [][]string
	records = append(records, []string{"Org", "Space", "Bound App Name", "Drain Name", "Drain URL", "Drain GUID", "Drain Service Last Operation"})
	for _, sd := range sds {
		for _, a := range sd.Apps {
			r := []string{
				sd.Organization.Name,
				sd.Space.Name,
				a.Name,
				sd.Name,
				sd.URL,
				sd.GUID,
				sd.LastOperation.State,
			}
			records = append(records, r)
		}
	}

	cw := csv.NewWriter(w)
	return cw.WriteAll(records)
}

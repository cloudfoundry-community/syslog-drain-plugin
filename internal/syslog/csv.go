package syslog

import (
	"encoding/csv"
	"io"
)

func WriteCSV(w io.Writer, sds []*Drain) error {
	var records [][]string

	for _, sd := range sds {
		for _, a := range sd.Apps {
			r := []string{
				sd.GUID,
				sd.Name,
				sd.URL,
				sd.LastOperation.State,
				a.Name,
				sd.Space.Name,
				sd.Organization.Name,
			}
			records = append(records, r)
		}
	}

	cw := csv.NewWriter(w)
	return cw.WriteAll(records)
}

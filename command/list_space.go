package command

import (
	"code.cloudfoundry.org/cli/plugin"
	"context"
	"github.com/cloudfoundry-community/syslog-drain-plugin/internal/syslog"
	"io"
)

func ListSpaceSyslogDrains(cliConnection plugin.CliConnection, log Logger, w io.Writer) error {
	cf, err := createCFClient(cliConnection)
	if err != nil {
		return err
	}
	l := syslog.NewDrainLister(cf, log)

	s, err := cliConnection.GetCurrentSpace()
	if err != nil {
		return err
	}

	sds, err := l.ListSpaceSyslogDrains(context.Background(), s.Guid)
	if err != nil {
		return err
	}

	return syslog.WriteCSV(w, sds)
}

package command

import (
	"code.cloudfoundry.org/cli/plugin"
	"context"
	"github.com/cloudfoundry-community/syslog-drain-plugin/internal/syslog"
	"io"
)

func ListSyslogDrains(cliConnection plugin.CliConnection, log Logger, w io.Writer) error {
	cf, err := createCFClient(cliConnection)
	if err != nil {
		return err
	}
	cf2, err := createCFv2Client(cliConnection)
	if err != nil {
		return err
	}
	l := syslog.NewDrainLister(cf, cf2, log)

	sds, err := l.ListSyslogDrains(context.Background())
	if err != nil {
		return err
	}

	return syslog.WriteCSV(w, sds)
}

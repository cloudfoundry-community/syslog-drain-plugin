package command

import (
	"code.cloudfoundry.org/cli/plugin"
	"context"
	"github.com/cloudfoundry-community/syslog-drain-plugin/internal/syslog"
	"io"
)

func ListOrgSyslogDrains(cliConnection plugin.CliConnection, log Logger, w io.Writer) error {
	cf, err := createCFClient(cliConnection)
	if err != nil {
		return err
	}
	l := syslog.NewDrainLister(cf, log)

	o, err := cliConnection.GetCurrentOrg()
	if err != nil {
		return err
	}

	sds, err := l.ListOrgSyslogDrains(context.Background(), o.Guid)
	if err != nil {
		return err
	}

	return syslog.WriteCSV(w, sds)
}

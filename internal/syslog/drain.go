package syslog

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"time"
)

type Drain struct {
	GUID          string
	Name          string
	URL           string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastOperation *resource.LastOperation
	Organization  *resource.Organization
	Space         *resource.Space
	Apps          []*resource.App
}

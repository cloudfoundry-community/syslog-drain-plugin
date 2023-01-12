package syslog

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type DrainLister struct {
	cf  *client.Client
	cf2 *cfclient.Client
	log Logger
}

func NewDrainLister(cf *client.Client, cf2 *cfclient.Client, log Logger) *DrainLister {
	return &DrainLister{
		cf:  cf,
		cf2: cf2,
		log: log,
	}
}

func (c *DrainLister) ListSyslogDrains(ctx context.Context) ([]*Drain, error) {
	var sds []*Drain

	opts := client.NewSpaceListOptions()
	for {
		spaces, pager, err := c.cf.Spaces.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		for _, space := range spaces {
			spaceDrains, err := c.ListSpaceSyslogDrains(ctx, space.GUID)
			if err != nil {
				return nil, err
			}
			sds = append(sds, spaceDrains...)
		}
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}

	return sds, nil
}

func (c *DrainLister) ListOrgSyslogDrains(ctx context.Context, orgGUID string) ([]*Drain, error) {
	var sds []*Drain

	opts := client.NewSpaceListOptions()
	opts.OrganizationGUIDs.EqualTo(orgGUID)
	for {
		spaces, pager, err := c.cf.Spaces.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		for _, space := range spaces {
			spaceDrains, err := c.ListSpaceSyslogDrains(ctx, space.GUID)
			if err != nil {
				return nil, err
			}
			sds = append(sds, spaceDrains...)
		}
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}

	return sds, nil
}

func (c *DrainLister) ListSpaceSyslogDrains(ctx context.Context, spaceGUID string) ([]*Drain, error) {
	var sds []*Drain

	space, org, err := c.cf.Spaces.GetIncludeOrganization(ctx, spaceGUID)
	if err != nil {
		return nil, err
	}

	opts := client.NewServiceInstanceListOptions()
	opts.Type = "user-provided"
	opts.SpaceGUIDs.EqualTo(space.GUID)
	for {
		sis, pager, err := c.cf.ServiceInstances.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		for _, si := range sis {
			if !isSyslogDrain(si) {
				continue
			}
			apps, err := c.listApps(ctx, si)
			if err != nil {
				return nil, err
			}
			sd := &Drain{
				GUID:          si.GUID,
				Name:          si.Name,
				URL:           *si.SyslogDrainURL,
				CreatedAt:     si.CreatedAt,
				UpdatedAt:     si.UpdatedAt,
				LastOperation: &si.LastOperation,
				Organization:  org,
				Space:         space,
				Apps:          apps,
			}
			sds = append(sds, sd)
		}
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return sds, nil
}

func (c *DrainLister) listApps(ctx context.Context, si *resource.ServiceInstance) ([]*resource.App, error) {
	var apps []*resource.App

	q := url.Values{}
	q.Set("results-per-page", "100")
	q.Set("q", fmt.Sprintf("service_instance_guid:%s", si.GUID))
	sbs, err := c.cf2.ListServiceBindingsByQuery(q)
	if err != nil {
		return nil, err
	}
	for _, sb := range sbs {
		app, err := c.cf.Applications.Get(ctx, sb.AppGuid)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func isSyslogDrain(si *resource.ServiceInstance) bool {
	return si.SyslogDrainURL != nil && *si.SyslogDrainURL != ""
}

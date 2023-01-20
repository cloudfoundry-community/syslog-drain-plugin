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
	spaceCount := 0
	for {
		spaces, pager, err := c.cf.Spaces.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		for _, space := range spaces {
			spaceCount++
			c.log.Debugf("Processing space %s (%d/%d)", space.Name, spaceCount, pager.TotalResults)
			spaceDrains, err := c.listSpaceSyslogDrains(ctx, space)
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

	c.debugPrintResults(sds)
	return sds, nil
}

func (c *DrainLister) ListOrgSyslogDrains(ctx context.Context, orgGUID string) ([]*Drain, error) {
	var sds []*Drain

	opts := client.NewSpaceListOptions()
	opts.OrganizationGUIDs.EqualTo(orgGUID)
	spaceCount := 0
	for {
		spaces, pager, err := c.cf.Spaces.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		for _, space := range spaces {
			spaceCount++
			c.log.Debugf("Processing space %s (%d/%d)", space.Name, spaceCount, pager.TotalResults)
			spaceDrains, err := c.listSpaceSyslogDrains(ctx, space)
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

	c.debugPrintResults(sds)
	return sds, nil
}

func (c *DrainLister) ListSpaceSyslogDrains(ctx context.Context, spaceGUID string) ([]*Drain, error) {
	space, err := c.cf.Spaces.Get(ctx, spaceGUID)
	if err != nil {
		return nil, err
	}

	c.log.Debugf("Processing space %s (1/1)", space.Name)
	drains, err := c.listSpaceSyslogDrains(ctx, space)
	if err != nil {
		return nil, err
	}

	c.debugPrintResults(drains)
	return drains, nil
}

func (c *DrainLister) listSpaceSyslogDrains(ctx context.Context, space *resource.Space) ([]*Drain, error) {
	var sds []*Drain

	org, err := c.cf.Organizations.Get(ctx, space.Relationships.Organization.Data.GUID)
	if err != nil {
		return nil, err
	}
	c.log.Debugf("Found org %s", org.Name)

	opts := client.NewServiceInstanceListOptions()
	opts.Type = "user-provided"
	opts.SpaceGUIDs.EqualTo(space.GUID)
	siCount := 0
	for {
		sis, pager, err := c.cf.ServiceInstances.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		for _, si := range sis {
			siCount++
			c.log.Debugf("Processing user-provided service instance %s (%d/%d)", si.Name, siCount, pager.TotalResults)
			if !isSyslogDrain(si) {
				c.log.Debugf("Skipping non-syslog service instance %s ", si.Name)
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

	for sbCount, sb := range sbs {
		c.log.Debugf("Processing service instance %s binding (%d/%d)", si.Name, sbCount+1, len(sbs)+1)
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

func (c *DrainLister) debugPrintResults(drains []*Drain) {
	spaces := map[string]*resource.Space{}
	orgs := map[string]*resource.Organization{}
	var appCount int

	for _, drain := range drains {
		appCount += len(drain.Apps)
		spaces[drain.Space.Name] = drain.Space
		orgs[drain.Organization.Name] = drain.Organization
	}

	c.log.Debugf("Found a total of %d syslog drains", len(drains))
	c.log.Debugf("  in %d org(s)", len(orgs))
	c.log.Debugf("  in %d space(s)", len(spaces))
	c.log.Debugf("  attached to %d app(s)", appCount)
}

package vcenter_client_factory

import (
	"context"
	"net/url"
	"time"

	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients/vcenter_manager"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/soap"
)

//go:generate counterfeiter . Vim25ClientCreator
type Vim25ClientCreator interface {
	NewClient(ctx context.Context, rt soap.RoundTripper) (*vim25.Client, error)
}

type ClientCreator struct {
}

func (g *ClientCreator) NewClient(ctx context.Context, rt soap.RoundTripper) (*vim25.Client, error) {

	vimClient, err := vim25.NewClient(ctx, rt)
	if err != nil {
		return nil, err
	}

	return vimClient, nil
}

//go:generate counterfeiter . FinderCreator
type FinderCreator interface {
	NewFinder(client *vim25.Client, all bool) *find.Finder
}

type GovmomiFinderCreator struct {
}

func (g *GovmomiFinderCreator) NewFinder(client *vim25.Client, all bool) *find.Finder {
	return find.NewFinder(client, all)
}

type ManagerFactory struct {
	VCenterServer      string
	Username           string
	Password           string
	InsecureConnection bool
	ClientCreator      Vim25ClientCreator
	FinderCreator      FinderCreator
}

func (f *ManagerFactory) VCenterManager(ctx context.Context) (*vcenter_manager.VCenterManager, error) {

	govmomiClient, err := f.govmomiClient(ctx)
	if err != nil {
		return nil, err
	}

	//TODO: understand the "all" parameter to new finder
	finder := f.FinderCreator.NewFinder(govmomiClient.Client, false)

	return vcenter_manager.NewVCenterManager(govmomiClient, govmomiClient.Client, finder, f.Username, f.Password)

}

func (f *ManagerFactory) govmomiClient(ctx context.Context) (*govmomi.Client, error) {

	sc, err := f.soapClient()
	if err != nil {
		return nil, err
	}

	vc, err := f.vimClient(ctx, sc)
	if err != nil {
		return nil, err
	}

	return &govmomi.Client{
		Client:         vc,
		SessionManager: session.NewManager(vc),
	}, nil

}

func (f *ManagerFactory) soapClient() (*soap.Client, error) {
	vCenterURL, err := soap.ParseURL(f.VCenterServer)
	if err != nil {
		return nil, err
	}
	credentials := url.UserPassword(f.Username, f.Password)
	vCenterURL.User = credentials

	soapClient := soap.NewClient(vCenterURL, f.InsecureConnection)
	//soapClient.SetRootCAs()
	//soapClient.SetRootCAs()

	return soapClient, nil
}

func (f *ManagerFactory) vimClient(ctx context.Context, soapClient *soap.Client) (*vim25.Client, error) {
	vimClient, err := f.ClientCreator.NewClient(ctx, soapClient)
	if err != nil {
		return nil, err
	}

	vimClient.RoundTripper = session.KeepAlive(vimClient.RoundTripper, 10*time.Minute)
	return vimClient, nil
}

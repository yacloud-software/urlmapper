package migrate

import (
	"context"
	"fmt"

	"golang.conradwood.net/go-easyops/authremote"
	"golang.yacloud.eu/apis/protomanager"
	"golang.yacloud.eu/urlmapper/db"
)

func Start() error {
	err := find_fqdn_for_jsonmappings()
	if err != nil {
		return err
	}

	err = find_fqdn_for_anymappings()
	if err != nil {
		return err
	}

	return nil
}

func find_fqdn_for_jsonmappings() error {
	ctx := context.Background()
	jmaps, err := db.DefaultDBJsonMapping().All(ctx)
	if err != nil {
		return err
	}
	for _, jmap := range jmaps {
		if jmap.ServiceName != "" && jmap.FQDNServiceName != "" {
			continue
		}
		ctx = authremote.Context()
		fmt.Printf("Need servicemapping for %v\n", jmap)
		serviceid := jmap.ServiceID
		p, err := byid(ctx, &protomanager.ID{ID: serviceid})
		if err != nil {
			return err
		}
		fqdn := p.PackageFQDN + "/" + p.Name
		fmt.Printf("Found: %s\n", fqdn)
		if jmap.FQDNServiceName == "" {
			jmap.FQDNServiceName = fqdn
		}
		if jmap.ServiceName == "" {
			jmap.ServiceName = p.Package + "." + p.Name
		}
		err = db.DefaultDBJsonMapping().Update(ctx, jmap)
		if err != nil {
			return err
		}

	}
	return nil
}
func find_fqdn_for_anymappings() error {
	ctx := context.Background()
	amaps, err := db.DefaultDBAnyHostMapping().All(ctx)
	if err != nil {
		return err
	}
	for _, amap := range amaps {
		if amap.ServiceName != "" && amap.FQDNServiceName != "" {
			continue
		}
		ctx = authremote.Context()
		fmt.Printf("Need fqdn for: %v\n", amap)
		serviceid := amap.ServiceID
		p, err := byid(ctx, &protomanager.ID{ID: serviceid})
		if err != nil {
			return err
		}
		fqdn := p.PackageFQDN + "/" + p.Name
		fmt.Printf("Found: %s\n", fqdn)
		if amap.FQDNServiceName == "" {
			amap.FQDNServiceName = fqdn
		}
		if amap.ServiceName == "" {
			amap.ServiceName = p.Package + "." + p.Name
		}
		err = db.DefaultDBAnyHostMapping().Update(ctx, amap)
		if err != nil {
			return err
		}

	}
	return nil
}

func byid(ctx context.Context, req *protomanager.ID) (*protomanager.Service, error) {
	panic("not implemented")
	//	p, err := protomanager.GetProtoManagerClient().FindServiceByID(ctx, &protomanager.ID{ID: serviceid})
	//	return p, err
	//	return &protomanager.Service{ID: req.ID}, nil
}

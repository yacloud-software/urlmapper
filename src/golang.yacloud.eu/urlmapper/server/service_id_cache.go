package main

import (
	"context"
	//	"fmt"
	//	"golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/cache"
	//	"golang.conradwood.net/go-easyops/errors"
	"golang.yacloud.eu/apis/protomanager"
	"time"
)

const (
	USE_PROTOMANAGER = true
)

var (
	scache = cache.New("serviceidcache", time.Duration(5)*time.Minute, 1000)
	fcache = cache.New("servicefqdncache", time.Duration(5)*time.Minute, 1000)
)

type fqdncacheentry struct {
	res *protomanager.Service
}

func getServiceByFQDN(ctx context.Context, fqdn string) (*fqdncacheentry, error) {
	obj := fcache.Get(fqdn)
	if obj != nil {
		sobj := obj.(*fqdncacheentry)
		return sobj, nil
	}
	p, err := protomanager.GetProtoManagerClient().FindServiceByFQDN(ctx, &protomanager.FindServiceFQDNRequest{FQDN: fqdn})
	if err != nil {
		return nil, err
	}
	fce := &fqdncacheentry{res: p}
	fcache.Put(fqdn, fce)
	return fce, nil
}

/*
func getServiceName(ctx context.Context, serviceid string) (string, error) {
	obj := scache.Get(serviceid)
	if obj != nil {
		sobj := obj.(string)
		return sobj, nil
	}
	servicename := ""
	if USE_PROTOMANAGER {
		p, err := protomanager.GetProtoManagerClient().FindServiceByID(ctx, &protomanager.ID{ID: serviceid})
		if err != nil {
			return "", err
		}
		servicename = fmt.Sprintf("%s.%s", p.Package, p.Name)
	} else {
		p, err := protorenderer.GetProtoRendererClient().FindServiceByID(ctx, &protorenderer.ID{ID: serviceid})
		if err != nil {
			return "", err
		}
		servicename = fmt.Sprintf("%s.%s", p.PackageName, p.Service.Name)

	}
	scache.Put(serviceid, servicename)
	return servicename, nil
}
*/

/*
func id_to_service_name(ctx context.Context, serviceid string) string {
	res, err := getServiceName(ctx, serviceid)
	if err != nil {
		return fmt.Sprintf("noservice(%s)", err)
	}
	return res
}
*/

// returns, e.g. "registry.Registry"
func (fce *fqdncacheentry) RegistryName() string {
	return fce.res.Package + "." + fce.res.Name
}

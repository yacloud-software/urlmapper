package main

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/cache"
	"golang.yacloud.eu/apis/protomanager"
	"time"
)

const (
	USE_PROTOMANAGER = false
)

var (
	scache = cache.New("serviceidcache", time.Duration(5)*time.Minute, 1000)
)

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
		servicename = p.Name
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

func id_to_service_name(ctx context.Context, serviceid string) string {
	res, err := getServiceName(ctx, serviceid)
	if err != nil {
		return fmt.Sprintf("noservice(%s)", err)
	}
	return res
}

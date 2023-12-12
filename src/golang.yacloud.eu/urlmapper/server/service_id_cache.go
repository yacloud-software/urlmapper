package main

import (
	"context"
	"fmt"
	prender "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/cache"
	"time"
)

var (
	scache = cache.NewResolvingCache("serviceidcache", time.Duration(5)*time.Minute, 1000)
)

func getServiceName(ctx context.Context, serviceid string) (string, error) {
	o, err := scache.Retrieve(serviceid, func(k string) (interface{}, error) {
		r, err := prender.GetProtoRendererClient().FindServiceByID(ctx, &prender.ID{ID: k})
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("%s.%s", r.PackageName, r.Service.Name), nil
	})
	if err != nil {
		return "", err
	}
	return o.(string), nil
}




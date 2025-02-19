package main

import (
	"context"
	"fmt"

	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	"golang.yacloud.eu/apis/protomanager"
	pb "golang.yacloud.eu/apis/urlmapper"
	"golang.yacloud.eu/urlmapper/db"
)

func (e *echoServer) SetRPCMapping(ctx context.Context, req *pb.RPCMappingRequest) (*common.Void, error) {
	if !auth.IsInGroup(ctx, "8") {
		return nil, errors.AccessDenied(ctx, "Access denied")
	}
	fmt.Printf("Expose service=%s, rpc=%s: %v\n", req.FQDNServiceName, req.RPCName, req.Expose)
	q := db.DefaultDBRPCMapping().NewQuery()
	q.AddEqual("fqdnservice", req.FQDNServiceName)
	q.AddEqual("rpcname", req.RPCName)
	rpcs, err := db.DefaultDBRPCMapping().ByDBQuery(ctx, q)
	if err != nil {
		return nil, err
	}
	id := uint64(0)
	if len(rpcs) != 0 {
		id = rpcs[0].ID
	}
	if req.Expose && (id != 0) {
		// nothing to do
		return &common.Void{}, nil
	}
	fsvc := &protomanager.FindServiceFQDNRequest{FQDN: req.FQDNServiceName}
	svc, err := protomanager.GetProtoManagerClient().FindServiceByFQDN(ctx, fsvc)
	if err != nil {
		return nil, err
	}
	if svc.RegistryName == "" {
		return nil, errors.InvalidArgs(ctx, "service has no registryname", "service has no registryname")
	}
	if id == 0 {
		r := &pb.RPCMapping{
			RPCName:     req.RPCName,
			FQDNService: req.FQDNServiceName,
			ServiceName: svc.RegistryName,
		}
		_, err = db.DefaultDBRPCMapping().Save(ctx, r)
		if err != nil {
			return nil, err
		}
	} else {
		err = db.DefaultDBRPCMapping().DeleteByID(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return &common.Void{}, nil
}

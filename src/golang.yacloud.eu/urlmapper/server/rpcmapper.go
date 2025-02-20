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
	prefix := fmt.Sprintf("[service=%s, rpc=%s] ", req.FQDNServiceName, req.RPCName)
	fmt.Printf("%sExpose: %v\n", prefix, req.Expose)
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
		fmt.Printf("%sNothing to do\n", prefix)
		// nothing to do
		return &common.Void{}, nil
	}
	fsvc := &protomanager.FindServiceFQDNRequest{FQDN: req.FQDNServiceName}
	svc, err := protomanager.GetProtoManagerClient().FindServiceByFQDN(ctx, fsvc)
	if err != nil {
		fmt.Printf("%sFailed to find service for %s: %s\n", prefix, fsvc.FQDN, errors.ErrorString(err))
		return nil, err
	}
	if svc.RegistryName == "" {
		return nil, errors.InvalidArgs(ctx, "service has no registryname", "service has no registryname")
	}
	if id == 0 {
		r := &pb.RPCMapping{
			Active:      true,
			RPCName:     req.RPCName,
			FQDNService: req.FQDNServiceName,
			ServiceName: svc.RegistryName,
		}
		_, err = db.DefaultDBRPCMapping().Save(ctx, r)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%sSaved\n", prefix)
	} else {
		err = db.DefaultDBRPCMapping().DeleteByID(ctx, id)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%sDeleted\n", prefix)
	}
	return &common.Void{}, nil
}

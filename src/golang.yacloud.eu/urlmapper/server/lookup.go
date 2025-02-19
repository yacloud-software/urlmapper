package main

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	pb "golang.yacloud.eu/apis/urlmapper"
	"golang.yacloud.eu/urlmapper/db"
)

func (e *echoServer) GetMapping(ctx context.Context, req *pb.MappingRequest) (*pb.Mapping, error) {
	fmt.Printf("New style mapping lookup for \"%s\"\n", req.Path)
	path := strings.TrimPrefix(req.Path, "/")
	url, err := url.Parse(req.Path)
	if err == nil {
		path = url.Path
		fmt.Printf("url detected, using path \"%s\"\n", path)
	}
	path = strings.TrimPrefix(path, "/")
	/*
	 lookup 3 things (in this order)
	 1. rpc mapping
	 2. service mapping
	 3. very old style
	*/
	if strings.HasPrefix(path, "_api/") {
		m, err := get_mapping_rpc_by_path(ctx, path)
		if m != nil {
			return &pb.Mapping{
				MappingFound:    true,
				RegistryName:    m.ServiceName,
				FQDNServiceName: m.FQDNService,
				RPCName:         m.RPCName,
			}, nil
		}
		if err != nil {
			return nil, err
		}
	}
	if strings.HasPrefix(path, "_api/") {
		m, err := get_mapping_service(ctx, path)
		if m != nil {
			return m, nil
		}
		if err != nil {
			return nil, err
		}
	}

	return &pb.Mapping{}, nil
}

// get an rpc mapping. path must be _api/fqdnservice/rpcname, e.g. "_api/golang.yacloud.eu/apis/vuehelper/Log"
func get_mapping_rpc_by_path(ctx context.Context, path string) (*pb.RPCMapping, error) {
	p := strings.TrimPrefix(path, "_api/")
	service := filepath.Dir(p)
	rpc := filepath.Base(p)
	fmt.Printf("Looking up service \"%s\", RPC \"%s\"\n", service, rpc)
	q := db.DefaultDBRPCMapping().NewQuery()
	q.AddEqual("fqdnservice", service)
	q.AddEqual("rpcname", rpc)
	rpcs, err := db.DefaultDBRPCMapping().ByDBQuery(ctx, q)
	if err != nil {
		return nil, err
	}
	if len(rpcs) == 0 {
		return nil, nil
	}
	r := rpcs[0]
	return r, nil
}

// get a service mapping. path must be _api/fqdnservice, e.g. "_api/golang.yacloud.eu/apis/vuehelper"
func get_mapping_service(ctx context.Context, path string) (*pb.Mapping, error) {
	p := strings.TrimPrefix(path, "_api/")
	service := filepath.Dir(p)
	fmt.Printf("Looking up service \"%s\"\n", service)
	res, err := db.DefaultDBAnyHostMapping().ByPath(ctx, service)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	r := res[0]
	rpcname := filepath.Base(p)
	return &pb.Mapping{
		MappingFound:    true,
		RegistryName:    r.ServiceName,
		FQDNServiceName: r.FQDNServiceName,
		RPCName:         rpcname,
	}, nil

}

// get a service mapping. path must be _api/fqdnservice, e.g. "_api/golang.yacloud.eu/apis/vuehelper"
func get_mapping_service_old(ctx context.Context, host string, path string) (*pb.Mapping, error) {
	// we don't do this any more. we do not need it any more (hopefully)
	return nil, nil
}

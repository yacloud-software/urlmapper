package main

import (
	"context"

	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	pb "golang.yacloud.eu/apis/urlmapper"
)

func (e *echoServer) SetRPCMapping(ctx context.Context, req *pb.RPCMappingRequest) (*common.Void, error) {
	if !auth.IsInGroup(ctx, "yacloud-admin") {
		return nil, errors.AccessDenied(ctx, "Access denied")
	}
	return &common.Void{}, nil
}

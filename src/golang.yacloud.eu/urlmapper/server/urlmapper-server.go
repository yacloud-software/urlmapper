package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	ipm "golang.conradwood.net/apis/ipmanager"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/sql"
	"golang.conradwood.net/go-easyops/utils"
	"golang.yacloud.eu/urlmapper/db"
	"google.golang.org/grpc"
	"os"
	"strings"
	pb "golang.yacloud.eu/apis/urlmapper"
)

var (
	jsonMapStore *db.DBJsonMapping
	debug        = flag.Bool("debug", false, "debug mode")
	port         = flag.Int("port", 4100, "The grpc server port")
)

type echoServer struct {
}

func main() {
	flag.Parse()
	fmt.Printf("Starting URLMapperServer...\n")
	psql, err := sql.Open()
	utils.Bail("failed to open sql", err)
	jsonMapStore = db.NewDBJsonMapping(psql)
	sd := server.NewServerDef()
	sd.Port = *port
	sd.Register = server.Register(
		func(server *grpc.Server) error {
			e := new(echoServer)
			pb.RegisterURLMapperServer(server, e)
			return nil
		},
	)
	err = server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
	os.Exit(0)
}

/************************************
* grpc functions
************************************/
func (e *echoServer) AddJsonMapping(ctx context.Context, req *pb.JsonMapping) (*common.Void, error) {
	req.Path = strings.Trim(req.Path, "/")
	req.Path = strings.ToLower(req.Path)
	req.Domain = strings.ToLower(req.Domain)
	jms, err := jsonMapStore.ByPath(ctx, req.Path)
	if err != nil {
		return nil, err
	}
	for _, j := range jms {
		if j.Domain == req.Domain {
			return nil, errors.InvalidArgs(ctx, "duplicate entry", "duplicate entry %s/%s", req.Domain, req.Path)
		}
	}
	ipc := ipm.GetIPManagerClient()
	ra, err := ipc.GetHostAccess(ctx, &ipm.HostAccessRequest{Host: req.Domain})
	if err != nil {
		return nil, err
	}
	if !ra.Response.Granted {
		return nil, errors.AccessDenied(ctx, "no access to domain %s (#%d)", ra.Domain.Name, ra.Domain.ID)
	}
	id, err := jsonMapStore.Save(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("stored jsonmapping: %d\n", id)
	return &common.Void{}, nil
}
func (e *echoServer) GetJsonMappings(ctx context.Context, req *common.Void) (*pb.JsonMappingResponseList, error) {

	jms, err := jsonMapStore.All(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.JsonMappingResponseList{}
	for _, r := range jms {
		jmr := &pb.JsonMappingResponse{Mapping: r}
		sname, err := getServiceName(ctx, jmr.Mapping.ServiceID)
		if err != nil {
			fmt.Printf("Service with ID \"%s\" not found: %s\n", jmr.Mapping.ServiceID, utils.ErrorString(err))
			return nil, err
		}
		jmr.GRPCService = sname
		res.Responses = append(res.Responses, jmr)
	}
	return res, nil
}
func (e *echoServer) GetJsonMappingWithUser(ctx context.Context, req *pb.GetJsonMappingRequest) (*pb.JsonMappingResponse, error) {
	rs, err := e.GetJsonMapping(ctx, req)
	if err != nil {
		return nil, err
	}
	if !auth.IsInGroup(ctx, rs.Mapping.GroupID) {
		return nil, errors.NotFound(ctx, "not found (%s/%s)", req.Domain, req.Path)
	}

	return rs, nil

}
func (e *echoServer) GetJsonMapping(ctx context.Context, req *pb.GetJsonMappingRequest) (*pb.JsonMappingResponse, error) {
	req.Path = strings.Trim(req.Path, "/")
	req.Path = strings.ToLower(req.Path)
	req.Domain = strings.ToLower(req.Domain)
	jms, err := jsonMapStore.ByPath(ctx, req.Path)
	if *debug {
		fmt.Printf("GetJsonMapping - path=\"%s\", domain=\"%s\"\n", req.Path, req.Domain)
	}
	if err != nil {
		return nil, err
	}

	var jm *pb.JsonMapping
	for _, j := range jms {
		if j.Domain == req.Domain {
			jm = j
			break
		}
	}
	if jm == nil {
		if *debug {
			fmt.Printf("Found no service for \"domain=%s, path=%s\"\n", req.Domain, req.Path)
		}
		return nil, errors.NotFound(ctx, "not found (%s/%s)", req.Domain, req.Path)
	}
	sname, err := getServiceName(ctx, jm.ServiceID)
	if err != nil {
		fmt.Printf("Service with ID \"%s\" not found: %s\n", jm.ServiceID, utils.ErrorString(err))
		return nil, err
	}
	res := &pb.JsonMappingResponse{
		Mapping:     jm,
		GRPCService: sname,
	}

	return res, nil

}
func (e *echoServer) GetJsonDomains(ctx context.Context, req *common.Void) (*pb.DomainList, error) {
	m := make(map[string]int)
	jms, err := jsonMapStore.All(ctx)
	if err != nil {
		return nil, err
	}

	for _, j := range jms {
		m[j.Domain]++
	}

	res := &pb.DomainList{}
	ipc := ipm.GetIPManagerClient()
	for k, _ := range m {
		ra, err := ipc.GetHostAccess(ctx, &ipm.HostAccessRequest{Host: k})
		if err != nil {
			return nil, err
		}
		if !ra.Response.Granted {
			continue
		}

		res.Domains = append(res.Domains, k)
	}
	return res, nil
}

func (e *echoServer) GetServiceMappings(ctx context.Context, req *pb.ServiceID) (*pb.JsonMappingResponseList, error) {
	jms, err := jsonMapStore.ByServiceID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	res := &pb.JsonMappingResponseList{}
	for _, jm := range jms {
		sname, err := getServiceName(ctx, jm.ServiceID)
		if err != nil {
			fmt.Printf("Service with ID \"%s\" not found: %s\n", jm.ServiceID, utils.ErrorString(err))
			return nil, err
		}
		r := &pb.JsonMappingResponse{Mapping: jm, GRPCService: sname}
		res.Responses = append(res.Responses, r)
	}
	return res, nil

}

package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	ipm "golang.conradwood.net/apis/ipmanager"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/sql"
	"golang.conradwood.net/go-easyops/utils"
	pb "golang.yacloud.eu/apis/urlmapper"
	"golang.yacloud.eu/urlmapper/db"
	"google.golang.org/grpc"
	"os"
	"strings"
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
	db.DefaultDBAnyHostMapping()
	sd := server.NewServerDef()
	sd.SetPort(*port)
	sd.SetRegister(server.Register(
		func(server *grpc.Server) error {
			e := new(echoServer)
			pb.RegisterURLMapperServer(server, e)
			return nil
		},
	))
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

	path := strings.Trim(req.Path, "/")

	if strings.HasPrefix(path, "_api/") {
		p := strings.TrimPrefix(path, "_api/")
		ahms, err := db.DefaultDBAnyHostMapping().ByPath(ctx, p)
		if err != nil {
			return nil, err
		}
		if len(ahms) > 0 {
			ah := ahms[0]
			sname, err := getServiceName(ctx, ah.ServiceID)
			if err != nil {
				fmt.Printf("Service with ID \"%s\" not found: %s\n", ah.ServiceID, utils.ErrorString(err))
				return nil, err
			}
			jm := &pb.JsonMapping{ID: 0, Domain: "*", Path: "/_api/" + ah.Path, ServiceID: ah.ServiceID, GroupID: ""}
			res := &pb.JsonMappingResponse{
				Mapping:     jm,
				GRPCService: sname,
			}
			return res, nil
		}
	}

	path = strings.ToLower(path)
	req.Domain = strings.ToLower(req.Domain)
	jms, err := jsonMapStore.ByPath(ctx, path)
	if *debug {
		fmt.Printf("GetJsonMapping - path=\"%s\", domain=\"%s\"\n", path, req.Domain)
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
	if jm != nil {
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

	if *debug {
		fmt.Printf("Found no service for \"domain=%s, path=%s\"\n", req.Domain, path)
	}
	return nil, errors.NotFound(ctx, "not found (%s/%s)", req.Domain, path)

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
	if *debug {
		fmt.Printf("Getting jsonmapping for service #%s\n", req.ID)
	}
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

	// add the "any" host mappings
	anys, err := db.DefaultDBAnyHostMapping().ByServiceID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	sname, err := getServiceName(ctx, req.ID)
	if err != nil {
		fmt.Printf("Service with ID \"%s\" not found: %s\n", req.ID, utils.ErrorString(err))
		return nil, err
	}
	for _, a := range anys {
		jm := &pb.JsonMapping{ID: 0, Domain: "*", Path: "/_api/" + a.Path, ServiceID: req.ID, GroupID: ""}
		r := &pb.JsonMappingResponse{Mapping: jm, GRPCService: sname}
		res.Responses = append(res.Responses, r)
	}
	return res, nil

}
func (e *echoServer) AddAnyHostMapping(ctx context.Context, req *pb.AnyMappingRequest) (*pb.AnyMappingResponse, error) {
	if req.ServiceName == "" {
		return nil, errors.InvalidArgs(ctx, "missing service name", "missing service name")
	}
	fsr := &pr.FindServiceByNameRequest{Name: req.ServiceName}
	sv, err := pr.GetProtoRendererClient().FindServiceByName(ctx, fsr)
	if err != nil {
		return nil, err
	}
	if len(sv.Services) == 0 {
		fmt.Printf("protorenderer knows no no such service: \"%s\"\n", fsr.Name)
		return nil, errors.NotFound(ctx, "no such service (%s)", fsr.Name)
	}
	if len(sv.Services) > 1 {
		return nil, errors.InvalidArgs(ctx, "ambigous service name", "protorenderer has multiple definitions of (%s)", fsr.Name)
	}
	srv := sv.Services[0]
	path := srv.PackageFQDN + "/" + srv.Service.Name
	//	path = strings.ToLower(path)
	fmt.Printf("Adding Service: %#v -> path=%s\n", srv, path)
	fmt.Printf("ProtoRenderService: %#v -> path=%s\n", srv.Service, path)
	fmt.Printf("ProtoRenderPackage: %#v -> path=%s\n", srv.Package, path)
	ahm := &pb.AnyHostMapping{
		Path:        path,
		ServiceID:   sv.Services[0].Service.ID,
		ServiceName: req.ServiceName,
	}
	_, err = db.DefaultDBAnyHostMapping().Save(ctx, ahm)
	if err != nil {
		return nil, err
	}
	res := &pb.AnyMappingResponse{Path: "_api/" + path}
	return res, nil
}

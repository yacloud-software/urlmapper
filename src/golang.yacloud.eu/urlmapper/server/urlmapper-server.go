package main

import (
	"context"
	"flag"
	"fmt"

	"golang.conradwood.net/apis/common"
	ipm "golang.conradwood.net/apis/ipmanager"

	//	pr "golang.conradwood.net/apis/protorenderer"
	"os"
	"sort"
	"strings"

	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/sql"
	"golang.conradwood.net/go-easyops/utils"
	"golang.yacloud.eu/apis/protomanager"
	pb "golang.yacloud.eu/apis/urlmapper"
	"golang.yacloud.eu/urlmapper/db"
	"golang.yacloud.eu/urlmapper/migrate"
	"google.golang.org/grpc"
)

var (
	panic_on_deprecated = flag.Bool("panic_on_deprecated", true, "panic on deprecated rpc calls")
	jsonMapStore        *db.DBJsonMapping
	debug               = flag.Bool("debug", false, "debug mode")
	port                = flag.Int("port", 4100, "The grpc server port")
)

type echoServer struct {
}

func main() {
	flag.Parse()
	fmt.Printf("Starting URLMapperServer...\n")
	utils.Bail("failed to migrate", migrate.Start())
	psql, err := sql.Open()
	utils.Bail("failed to open sql", err)
	jsonMapStore = db.NewDBJsonMapping(psql)
	db.DefaultDBAnyHostMapping()
	db.DefaultDBRPCMapping()
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
	// *************** DEPRECATED ******************
	if *panic_on_deprecated {
		panic("deprecated code path")
	}
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
	// *************** DEPRECATED ******************
	if *panic_on_deprecated {
		panic("deprecated code path")
	}
	jms, err := jsonMapStore.All(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.JsonMappingResponseList{}
	for _, r := range jms {
		jmr := &pb.JsonMappingResponse{Mapping: r}
		sname, err := getServiceByFQDN(ctx, jmr.Mapping.FQDNServiceName)
		if err != nil {
			fmt.Printf("Service with ID \"%s\" not found: %s\n", jmr.Mapping.FQDNServiceName, utils.ErrorString(err))
			return nil, err
		}
		jmr.GRPCService = sname.RegistryName()
		res.Responses = append(res.Responses, jmr)
	}
	return res, nil
}
func (e *echoServer) GetJsonMappingWithUser(ctx context.Context, req *pb.GetJsonMappingRequest) (*pb.JsonMappingResponse, error) {
	// *************** DEPRECATED ******************
	if *panic_on_deprecated {
		panic("deprecated code path")
	}
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
	fmt.Printf("Old codepath, looking up mapping for path \"%s\"\n", req.Path)
	path := strings.Trim(req.Path, "/")

	if strings.HasPrefix(path, "_api/") {
		p := strings.TrimPrefix(path, "_api/")
		// try rpc mapping
		rhms, err := db.DefaultDBRPCMapping().ByFQDNService(ctx, p)
		if err != nil {
			return nil, err
		}
		if len(rhms) > 0 {
			ah := rhms[0]
			jm := &pb.JsonMapping{ID: 0, Domain: "*", Path: "/_api/" + ah.FQDNService, GroupID: "", RPC: ah.RPCName}
			res := &pb.JsonMappingResponse{
				Mapping:     jm,
				GRPCService: ah.ServiceName,
			}
			return res, nil
		}
		// try service mapping
		ahms, err := db.DefaultDBAnyHostMapping().ByPath(ctx, p)
		if err != nil {
			return nil, err
		}
		if len(ahms) > 0 {
			ah := ahms[0]
			sname, err := getServiceByFQDN(ctx, ah.FQDNServiceName)
			if err != nil {
				fmt.Printf("Service with ID \"%s\" not found: %s\n", ah.FQDNServiceName, utils.ErrorString(err))
				return nil, err
			}
			jm := &pb.JsonMapping{ID: 0, Domain: "*", Path: "/_api/" + ah.Path, GroupID: ""}
			res := &pb.JsonMappingResponse{
				Mapping:     jm,
				GRPCService: sname.RegistryName(),
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
		sname, err := getServiceByFQDN(ctx, jm.FQDNServiceName)
		if err != nil {
			fmt.Printf("Service with ID \"%s\" not found: %s\n", jm.FQDNServiceName, utils.ErrorString(err))
			return nil, err
		}
		res := &pb.JsonMappingResponse{
			Mapping:     jm,
			GRPCService: sname.RegistryName(),
		}
		return res, nil
	}

	if *debug {
		fmt.Printf("Found no service for \"domain=%s, path=%s\"\n", req.Domain, path)
	}
	return nil, errors.NotFound(ctx, "not found (%s/%s)", req.Domain, path)

}
func (e *echoServer) GetJsonDomains(ctx context.Context, req *common.Void) (*pb.DomainList, error) {
	// *************** DEPRECATED ******************
	if *panic_on_deprecated {
		panic("deprecated code path")
	}
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

/*
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
		sname, err := getServiceByFQDN(ctx, jm.FQDNServiceName)
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
		jm := &pb.JsonMapping{ID: 0, Domain: "*", Path: "/_api/" + a.Path, GroupID: ""}
		r := &pb.JsonMappingResponse{Mapping: jm, GRPCService: sname}
		res.Responses = append(res.Responses, r)
	}
	return res, nil
}
*/

func (e *echoServer) AddAnyHostMapping(ctx context.Context, req *pb.AnyMappingRequest) (*pb.AnyMappingResponse, error) {
	if req.ServiceName == "" {
		return nil, errors.InvalidArgs(ctx, "missing service name", "missing service name")
	}
	fsr := &protomanager.ServicesByNameRequest{Name: req.ServiceName, ExactMatchOnly: true}
	sv, err := protomanager.GetProtoManagerClient().FindServicesByName(ctx, fsr)
	if err != nil {
		return nil, err
	}
	if len(sv.Services) == 0 {
		fmt.Printf("protomanager knows no no such service: \"%s\"\n", fsr.Name)
		return nil, errors.NotFound(ctx, "no such service (%s)", fsr.Name)
	}
	if len(sv.Services) > 1 {
		return nil, errors.InvalidArgs(ctx, "ambigous service name", "protorenderer has multiple definitions of (%s)", fsr.Name)
	}
	srv := sv.Services[0]
	path := srv.PackageFQDN + "/" + srv.Name
	//serviceid := srv.ID

	fmt.Printf("Adding Service: %#v -> path=%s\n", srv, path)
	ahm := &pb.AnyHostMapping{
		Path:            path,
		ServiceName:     req.ServiceName,
		FQDNServiceName: path,
	}
	_, err = db.DefaultDBAnyHostMapping().Save(ctx, ahm)
	if err != nil {
		return nil, err
	}
	res := &pb.AnyMappingResponse{Path: "_api/" + path}
	return res, nil
}

func (e *echoServer) GetAllMappings(ctx context.Context, req *common.Void) (*pb.AllMappingList, error) {
	res := &pb.AllMappingList{}

	// *** do json mappings *** /
	jms, err := jsonMapStore.All(ctx)
	if err != nil {
		return nil, err
	}
	for _, jm := range jms {

		a := &pb.AllMapping{
			ServiceID:   "",
			ServiceName: jm.ServiceName,
			Domain:      jm.Domain,
			Path:        jm.Path,
			RPC:         "*",
		}
		res.AllMappings = append(res.AllMappings, a)
	}
	// *** do any mappings *** /
	ams, err := db.DefaultDBAnyHostMapping().All(ctx)
	if err != nil {
		return nil, err
	}
	for _, am := range ams {
		a := &pb.AllMapping{
			//	ServiceID:   am.ServiceID,
			ServiceName: am.ServiceName,
			Domain:      "*",
			Path:        am.Path,
			RPC:         "*",
		}
		res.AllMappings = append(res.AllMappings, a)
	}

	// *** do rpc mappings ***//
	rms, err := db.DefaultDBRPCMapping().All(ctx)
	if err != nil {
		return nil, err
	}
	for _, rm := range rms {
		a := &pb.AllMapping{
			//	ServiceID:   am.ServiceID,
			ServiceName: rm.ServiceName,
			Domain:      "*",
			Path:        rm.FQDNService,
			RPC:         rm.RPCName,
		}
		res.AllMappings = append(res.AllMappings, a)
	}
	sort.Slice(res.AllMappings, func(i, j int) bool {
		return res.AllMappings[i].ServiceName < res.AllMappings[j].ServiceName
	})
	return res, nil

}

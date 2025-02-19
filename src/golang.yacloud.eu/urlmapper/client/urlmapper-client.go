package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"golang.conradwood.net/apis/common"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"golang.yacloud.eu/apis/urlmapper"
	pb "golang.yacloud.eu/apis/urlmapper"
)

var (
	addflag           = flag.String("add", "", "`servicename` - add this service as [anyhost]/_api/[domain]/[path]")
	add_specific_flag = flag.Bool("add_specific", false, "add a specific mapping (with host)")
	findflag          = flag.String("find", "", "find a mapping by path, e.g /_api/foo/bar")
	mapurl            = flag.String("url", "", "url to serve a mapping on")
	grpcservice       = flag.String("service", "", "`grpc service` name, e.g. \"urlmapper.URLMapper\"")
	all               = flag.Bool("all", false, "list all")

	echoClient pb.URLMapperClient
)

func main() {
	flag.Parse()

	echoClient = pb.GetURLMapperClient()
	if *findflag != "" {
		utils.Bail("failed to find", Find())
		os.Exit(0)
	}
	if *all {
		utils.Bail("failed to list all", All())
		os.Exit(0)
	}
	if *addflag != "" {
		Add()
		os.Exit(0)
	}
	if *add_specific_flag {
		AddSpecific()
		os.Exit(0)
	}
	// a context with authentication
	ctx := authremote.Context()

	empty := &common.Void{}
	response, err := echoClient.GetAllMappings(ctx, empty)
	utils.Bail("Failed to ping server", err)
	t := utils.Table{}
	t.AddHeaders("serviceid", "servicename", "domain", "path")
	for _, am := range response.AllMappings {
		t.AddString(am.ServiceID)
		t.AddString(am.ServiceName)
		t.AddString(am.Domain)
		t.AddString(am.Path)
		t.NewRow()
	}
	fmt.Println(t.ToPrettyString())
	fmt.Printf("Done.\n")
	os.Exit(0)
}
func Add() {
	req := &pb.AnyMappingRequest{Path: "", ServiceName: *addflag}
	ctx := authremote.Context()
	r, err := echoClient.AddAnyHostMapping(ctx, req)
	utils.Bail("failed to add", err)
	fmt.Printf("Added - serving on [anyhost]/%s\n", r.Path)
	fmt.Printf("Done\n")
	return
}
func AddSpecific() {
	req := getJsonMap()
	fmt.Printf("Adding %v\n", req)
	ctx := authremote.Context()
	_, err := echoClient.AddJsonMapping(ctx, req)
	utils.Bail("failed to add mapping", err)

}

// get json mapping from parameters
func getJsonMap() *pb.JsonMapping {
	if *mapurl == "" {
		fmt.Printf("Missing -url\n")
		os.Exit(10)
	}
	if *grpcservice == "" {
		fmt.Printf("Missing -service\n")
		os.Exit(10)
	}
	u, err := url.Parse(*mapurl)
	utils.Bail("not a valid url", err)
	fmt.Printf("Scheme: \"%s\"\n", u.Scheme)
	if u.Scheme == "" {
		fmt.Printf("Please specify protocol (http/https)\n")
		os.Exit(10)
	}
	fsr := &pr.FindServiceByNameRequest{Name: *grpcservice}
	ctx := authremote.Context()
	sv, err := pr.GetProtoRendererClient().FindServiceByName(ctx, fsr)
	utils.Bail("failed to get services", err)
	if len(sv.Services) == 0 {
		fmt.Printf("protorenderer knows no such service: \"%s\"\n", fsr.Name)
		os.Exit(10)
	}
	if len(sv.Services) > 1 {
		fmt.Printf("protorenderer has multiple (%d) services with name: \"%s\"\n", len(sv.Services), fsr.Name)
		os.Exit(10)
	}
	res := &pb.JsonMapping{
		Domain:    u.Host,
		Path:      u.Path,
		ServiceID: sv.Services[0].Service.ID,
	}
	return res
}

func All() error {
	ctx := authremote.Context()
	um, err := urlmapper.GetURLMapperClient().GetAllMappings(ctx, &common.Void{})
	if err != nil {
		return err
	}
	t := utils.Table{}
	t.AddHeaders("ServiceName", "Path", "Domain", "RPC")
	for _, m := range um.AllMappings {
		t.AddString(m.ServiceName)
		t.AddString(m.Path)
		t.AddString(m.Domain)
		t.AddString(m.RPC)
		t.NewRow()
	}
	fmt.Println(t.ToPrettyString())
	return nil
}

func Find() error {
	ctx := authremote.Context()
	mr := &urlmapper.MappingRequest{Path: *findflag}
	mapping, err := urlmapper.GetURLMapperClient().GetMapping(ctx, mr)
	if err != nil {
		return err
	}
	if !mapping.MappingFound {
		fmt.Printf("No mapping found for \"%s\"\n", mr.Path)
		return nil
	}
	fmt.Printf("Mapping for \"%s\":\n", mr.Path)
	fmt.Printf("  RegistryName    : %s\n", mapping.RegistryName)
	fmt.Printf("  FQDNService     : %s\n", mapping.FQDNServiceName)
	fmt.Printf("  RPC             : %s\n", mapping.RPCName)
	return nil
}

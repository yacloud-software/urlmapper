package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	pb "golang.yacloud.eu/apis/urlmapper"
	"net/url"
	"os"
)

var (
	addflag           = flag.String("add", "", "`servicename` - add this service as [anyhost]/_api/[domain]/[path]")
	add_specific_flag = flag.Bool("add_specific", false, "add a specific mapping (with host)")
	findflag          = flag.Bool("find", false, "find a mapping")
	mapurl            = flag.String("url", "", "url to serve a mapping on")
	grpcservice       = flag.String("service", "", "`grpc service` name, e.g. \"urlmapper.URLMapper\"")
	echoClient        pb.URLMapperClient
)

func main() {
	flag.Parse()

	echoClient = pb.GetURLMapperClient()
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
	response, err := echoClient.GetJsonMappings(ctx, empty)
	utils.Bail("Failed to ping server", err)
	fmt.Printf("%d mappings:\n", len(response.Responses))
	for _, r := range response.Responses {
		m := r.Mapping
		fmt.Printf("%d %s/%s -> %s\n", m.ID, m.Domain, m.Path, r.GRPCService)
	}

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

//get json mapping from parameters
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
		fmt.Printf("No such service: \"%s\"\n", fsr.Name)
		os.Exit(10)
	}
	if len(sv.Services) > 1 {
		fmt.Printf("Multiple (%d) services with name: \"%s\"\n", len(sv.Services), fsr.Name)
		os.Exit(10)
	}
	res := &pb.JsonMapping{
		Domain:    u.Host,
		Path:      u.Path,
		ServiceID: sv.Services[0].Service.ID,
	}
	return res
}






// client create: URLMapperClient
/*
  Created by /home/cnw/devel/go/yatools/src/golang.yacloud.eu/yatools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   filename  : protos/golang.yacloud.eu/apis/urlmapper/urlmapper.proto
   gopackage : golang.yacloud.eu/apis/urlmapper
   importname: ai_0
   clientfunc: GetURLMapper
   serverfunc: NewURLMapper
   lookupfunc: URLMapperLookupID
   varname   : client_URLMapperClient_0
   clientname: URLMapperClient
   servername: URLMapperServer
   gsvcname  : urlmapper.URLMapper
   lockname  : lock_URLMapperClient_0
   activename: active_URLMapperClient_0
*/

package urlmapper

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_URLMapperClient_0 sync.Mutex
  client_URLMapperClient_0 URLMapperClient
)

func GetURLMapperClient() URLMapperClient { 
    if client_URLMapperClient_0 != nil {
        return client_URLMapperClient_0
    }

    lock_URLMapperClient_0.Lock() 
    if client_URLMapperClient_0 != nil {
       lock_URLMapperClient_0.Unlock()
       return client_URLMapperClient_0
    }

    client_URLMapperClient_0 = NewURLMapperClient(client.Connect(URLMapperLookupID()))
    lock_URLMapperClient_0.Unlock()
    return client_URLMapperClient_0
}

func URLMapperLookupID() string { return "urlmapper.URLMapper" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.

func init() {
   client.RegisterDependency("urlmapper.URLMapper")
   AddService("urlmapper.URLMapper")
}


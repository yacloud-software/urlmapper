// client create: GDriveServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/gdrive/gdrive.proto
   gopackage : golang.conradwood.net/apis/gdrive
   importname: ai_0
   varname   : client_GDriveServiceClient_0
   clientname: GDriveServiceClient
   servername: GDriveServiceServer
   gscvname  : gdrive.GDriveService
   lockname  : lock_GDriveServiceClient_0
   activename: active_GDriveServiceClient_0
*/

package gdrive

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GDriveServiceClient_0 sync.Mutex
  client_GDriveServiceClient_0 GDriveServiceClient
)

func GetGDriveClient() GDriveServiceClient { 
    if client_GDriveServiceClient_0 != nil {
        return client_GDriveServiceClient_0
    }

    lock_GDriveServiceClient_0.Lock() 
    if client_GDriveServiceClient_0 != nil {
       lock_GDriveServiceClient_0.Unlock()
       return client_GDriveServiceClient_0
    }

    client_GDriveServiceClient_0 = NewGDriveServiceClient(client.Connect("gdrive.GDriveService"))
    lock_GDriveServiceClient_0.Unlock()
    return client_GDriveServiceClient_0
}

func GetGDriveServiceClient() GDriveServiceClient { 
    if client_GDriveServiceClient_0 != nil {
        return client_GDriveServiceClient_0
    }

    lock_GDriveServiceClient_0.Lock() 
    if client_GDriveServiceClient_0 != nil {
       lock_GDriveServiceClient_0.Unlock()
       return client_GDriveServiceClient_0
    }

    client_GDriveServiceClient_0 = NewGDriveServiceClient(client.Connect("gdrive.GDriveService"))
    lock_GDriveServiceClient_0.Unlock()
    return client_GDriveServiceClient_0
}


.PHONY:	client server
all:	client server
client:
	cd client && go install ${LDFLAGS} urlmapper-client.go
server:
	cd server && go install ${LDFLAGS} urlmapper-server.go `ls -1 *.go|grep -v urlmapper-server.go`

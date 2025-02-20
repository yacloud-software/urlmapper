package db

/*
 This file was created by mkdb-client.
 The intention is not to modify this file, but you may extend the struct DBJsonMapping
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence jsonmapping_seq;

Main Table:

 CREATE TABLE jsonmapping (id integer primary key default nextval('jsonmapping_seq'),domain text not null  ,path text not null  ,serviceid text not null  ,groupid text not null  ,fqdnservicename text not null  ,servicename text not null  ,rpc text not null  ,active boolean not null  );

Alter statements:
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS domain text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS path text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS serviceid text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS groupid text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS fqdnservicename text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS servicename text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS rpc text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS active boolean not null default false;


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE jsonmapping_archive (id integer unique not null,domain text not null,path text not null,serviceid text not null,groupid text not null,fqdnservicename text not null,servicename text not null,rpc text not null,active boolean not null);
*/

import (
	"context"
	gosql "database/sql"
	"fmt"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/sql"
	savepb "golang.yacloud.eu/apis/urlmapper"
	"os"
	"sync"
)

var (
	default_def_DBJsonMapping *DBJsonMapping
)

type DBJsonMapping struct {
	DB                   *sql.DB
	SQLTablename         string
	SQLArchivetablename  string
	customColumnHandlers []CustomColumnHandler
	lock                 sync.Mutex
}

func DefaultDBJsonMapping() *DBJsonMapping {
	if default_def_DBJsonMapping != nil {
		return default_def_DBJsonMapping
	}
	psql, err := sql.Open()
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err)
		os.Exit(10)
	}
	res := NewDBJsonMapping(psql)
	ctx := context.Background()
	err = res.CreateTable(ctx)
	if err != nil {
		fmt.Printf("Failed to create table: %s\n", err)
		os.Exit(10)
	}
	default_def_DBJsonMapping = res
	return res
}
func NewDBJsonMapping(db *sql.DB) *DBJsonMapping {
	foo := DBJsonMapping{DB: db}
	foo.SQLTablename = "jsonmapping"
	foo.SQLArchivetablename = "jsonmapping_archive"
	return &foo
}

func (a *DBJsonMapping) GetCustomColumnHandlers() []CustomColumnHandler {
	return a.customColumnHandlers
}
func (a *DBJsonMapping) AddCustomColumnHandler(w CustomColumnHandler) {
	a.lock.Lock()
	a.customColumnHandlers = append(a.customColumnHandlers, w)
	a.lock.Unlock()
}

func (a *DBJsonMapping) NewQuery() *Query {
	return newQuery(a)
}

// archive. It is NOT transactionally save.
func (a *DBJsonMapping) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBJsonMapping", "insert into "+a.SQLArchivetablename+" (id,domain, path, serviceid, groupid, fqdnservicename, servicename, rpc, active) values ($1,$2, $3, $4, $5, $6, $7, $8, $9) ", p.ID, p.Domain, p.Path, p.ServiceID, p.GroupID, p.FQDNServiceName, p.ServiceName, p.RPC, p.Active)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// return a map with columnname -> value_from_proto
func (a *DBJsonMapping) buildSaveMap(ctx context.Context, p *savepb.JsonMapping) (map[string]interface{}, error) {
	extra, err := extraFieldsToStore(ctx, a, p)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	res["id"] = a.get_col_from_proto(p, "id")
	res["domain"] = a.get_col_from_proto(p, "domain")
	res["path"] = a.get_col_from_proto(p, "path")
	res["serviceid"] = a.get_col_from_proto(p, "serviceid")
	res["groupid"] = a.get_col_from_proto(p, "groupid")
	res["fqdnservicename"] = a.get_col_from_proto(p, "fqdnservicename")
	res["servicename"] = a.get_col_from_proto(p, "servicename")
	res["rpc"] = a.get_col_from_proto(p, "rpc")
	res["active"] = a.get_col_from_proto(p, "active")
	if extra != nil {
		for k, v := range extra {
			res[k] = v
		}
	}
	return res, nil
}

func (a *DBJsonMapping) Save(ctx context.Context, p *savepb.JsonMapping) (uint64, error) {
	qn := "save_DBJsonMapping"
	smap, err := a.buildSaveMap(ctx, p)
	if err != nil {
		return 0, err
	}
	delete(smap, "id") // save without id
	return a.saveMap(ctx, qn, smap, p)
}

// Save using the ID specified
func (a *DBJsonMapping) SaveWithID(ctx context.Context, p *savepb.JsonMapping) error {
	qn := "insert_DBJsonMapping"
	smap, err := a.buildSaveMap(ctx, p)
	if err != nil {
		return err
	}
	_, err = a.saveMap(ctx, qn, smap, p)
	return err
}

// use a hashmap of columnname->values to store to database (see buildSaveMap())
func (a *DBJsonMapping) saveMap(ctx context.Context, queryname string, smap map[string]interface{}, p *savepb.JsonMapping) (uint64, error) {
	// Save (and use database default ID generation)

	var rows *gosql.Rows
	var e error

	q_cols := ""
	q_valnames := ""
	q_vals := make([]interface{}, 0)
	deli := ""
	i := 0
	// build the 2 parts of the query (column names and value names) as well as the values themselves
	for colname, val := range smap {
		q_cols = q_cols + deli + colname
		i++
		q_valnames = q_valnames + deli + fmt.Sprintf("$%d", i)
		q_vals = append(q_vals, val)
		deli = ","
	}
	rows, e = a.DB.QueryContext(ctx, queryname, "insert into "+a.SQLTablename+" ("+q_cols+") values ("+q_valnames+") returning id", q_vals...)
	if e != nil {
		return 0, a.Error(ctx, queryname, e)
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, a.Error(ctx, queryname, errors.Errorf("No rows after insert"))
	}
	var id uint64
	e = rows.Scan(&id)
	if e != nil {
		return 0, a.Error(ctx, queryname, errors.Errorf("failed to scan id after insert: %s", e))
	}
	p.ID = id
	return id, nil
}

func (a *DBJsonMapping) Update(ctx context.Context, p *savepb.JsonMapping) error {
	qn := "DBJsonMapping_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set domain=$1, path=$2, serviceid=$3, groupid=$4, fqdnservicename=$5, servicename=$6, rpc=$7, active=$8 where id = $9", a.get_Domain(p), a.get_Path(p), a.get_ServiceID(p), a.get_GroupID(p), a.get_FQDNServiceName(p), a.get_ServiceName(p), a.get_RPC(p), a.get_Active(p), p.ID)

	return a.Error(ctx, qn, e)
}

// delete by id field
func (a *DBJsonMapping) DeleteByID(ctx context.Context, p uint64) error {
	qn := "deleteDBJsonMapping_ByID"
	_, e := a.DB.ExecContext(ctx, qn, "delete from "+a.SQLTablename+" where id = $1", p)
	return a.Error(ctx, qn, e)
}

// get it by primary id
func (a *DBJsonMapping) ByID(ctx context.Context, p uint64) (*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByID"
	l, e := a.fromQuery(ctx, qn, "id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, errors.Errorf("No JsonMapping with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, errors.Errorf("Multiple (%d) JsonMapping with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBJsonMapping) TryByID(ctx context.Context, p uint64) (*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_TryByID"
	l, e := a.fromQuery(ctx, qn, "id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("TryByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, nil
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, errors.Errorf("Multiple (%d) JsonMapping with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by multiple primary ids
func (a *DBJsonMapping) ByIDs(ctx context.Context, p []uint64) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByIDs"
	l, e := a.fromQuery(ctx, qn, "id in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("TryByID: error scanning (%s)", e))
	}
	return l, nil
}

// get all rows
func (a *DBJsonMapping) All(ctx context.Context) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_all"
	l, e := a.fromQuery(ctx, qn, "true")
	if e != nil {
		return nil, errors.Errorf("All: error scanning (%s)", e)
	}
	return l, nil
}

/**********************************************************************
* GetBy[FIELD] functions
**********************************************************************/

// get all "DBJsonMapping" rows with matching Domain
func (a *DBJsonMapping) ByDomain(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByDomain"
	l, e := a.fromQuery(ctx, qn, "domain = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByDomain: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with multiple matching Domain
func (a *DBJsonMapping) ByMultiDomain(ctx context.Context, p []string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByDomain"
	l, e := a.fromQuery(ctx, qn, "domain in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByDomain: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeDomain(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeDomain"
	l, e := a.fromQuery(ctx, qn, "domain ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByDomain: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching Path
func (a *DBJsonMapping) ByPath(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByPath"
	l, e := a.fromQuery(ctx, qn, "path = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPath: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with multiple matching Path
func (a *DBJsonMapping) ByMultiPath(ctx context.Context, p []string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByPath"
	l, e := a.fromQuery(ctx, qn, "path in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPath: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikePath(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikePath"
	l, e := a.fromQuery(ctx, qn, "path ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPath: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching ServiceID
func (a *DBJsonMapping) ByServiceID(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByServiceID"
	l, e := a.fromQuery(ctx, qn, "serviceid = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceID: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with multiple matching ServiceID
func (a *DBJsonMapping) ByMultiServiceID(ctx context.Context, p []string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByServiceID"
	l, e := a.fromQuery(ctx, qn, "serviceid in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceID: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeServiceID(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeServiceID"
	l, e := a.fromQuery(ctx, qn, "serviceid ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceID: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching GroupID
func (a *DBJsonMapping) ByGroupID(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByGroupID"
	l, e := a.fromQuery(ctx, qn, "groupid = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByGroupID: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with multiple matching GroupID
func (a *DBJsonMapping) ByMultiGroupID(ctx context.Context, p []string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByGroupID"
	l, e := a.fromQuery(ctx, qn, "groupid in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByGroupID: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeGroupID(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeGroupID"
	l, e := a.fromQuery(ctx, qn, "groupid ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByGroupID: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching FQDNServiceName
func (a *DBJsonMapping) ByFQDNServiceName(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByFQDNServiceName"
	l, e := a.fromQuery(ctx, qn, "fqdnservicename = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByFQDNServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with multiple matching FQDNServiceName
func (a *DBJsonMapping) ByMultiFQDNServiceName(ctx context.Context, p []string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByFQDNServiceName"
	l, e := a.fromQuery(ctx, qn, "fqdnservicename in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByFQDNServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeFQDNServiceName(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeFQDNServiceName"
	l, e := a.fromQuery(ctx, qn, "fqdnservicename ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByFQDNServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching ServiceName
func (a *DBJsonMapping) ByServiceName(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByServiceName"
	l, e := a.fromQuery(ctx, qn, "servicename = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with multiple matching ServiceName
func (a *DBJsonMapping) ByMultiServiceName(ctx context.Context, p []string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByServiceName"
	l, e := a.fromQuery(ctx, qn, "servicename in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeServiceName(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeServiceName"
	l, e := a.fromQuery(ctx, qn, "servicename ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching RPC
func (a *DBJsonMapping) ByRPC(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByRPC"
	l, e := a.fromQuery(ctx, qn, "rpc = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByRPC: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with multiple matching RPC
func (a *DBJsonMapping) ByMultiRPC(ctx context.Context, p []string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByRPC"
	l, e := a.fromQuery(ctx, qn, "rpc in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByRPC: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeRPC(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeRPC"
	l, e := a.fromQuery(ctx, qn, "rpc ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByRPC: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching Active
func (a *DBJsonMapping) ByActive(ctx context.Context, p bool) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByActive"
	l, e := a.fromQuery(ctx, qn, "active = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByActive: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with multiple matching Active
func (a *DBJsonMapping) ByMultiActive(ctx context.Context, p []bool) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByActive"
	l, e := a.fromQuery(ctx, qn, "active in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByActive: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeActive(ctx context.Context, p bool) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeActive"
	l, e := a.fromQuery(ctx, qn, "active ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByActive: error scanning (%s)", e))
	}
	return l, nil
}

/**********************************************************************
* The field getters
**********************************************************************/

// getter for field "ID" (ID) [uint64]
func (a *DBJsonMapping) get_ID(p *savepb.JsonMapping) uint64 {
	return uint64(p.ID)
}

// getter for field "Domain" (Domain) [string]
func (a *DBJsonMapping) get_Domain(p *savepb.JsonMapping) string {
	return string(p.Domain)
}

// getter for field "Path" (Path) [string]
func (a *DBJsonMapping) get_Path(p *savepb.JsonMapping) string {
	return string(p.Path)
}

// getter for field "ServiceID" (ServiceID) [string]
func (a *DBJsonMapping) get_ServiceID(p *savepb.JsonMapping) string {
	return string(p.ServiceID)
}

// getter for field "GroupID" (GroupID) [string]
func (a *DBJsonMapping) get_GroupID(p *savepb.JsonMapping) string {
	return string(p.GroupID)
}

// getter for field "FQDNServiceName" (FQDNServiceName) [string]
func (a *DBJsonMapping) get_FQDNServiceName(p *savepb.JsonMapping) string {
	return string(p.FQDNServiceName)
}

// getter for field "ServiceName" (ServiceName) [string]
func (a *DBJsonMapping) get_ServiceName(p *savepb.JsonMapping) string {
	return string(p.ServiceName)
}

// getter for field "RPC" (RPC) [string]
func (a *DBJsonMapping) get_RPC(p *savepb.JsonMapping) string {
	return string(p.RPC)
}

// getter for field "Active" (Active) [bool]
func (a *DBJsonMapping) get_Active(p *savepb.JsonMapping) bool {
	return bool(p.Active)
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBJsonMapping) ByDBQuery(ctx context.Context, query *Query) ([]*savepb.JsonMapping, error) {
	extra_fields, err := extraFieldsToQuery(ctx, a)
	if err != nil {
		return nil, err
	}
	i := 0
	for col_name, value := range extra_fields {
		i++
		efname := fmt.Sprintf("EXTRA_FIELD_%d", i)
		query.Add(col_name+" = "+efname, QP{efname: value})
	}

	gw, paras := query.ToPostgres()
	queryname := "custom_dbquery"
	rows, err := a.DB.QueryContext(ctx, queryname, "select "+a.SelectCols()+" from "+a.Tablename()+" where "+gw, paras...)
	if err != nil {
		return nil, err
	}
	res, err := a.FromRows(ctx, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (a *DBJsonMapping) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.JsonMapping, error) {
	return a.fromQuery(ctx, "custom_query_"+a.Tablename(), query_where, args...)
}

// from a query snippet (the part after WHERE)
func (a *DBJsonMapping) fromQuery(ctx context.Context, queryname string, query_where string, args ...interface{}) ([]*savepb.JsonMapping, error) {
	extra_fields, err := extraFieldsToQuery(ctx, a)
	if err != nil {
		return nil, err
	}
	eq := ""
	if extra_fields != nil && len(extra_fields) > 0 {
		eq = " AND ("
		// build the extraquery "eq"
		i := len(args)
		deli := ""
		for col_name, value := range extra_fields {
			i++
			eq = eq + deli + col_name + fmt.Sprintf(" = $%d", i)
			deli = " AND "
			args = append(args, value)
		}
		eq = eq + ")"
	}
	rows, err := a.DB.QueryContext(ctx, queryname, "select "+a.SelectCols()+" from "+a.Tablename()+" where ( "+query_where+") "+eq, args...)
	if err != nil {
		return nil, err
	}
	res, err := a.FromRows(ctx, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}
	return res, nil
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBJsonMapping) get_col_from_proto(p *savepb.JsonMapping, colname string) interface{} {
	if colname == "id" {
		return a.get_ID(p)
	} else if colname == "domain" {
		return a.get_Domain(p)
	} else if colname == "path" {
		return a.get_Path(p)
	} else if colname == "serviceid" {
		return a.get_ServiceID(p)
	} else if colname == "groupid" {
		return a.get_GroupID(p)
	} else if colname == "fqdnservicename" {
		return a.get_FQDNServiceName(p)
	} else if colname == "servicename" {
		return a.get_ServiceName(p)
	} else if colname == "rpc" {
		return a.get_RPC(p)
	} else if colname == "active" {
		return a.get_Active(p)
	}
	panic(fmt.Sprintf("in table \"%s\", column \"%s\" cannot be resolved to proto field name", a.Tablename(), colname))
}

func (a *DBJsonMapping) Tablename() string {
	return a.SQLTablename
}

func (a *DBJsonMapping) SelectCols() string {
	return "id,domain, path, serviceid, groupid, fqdnservicename, servicename, rpc, active"
}
func (a *DBJsonMapping) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".domain, " + a.SQLTablename + ".path, " + a.SQLTablename + ".serviceid, " + a.SQLTablename + ".groupid, " + a.SQLTablename + ".fqdnservicename, " + a.SQLTablename + ".servicename, " + a.SQLTablename + ".rpc, " + a.SQLTablename + ".active"
}

func (a *DBJsonMapping) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.JsonMapping, error) {
	var res []*savepb.JsonMapping
	for rows.Next() {
		// SCANNER:
		foo := &savepb.JsonMapping{}
		// create the non-nullable pointers
		// create variables for scan results
		scanTarget_0 := &foo.ID
		scanTarget_1 := &foo.Domain
		scanTarget_2 := &foo.Path
		scanTarget_3 := &foo.ServiceID
		scanTarget_4 := &foo.GroupID
		scanTarget_5 := &foo.FQDNServiceName
		scanTarget_6 := &foo.ServiceName
		scanTarget_7 := &foo.RPC
		scanTarget_8 := &foo.Active
		err := rows.Scan(scanTarget_0, scanTarget_1, scanTarget_2, scanTarget_3, scanTarget_4, scanTarget_5, scanTarget_6, scanTarget_7, scanTarget_8)
		// END SCANNER

		if err != nil {
			return nil, a.Error(ctx, "fromrow-scan", err)
		}
		res = append(res, foo)
	}
	return res, nil
}

/**********************************************************************
* Helper to create table and columns
**********************************************************************/
func (a *DBJsonMapping) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),domain text not null ,path text not null ,serviceid text not null ,groupid text not null ,fqdnservicename text not null ,servicename text not null ,rpc text not null ,active boolean not null );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),domain text not null ,path text not null ,serviceid text not null ,groupid text not null ,fqdnservicename text not null ,servicename text not null ,rpc text not null ,active boolean not null );`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS domain text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS path text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS serviceid text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS groupid text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS fqdnservicename text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS servicename text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS rpc text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS active boolean not null default false;`,

		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS domain text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS path text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS serviceid text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS groupid text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS fqdnservicename text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS servicename text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS rpc text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS active boolean not null  default false;`,
	}

	for i, c := range csql {
		_, e := a.DB.ExecContext(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
		if e != nil {
			return e
		}
	}

	// these are optional, expected to fail
	csql = []string{
		// Indices:

		// Foreign keys:

	}
	for i, c := range csql {
		a.DB.ExecContextQuiet(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
	}
	return nil
}

/**********************************************************************
* Helper to meaningful errors
**********************************************************************/
func (a *DBJsonMapping) Error(ctx context.Context, q string, e error) error {
	if e == nil {
		return nil
	}
	return errors.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}


package db

/*
 This file was created by mkdb-client.
 The intention is not to modify this file, but you may extend the struct DBAnyHostMapping
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence anyhostmapping_seq;

Main Table:

 CREATE TABLE anyhostmapping (id integer primary key default nextval('anyhostmapping_seq'),path text not null  unique  ,serviceid text not null  ,servicename text not null  ,fqdnservicename text not null  ,active boolean not null  );

Alter statements:
ALTER TABLE anyhostmapping ADD COLUMN IF NOT EXISTS path text not null unique  default '';
ALTER TABLE anyhostmapping ADD COLUMN IF NOT EXISTS serviceid text not null default '';
ALTER TABLE anyhostmapping ADD COLUMN IF NOT EXISTS servicename text not null default '';
ALTER TABLE anyhostmapping ADD COLUMN IF NOT EXISTS fqdnservicename text not null default '';
ALTER TABLE anyhostmapping ADD COLUMN IF NOT EXISTS active boolean not null default false;


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE anyhostmapping_archive (id integer unique not null,path text not null,serviceid text not null,servicename text not null,fqdnservicename text not null,active boolean not null);
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
	default_def_DBAnyHostMapping *DBAnyHostMapping
)

type DBAnyHostMapping struct {
	DB                   *sql.DB
	SQLTablename         string
	SQLArchivetablename  string
	customColumnHandlers []CustomColumnHandler
	lock                 sync.Mutex
}

func DefaultDBAnyHostMapping() *DBAnyHostMapping {
	if default_def_DBAnyHostMapping != nil {
		return default_def_DBAnyHostMapping
	}
	psql, err := sql.Open()
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err)
		os.Exit(10)
	}
	res := NewDBAnyHostMapping(psql)
	ctx := context.Background()
	err = res.CreateTable(ctx)
	if err != nil {
		fmt.Printf("Failed to create table: %s\n", err)
		os.Exit(10)
	}
	default_def_DBAnyHostMapping = res
	return res
}
func NewDBAnyHostMapping(db *sql.DB) *DBAnyHostMapping {
	foo := DBAnyHostMapping{DB: db}
	foo.SQLTablename = "anyhostmapping"
	foo.SQLArchivetablename = "anyhostmapping_archive"
	return &foo
}

func (a *DBAnyHostMapping) GetCustomColumnHandlers() []CustomColumnHandler {
	return a.customColumnHandlers
}
func (a *DBAnyHostMapping) AddCustomColumnHandler(w CustomColumnHandler) {
	a.lock.Lock()
	a.customColumnHandlers = append(a.customColumnHandlers, w)
	a.lock.Unlock()
}

func (a *DBAnyHostMapping) NewQuery() *Query {
	return newQuery(a)
}

// archive. It is NOT transactionally save.
func (a *DBAnyHostMapping) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBAnyHostMapping", "insert into "+a.SQLArchivetablename+" (id,path, serviceid, servicename, fqdnservicename, active) values ($1,$2, $3, $4, $5, $6) ", p.ID, p.Path, p.ServiceID, p.ServiceName, p.FQDNServiceName, p.Active)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// return a map with columnname -> value_from_proto
func (a *DBAnyHostMapping) buildSaveMap(ctx context.Context, p *savepb.AnyHostMapping) (map[string]interface{}, error) {
	extra, err := extraFieldsToStore(ctx, a, p)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	res["id"] = a.get_col_from_proto(p, "id")
	res["path"] = a.get_col_from_proto(p, "path")
	res["serviceid"] = a.get_col_from_proto(p, "serviceid")
	res["servicename"] = a.get_col_from_proto(p, "servicename")
	res["fqdnservicename"] = a.get_col_from_proto(p, "fqdnservicename")
	res["active"] = a.get_col_from_proto(p, "active")
	if extra != nil {
		for k, v := range extra {
			res[k] = v
		}
	}
	return res, nil
}

func (a *DBAnyHostMapping) Save(ctx context.Context, p *savepb.AnyHostMapping) (uint64, error) {
	qn := "save_DBAnyHostMapping"
	smap, err := a.buildSaveMap(ctx, p)
	if err != nil {
		return 0, err
	}
	delete(smap, "id") // save without id
	return a.saveMap(ctx, qn, smap, p)
}

// Save using the ID specified
func (a *DBAnyHostMapping) SaveWithID(ctx context.Context, p *savepb.AnyHostMapping) error {
	qn := "insert_DBAnyHostMapping"
	smap, err := a.buildSaveMap(ctx, p)
	if err != nil {
		return err
	}
	_, err = a.saveMap(ctx, qn, smap, p)
	return err
}

// use a hashmap of columnname->values to store to database (see buildSaveMap())
func (a *DBAnyHostMapping) saveMap(ctx context.Context, queryname string, smap map[string]interface{}, p *savepb.AnyHostMapping) (uint64, error) {
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

func (a *DBAnyHostMapping) Update(ctx context.Context, p *savepb.AnyHostMapping) error {
	qn := "DBAnyHostMapping_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set path=$1, serviceid=$2, servicename=$3, fqdnservicename=$4, active=$5 where id = $6", a.get_Path(p), a.get_ServiceID(p), a.get_ServiceName(p), a.get_FQDNServiceName(p), a.get_Active(p), p.ID)

	return a.Error(ctx, qn, e)
}

// delete by id field
func (a *DBAnyHostMapping) DeleteByID(ctx context.Context, p uint64) error {
	qn := "deleteDBAnyHostMapping_ByID"
	_, e := a.DB.ExecContext(ctx, qn, "delete from "+a.SQLTablename+" where id = $1", p)
	return a.Error(ctx, qn, e)
}

// get it by primary id
func (a *DBAnyHostMapping) ByID(ctx context.Context, p uint64) (*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByID"
	l, e := a.fromQuery(ctx, qn, "id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, errors.Errorf("No AnyHostMapping with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, errors.Errorf("Multiple (%d) AnyHostMapping with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBAnyHostMapping) TryByID(ctx context.Context, p uint64) (*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_TryByID"
	l, e := a.fromQuery(ctx, qn, "id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("TryByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, nil
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, errors.Errorf("Multiple (%d) AnyHostMapping with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by multiple primary ids
func (a *DBAnyHostMapping) ByIDs(ctx context.Context, p []uint64) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByIDs"
	l, e := a.fromQuery(ctx, qn, "id in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("TryByID: error scanning (%s)", e))
	}
	return l, nil
}

// get all rows
func (a *DBAnyHostMapping) All(ctx context.Context) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_all"
	l, e := a.fromQuery(ctx, qn, "true")
	if e != nil {
		return nil, errors.Errorf("All: error scanning (%s)", e)
	}
	return l, nil
}

/**********************************************************************
* GetBy[FIELD] functions
**********************************************************************/

// get all "DBAnyHostMapping" rows with matching Path
func (a *DBAnyHostMapping) ByPath(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByPath"
	l, e := a.fromQuery(ctx, qn, "path = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPath: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBAnyHostMapping" rows with multiple matching Path
func (a *DBAnyHostMapping) ByMultiPath(ctx context.Context, p []string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByPath"
	l, e := a.fromQuery(ctx, qn, "path in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPath: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBAnyHostMapping) ByLikePath(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByLikePath"
	l, e := a.fromQuery(ctx, qn, "path ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPath: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBAnyHostMapping" rows with matching ServiceID
func (a *DBAnyHostMapping) ByServiceID(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByServiceID"
	l, e := a.fromQuery(ctx, qn, "serviceid = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceID: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBAnyHostMapping" rows with multiple matching ServiceID
func (a *DBAnyHostMapping) ByMultiServiceID(ctx context.Context, p []string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByServiceID"
	l, e := a.fromQuery(ctx, qn, "serviceid in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceID: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBAnyHostMapping) ByLikeServiceID(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByLikeServiceID"
	l, e := a.fromQuery(ctx, qn, "serviceid ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceID: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBAnyHostMapping" rows with matching ServiceName
func (a *DBAnyHostMapping) ByServiceName(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByServiceName"
	l, e := a.fromQuery(ctx, qn, "servicename = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBAnyHostMapping" rows with multiple matching ServiceName
func (a *DBAnyHostMapping) ByMultiServiceName(ctx context.Context, p []string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByServiceName"
	l, e := a.fromQuery(ctx, qn, "servicename in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBAnyHostMapping) ByLikeServiceName(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByLikeServiceName"
	l, e := a.fromQuery(ctx, qn, "servicename ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBAnyHostMapping" rows with matching FQDNServiceName
func (a *DBAnyHostMapping) ByFQDNServiceName(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByFQDNServiceName"
	l, e := a.fromQuery(ctx, qn, "fqdnservicename = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByFQDNServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBAnyHostMapping" rows with multiple matching FQDNServiceName
func (a *DBAnyHostMapping) ByMultiFQDNServiceName(ctx context.Context, p []string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByFQDNServiceName"
	l, e := a.fromQuery(ctx, qn, "fqdnservicename in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByFQDNServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBAnyHostMapping) ByLikeFQDNServiceName(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByLikeFQDNServiceName"
	l, e := a.fromQuery(ctx, qn, "fqdnservicename ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByFQDNServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBAnyHostMapping" rows with matching Active
func (a *DBAnyHostMapping) ByActive(ctx context.Context, p bool) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByActive"
	l, e := a.fromQuery(ctx, qn, "active = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByActive: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBAnyHostMapping" rows with multiple matching Active
func (a *DBAnyHostMapping) ByMultiActive(ctx context.Context, p []bool) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByActive"
	l, e := a.fromQuery(ctx, qn, "active in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByActive: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBAnyHostMapping) ByLikeActive(ctx context.Context, p bool) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByLikeActive"
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
func (a *DBAnyHostMapping) get_ID(p *savepb.AnyHostMapping) uint64 {
	return uint64(p.ID)
}

// getter for field "Path" (Path) [string]
func (a *DBAnyHostMapping) get_Path(p *savepb.AnyHostMapping) string {
	return string(p.Path)
}

// getter for field "ServiceID" (ServiceID) [string]
func (a *DBAnyHostMapping) get_ServiceID(p *savepb.AnyHostMapping) string {
	return string(p.ServiceID)
}

// getter for field "ServiceName" (ServiceName) [string]
func (a *DBAnyHostMapping) get_ServiceName(p *savepb.AnyHostMapping) string {
	return string(p.ServiceName)
}

// getter for field "FQDNServiceName" (FQDNServiceName) [string]
func (a *DBAnyHostMapping) get_FQDNServiceName(p *savepb.AnyHostMapping) string {
	return string(p.FQDNServiceName)
}

// getter for field "Active" (Active) [bool]
func (a *DBAnyHostMapping) get_Active(p *savepb.AnyHostMapping) bool {
	return bool(p.Active)
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBAnyHostMapping) ByDBQuery(ctx context.Context, query *Query) ([]*savepb.AnyHostMapping, error) {
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

func (a *DBAnyHostMapping) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.AnyHostMapping, error) {
	return a.fromQuery(ctx, "custom_query_"+a.Tablename(), query_where, args...)
}

// from a query snippet (the part after WHERE)
func (a *DBAnyHostMapping) fromQuery(ctx context.Context, queryname string, query_where string, args ...interface{}) ([]*savepb.AnyHostMapping, error) {
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
func (a *DBAnyHostMapping) get_col_from_proto(p *savepb.AnyHostMapping, colname string) interface{} {
	if colname == "id" {
		return a.get_ID(p)
	} else if colname == "path" {
		return a.get_Path(p)
	} else if colname == "serviceid" {
		return a.get_ServiceID(p)
	} else if colname == "servicename" {
		return a.get_ServiceName(p)
	} else if colname == "fqdnservicename" {
		return a.get_FQDNServiceName(p)
	} else if colname == "active" {
		return a.get_Active(p)
	}
	panic(fmt.Sprintf("in table \"%s\", column \"%s\" cannot be resolved to proto field name", a.Tablename(), colname))
}

func (a *DBAnyHostMapping) Tablename() string {
	return a.SQLTablename
}

func (a *DBAnyHostMapping) SelectCols() string {
	return "id,path, serviceid, servicename, fqdnservicename, active"
}
func (a *DBAnyHostMapping) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".path, " + a.SQLTablename + ".serviceid, " + a.SQLTablename + ".servicename, " + a.SQLTablename + ".fqdnservicename, " + a.SQLTablename + ".active"
}

func (a *DBAnyHostMapping) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.AnyHostMapping, error) {
	var res []*savepb.AnyHostMapping
	for rows.Next() {
		// SCANNER:
		foo := &savepb.AnyHostMapping{}
		// create the non-nullable pointers
		// create variables for scan results
		scanTarget_0 := &foo.ID
		scanTarget_1 := &foo.Path
		scanTarget_2 := &foo.ServiceID
		scanTarget_3 := &foo.ServiceName
		scanTarget_4 := &foo.FQDNServiceName
		scanTarget_5 := &foo.Active
		err := rows.Scan(scanTarget_0, scanTarget_1, scanTarget_2, scanTarget_3, scanTarget_4, scanTarget_5)
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
func (a *DBAnyHostMapping) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),path text not null ,serviceid text not null ,servicename text not null ,fqdnservicename text not null ,active boolean not null );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),path text not null ,serviceid text not null ,servicename text not null ,fqdnservicename text not null ,active boolean not null );`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS path text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS serviceid text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS servicename text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS fqdnservicename text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS active boolean not null default false;`,

		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS path text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS serviceid text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS servicename text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS fqdnservicename text not null  default '';`,
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
		`create unique index if not exists uniq_anyhostmapping_path on anyhostmapping (path);`,
		`alter table anyhostmapping add constraint uniq_anyhostmapping_path unique using index uniq_anyhostmapping_path;`,

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
func (a *DBAnyHostMapping) Error(ctx context.Context, q string, e error) error {
	if e == nil {
		return nil
	}
	return errors.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}


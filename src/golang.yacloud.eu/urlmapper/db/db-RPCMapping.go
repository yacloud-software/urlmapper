package db

/*
 This file was created by mkdb-client.
 The intention is not to modify this file, but you may extend the struct DBRPCMapping
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence rpcmapping_seq;

Main Table:

 CREATE TABLE rpcmapping (id integer primary key default nextval('rpcmapping_seq'),servicename text not null  ,fqdnservice text not null  ,rpcname text not null  ,active boolean not null  );

Alter statements:
ALTER TABLE rpcmapping ADD COLUMN IF NOT EXISTS servicename text not null default '';
ALTER TABLE rpcmapping ADD COLUMN IF NOT EXISTS fqdnservice text not null default '';
ALTER TABLE rpcmapping ADD COLUMN IF NOT EXISTS rpcname text not null default '';
ALTER TABLE rpcmapping ADD COLUMN IF NOT EXISTS active boolean not null default false;


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE rpcmapping_archive (id integer unique not null,servicename text not null,fqdnservice text not null,rpcname text not null,active boolean not null);
*/

import (
	"context"
	gosql "database/sql"
	"fmt"
	"os"
	"sync"

	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/sql"
	savepb "golang.yacloud.eu/apis/urlmapper"
)

var (
	default_def_DBRPCMapping *DBRPCMapping
)

type DBRPCMapping struct {
	DB                   *sql.DB
	SQLTablename         string
	SQLArchivetablename  string
	customColumnHandlers []CustomColumnHandler
	lock                 sync.Mutex
}

func DefaultDBRPCMapping() *DBRPCMapping {
	if default_def_DBRPCMapping != nil {
		return default_def_DBRPCMapping
	}
	psql, err := sql.Open()
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err)
		os.Exit(10)
	}
	res := NewDBRPCMapping(psql)
	ctx := context.Background()
	err = res.CreateTable(ctx)
	if err != nil {
		fmt.Printf("Failed to create table: %s\n", err)
		os.Exit(10)
	}
	default_def_DBRPCMapping = res
	return res
}
func NewDBRPCMapping(db *sql.DB) *DBRPCMapping {
	foo := DBRPCMapping{DB: db}
	foo.SQLTablename = "rpcmapping"
	foo.SQLArchivetablename = "rpcmapping_archive"
	return &foo
}

func (a *DBRPCMapping) GetCustomColumnHandlers() []CustomColumnHandler {
	return a.customColumnHandlers
}
func (a *DBRPCMapping) AddCustomColumnHandler(w CustomColumnHandler) {
	a.lock.Lock()
	a.customColumnHandlers = append(a.customColumnHandlers, w)
	a.lock.Unlock()
}

func (a *DBRPCMapping) NewQuery() *Query {
	return newQuery(a)
}

// archive. It is NOT transactionally save.
func (a *DBRPCMapping) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBRPCMapping", "insert into "+a.SQLArchivetablename+" (id,servicename, fqdnservice, rpcname, active) values ($1,$2, $3, $4, $5) ", p.ID, p.ServiceName, p.FQDNService, p.RPCName, p.Active)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// return a map with columnname -> value_from_proto
func (a *DBRPCMapping) buildSaveMap(ctx context.Context, p *savepb.RPCMapping) (map[string]interface{}, error) {
	extra, err := extraFieldsToStore(ctx, a, p)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	res["id"] = a.get_col_from_proto(p, "id")
	res["servicename"] = a.get_col_from_proto(p, "servicename")
	res["fqdnservice"] = a.get_col_from_proto(p, "fqdnservice")
	res["rpcname"] = a.get_col_from_proto(p, "rpcname")
	res["active"] = a.get_col_from_proto(p, "active")
	if extra != nil {
		for k, v := range extra {
			res[k] = v
		}
	}
	return res, nil
}

func (a *DBRPCMapping) Save(ctx context.Context, p *savepb.RPCMapping) (uint64, error) {
	qn := "save_DBRPCMapping"
	smap, err := a.buildSaveMap(ctx, p)
	if err != nil {
		return 0, err
	}
	delete(smap, "id") // save without id
	return a.saveMap(ctx, qn, smap, p)
}

// Save using the ID specified
func (a *DBRPCMapping) SaveWithID(ctx context.Context, p *savepb.RPCMapping) error {
	qn := "insert_DBRPCMapping"
	smap, err := a.buildSaveMap(ctx, p)
	if err != nil {
		return err
	}
	_, err = a.saveMap(ctx, qn, smap, p)
	return err
}

// use a hashmap of columnname->values to store to database (see buildSaveMap())
func (a *DBRPCMapping) saveMap(ctx context.Context, queryname string, smap map[string]interface{}, p *savepb.RPCMapping) (uint64, error) {
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

func (a *DBRPCMapping) Update(ctx context.Context, p *savepb.RPCMapping) error {
	qn := "DBRPCMapping_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set servicename=$1, fqdnservice=$2, rpcname=$3, active=$4 where id = $5", a.get_ServiceName(p), a.get_FQDNService(p), a.get_RPCName(p), a.get_Active(p), p.ID)

	return a.Error(ctx, qn, e)
}

// delete by id field
func (a *DBRPCMapping) DeleteByID(ctx context.Context, p uint64) error {
	qn := "deleteDBRPCMapping_ByID"
	_, e := a.DB.ExecContext(ctx, qn, "delete from "+a.SQLTablename+" where id = $1", p)
	return a.Error(ctx, qn, e)
}

// get it by primary id
func (a *DBRPCMapping) ByID(ctx context.Context, p uint64) (*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByID"
	l, e := a.fromQuery(ctx, qn, "id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, errors.Errorf("No RPCMapping with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, errors.Errorf("Multiple (%d) RPCMapping with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBRPCMapping) TryByID(ctx context.Context, p uint64) (*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_TryByID"
	l, e := a.fromQuery(ctx, qn, "id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("TryByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, nil
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, errors.Errorf("Multiple (%d) RPCMapping with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by multiple primary ids
func (a *DBRPCMapping) ByIDs(ctx context.Context, p []uint64) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByIDs"
	l, e := a.fromQuery(ctx, qn, "id in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("TryByID: error scanning (%s)", e))
	}
	return l, nil
}

// get all rows
func (a *DBRPCMapping) All(ctx context.Context) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_all"
	l, e := a.fromQuery(ctx, qn, "true")
	if e != nil {
		return nil, errors.Errorf("All: error scanning (%s)", e)
	}
	return l, nil
}

/**********************************************************************
* GetBy[FIELD] functions
**********************************************************************/

// get all "DBRPCMapping" rows with matching ServiceName
func (a *DBRPCMapping) ByServiceName(ctx context.Context, p string) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByServiceName"
	l, e := a.fromQuery(ctx, qn, "servicename = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBRPCMapping" rows with multiple matching ServiceName
func (a *DBRPCMapping) ByMultiServiceName(ctx context.Context, p []string) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByServiceName"
	l, e := a.fromQuery(ctx, qn, "servicename in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBRPCMapping) ByLikeServiceName(ctx context.Context, p string) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByLikeServiceName"
	l, e := a.fromQuery(ctx, qn, "servicename ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBRPCMapping" rows with matching FQDNService
func (a *DBRPCMapping) ByFQDNService(ctx context.Context, p string) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByFQDNService"
	l, e := a.fromQuery(ctx, qn, "fqdnservice = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByFQDNService: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBRPCMapping" rows with multiple matching FQDNService
func (a *DBRPCMapping) ByMultiFQDNService(ctx context.Context, p []string) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByFQDNService"
	l, e := a.fromQuery(ctx, qn, "fqdnservice in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByFQDNService: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBRPCMapping) ByLikeFQDNService(ctx context.Context, p string) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByLikeFQDNService"
	l, e := a.fromQuery(ctx, qn, "fqdnservice ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByFQDNService: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBRPCMapping" rows with matching RPCName
func (a *DBRPCMapping) ByRPCName(ctx context.Context, p string) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByRPCName"
	l, e := a.fromQuery(ctx, qn, "rpcname = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByRPCName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBRPCMapping" rows with multiple matching RPCName
func (a *DBRPCMapping) ByMultiRPCName(ctx context.Context, p []string) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByRPCName"
	l, e := a.fromQuery(ctx, qn, "rpcname in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByRPCName: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBRPCMapping) ByLikeRPCName(ctx context.Context, p string) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByLikeRPCName"
	l, e := a.fromQuery(ctx, qn, "rpcname ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByRPCName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBRPCMapping" rows with matching Active
func (a *DBRPCMapping) ByActive(ctx context.Context, p bool) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByActive"
	l, e := a.fromQuery(ctx, qn, "active = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByActive: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBRPCMapping" rows with multiple matching Active
func (a *DBRPCMapping) ByMultiActive(ctx context.Context, p []bool) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByActive"
	l, e := a.fromQuery(ctx, qn, "active in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByActive: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBRPCMapping) ByLikeActive(ctx context.Context, p bool) ([]*savepb.RPCMapping, error) {
	qn := "DBRPCMapping_ByLikeActive"
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
func (a *DBRPCMapping) get_ID(p *savepb.RPCMapping) uint64 {
	return uint64(p.ID)
}

// getter for field "ServiceName" (ServiceName) [string]
func (a *DBRPCMapping) get_ServiceName(p *savepb.RPCMapping) string {
	return string(p.ServiceName)
}

// getter for field "FQDNService" (FQDNService) [string]
func (a *DBRPCMapping) get_FQDNService(p *savepb.RPCMapping) string {
	return string(p.FQDNService)
}

// getter for field "RPCName" (RPCName) [string]
func (a *DBRPCMapping) get_RPCName(p *savepb.RPCMapping) string {
	return string(p.RPCName)
}

// getter for field "Active" (Active) [bool]
func (a *DBRPCMapping) get_Active(p *savepb.RPCMapping) bool {
	return bool(p.Active)
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBRPCMapping) ByDBQuery(ctx context.Context, query *Query) ([]*savepb.RPCMapping, error) {
	extra_fields, err := extraFieldsToQuery(ctx, a)
	if err != nil {
		return nil, err
	}
	i := 0
	for col_name, value := range extra_fields {
		i++
		//		efname := fmt.Sprintf("EXTRA_FIELD_%d", i)
		query.AddEqual(col_name, value)
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

func (a *DBRPCMapping) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.RPCMapping, error) {
	return a.fromQuery(ctx, "custom_query_"+a.Tablename(), query_where, args...)
}

// from a query snippet (the part after WHERE)
func (a *DBRPCMapping) fromQuery(ctx context.Context, queryname string, query_where string, args ...interface{}) ([]*savepb.RPCMapping, error) {
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
func (a *DBRPCMapping) get_col_from_proto(p *savepb.RPCMapping, colname string) interface{} {
	if colname == "id" {
		return a.get_ID(p)
	} else if colname == "servicename" {
		return a.get_ServiceName(p)
	} else if colname == "fqdnservice" {
		return a.get_FQDNService(p)
	} else if colname == "rpcname" {
		return a.get_RPCName(p)
	} else if colname == "active" {
		return a.get_Active(p)
	}
	panic(fmt.Sprintf("in table \"%s\", column \"%s\" cannot be resolved to proto field name", a.Tablename(), colname))
}

func (a *DBRPCMapping) Tablename() string {
	return a.SQLTablename
}

func (a *DBRPCMapping) SelectCols() string {
	return "id,servicename, fqdnservice, rpcname, active"
}
func (a *DBRPCMapping) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".servicename, " + a.SQLTablename + ".fqdnservice, " + a.SQLTablename + ".rpcname, " + a.SQLTablename + ".active"
}

func (a *DBRPCMapping) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.RPCMapping, error) {
	var res []*savepb.RPCMapping
	for rows.Next() {
		// SCANNER:
		foo := &savepb.RPCMapping{}
		// create the non-nullable pointers
		// create variables for scan results
		scanTarget_0 := &foo.ID
		scanTarget_1 := &foo.ServiceName
		scanTarget_2 := &foo.FQDNService
		scanTarget_3 := &foo.RPCName
		scanTarget_4 := &foo.Active
		err := rows.Scan(scanTarget_0, scanTarget_1, scanTarget_2, scanTarget_3, scanTarget_4)
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
func (a *DBRPCMapping) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),servicename text not null ,fqdnservice text not null ,rpcname text not null ,active boolean not null );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),servicename text not null ,fqdnservice text not null ,rpcname text not null ,active boolean not null );`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS servicename text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS fqdnservice text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS rpcname text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS active boolean not null default false;`,

		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS servicename text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS fqdnservice text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS rpcname text not null  default '';`,
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
func (a *DBRPCMapping) Error(ctx context.Context, q string, e error) error {
	if e == nil {
		return nil
	}
	return errors.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}

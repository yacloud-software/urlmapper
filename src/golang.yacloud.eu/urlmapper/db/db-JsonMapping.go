package db

/*
 This file was created by mkdb-client.
 The intention is not to modify thils file, but you may extend the struct DBJsonMapping
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence jsonmapping_seq;

Main Table:

 CREATE TABLE jsonmapping (id integer primary key default nextval('jsonmapping_seq'),domain text not null  ,path text not null  ,serviceid text not null  ,groupid text not null  ,fqdnservicename text not null  ,servicename text not null  );

Alter statements:
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS domain text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS path text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS serviceid text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS groupid text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS fqdnservicename text not null default '';
ALTER TABLE jsonmapping ADD COLUMN IF NOT EXISTS servicename text not null default '';


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE jsonmapping_archive (id integer unique not null,domain text not null,path text not null,serviceid text not null,groupid text not null,fqdnservicename text not null,servicename text not null);
*/

import (
	"context"
	gosql "database/sql"
	"fmt"
	"golang.conradwood.net/go-easyops/sql"
	savepb "golang.yacloud.eu/apis/urlmapper"
	"os"
)

var (
	default_def_DBJsonMapping *DBJsonMapping
)

type DBJsonMapping struct {
	DB                  *sql.DB
	SQLTablename        string
	SQLArchivetablename string
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

// archive. It is NOT transactionally save.
func (a *DBJsonMapping) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBJsonMapping", "insert into "+a.SQLArchivetablename+" (id,domain, path, serviceid, groupid, fqdnservicename, servicename) values ($1,$2, $3, $4, $5, $6, $7) ", p.ID, p.Domain, p.Path, p.ServiceID, p.GroupID, p.FQDNServiceName, p.ServiceName)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBJsonMapping) Save(ctx context.Context, p *savepb.JsonMapping) (uint64, error) {
	qn := "DBJsonMapping_Save"
	rows, e := a.DB.QueryContext(ctx, qn, "insert into "+a.SQLTablename+" (domain, path, serviceid, groupid, fqdnservicename, servicename) values ($1, $2, $3, $4, $5, $6) returning id", a.get_Domain(p), a.get_Path(p), a.get_ServiceID(p), a.get_GroupID(p), a.get_FQDNServiceName(p), a.get_ServiceName(p))
	if e != nil {
		return 0, a.Error(ctx, qn, e)
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, a.Error(ctx, qn, fmt.Errorf("No rows after insert"))
	}
	var id uint64
	e = rows.Scan(&id)
	if e != nil {
		return 0, a.Error(ctx, qn, fmt.Errorf("failed to scan id after insert: %s", e))
	}
	p.ID = id
	return id, nil
}

// Save using the ID specified
func (a *DBJsonMapping) SaveWithID(ctx context.Context, p *savepb.JsonMapping) error {
	qn := "insert_DBJsonMapping"
	_, e := a.DB.ExecContext(ctx, qn, "insert into "+a.SQLTablename+" (id,domain, path, serviceid, groupid, fqdnservicename, servicename) values ($1,$2, $3, $4, $5, $6, $7) ", p.ID, p.Domain, p.Path, p.ServiceID, p.GroupID, p.FQDNServiceName, p.ServiceName)
	return a.Error(ctx, qn, e)
}

func (a *DBJsonMapping) Update(ctx context.Context, p *savepb.JsonMapping) error {
	qn := "DBJsonMapping_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set domain=$1, path=$2, serviceid=$3, groupid=$4, fqdnservicename=$5, servicename=$6 where id = $7", a.get_Domain(p), a.get_Path(p), a.get_ServiceID(p), a.get_GroupID(p), a.get_FQDNServiceName(p), a.get_ServiceName(p), p.ID)

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
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, fmt.Errorf("No JsonMapping with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) JsonMapping with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBJsonMapping) TryByID(ctx context.Context, p uint64) (*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_TryByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("TryByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("TryByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, nil
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) JsonMapping with id %v", len(l), p))
	}
	return l[0], nil
}

// get all rows
func (a *DBJsonMapping) All(ctx context.Context) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_all"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" order by id")
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("All: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("All: error scanning (%s)", e)
	}
	return l, nil
}

/**********************************************************************
* GetBy[FIELD] functions
**********************************************************************/

// get all "DBJsonMapping" rows with matching Domain
func (a *DBJsonMapping) ByDomain(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByDomain"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where domain = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByDomain: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByDomain: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeDomain(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeDomain"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where domain ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByDomain: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByDomain: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching Path
func (a *DBJsonMapping) ByPath(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByPath"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where path = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPath: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPath: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikePath(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikePath"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where path ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPath: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPath: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching ServiceID
func (a *DBJsonMapping) ByServiceID(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByServiceID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where serviceid = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByServiceID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByServiceID: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeServiceID(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeServiceID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where serviceid ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByServiceID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByServiceID: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching GroupID
func (a *DBJsonMapping) ByGroupID(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByGroupID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where groupid = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByGroupID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByGroupID: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeGroupID(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeGroupID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where groupid ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByGroupID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByGroupID: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching FQDNServiceName
func (a *DBJsonMapping) ByFQDNServiceName(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByFQDNServiceName"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where fqdnservicename = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByFQDNServiceName: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByFQDNServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeFQDNServiceName(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeFQDNServiceName"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where fqdnservicename ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByFQDNServiceName: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByFQDNServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBJsonMapping" rows with matching ServiceName
func (a *DBJsonMapping) ByServiceName(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByServiceName"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where servicename = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByServiceName: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByServiceName: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBJsonMapping) ByLikeServiceName(ctx context.Context, p string) ([]*savepb.JsonMapping, error) {
	qn := "DBJsonMapping_ByLikeServiceName"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, path, serviceid, groupid, fqdnservicename, servicename from "+a.SQLTablename+" where servicename ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByServiceName: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByServiceName: error scanning (%s)", e))
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

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBJsonMapping) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.JsonMapping, error) {
	rows, err := a.DB.QueryContext(ctx, "custom_query_"+a.Tablename(), "select "+a.SelectCols()+" from "+a.Tablename()+" where "+query_where, args...)
	if err != nil {
		return nil, err
	}
	return a.FromRows(ctx, rows)
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBJsonMapping) Tablename() string {
	return a.SQLTablename
}

func (a *DBJsonMapping) SelectCols() string {
	return "id,domain, path, serviceid, groupid, fqdnservicename, servicename"
}
func (a *DBJsonMapping) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".domain, " + a.SQLTablename + ".path, " + a.SQLTablename + ".serviceid, " + a.SQLTablename + ".groupid, " + a.SQLTablename + ".fqdnservicename, " + a.SQLTablename + ".servicename"
}

func (a *DBJsonMapping) FromRowsOld(ctx context.Context, rows *gosql.Rows) ([]*savepb.JsonMapping, error) {
	var res []*savepb.JsonMapping
	for rows.Next() {
		foo := savepb.JsonMapping{}
		err := rows.Scan(&foo.ID, &foo.Domain, &foo.Path, &foo.ServiceID, &foo.GroupID, &foo.FQDNServiceName, &foo.ServiceName)
		if err != nil {
			return nil, a.Error(ctx, "fromrow-scan", err)
		}
		res = append(res, &foo)
	}
	return res, nil
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
		err := rows.Scan(scanTarget_0, scanTarget_1, scanTarget_2, scanTarget_3, scanTarget_4, scanTarget_5, scanTarget_6)
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
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),domain text not null ,path text not null ,serviceid text not null ,groupid text not null ,fqdnservicename text not null ,servicename text not null );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),domain text not null ,path text not null ,serviceid text not null ,groupid text not null ,fqdnservicename text not null ,servicename text not null );`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS domain text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS path text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS serviceid text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS groupid text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS fqdnservicename text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS servicename text not null default '';`,

		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS domain text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS path text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS serviceid text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS groupid text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS fqdnservicename text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS servicename text not null  default '';`,
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
	return fmt.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}


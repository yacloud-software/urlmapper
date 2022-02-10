package db

/*
 This file was created by mkdb-client.
 The intention is not to modify thils file, but you may extend the struct DBAnyHostMapping
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence anyhostmapping_seq;

Main Table:

 CREATE TABLE anyhostmapping (id integer primary key default nextval('anyhostmapping_seq'),path text not null  unique  ,serviceid text not null  );

Alter statements:
ALTER TABLE anyhostmapping ADD COLUMN path text not null unique  default '';
ALTER TABLE anyhostmapping ADD COLUMN serviceid text not null default '';


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE anyhostmapping_archive (id integer unique not null,path text not null,serviceid text not null);
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
	default_def_DBAnyHostMapping *DBAnyHostMapping
)

type DBAnyHostMapping struct {
	DB                  *sql.DB
	SQLTablename        string
	SQLArchivetablename string
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

// archive. It is NOT transactionally save.
func (a *DBAnyHostMapping) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBAnyHostMapping", "insert into "+a.SQLArchivetablename+"+ (id,path, serviceid) values ($1,$2, $3) ", p.ID, p.Path, p.ServiceID)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBAnyHostMapping) Save(ctx context.Context, p *savepb.AnyHostMapping) (uint64, error) {
	qn := "DBAnyHostMapping_Save"
	rows, e := a.DB.QueryContext(ctx, qn, "insert into "+a.SQLTablename+" (path, serviceid) values ($1, $2) returning id", p.Path, p.ServiceID)
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
func (a *DBAnyHostMapping) SaveWithID(ctx context.Context, p *savepb.AnyHostMapping) error {
	qn := "insert_DBAnyHostMapping"
	_, e := a.DB.ExecContext(ctx, qn, "insert into "+a.SQLTablename+" (id,path, serviceid) values ($1,$2, $3) ", p.ID, p.Path, p.ServiceID)
	return a.Error(ctx, qn, e)
}

func (a *DBAnyHostMapping) Update(ctx context.Context, p *savepb.AnyHostMapping) error {
	qn := "DBAnyHostMapping_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set path=$1, serviceid=$2 where id = $3", p.Path, p.ServiceID, p.ID)

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
	rows, e := a.DB.QueryContext(ctx, qn, "select id,path, serviceid from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, fmt.Errorf("No AnyHostMapping with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) AnyHostMapping with id %v", len(l), p))
	}
	return l[0], nil
}

// get all rows
func (a *DBAnyHostMapping) All(ctx context.Context) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_all"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,path, serviceid from "+a.SQLTablename+" order by id")
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

// get all "DBAnyHostMapping" rows with matching Path
func (a *DBAnyHostMapping) ByPath(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByPath"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,path, serviceid from "+a.SQLTablename+" where path = $1", p)
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
func (a *DBAnyHostMapping) ByLikePath(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByLikePath"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,path, serviceid from "+a.SQLTablename+" where path ilike $1", p)
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

// get all "DBAnyHostMapping" rows with matching ServiceID
func (a *DBAnyHostMapping) ByServiceID(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByServiceID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,path, serviceid from "+a.SQLTablename+" where serviceid = $1", p)
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
func (a *DBAnyHostMapping) ByLikeServiceID(ctx context.Context, p string) ([]*savepb.AnyHostMapping, error) {
	qn := "DBAnyHostMapping_ByLikeServiceID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,path, serviceid from "+a.SQLTablename+" where serviceid ilike $1", p)
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

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBAnyHostMapping) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.AnyHostMapping, error) {
	rows, err := a.DB.QueryContext(ctx, "custom_query_"+a.Tablename(), "select "+a.SelectCols()+" from "+a.Tablename()+" where "+query_where, args...)
	if err != nil {
		return nil, err
	}
	return a.FromRows(ctx, rows)
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBAnyHostMapping) Tablename() string {
	return a.SQLTablename
}

func (a *DBAnyHostMapping) SelectCols() string {
	return "id,path, serviceid"
}
func (a *DBAnyHostMapping) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".path, " + a.SQLTablename + ".serviceid"
}

func (a *DBAnyHostMapping) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.AnyHostMapping, error) {
	var res []*savepb.AnyHostMapping
	for rows.Next() {
		foo := savepb.AnyHostMapping{}
		err := rows.Scan(&foo.ID, &foo.Path, &foo.ServiceID)
		if err != nil {
			return nil, a.Error(ctx, "fromrow-scan", err)
		}
		res = append(res, &foo)
	}
	return res, nil
}

/**********************************************************************
* Helper to create table and columns
**********************************************************************/
func (a *DBAnyHostMapping) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),path text not null  unique  ,serviceid text not null  );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),path text not null  unique  ,serviceid text not null  );`,
	}
	for i, c := range csql {
		_, e := a.DB.ExecContext(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
		if e != nil {
			return e
		}
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
	return fmt.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}

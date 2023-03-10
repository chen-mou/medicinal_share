// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"medicinal_share/main/entity"
)

func newFile(db *gorm.DB, opts ...gen.DOOption) file {
	_file := file{}

	_file.fileDo.UseDB(db, opts...)
	_file.fileDo.UseModel(&entity.File{})

	tableName := _file.fileDo.TableName()
	_file.ALL = field.NewAsterisk(tableName)
	_file.Id = field.NewInt64(tableName, "id")
	_file.Hash = field.NewString(tableName, "hash")
	_file.Machine = field.NewString(tableName, "machine")
	_file.Status = field.NewString(tableName, "status")
	_file.Path = field.NewString(tableName, "path")
	_file.Uri = field.NewString(tableName, "uri")

	_file.fillFieldMap()

	return _file
}

type file struct {
	fileDo fileDo

	ALL     field.Asterisk
	Id      field.Int64
	Hash    field.String
	Machine field.String
	Status  field.String
	Path    field.String
	Uri     field.String

	fieldMap map[string]field.Expr
}

func (f file) Table(newTableName string) *file {
	f.fileDo.UseTable(newTableName)
	return f.updateTableName(newTableName)
}

func (f file) As(alias string) *file {
	f.fileDo.DO = *(f.fileDo.As(alias).(*gen.DO))
	return f.updateTableName(alias)
}

func (f *file) updateTableName(table string) *file {
	f.ALL = field.NewAsterisk(table)
	f.Id = field.NewInt64(table, "id")
	f.Hash = field.NewString(table, "hash")
	f.Machine = field.NewString(table, "machine")
	f.Status = field.NewString(table, "status")
	f.Path = field.NewString(table, "path")
	f.Uri = field.NewString(table, "uri")

	f.fillFieldMap()

	return f
}

func (f *file) WithContext(ctx context.Context) IFileDo { return f.fileDo.WithContext(ctx) }

func (f file) TableName() string { return f.fileDo.TableName() }

func (f file) Alias() string { return f.fileDo.Alias() }

func (f *file) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := f.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (f *file) fillFieldMap() {
	f.fieldMap = make(map[string]field.Expr, 6)
	f.fieldMap["id"] = f.Id
	f.fieldMap["hash"] = f.Hash
	f.fieldMap["machine"] = f.Machine
	f.fieldMap["status"] = f.Status
	f.fieldMap["path"] = f.Path
	f.fieldMap["uri"] = f.Uri
}

func (f file) clone(db *gorm.DB) file {
	f.fileDo.ReplaceConnPool(db.Statement.ConnPool)
	return f
}

func (f file) replaceDB(db *gorm.DB) file {
	f.fileDo.ReplaceDB(db)
	return f
}

type fileDo struct{ gen.DO }

type IFileDo interface {
	gen.SubQuery
	Debug() IFileDo
	WithContext(ctx context.Context) IFileDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IFileDo
	WriteDB() IFileDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IFileDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IFileDo
	Not(conds ...gen.Condition) IFileDo
	Or(conds ...gen.Condition) IFileDo
	Select(conds ...field.Expr) IFileDo
	Where(conds ...gen.Condition) IFileDo
	Order(conds ...field.Expr) IFileDo
	Distinct(cols ...field.Expr) IFileDo
	Omit(cols ...field.Expr) IFileDo
	Join(table schema.Tabler, on ...field.Expr) IFileDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IFileDo
	RightJoin(table schema.Tabler, on ...field.Expr) IFileDo
	Group(cols ...field.Expr) IFileDo
	Having(conds ...gen.Condition) IFileDo
	Limit(limit int) IFileDo
	Offset(offset int) IFileDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IFileDo
	Unscoped() IFileDo
	Create(values ...*entity.File) error
	CreateInBatches(values []*entity.File, batchSize int) error
	Save(values ...*entity.File) error
	First() (*entity.File, error)
	Take() (*entity.File, error)
	Last() (*entity.File, error)
	Find() ([]*entity.File, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.File, err error)
	FindInBatches(result *[]*entity.File, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*entity.File) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IFileDo
	Assign(attrs ...field.AssignExpr) IFileDo
	Joins(fields ...field.RelationField) IFileDo
	Preload(fields ...field.RelationField) IFileDo
	FirstOrInit() (*entity.File, error)
	FirstOrCreate() (*entity.File, error)
	FindByPage(offset int, limit int) (result []*entity.File, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IFileDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	FilterById(id int64) (result *entity.File, err error)
}

//FilterById
//
//select * from @@table where id = @id
func (f fileDo) FilterById(id int64) (result *entity.File, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("select * from file where id = ? ")

	var executeSQL *gorm.DB

	executeSQL = f.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result)
	err = executeSQL.Error
	return
}

func (f fileDo) Debug() IFileDo {
	return f.withDO(f.DO.Debug())
}

func (f fileDo) WithContext(ctx context.Context) IFileDo {
	return f.withDO(f.DO.WithContext(ctx))
}

func (f fileDo) ReadDB() IFileDo {
	return f.Clauses(dbresolver.Read)
}

func (f fileDo) WriteDB() IFileDo {
	return f.Clauses(dbresolver.Write)
}

func (f fileDo) Session(config *gorm.Session) IFileDo {
	return f.withDO(f.DO.Session(config))
}

func (f fileDo) Clauses(conds ...clause.Expression) IFileDo {
	return f.withDO(f.DO.Clauses(conds...))
}

func (f fileDo) Returning(value interface{}, columns ...string) IFileDo {
	return f.withDO(f.DO.Returning(value, columns...))
}

func (f fileDo) Not(conds ...gen.Condition) IFileDo {
	return f.withDO(f.DO.Not(conds...))
}

func (f fileDo) Or(conds ...gen.Condition) IFileDo {
	return f.withDO(f.DO.Or(conds...))
}

func (f fileDo) Select(conds ...field.Expr) IFileDo {
	return f.withDO(f.DO.Select(conds...))
}

func (f fileDo) Where(conds ...gen.Condition) IFileDo {
	return f.withDO(f.DO.Where(conds...))
}

func (f fileDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IFileDo {
	return f.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (f fileDo) Order(conds ...field.Expr) IFileDo {
	return f.withDO(f.DO.Order(conds...))
}

func (f fileDo) Distinct(cols ...field.Expr) IFileDo {
	return f.withDO(f.DO.Distinct(cols...))
}

func (f fileDo) Omit(cols ...field.Expr) IFileDo {
	return f.withDO(f.DO.Omit(cols...))
}

func (f fileDo) Join(table schema.Tabler, on ...field.Expr) IFileDo {
	return f.withDO(f.DO.Join(table, on...))
}

func (f fileDo) LeftJoin(table schema.Tabler, on ...field.Expr) IFileDo {
	return f.withDO(f.DO.LeftJoin(table, on...))
}

func (f fileDo) RightJoin(table schema.Tabler, on ...field.Expr) IFileDo {
	return f.withDO(f.DO.RightJoin(table, on...))
}

func (f fileDo) Group(cols ...field.Expr) IFileDo {
	return f.withDO(f.DO.Group(cols...))
}

func (f fileDo) Having(conds ...gen.Condition) IFileDo {
	return f.withDO(f.DO.Having(conds...))
}

func (f fileDo) Limit(limit int) IFileDo {
	return f.withDO(f.DO.Limit(limit))
}

func (f fileDo) Offset(offset int) IFileDo {
	return f.withDO(f.DO.Offset(offset))
}

func (f fileDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IFileDo {
	return f.withDO(f.DO.Scopes(funcs...))
}

func (f fileDo) Unscoped() IFileDo {
	return f.withDO(f.DO.Unscoped())
}

func (f fileDo) Create(values ...*entity.File) error {
	if len(values) == 0 {
		return nil
	}
	return f.DO.Create(values)
}

func (f fileDo) CreateInBatches(values []*entity.File, batchSize int) error {
	return f.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (f fileDo) Save(values ...*entity.File) error {
	if len(values) == 0 {
		return nil
	}
	return f.DO.Save(values)
}

func (f fileDo) First() (*entity.File, error) {
	if result, err := f.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*entity.File), nil
	}
}

func (f fileDo) Take() (*entity.File, error) {
	if result, err := f.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*entity.File), nil
	}
}

func (f fileDo) Last() (*entity.File, error) {
	if result, err := f.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*entity.File), nil
	}
}

func (f fileDo) Find() ([]*entity.File, error) {
	result, err := f.DO.Find()
	return result.([]*entity.File), err
}

func (f fileDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.File, err error) {
	buf := make([]*entity.File, 0, batchSize)
	err = f.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (f fileDo) FindInBatches(result *[]*entity.File, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return f.DO.FindInBatches(result, batchSize, fc)
}

func (f fileDo) Attrs(attrs ...field.AssignExpr) IFileDo {
	return f.withDO(f.DO.Attrs(attrs...))
}

func (f fileDo) Assign(attrs ...field.AssignExpr) IFileDo {
	return f.withDO(f.DO.Assign(attrs...))
}

func (f fileDo) Joins(fields ...field.RelationField) IFileDo {
	for _, _f := range fields {
		f = *f.withDO(f.DO.Joins(_f))
	}
	return &f
}

func (f fileDo) Preload(fields ...field.RelationField) IFileDo {
	for _, _f := range fields {
		f = *f.withDO(f.DO.Preload(_f))
	}
	return &f
}

func (f fileDo) FirstOrInit() (*entity.File, error) {
	if result, err := f.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*entity.File), nil
	}
}

func (f fileDo) FirstOrCreate() (*entity.File, error) {
	if result, err := f.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*entity.File), nil
	}
}

func (f fileDo) FindByPage(offset int, limit int) (result []*entity.File, count int64, err error) {
	result, err = f.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = f.Offset(-1).Limit(-1).Count()
	return
}

func (f fileDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = f.Count()
	if err != nil {
		return
	}

	err = f.Offset(offset).Limit(limit).Scan(result)
	return
}

func (f fileDo) Scan(result interface{}) (err error) {
	return f.DO.Scan(result)
}

func (f fileDo) Delete(models ...*entity.File) (result gen.ResultInfo, err error) {
	return f.DO.Delete(models)
}

func (f *fileDo) withDO(do gen.Dao) *fileDo {
	f.DO = *do.(*gen.DO)
	return f
}

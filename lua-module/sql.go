package lua_module

import (
	"database/sql"
	"fmt"
	. "github.com/Icy2010/TBys"
	lua "github.com/yuin/gopher-lua"
	"gorm.io/gorm"
	"reflect"
)

type TLuaSQL struct {
	db      *gorm.DB
	sqlType int
}

func (this *TLuaSQL) rowsToLuaTable(rows *sql.Rows, L *lua.LState, OnTable func(data *lua.LTable)) int {
	count := 0
	err := DoSqlData(rows, func(result TSqlData) {
		data := L.NewTable()
		for s, i := range result {
			val := reflect.ValueOf(i)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				data.RawSet(lua.LString(s), lua.LNumber(val.Int()))
			case reflect.Float64, reflect.Float32:
				data.RawSet(lua.LString(s), lua.LNumber(val.Float()))
			case reflect.String:
				data.RawSet(lua.LString(s), lua.LString(val.String()))
			case reflect.Bool:
				data.RawSet(lua.LString(s), lua.LBool(val.Bool()))
			default:
				continue
			}
		}

		count++
		if OnTable != nil {
			OnTable(data)
		}
	})
	if err != nil {
		Logger().Error(err)
	}

	return count
}

func (this *TLuaSQL) queryToLuaTable(rows *sql.Rows, L *lua.LState) *lua.LTable {
	table := L.NewTable()
	this.rowsToLuaTable(rows, L, func(data *lua.LTable) {
		table.Append(data)
	})

	return table
}

/*
Author: icy
Description: 数据库查询
Date:  2023/1/4 下午9:03
Param: 语句  参数2 空=返回表
参数操作方式 字符串
'table' 参数 3 回调 function(data)
'json' 参数3 文件路径  参数4 进度回调 function(index)
return: 参数2 为空返回表  反之返回记录数量
*/

func (this *TLuaSQL) query(L *lua.LState) int {
	sql := L.ToString(1)
	PDB := this.db
	if rows, err := PDB.Raw(sql).Rows(); err == nil {
		doData := L.ToFunction(2)
		c := 0
		if doData != nil {
			c = this.rowsToLuaTable(rows, L, func(data *lua.LTable) {
				err = L.CallByParam(lua.P{
					Fn:      doData,
					NRet:    0,
					Protect: true,
				}, data)
				if err != nil {
					Logger().Error(err)
				}
			})

			L.Push(lua.LNumber(c))
		} else {
			L.Push(this.queryToLuaTable(rows, L))
		}

	} else {
		Logger().Error(err)
	}

	return 1
}

func (this *TLuaSQL) execute(L *lua.LState) int {
	sql := L.ToString(1)
	n := this.db.Exec(sql).RowsAffected
	L.Push(lua.LNumber(n))
	return 1
}

/*
Author: icy
Description: 数据库插入
Date:  2022/12/16 23:39
Param: 1 表名 2 lua表{name=value}
return: 返回自增 如果有自增的话 没有就是 成功就是返回1反之0
*/

func (this *TLuaSQL) insert(L *lua.LState) int {
	Result := 0
	tableName := L.ToString(1)
	table := L.ToTable(2)
	if table != nil {
		db := this.db
		sql, values := LuaTableToSqlInsert(this.sqlType, tableName, table)
		ra := db.Exec(sql, values...).RowsAffected
		if ra > 0 {
			switch this.sqlType {
			case CDB_MYSQL:
				db.Raw(`select LAST_INSERT_ID() as id `).Row().Scan(&Result)
			case CDB_SQLITE:
				db.Raw(`select last_insert_rowid() as id `).Row().Scan(&Result)
			case CDB_MSSQL:
				db.Raw(fmt.Sprintf(`select ident_current('%s') as id`, tableName)).Row().Scan(&Result)
			default:
				Result = int(db.RowsAffected)
			}

			if Result == 0 {
				Result = int(ra)
			}
		}
	}

	L.Push(lua.LNumber(Result))
	return 1
}

/*
Author: icy
Description: 数据库更新
Date:  2022/12/12 23:35
Param: 1 数据库表名  2 lua表{name=value} 3 条件
return: 成功返回1 反之0
*/

func (this *TLuaSQL) update(L *lua.LState) int {
	Result := 0

	tableName := L.ToString(1)
	table := L.ToTable(2)
	where := L.ToString(3)
	if table != nil {
		data := LuaTableToMap(table)
		Result = int(this.db.Table(tableName).Where(where).Updates(data).RowsAffected)
	}

	L.Push(lua.LNumber(Result))
	return 1
}

/*
Author: icy
Description: Mysql数据库连接
Date:  2022/12/13 23:24
Param: 是一个lua表  {ip = "",name="",pw = "",db="",port="",charset = ""}
return: 布尔值 是否成功
*/

func GetMySqlOptions(table *lua.LTable) TOptSQL {
	return TOptSQL{
		Host:     lua.LVAsString(table.RawGetString(`host`)),
		UserName: lua.LVAsString(table.RawGetString(`name`)),
		PassWord: lua.LVAsString(table.RawGetString(`pw`)),
		DataBase: lua.LVAsString(table.RawGetString(`db`)),
		Charset: func() string {
			val := lua.LVAsString(table.RawGetString(`charset`))
			if val == "" {
				val = "utf8mb4"
			}
			return val
		}(),
		Port: func() string {
			val := lua.LVAsString(table.RawGetString(`port`))
			if val == "" {
				val = "3306"
			}
			return val
		}(),
	}
}

func (this *TLuaSQL) mysql(L *lua.LState) int {
	this.sqlType = CDB_MYSQL
	table := L.ToTable(1)
	succ := false
	if table != nil {
		data := GetMySqlOptions(table)

		var err error = nil
		this.db, err = CreateMySQLDB(data)

		succ = err == nil
	}

	L.Push(lua.LBool(succ))
	return 1
}

/*
Author: icy
Description: sqlite 连接
Date:  2022/12/13 23:20
Param: 文件路径 字符串
return: 布尔值是否成功
*/

func (this *TLuaSQL) sqlite(L *lua.LState) int {
	fileName := L.ToString(1)
	succ := false
	this.sqlType = CDB_SQLITE
	if PathExist(fileName) {
		var err error = nil
		this.db, err = CreateSQLiteDB(fileName)

		succ = err == nil
	}

	L.Push(lua.LBool(succ))
	return 1
}

/*
Author: icy
Description: sqlServer 连接函数
Date:  2022/12/13 23:18
Param: 是一个lua表  {ip = "",name="",pw = "",db="",port=""}
return: 布尔值 是否成功
*/

func (this *TLuaSQL) mssql(L *lua.LState) int {
	this.sqlType = CDB_MSSQL
	table := L.ToTable(1)
	succ := false
	if table != nil {
		data := TOptSQL{
			Host:     lua.LVAsString(table.RawGetString(`host`)),
			UserName: lua.LVAsString(table.RawGetString(`name`)),
			PassWord: lua.LVAsString(table.RawGetString(`pw`)),
			DataBase: lua.LVAsString(table.RawGetString(`db`)),
			Port: func() string {
				val := lua.LVAsString(table.RawGetString(`port`))
				if val == "" {
					val = "1433"
				}
				return val
			}(),
		}

		var err error = nil
		this.db, err = CreateMSSQLDB(data)

		succ = err == nil
	}

	L.Push(lua.LBool(succ))
	return 1
}

/*--------------------------------------------------------------------------------------------------------------------*/

func SqlRegisterModule(L *lua.LState) {
	sql := &TLuaSQL{}
	L.RegisterModule(`sql`, map[string]lua.LGFunction{
		`mysql`:  sql.mysql,
		`sqlite`: sql.sqlite,
		`mssql`:  sql.mssql,

		`query`:   sql.query,
		`execute`: sql.execute,
		`insert`:  sql.insert,
		`update`:  sql.update,
	})
}

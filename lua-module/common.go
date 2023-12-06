package lua_module

import (
	"encoding/json"
	"fmt"
	. "github.com/Icy2010/TBys"
	"github.com/tidwall/gjson"
	lua "github.com/yuin/gopher-lua"
	"reflect"
	"strings"
)

func LuaLoggedError(err error) {
	Logger().Error(`[Lua-Error]`, err)
}

func LuaTableToMySQLInsert(TableName string, table *lua.LTable) (string, []any) {
	Result := "insert into `" + TableName + "`("
	Keys := ""
	Values := ""
	list := make([]any, 0)
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		if key.Type() == lua.LTString {
			Keys += fmt.Sprintf("`%s`,", key.String())
			Values += `?,`
			switch value.Type() {
			case lua.LTString:
				list = append(list, lua.LVAsString(value))
			case lua.LTNumber:
				list = append(list, lua.LVAsNumber(value))
			case lua.LTBool:
				list = append(list, lua.LVAsBool(value))
			default:
				list = append(list, 0)
			}

		}
	})

	Keys = strings.Trim(Keys, `,`)
	Values = strings.Trim(Values, `,`)
	Result += fmt.Sprintf(`%s)VALUES(%s)`, Keys, Values)

	return Result, list
}

func LuaTableToSqliteInsert(TableName string, table *lua.LTable) (string, []any) {
	Result := `insert into " + TableName + "(`
	Keys := ""
	Values := ""
	list := make([]any, 0)
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		if key.Type() == lua.LTString {
			Keys += fmt.Sprintf(`"%s",`, key.String())
			Values += `?,`
			switch value.Type() {
			case lua.LTString:
				list = append(list, lua.LVAsString(value))
			case lua.LTNumber:
				list = append(list, lua.LVAsNumber(value))
			case lua.LTBool:
				list = append(list, lua.LVAsBool(value))
			default:
				list = append(list, 0)
			}

		}
	})

	Keys = strings.Trim(Keys, `,`)
	Values = strings.Trim(Values, `,`)
	Result += fmt.Sprintf(`%s)VALUES(%s)`, Keys, Values)

	return Result, list
}

func LuaTableToSqlInsert(way int, tableName string, table *lua.LTable) (string, []any) {
	switch way {
	case CDB_MYSQL:
		return LuaTableToMySQLInsert(tableName, table)
	default:
		return LuaTableToSqliteInsert(tableName, table)
	}
}

func LuaTableToMap(table *lua.LTable) map[string]interface{} {
	data := make(map[string]interface{}, 0)
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		if key.Type() == lua.LTString {
			switch value.Type() {
			case lua.LTString:
				data[lua.LVAsString(key)] = lua.LVAsString(value)
			case lua.LTNumber:
				data[lua.LVAsString(key)] = lua.LVAsNumber(value)
			case lua.LTBool:
				data[lua.LVAsString(key)] = lua.LVAsBool(value)
			default:
				break
			}

		}
	})

	return data
}

func LuaTableToJSON(table *lua.LTable) string {
	data := LuaTableToMap(table)

	if len(data) > 0 {
		if bytes, err := json.Marshal(data); err == nil {
			return string(bytes)
		}

	}

	return "{}"
}

func LuaTableToStrings(table *lua.LTable) []string {
	data := make([]string, 0)
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		if value.Type() == lua.LTString {
			data = append(data, lua.LVAsString(value))
		}
	})
	return data
}

func LauTableView(table *lua.LTable) string {
	buffer := "{\n"
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		if key.Type() == lua.LTString || key.Type() == lua.LTNumber {
			buffer += key.String() + ` = ` + value.String() + " ,\n"
		}
	})
	if len(buffer) > 2 {
		return buffer[:len(buffer)-2] + "\n}"
	}

	return `{}`
}

func LuaTableToMapString(table *lua.LTable) map[string]string {
	data := make(map[string]string, 0)
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		if key.Type() == lua.LTString {
			switch value.Type() {
			case lua.LTString:
				data[lua.LVAsString(key)] = lua.LVAsString(value)
			case lua.LTNumber:
				data[lua.LVAsString(key)] = value.String()
			case lua.LTBool:
				data[lua.LVAsString(key)] = value.String()
			default:
				break
			}

		}
	})

	return data
}

func ToLuaValue(val any) lua.LValue {
	v := reflect.ValueOf(val)
	switch v.Type().Kind() {
	case reflect.String:
		return lua.LString(val.(string))
	case reflect.Int:
		return lua.LNumber(val.(int))
	case reflect.Int16:
		return lua.LNumber(val.(int16))
	case reflect.Int8:
		return lua.LNumber(val.(int8))
	case reflect.Int32:
		return lua.LNumber(val.(int32))
	case reflect.Int64:
		return lua.LNumber(val.(int64))
	case reflect.Float32:
		return lua.LNumber(val.(float32))
	case reflect.Float64:
		return lua.LNumber(val.(float64))
	case reflect.Uint:
		return lua.LNumber(val.(uint))
	case reflect.Uint8:
		return lua.LNumber(val.(uint8))
	case reflect.Uint16:
		return lua.LNumber(val.(uint16))
	case reflect.Uint32:
		return lua.LNumber(val.(uint32))
	case reflect.Uint64:
		return lua.LNumber(val.(uint64))
	case reflect.Bool:
		return lua.LNumber(val.(uint64))
	}
	return lua.LNil
}

func SetLuaTableValue(table *lua.LTable, key string, val any) {
	v := reflect.ValueOf(val)
	switch v.Type().Kind() {
	case reflect.String:
		table.RawSetString(key, lua.LString(val.(string)))
	case reflect.Int:
		table.RawSetString(key, lua.LNumber(val.(int)))
	case reflect.Int16:
		table.RawSetString(key, lua.LNumber(val.(int16)))
	case reflect.Int8:
		table.RawSetString(key, lua.LNumber(val.(int8)))
	case reflect.Int32:
		table.RawSetString(key, lua.LNumber(val.(int32)))
	case reflect.Int64:
		table.RawSetString(key, lua.LNumber(val.(int64)))
	case reflect.Float32:
		table.RawSetString(key, lua.LNumber(val.(float32)))
	case reflect.Float64:
		table.RawSetString(key, lua.LNumber(val.(float64)))
	case reflect.Uint:
		table.RawSetString(key, lua.LNumber(val.(uint)))
	case reflect.Uint8:
		table.RawSetString(key, lua.LNumber(val.(uint8)))
	case reflect.Uint16:
		table.RawSetString(key, lua.LNumber(val.(uint16)))
	case reflect.Uint32:
		table.RawSetString(key, lua.LNumber(val.(uint32)))
	case reflect.Uint64:
		table.RawSetString(key, lua.LNumber(val.(uint64)))
	case reflect.Bool:
		table.RawSetString(key, lua.LBool(val.(bool)))
	}
}

//json

func JSONObjectToSqlInsert(jo gjson.Result, tableName string) (string, []any) {
	if jo.IsObject() {
		Names := ``
		Values := ``
		list := make([]any, 0)

		for k, v := range jo.Map() {
			Names += k
			Names += `,`
			Values += `?,`

			list = append(list, v)
		}

		return fmt.Sprintf(`insert into %s (%s)values(%s)`, tableName, Names, Values), list

	}
	return ``, nil
}

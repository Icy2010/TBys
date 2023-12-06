package TBys

import (
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/glebarez/sqlite"
	"github.com/tidwall/gjson"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type TResultBase struct {
	State   int    `json:"state"`
	Message string `json:"msg"`
}

type TResult struct {
	State   int         `json:"state"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type TResultError TResultBase

type TResultCommon struct {
	State   int           `json:"state"`
	Message string        `json:"msg"`
	Data    []interface{} `json:"data"`
}

type TFormLimit struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type TSqlData map[string]interface{}

func StrToInt(s string) int {
	i, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return 0
	}

	return int(i)
}

func StrToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0
	}

	return i
}

func StrToFloat(s string) float64 {
	i, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0
	}

	return i
}

func StrToFloat64(s string) float64 {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return i
}

func Substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ExtractFileName(path string) string {
	_, fileName := filepath.Split(path)
	return fileName
}

func IncludeTrailingPathDelimiter(path string) string {
	if len(path) > 0 {
		d := PathSeparators()
		s := path[len(path)-1:]
		if d != s {
			return path + d
		}
	}

	return path
}

func SHA1(val []byte) string {
	h := sha1.New()
	h.Write(val)
	return hex.EncodeToString(h.Sum(nil))
}

func MD5(str string) string {
	_md5 := md5.New()
	_md5.Write([]byte(str))
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func MD5Data(reader io.Reader) string {
	_md5 := md5.New()
	if _, err := io.Copy(_md5, reader); err == nil {
		return hex.EncodeToString(_md5.Sum(nil))
	}
	return ""
}

func MD5File(fileName string) string {
	f, err := os.Open(fileName)
	if err == nil {
		return MD5Data(f)
	}
	return ""
}

func GetFileBytes(fileName string) []byte {
	if PathExist(fileName) {
		if file, err := os.Open(fileName); err == nil {
			defer file.Close()
			info, _ := file.Stat()
			buf := make([]byte, info.Size())
			file.Read(buf)

			return buf
		}
	}
	return nil
}

const (
	//time format

	FormatDay      = "2006-01-02"
	FormatSecond   = "2006-01-02 15:04:05"
	FormatMinute   = "2006-01-02 15:04"
	FormatOnlyHour = "15:04"
	FromatYMDHMS   = `20060102150405`
)

/*--------------------------------------------------------------------------------------------------------------------*/

var (
	cacheDefaultExpiration time.Duration
)

/*--------------------------------------------------------------------------------------------------------------------*/

func CacheDefExpiration() time.Duration {
	return cacheDefaultExpiration
}

func MakeDir(path string) error {
	if !PathExist(path) {
		err := os.Mkdir(path, 0777)
		if err == nil {
			err = os.Chmod(path, 0777) // 再修改权限
		}
		return err
	}
	return nil
}

func JoinPath(p1, p2 string) string {
	sep := PathSeparators()
	p := p1
	A := string([]byte(p1)[len(p1)-1])
	B := string([]byte(p2)[0])
	if A == sep && B != sep {
		return p1 + p2
	}

	if A != sep {
		p += sep
	}

	if B == sep {
		p += strings.TrimLeft(p2, sep)
	} else {
		p += p2
	}

	A = string([]byte(p1)[len(p1)-1])
	if A != sep {
		p += sep
	}

	return p
}

func MakeFolderMonth(base_path, topfolder string) string {
	folder := fmt.Sprintf(`%s/%s/`, topfolder, time.Now().Format("200601"))
	p := JoinPath(base_path, folder)
	if !PathExist(p) {
		MakeDir(p)
	}

	return folder
}

func PathSeparators() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}

func GetWorkPath(FolderName string) string {
	p, _ := os.Getwd()

	if FolderName != "" {
		return p + PathSeparators() + FolderName + PathSeparators()
	}

	return p + PathSeparators()
}

func JoinFileName(p1, p2 string) string {
	A := string([]byte(p1)[len(p1)-1])
	if A == PathSeparators() {
		return p1 + p2
	} else {
		return p1 + PathSeparators() + p2
	}
}

func GetFileSize(Filename string) int64 {
	fi, err := os.Stat(Filename)
	if err != nil {
		return 0
	}

	return fi.Size()
}

func HasValue[T comparable](val T, list []T) bool {
	if len(list) > 0 {
		for _, p := range list {
			if val == p {
				return true
			}
		}
	}

	return false
}

const (
	CDB_MYSQL  = 0
	CDB_SQLITE = 1
	CDB_MSSQL  = 2
)

type TOptSQL struct {
	Host     string `json:"host,omitempty"`
	UserName string `json:"userName,omitempty"`
	PassWord string `json:"passWord,omitempty"`
	DataBase string `json:"dataBase,omitempty"`
	Charset  string `json:"charset,omitempty"`
	Port     string `json:"port,omitempty"`
}

func (this *TOptSQL) MySQLConnection() string {
	return fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local`, this.UserName, this.PassWord, this.Host, this.Port, this.DataBase, this.Charset)
}
func (this *TOptSQL) MSSQLConnection() string {
	return fmt.Sprintf(`sqlserver://%s:%s@%s:%s?database=%s`, this.UserName, this.PassWord, this.Host, this.Port, this.DataBase)
}

func CreateMySQLDB(opt TOptSQL) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(opt.MySQLConnection()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf(`[MYSQL-连接失败]%s`, err.Error())
	}

	if sqlDB, err := db.DB(); err == nil {

		sqlDB.SetMaxIdleConns(50)           //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
		sqlDB.SetMaxOpenConns(100)          //设置数据库连接池最大连接数
		sqlDB.SetConnMaxLifetime(time.Hour) //设置了连接可复用的最大时间
	} else {
		return nil, err
	}

	return db, nil
}

func CreateSQLiteDB(FileName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(FileName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("[SQLITE-连接失败]%s", err.Error())
	}

	if sqlDB, err := db.DB(); err == nil {

		sqlDB.SetMaxIdleConns(50)           //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
		sqlDB.SetMaxOpenConns(100)          //设置数据库连接池最大连接数
		sqlDB.SetConnMaxLifetime(time.Hour) //设置了连接可复用的最大时间
	} else {
		return nil, err
	}

	return db, nil
}

func CreateMSSQLDB(Opt TOptSQL) (*gorm.DB, error) {
	db, err := gorm.Open(sqlserver.Open(Opt.MSSQLConnection()), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("[MSSQL-连接失败]%s", err.Error())
	}

	if sqlDB, err := db.DB(); err == nil {

		sqlDB.SetMaxIdleConns(50)           //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
		sqlDB.SetMaxOpenConns(100)          //设置数据库连接池最大连接数
		sqlDB.SetConnMaxLifetime(time.Hour) //设置了连接可复用的最大时间
	} else {
		return nil, err
	}

	return db, nil
}

func ToSqlData(rows *sql.Rows) []TSqlData {
	columns, err := rows.Columns()
	if err != nil {
		return nil
	}
	defer rows.Close()

	count := len(columns)
	rc := 0
	tableData := make([]TSqlData, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		rc++
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}

	return tableData
}

/**
 * @Author: icy
 * @Description: 解析数据查询返回结果集 此为带回调的
 * @Param 数据库返回结果集，回调函数
 * @return 错误类型
 * @Date: 12/5/23 10:43 PM
 */

func DoSqlData(rows *sql.Rows, OnData func(result TSqlData)) error {
	if OnData == nil {
		return fmt.Errorf("此方法数据回调不能为空")
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil
	}
	count := len(columns)
	rc := 0
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		rc++
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		OnData(entry)
	}

	return err
}

/**
 * @Author: icy
 * @Description: 复制文件
 * @Param 源文件路径,目标路径
 * @return 错误类型
 * @Date: 12/5/23 10:45 PM
 */

func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

/**
 * @Author: icy
 * @Description: 复制文件夹
 * @Param 目标路径,目标路径
 * @return 错误类型
 * @Date: 12/5/23 10:45 PM
 */

func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

const (
	//通用翻译API HTTP地址：
	CBDAPI_HTTP = `http://api.fanyi.baidu.com/api/trans/vip/translate`
	//通用翻译API HTTPS地址：
	CBDAPI_HTTPS = `https://fanyi-api.baidu.com/api/trans/vip/translate`
)

func BaiduTranslate(appid string, appkey string, fr string, to string, query string) string {
	client := &http.Client{Timeout: 5 * time.Second}

	rand.Seed(int64(time.Now().UnixNano()))
	salt := strconv.Itoa(rand.Intn(32768) + (65536 - 32768))
	sign := MD5(appid + query + salt + appkey)

	payload := url.Values{"appid": {appid}, "q": {query}, "from": {fr}, "to": {to}, "salt": {salt}, "sign": {sign}}
	resp, err := client.Post(CBDAPI_HTTPS,
		"application/x-www-form-urlencoded",
		strings.NewReader(payload.Encode()))

	if err == nil {
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		JO := gjson.ParseBytes(data)
		if JO.Exists() {
			if JO.Get("error_code").Int() > 0 { // 如果存在 这个字段肯定不会是零的咯
				return JO.Get("error_msg").String()
			}

			if JO.Get("trans_result").IsArray() {
				return JO.Get("trans_result").Array()[0].Get("dst").String()
			}
		}
	}

	return ""
}

func CallTableFunction(db *gorm.DB, in any, name string) {
	p := reflect.TypeOf(in)
	if p.Kind() == reflect.Pointer {
		p = p.Elem()
	}

	if p.Kind() != reflect.Struct {
		return
	}

	obj := reflect.ValueOf(in)
	newMethod := obj.MethodByName(name)
	if !newMethod.IsValid() {
		return
	}

	val := make([]reflect.Value, 0)
	val = append(val, reflect.ValueOf(db))
	newMethod.Call(val)
}

func valueIsFloat(val any) bool {
	k := reflect.TypeOf(val)
	switch k.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	}

	return false
}

func valueIsInt(val any) bool {
	k := reflect.TypeOf(val)
	switch k.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	}

	return false
}

func valueIsUInt(val any) bool {
	k := reflect.TypeOf(val)
	switch k.Kind() {
	case reflect.Uint, reflect.Int8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}

	return false
}

/**
 * @Author: icy
 * @Description: 数据查询数据转换到结构体
 * @Param 数据，输出结构
 * @return 错误类型
 * @Date: 12/5/23 10:47 PM
 */

func SqlDataToStruct(data TSqlData, out any) error {
	p := reflect.TypeOf(out)
	if p.Kind() == reflect.Pointer {
		p = p.Elem()
		if p.Kind() != reflect.Struct {
			return fmt.Errorf(`非结构体`)
		}

		size := p.NumField()
		if size == 0 {
			return fmt.Errorf(`这是一个空的结构体`)
		}

		value := reflect.ValueOf(out)
		for k, v := range data {
			val := value.FieldByName(k)
			if !val.IsNil() {
				switch reflect.TypeOf(v).Kind() {
				case reflect.Bool:
					if reflect.TypeOf(val).Kind() == reflect.Bool {
						value.Set(reflect.ValueOf(v))
					}

				case reflect.String:
					if reflect.TypeOf(val).Kind() == reflect.String {
						value.Set(reflect.ValueOf(v))
					}

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if valueIsInt(val) {
						value.Set(reflect.ValueOf(v))
					}

				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if valueIsUInt(val) {
						value.Set(reflect.ValueOf(v))
					}

				case reflect.Float32, reflect.Float64:
					if valueIsFloat(val) {
						value.Set(reflect.ValueOf(v))
					}
				}
			}
		}
	} else {
		return fmt.Errorf(`不正确的结构变体`)
	}

	return nil
}

/*--------------------------------------------------------------------------------------------------------------------*/

const (
	CLOG_LEVEL_INFO  = 1
	CLOG_LEVEL_ERROR = 2
	CLOG_LEVEL_WARN  = 3
	CLOG_LEVEL_DEBUG = 4
)

type TBasic struct {
	db    *gorm.DB
	Print bool
}

func NewBasic(db *gorm.DB) TBasic {
	return TBasic{
		db:    db,
		Print: true,
	}
}

func (this *TBasic) levelString(level int) string {
	switch level {
	case CLOG_LEVEL_ERROR:
		return `error`
	case CLOG_LEVEL_WARN:
		return `warn`
	case CLOG_LEVEL_DEBUG:
		return `debug`
	default:
		return `info`
	}
}

func (this *TBasic) DB(dst ...any) *gorm.DB {
	if len(dst) > 0 {
		p := reflect.TypeOf(dst[0])

		if p == nil {
			return this.db
		}

		if p.Kind() == reflect.String {
			if v, ok := dst[0].(string); ok {
				return this.db.Table(v)
			}
		}

		if p.Kind() == reflect.Pointer {
			p = p.Elem()
		}

		if p.Kind() != reflect.Struct {
			return this.db
		}

		obj := reflect.ValueOf(dst[0])
		Method := obj.MethodByName(`TableName`)
		if !Method.IsValid() {
			return this.db
		}

		val := Method.Call(nil)
		if len(val) > 0 {
			return this.db.Table(val[0].String())
		}
	}

	return this.db
}

func (this *TBasic) Table(TableName string) *gorm.DB {
	return this.db.Table(TableName)
}

func (this *TBasic) logger(level int, title, value string) error {
	s := fmt.Sprintf("[%s] [%s] [%s] %s\n",
		time.Now().Format(`2006-01-02 15:04:05`),
		this.levelString(level),
		title,
		value)
	if this.Print {
		switch level {
		case CLOG_LEVEL_ERROR:
			{
				_, _ = color.New(color.FgRed).Println(s)
				return errors.New(s)
			}
		case CLOG_LEVEL_WARN:
			_, _ = color.New(color.BgYellow).Println(s)
		case CLOG_LEVEL_DEBUG:
			_, _ = color.New(color.FgGreen).Println(s)
		default:
			fmt.Println(s)
		}
	}

	return nil
}

func (this *TBasic) LogInfo(title, value string) {
	_ = this.logger(CLOG_LEVEL_INFO, title, value)
}

func (this *TBasic) LogError(title, value string) error {
	return this.logger(CLOG_LEVEL_ERROR, title, value)
}

func (this *TBasic) LogWarn(title, value string) {
	_ = this.logger(CLOG_LEVEL_WARN, title, value)
}

func (this *TBasic) LogDebug(title, value string) {
	_ = this.logger(CLOG_LEVEL_DEBUG, title, value)
}

/*--------------------------------------------------------------------------------------------------------------------*/
/**
 * @Author: icy
 * @Description: 随机休眠
 * @Param 休眠种子
 * @return 无
 * @Date: 12/5/23 11:06 PM
 */

func RandSleep(seed uint) {
	n := rand.Intn(int(seed))
	if n <= 0 {
		n = 10
	}

	time.Sleep(time.Duration(n) * time.Second)
}

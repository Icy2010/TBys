package TBys

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/**
 * @Author: icy
 * @Description: 解压Gzip数据到字符串
 * @Param 缓冲字节
 * @return 结果字符串 错误类型
 * @Date: 12/5/22 7:31 PM
 */

func UnGZipToString(buff []byte) (string, error) {
	reader := bytes.NewReader(buff)
	body, err := gzip.NewReader(reader)
	if err == nil {
		defer body.Close()
		data, err := ioutil.ReadAll(body)
		return string(data), err
	}

	return "", err
}

/**
 * @Author: icy
 * @Description: 解压GZip 数据 到 字节数组
 * @Param 缓冲字节
 * @return 缓冲字节 错误类型
 * @Date: 12/5/22 7:29 PM
 */

func UnGzipToBuff(buff []byte) ([]byte, error) {
	reader := bytes.NewReader(buff)
	body, err := gzip.NewReader(reader)
	if err == nil {
		defer body.Close()
		data, e := ioutil.ReadAll(body)
		return data, e
	}
	return nil, err
}

/*
Author: icy
Description: 获取代理客户端 不验证安全
Date:  2023/1/26 下午9:10
Param: 无
return: http请求客户端
*/

func GetProxyClient(ProxyUrl string) (*http.Client, error) {
	if proxy, err := url.Parse(ProxyUrl); err == nil {
		tr := &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		return &http.Client{
			Transport: tr,
		}, nil
	} else {
		return nil, err
	}
}

type THttp struct {
	client      *http.Client
	Timeout     time.Duration
	GZip        bool
	headers     map[string]string
	ContentType string

	OnResult func(resp *http.Response) error
}

func NewHTTP(client *http.Client) *THttp {
	return &THttp{
		client: func() *http.Client {
			if client != nil {
				return client
			}

			return &http.Client{}
		}(),
		Timeout: 5000,
		GZip:    false,
		headers: map[string]string{
			"User-Agent":      "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0",
			"Connection":      "keep-alive",
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"Accept-Encoding": "gzip, deflate, br",
			"Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
			"cache-control":   "max-age=0",
			"token":           "",
		},
		ContentType: "application/json",
	}
}

func (this *THttp) HasHeader(field string) (string, bool) {
	for k, _ := range this.headers {
		if strings.EqualFold(k, field) {
			return k, true
		}
	}
	return "", false
}

func (this *THttp) SetHeader(Field, Value string) {

	if key, Has := this.HasHeader(Field); Has {
		this.headers[key] = Value
	} else {
		this.headers[Field] = Value
	}
}

func (this *THttp) GetHeader(Field string) string {
	if key, Has := this.HasHeader(Field); Has {
		return this.headers[key]
	}

	return ``
}

func (this *THttp) newForm(body map[string]interface{}) string {
	Buff := strings.Builder{}

	if len(body) > 0 {
		for k, v := range body {
			switch v.(type) {
			case string:
				Buff.WriteString(fmt.Sprintf(`%s=%s;`, k, v.(string)))
			case int:
				Buff.WriteString(fmt.Sprintf(`%s=%d;`, k, v.(int)))
			case int64:
				Buff.WriteString(fmt.Sprintf(`%s=%d;`, k, v.(int64)))
			case float64:
				Buff.WriteString(fmt.Sprintf(`%s=%f;`, k, v.(float64)))
			case float32:
				Buff.WriteString(fmt.Sprintf(`%s=%f;`, k, v.(float32)))
			}
		}
	}

	return Buff.String()
}

func (this *THttp) defResult() gjson.Result {
	return gjson.Parse(`{"state": -1, "msg":"未知的错误."}`)
}

func (this *THttp) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	res, err := http.NewRequest(method, url, body)
	if len(this.headers) > 0 {
		for k, v := range this.headers {
			if k != "" {
				res.Header.Add(k, v)
			}
		}
	}
	return res, err
}

/**
 * @Author: icy
 * @Description: 所有不带数据提交的请求鼻祖
 * @Param 方法,url GET/DELETE/...
 * @return 结果字节数组,错误类型
 * @Date: 12/5/23 10:38 PM
 */

func (this *THttp) requestR(method, url string) ([]byte, error) {

	client := &http.Client{}

	res, err := this.newRequest(method, url, nil)
	if err == nil {
		var resp *http.Response
		resp, err = client.Do(res)
		if err == nil {
			defer resp.Body.Close()

			if this.OnResult != nil {
				err = this.OnResult(resp)
				return nil, err
			}

			var body []byte
			body, err = io.ReadAll(resp.Body)
			if err == nil {
				if err == nil && this.GZip {
					var data []byte
					data, err = UnGzipToBuff(body)
					if err != nil {
						return body, nil
					}
					return data, err
				} else {
					return body, err
				}
			}
		}
	}

	return nil, err
}

func (this *THttp) Get(url string) ([]byte, error) {
	return this.requestR(http.MethodGet, url)
}

func (this *THttp) Delete(url string) ([]byte, error) {
	return this.requestR(http.MethodDelete, url)
}

/**
 * @Author: icy
 * @Description: 所有带数据提交的请求鼻祖
 * @Param 方法,url,请求数据
 * @return 结果字节数组,错误类型 POST/PUT/....
 * @Date: 12/5/23 10:58 PM
 */

func (this *THttp) requestPP(method, url string, data []byte) ([]byte, error) {
	req, err := this.newRequest(method, url, bytes.NewReader(data))
	if err == nil {
		var resp *http.Response
		resp, err = this.client.Do(req)
		if err == nil {
			defer resp.Body.Close()

			if this.OnResult != nil {
				err = this.OnResult(resp)
				return nil, err
			}

			var body []byte
			body, err = io.ReadAll(resp.Body)
			if err == nil {

				if err == nil && this.GZip {
					var buff []byte
					buff, err = UnGzipToBuff(body)
					return buff, err
				} else {
					return body, err
				}
			}
		}
	}
	return nil, err
}

// post

func (this *THttp) Post(url string, data map[string]any) ([]byte, error) {
	if this.ContentType != "" {
		this.SetHeader("Content-Type", this.ContentType)
	} else {
		this.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	}

	return this.PostBuff(url, []byte(this.newForm(data)))
}

func (this *THttp) PostBuff(url string, data []byte) ([]byte, error) {
	return this.requestPP(http.MethodPost, url, data)
}

// put

func (this *THttp) Put(url string, data map[string]any) ([]byte, error) {
	if this.ContentType != "" {
		this.SetHeader("Content-Type", this.ContentType)
	} else {
		this.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	}

	return this.PutBuff(url, []byte(this.newForm(data)))
}

func (this *THttp) PutBuff(url string, data []byte) ([]byte, error) {
	return this.requestPP(http.MethodPost, url, data)
}

func (this *THttp) BytesToString(data []byte) string {
	return string(data)
}

func (this *THttp) BytesToJSON(data []byte) gjson.Result {
	return gjson.Parse(string(data))
}

func (this *THttp) Upload(url string, FileName string, Form map[string]any) ([]byte, error) {
	file, err := os.Open(FileName)
	if err == nil {
		defer file.Close()
		bodyBuf := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuf)

		if len(Form) > 0 {
			for k, v := range Form {
				switch v.(type) {
				case string:
					_ = bodyWriter.WriteField(k, v.(string))
				case int:
					_ = bodyWriter.WriteField(k, fmt.Sprintf(`%d`, v.(int)))
				case int64:
					_ = bodyWriter.WriteField(k, fmt.Sprintf(`%d`, v.(int64)))
				case float64:
					_ = bodyWriter.WriteField(k, fmt.Sprintf(`%f`, v.(float64)))
				case float32:
					_ = bodyWriter.WriteField(k, fmt.Sprintf(`%f`, v.(float32)))

				}

			}
		}

		//关键的一步操作
		var fileWriter io.Writer
		fileWriter, err = bodyWriter.CreateFormFile("file", filepath.Base(FileName))
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(fileWriter, file)
		if err != nil {
			return nil, err
		}

		this.SetHeader("Content-Type", bodyWriter.FormDataContentType())
		err = bodyWriter.Close() // 不能延迟执行哦 傻逼的哦 会导致未关闭 文件发送失败

		return this.PostBuff(url, bodyBuf.Bytes())
	}

	return nil, err
}

func (this *THttp) Header() map[string]string {
	return this.headers
}

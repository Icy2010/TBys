package TBys

import (
	"errors"
	htmlquery "github.com/antchfx/xquery/html"
	"github.com/dop251/goja"
	"gorm.io/gorm"
	"strings"
)

type TBysBasic struct {
	TBasic
}

func NewBysBasic(db *gorm.DB) TBysBasic {
	return TBysBasic{NewBasic(db)}
}

/**
 * @Author: icy
 * @Description: 获取HTML内容中的 JS脚本内容
 * @Param 网页源码 , goja的JS运行时，脚本源的回调
 * @return 字符串结果，是否成功
 * @Date: 12/5/23 11:29 PM
 */

func (this *TBysBasic) GetHTMLScript(content string, js *goja.Runtime, onScript func(value string, size uint) (string, bool)) (string, bool) {
	if onScript == nil {
		return "", false
	}

	root, _ := htmlquery.Parse(strings.NewReader(content))
	list := htmlquery.Find(root, `//script`)
	if len(list) > 0 {
		for _, v := range list {
			text := htmlquery.InnerText(v)
			if value, next := onScript(text, uint(len(text))); next {
				if val, err := js.RunString(value); err == nil {
					return val.String(), true
				}
			}
		}
	}

	return "", false
}

/**
 * @Author: icy
 * @Description: 获取HTML源码中的任何内容
 * @Param 网页源码 , XPATH，输出回调
 * @return 错误类型
 * @Date: 12/5/23 11:33 PM
 */

func (this *TBysBasic) GetHTMLContent(src, expr string, OnText func(val string) bool) error {
	if OnText != nil {
		return errors.New(`请传入回调函数`)
	}

	root, _ := htmlquery.Parse(strings.NewReader(src))
	list := htmlquery.Find(root, expr)
	if len(list) > 0 {
		for _, v := range list {
			text := htmlquery.InnerText(v)
			if OnText(text) {
				return nil
			}
		}
	} else {
		return errors.New(`未查找到内容`)
	}

	return nil
}

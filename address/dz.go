package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"regexp"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pupuk/addr"
)

// 判断一个字符串是否在一个数组中
func in(target string, str_array []string) bool {
	for _, element := range str_array {
		if strings.Contains(element, target) {
			return true
		}
	}
	return false
}

func check(each_array []string, array []string) bool {
	for _, element := range each_array {
		if in(element, array) {
			return true
		}
	}
	return false
}

func checkAddress(inputAdd string) string {
	parse := addr.Smart(inputAdd)
	var gvCity = []string{"北京市", "上海市", "天津市", "重庆市"}
	var ret string
	var bt bytes.Buffer

	// 输出解析结果
	fmt.Println(parse.Province) // 广东省
	fmt.Println(parse.City)     // 深圳市
	fmt.Println(parse.Region)   // 龙华区
	fmt.Println(parse.Street)   // 龙华街道1980科技文化产业园3栋317
	fmt.Println(parse.Address)  // 深圳市龙华区龙华街道1980科技文化产业园3栋317

	// 匹配所有省级
	pReg := regexp.MustCompile(`.+?(省|市|自治区|特别行政区|区)`)
	pArr := pReg.FindAllString(inputAdd, -1)

	// 匹配所有市级
	// 由于该匹配可能会遗漏部分，所以合并省级匹配
	cReg := regexp.MustCompile(`.+?(省|市|自治州|州|地区|盟|县|自治县|区|林区)`)
	cArr := append(cReg.FindAllString(inputAdd, -1), pArr...)

	// 匹配所有区县级
	// 由于该匹配可能会遗漏部分(如：东乡区)所以合并市级匹配
	rReg := regexp.MustCompile(`.+?(市|县|自治县|旗|自治旗|区|林区|特区|街道|镇|乡)`)
	rArr := append(rReg.FindAllString(inputAdd, -1), cArr...)

	if parse.Province == "" && parse.City == "" && parse.Region == "" && parse.Street == "" {
		return "输入地址有误！请输入详细地址"
	}

	if !strings.Contains(inputAdd, parse.Province) {
		if in("省", rArr) {
			bt.WriteString("地址输入可能有误！请核对是否为 " + parse.Province + "\n")
		} else if !check(gvCity, rArr) {
			bt.WriteString("缺少省级信息！\n")
		}
	}

	if !strings.Contains(inputAdd, parse.City) {
		if in("市", rArr) {
			bt.WriteString("地址输入可能有误！请核对是否为 " + parse.City + "\n")
		} else {
			bt.WriteString("缺少市级信息！\n")
		}
	}

	if !strings.Contains(inputAdd, parse.Region) {
		if in("县", rArr) {
			bt.WriteString("地址输入可能有误！请核对是否为 " + parse.Region + "\n")
		} else {
			bt.WriteString("缺少县级信息！\n")
		}
	}
	ret = bt.String()
	//fmt.Println(ret)
	return ret
}

type Site struct {
	Msg string
}

func main() {

	r := gin.Default()
	r.Static("xxx", "./static")
	r.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.LoadHTMLFiles("form.tmpl")

	r.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/posts/index", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(http.StatusOK, "get/index.tmpl", nil)
	})

	r.GET("/address", func(c *gin.Context) {
		//获取前端传过来的数据
		var site Site
		var data string
		err := c.ShouldBind(&site)
		if err != nil {
			panic(err)
		}
		data = site.Msg
		sta := checkAddress(data) //地址状态信息
		fmt.Println(sta)
		c.JSON(http.StatusOK, gin.H{
			"Msg": sta,
		})
	})

	r.Run(":9090")
}

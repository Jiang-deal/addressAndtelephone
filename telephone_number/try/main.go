/*对一个电话号码进行判断，判断是否为空号，使用gin框架，
接收前端传过来的电话号码，首先使用正则表达式对号码进行一个判断，
在调用一个接口判断号码状态*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type mes struct {
	Code   int
	Msg    string
	Result []mes2
}

type mes2 struct {
	Mobile string
	State  int
}

//号码格式验证
func VerifyMobileFormat(mobileNum string) bool {
	mobile := "^1[3|4|5|7|8][0-9]{9}$"
	zuoji := "^[0][1-9]{2,3}-[0-9]{5,10}$"
	len := len([]byte(mobileNum)) //获取电话的长度
	var reg *regexp.Regexp
	if len == 11 {
		reg = regexp.MustCompile(mobile)

	}
	if len == 7 || len == 8 {
		reg = regexp.MustCompile(zuoji)
	}
	return reg.MatchString(mobileNum)
}

//结构体
type Telephone struct {
	Number string
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

func main() {
	// var telephonenumber string
	// telephonenumber = getInput()

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
	//r.LoadHTMLGlob("s")

	r.GET("/posts/index", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(http.StatusOK, "get/index.tmpl", nil)
	})

	r.GET("/telephonenumber", func(c *gin.Context) {

		var number Telephone
		err1 := c.ShouldBind(&number) //从前端获取数据，并把值赋给number
		data := number.Number
		ri := VerifyMobileFormat(data)
		fmt.Println("电话号码格式正确与否", ri)
		//如果电话格式正确就继续判断是否为空号
		if ri == true {
			if err1 == nil {
				s1 := "http://139.196.108.241:8081/Api/Detect.ashx?"
				s2 := "account=17318572578&pswd=19981013415Xx&mobile="
				var bt bytes.Buffer
				bt.WriteString(s1)
				bt.WriteString(s2)
				bt.WriteString(data)
				s3 := bt.String()
				fmt.Println(s3)
				response, err := http.Get(s3)
				if err != nil {
					fmt.Print("网站访问错误")
				}
				defer response.Body.Close()

				body, err := ioutil.ReadAll(response.Body)
				if err == nil {
					var data_re = string(body)
					fmt.Println(data_re) //返回的书json数据.我们要把这个数据中的state取出来
					//这里是获得json中state的value
					var data mes
					json.Unmarshal([]byte(data_re), &data)

					fmt.Println("号码状态", data.Result[0].State) //-1 风险 0 空号 1 实号

					c.JSON(http.StatusOK, gin.H{
						"State": data.Result[0].State,
					})
				}

			} else {
				panic("接口出错")
			}

		} else {
			c.JSON(http.StatusOK, gin.H{
				"State": 2,
			})
		}

	})

	r.Run()

}

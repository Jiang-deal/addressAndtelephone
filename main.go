/*对一个电话号码进行判断，判断是否为空号，使用gin框架，
接收前端传过来的电话号码，首先使用正则表达式对号码进行一个判断，
在调用一个接口判断号码状态*/

/*
保存电话号码和号码状态在数据库，并在前端把这些数据显示出来
1.创建一个数据库 telephone_project 库中建表 make_telephone 表中包含字段 id time telephone state
2.在向前端返回用户状态时同时向数据库中保存这条信息
3.输入电话号码之后，先判断格式，在从数据库中查询在七天之内是否存在这个电话信息，如果存在电话信息把七天内的数据显示在前端，如果没有在调用接口去查询，并把查询信息保存在数据库

*/

///
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

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

// 数据库绑定电话信息
type MakeTelephone struct {
	ID        int
	Time      string
	Telephone string
	State     string
}

// 数据库绑定电话信息
type display_data struct {
	Time      string
	Telephone string
	State     string
}

//号码格式验证
func VerifyMobileFormat(mobileNum string) bool {
	mobile := "^1[3|4|5|7|8][0-9]{9}$"
	zuoji := "^[0][1-9]{2,3}-[0-9]{5,10}$"
	len := len([]byte(mobileNum)) //获取电话的长度
	var reg *regexp.Regexp
	if len == 11 {
		reg = regexp.MustCompile(mobile)
		return reg.MatchString(mobileNum)

	} else if len == 7 || len == 8 {
		reg = regexp.MustCompile(zuoji)
		return reg.MatchString(mobileNum)
	} else {
		return false
	}

}

//连接数据库创建表
func CreateDb(id int, tel string, sta int, t string) (bool, string) {
	db, err := gorm.Open("mysql", "root:123456@tcp(localhost:3306)/telephone_project?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println(id)

	// 自动迁移
	db.AutoMigrate(&MakeTelephone{})

	var u1 MakeTelephone
	if sta == 0 {
		u1 = MakeTelephone{Time: t, Telephone: tel, State: "空号"}
		db.Create(&u1)
		fmt.Println("数据库数据添加成功")
		return true, "空号"
	} else if sta == -1 {
		u1 = MakeTelephone{Time: t, Telephone: tel, State: "危险号"}
		db.Create(&u1)
		fmt.Println("数据库数据添加成功")
		return true, "危险号"
	} else if sta == 1 {
		u1 = MakeTelephone{Time: t, Telephone: tel, State: "实号"}
		db.Create(&u1)
		fmt.Println("数据库数据添加成功")
		return true, "实号"
	} else {
		fmt.Println("数据库数据添加失败")
		return false, "失败"
	}

}

//根据电话号码查询数据库数据,如果找到了，就返回true否则为false
func FindTe(tele string, data *[]MakeTelephone, currentTime, oldTime string) bool {
	db1, err := gorm.Open("mysql", "root:123456@tcp(localhost:3306)/telephone_project?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db1.Close()
	var results []MakeTelephone
	// db1.Find(&results) //查询所有的记录
	// fmt.Println("Find的值：", results)
	// fmt.Println(results[0].ID)
	// fmt.Println(len(results))
	db1.Where("Time BETWEEN ? AND ? AND Telephone=?", oldTime, currentTime, tele).Find(&results)
	//db1.Where("Telephone=?", tele).Debug().First(&results)
	if len(results) == 0 {
		return false
	} else {
		for _, value := range results {
			fmt.Println(value.ID)
		}
		*data = results
		return true
	}

}

//结构体
type Telephone struct {
	Number string
}

//创建结构体的目的：方便向前端传数据
type Redata struct {
	Telephone string
	State     string
	Source    string
	Time      string
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

	id := 1

	r := gin.Default()
	r.Static("xxx", "./static")
	r.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.LoadHTMLFiles("form2.tmpl")
	//r.LoadHTMLGlob("s")

	r.GET("/posts/index", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(http.StatusOK, "get/index.tmpl", nil)
	})
	r.GET("/database", func(c *gin.Context) {
		//做一个数据库的查询
		// root:root@tcp(localhost:3306)
		db1, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/telephone_project?charset=utf8mb4&parseTime=True&loc=Local")
		if err != nil {
			panic(err)
		}
		defer db1.Close()
		var results []MakeTelephone
		db1.Find(&results) //查询所有的记录
		fmt.Println("Find的值：", results)
		fmt.Println(results[0].ID)
		fmt.Println(len(results))

		var dis_datas display_data
		for i := 0; i < len(results); i++ {

		}

		//封装一下数据库的数据成json
		// c.JSON：返回JSON格式的数据
		c.JSON(http.StatusOK, dis_datas)
	})

	r.GET("/telephonenumber", func(c *gin.Context) {

		var number Telephone
		err1 := c.ShouldBind(&number) //从前端获取数据，并把值赋给number
		data := number.Number
		fmt.Println("电话号码", data)

		ri := VerifyMobileFormat(data)
		cli_tele := data
		fmt.Println("电话号码格式正确与否", ri)
		//如果电话格式正确就继续判断是否为空号
		if ri == true {
			//根据电话号码获取数据库中的数据，若存在则显示在前端
			if err1 == nil {
				//获取当前时间
				t1 := time.Now()
				//var u1 string
				//u1 = t1.Format("20060102150405")
				u1_1 := t1.Format("2006-01-02 15:04:05")

				t2 := t1.AddDate(0, 0, -7) //获取的是前七天的时间
				//var u2 string
				//u2 = t2.Format("20060102150405")
				u2_1 := t2.Format("2006-01-02 15:04:05")

				db_data := make([]MakeTelephone, 0, 10)

				dbtele := FindTe(data, &db_data, u1_1, u2_1) //若存在数据则返回true
				if !dbtele {
					s1 := "http://139.196.108.241:8081/Api/Detect.ashx?"
					s2 := "account=15797834583&pswd=Zhaohan123&mobile="
					var bt bytes.Buffer
					bt.WriteString(s1)
					bt.WriteString(s2)
					bt.WriteString(data)
					s3 := bt.String()
					fmt.Println(s3)

					//记录访问网站的时间
					t := time.Now()
					var k string
					k = t.Format("2006-01-02 15:04:05")

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

						//向数据库保存电话信息以及电话状态

						tele_state := data.Result[0].State
						res, s := CreateDb(id, cli_tele, tele_state, k)
						if res == true {
							fmt.Println("号码信息添加数据库成功")
							id++
						}
						//返回前端只是返回一个，所以就没必要用结构体的方法
						u := make([]Redata, 0, 10)
						var temp Redata

						temp.Telephone = cli_tele
						temp.State = s
						temp.Source = "接口"
						temp.Time = k
						u = append(u, temp)
						c.JSON(http.StatusOK, u)
					}
				} else { //数据库中的数据返回并返回给前端
					fmt.Println("数据库中存在这条数据")
					//数据库中添加数据
					var i int
					if db_data[len(db_data)-1].State == "空号" {
						i = 0
					}
					if db_data[len(db_data)-1].State == "危险号" {
						i = -1
					}
					if db_data[len(db_data)-1].State == "实号" {
						i = 1
					}
					res, _ := CreateDb(id, cli_tele, i, u1_1)
					if res == true {
						fmt.Println("号码信息添加数据库成功")
						id++
					}

					//声明一个结构化数组，作为返回数据
					s := make([]Redata, 0, 10)
					for t := 0; t < len(db_data); t++ {
						var temp Redata

						temp.Telephone = db_data[t].Telephone
						temp.State = db_data[t].State
						temp.Source = "数据库"
						temp.Time = db_data[t].Time

						s = append(s, temp)

					}
					//fmt.Println(s)
					//r, _ := json.Marshal(s) //结构体转为json
					//fmt.Println(r)

					fmt.Println(db_data[len(db_data)-1].State) //返回最近的一条记录
					fmt.Println("近七天的查询次数", len(db_data)+1)
					c.JSON(http.StatusOK, s)
				}

			} else {
				panic("从前端获取数据失败")
			}

		} else {
			t := time.Now()
			var k string
			k = t.Format("2006-01-02 15:04:05")
			u := make([]Redata, 0, 10)
			var temp Redata

			temp.Telephone = cli_tele
			temp.State = "格式错误"
			temp.Source = "前端"
			temp.Time = k
			u = append(u, temp)
			c.JSON(http.StatusOK, u)
		}

	})

	r.Run()

}

package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	var str = `{"code":true,"msg":"OK","result":[
        {
            "mobile":"13657087926",
            "state":0
        }
    ]}`

	fmt.Println(str)

	var data mes
	json.Unmarshal([]byte(str), &data)

	fmt.Println(data)
	fmt.Println(data.Code)
	fmt.Println(data.Msg)
	fmt.Println(data.Result[0].Mobile)

}

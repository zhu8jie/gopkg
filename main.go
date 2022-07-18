package main

import (
	"fmt"

	"github.com/zhu8jie/gopkg/ruleengine"
)

/**
 * @Description
 * @Author weiyanwei
 * @Date 2022/7/11 21:02
 **/

func main() {

	ruleStr := `!(a == 1 && b == 2 && c == "test" && d == false)`

	// 匹配变量
	params := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": "test",
		"d": true,
	}

	result, err := ruleengine.Match(ruleStr, params)

	fmt.Println(result, err)
}

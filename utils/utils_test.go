package utils

import (
	"fmt"
	"testing"
	"time"
)

type Person struct {
	Name     string
	Age      int
	Active   bool
	Flag     bool
	Num      float64
	TimeData time.Time
	S        []string
}

func TestSliceToStructs(t *testing.T) {
	// 准备测试数据
	data := [][]any{
		{"张三", "25", true, "on", 18.5, "2025-01-20T15:06:25.85", []int{1, 2, 3}},
		{"李四", 30, false, "true", "19.5", "2022-01-02 00:00:00", []string{"d", "e", "f"}},
	}
	fmt.Println(data)
	// 转换数据
	persons, err := SliceToStructs[Person](data)
	if err != nil {
		fmt.Println("错误:", err)
		return
	}

	// 打印结果
	for _, p := range persons {
		fmt.Printf("%+v\n", p)
	}
}

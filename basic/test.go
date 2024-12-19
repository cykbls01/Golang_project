package main

import (
	"fmt"
	"reflect"
	"strings"
)

// FilterStructsByFieldName 函数接受一个切片和一个字段名，返回一个新的切片，其中只包含字段值包含指定子字符串的元素
func FilterStructsByFieldName(slice interface{}, fieldName, subStr string) interface{} {
	// 获取切片的反射类型和值
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		panic("input is not a slice")
	}

	// 创建一个新的切片来存储结果
	resultSlice := reflect.MakeSlice(sliceValue.Type(), 0, sliceValue.Len())

	// 遍历原始切片
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i)
		if elem.Kind() == reflect.Struct {
			// 获取结构体字段的反射值
			fieldValue := elem.FieldByName(fieldName)
			if fieldValue.IsValid() && fieldValue.Kind() == reflect.String && strings.Contains(fieldValue.String(), subStr) {
				// 如果字段值包含指定子字符串，则将该元素添加到结果切片中
				resultSlice = reflect.Append(resultSlice, elem)
			}
		}
	}

	return resultSlice.Interface()
}

type Person struct {
	Name string
	Age  int
}

func main() {
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bobccekjy", Age: 25},
		{Name: "Charlie", Age: 35},
		{Name: "DavidccekjySmith", Age: 40},
	}

	// 调用 FilterStructsByFieldName 函数，注意这里返回的是一个 interface{} 类型，需要类型断言转换回 []Person
	filteredPeople, ok := FilterStructsByFieldName(people, "Name", "ccekjy").([]Person)
	if !ok {
		panic("type assertion failed")
	}

	// 打印筛选后的结果
	for _, person := range filteredPeople {
		fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
	}
}

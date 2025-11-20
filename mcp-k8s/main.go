package main

import (
	"fmt"
	"mcp-k8s/Util"
)

func main() {
	//Util.Pre()
	//Util.InitDB()
	//task.UpdateServices(os.Args[1])
	//Util.Post()

	type Address struct {
		City   string // 普通字段
		Zip    int    // 非字符串字段
		detail string // 未导出字段（首字母小写）
	}
	type User struct {
		Name     string   // 直接字段
		Age      int      // 数值类型
		Addr     Address  // 嵌套结构体
		Phone    *string  // 指针字段
		EmptyPtr *Address // 可能为nil的指针
	}

	// 测试数据准备
	phone1 := "13800138000"
	phone2 := "13900139000"
	addr1 := Address{City: "Beijing", Zip: 100000, detail: "xxx街道"} // detail未导出
	addr2 := Address{City: "Shanghai", Zip: 200000, detail: "yyy街道"}
	userSlice := []User{
		{Name: "Alice", Age: 25, Addr: addr1, Phone: &phone1, EmptyPtr: &addr1},
		{Name: "Bob", Age: 30, Addr: addr2, Phone: &phone2, EmptyPtr: nil},        // EmptyPtr为nil
		{Name: "Charlie", Age: 35, Addr: Address{}, Phone: nil, EmptyPtr: &addr2}, // Phone为nil，Addr字段为空
	}

	// 测试1：直接属性（Name）
	fmt.Println("测试1（直接属性Name）：")
	fmt.Println(Util.PluckAndJoin(userSlice, "Name", ",")) // 输出：Alice,Bob,Charlie

	// 测试2：数值类型属性（Age）
	fmt.Println("\n测试2（数值属性Age）：")
	fmt.Println(Util.PluckAndJoin(userSlice, "Age", "|")) // 输出：25|30|35

	// 测试3：嵌套属性（Addr.City）
	fmt.Println("\n测试3（嵌套属性Addr.City）：")
	fmt.Println(Util.PluckAndJoin(userSlice, "Addr.City", "-")) // 输出：Beijing-Shanghai-

	// 测试4：嵌套数值属性（Addr.Zip）
	fmt.Println("\n测试4（嵌套数值属性Addr.Zip）：")
	fmt.Println(Util.PluckAndJoin(userSlice, "Addr.Zip", " ")) // 输出：100000 200000 0

	// 测试5：指针字段（Phone）
	fmt.Println("\n测试5（指针字段Phone）：")
	fmt.Println(Util.PluckAndJoin(userSlice, "Phone", ",")) // 输出：13800138000,13900139000,

	// 测试6：嵌套指针字段（EmptyPtr.City）
	fmt.Println("\n测试6（嵌套指针字段EmptyPtr.City）：")
	fmt.Println(Util.PluckAndJoin(userSlice, "EmptyPtr.City", "|")) // 输出：Beijing||Shanghai

	// 测试7：不存在的属性（NoExist）
	fmt.Println("\n测试7（不存在的属性NoExist）：")
	fmt.Println(Util.PluckAndJoin(userSlice, "NoExist", ",")) // 输出：,,

	// 测试8：未导出字段（Addr.detail）
	fmt.Println("\n测试8（未导出字段Addr.detail）：")
	fmt.Println(Util.PluckAndJoin(userSlice, "Addr.detail", ",")) // 输出：,,
}

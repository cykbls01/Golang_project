package Util

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func SaveFromFunc[T any](path string, f func(string, string) []T) error {
	data := f(path, "")
	if len(data) == 0 {
		fmt.Println("⚠️  没有需要保存的数据")
		return nil
	}

	var saveErrs []error
	for idx, item := range data {
		if err := DB.Save(&item).Error; err != nil {
			saveErrs = append(saveErrs, fmt.Errorf("第%d条数据保存失败: %w", idx+1, err))
			continue
		}
		fmt.Printf("✅ 成功保存数据: %+v\n", item)
	}

	if len(saveErrs) > 0 {
		return fmt.Errorf("总共有%d条数据保存失败: %w", len(saveErrs), errors.Join(saveErrs...))
	}
	return nil
}

func PluckAndJoin[T any](arr []T, fieldName string, separator ...string) string {
	// 处理分隔符默认值
	sep := ""
	if len(separator) > 0 {
		sep = separator[0]
	}

	// 1. 检查数组是否为空
	if len(arr) == 0 {
		return "" // 空数组返回空字符串，无错误
	}

	// 2. 验证数组元素是结构体（避免传入 []int 等非结构体数组）
	elemType := reflect.TypeOf(arr).Elem()
	if elemType.Kind() != reflect.Struct {
		fmt.Errorf("数组元素必须是结构体，当前是 %s", elemType.Kind())
		return ""
	}

	// 3. 检查属性是否存在且可访问（导出字段：首字母大写）
	field, ok := elemType.FieldByName(fieldName)
	if !ok {
		fmt.Errorf("结构体 %s 不存在属性 %s", elemType.Name(), fieldName)
		return ""
	}
	if field.PkgPath != "" {
		fmt.Errorf("属性 %s 未导出（首字母需大写）", fieldName) // 非导出字段（首字母小写）会有 PkgPath
		return ""
	}

	// 4. 反射遍历数组，提取属性值并转为字符串
	var builder strings.Builder
	for i, item := range arr {
		// 获取当前元素的反射值
		itemValue := reflect.ValueOf(item)
		// 获取指定属性的值
		fieldValue := itemValue.FieldByName(fieldName)

		// 5. 将属性值转为字符串（兼容大部分类型：string、int、bool、time.Time 等）
		strVal, err := fieldValueToString(fieldValue)
		if err != nil {
			fmt.Errorf("元素索引 %d 的属性 %s 转换字符串失败：%v", i, fieldName, err)
			return ""
		}

		// 6. 拼接（跳过最后一个元素的分隔符）
		if i > 0 {
			builder.WriteString(sep)
		}
		builder.WriteString(strVal)
	}

	return builder.String()
}

func PluckAndJoinNested[T any](slice []T, fieldPath string, sep string) string {
	// 处理空输入
	if len(slice) == 0 || fieldPath == "" {
		return ""
	}

	// 解析属性路径（按"."拆分，支持嵌套）
	pathParts := strings.Split(fieldPath, ".")
	if len(pathParts) == 0 {
		return ""
	}

	var result []string
	// 遍历切片中每个元素
	for _, elem := range slice {
		// 提取当前元素的目标属性值（转为字符串）
		valStr := getNestedFieldValue(elem, pathParts)
		result = append(result, valStr)
	}

	// 用分隔符拼接所有属性值
	return strings.Join(result, sep)
}

// getNestedFieldValue 反射获取元素的嵌套属性值，转为字符串
// 支持处理：指针类型（*T、**T等）、嵌套结构体、导出字段校验
func getNestedFieldValue[T any](elem T, pathParts []string) string {
	// 1. 初始化反射值（处理元素本身是指针的情况）
	val := reflect.ValueOf(elem)

	// 循环解析指针（支持多重指针，如***Inner）
	for val.Kind() == reflect.Ptr {
		if val.IsNil() { // 指针为nil，无法获取属性
			return ""
		}
		val = val.Elem()
	}

	// 2. 逐层遍历属性路径
	for _, fieldName := range pathParts {
		// 当前值必须是结构体才能继续获取字段
		if val.Kind() != reflect.Struct {
			return ""
		}

		// 获取当前层字段（必须是导出字段，首字母大写）
		field := val.FieldByName(fieldName)
		if !field.IsValid() { // 字段不存在
			return ""
		}
		if !field.CanInterface() { // 字段未导出（首字母小写），反射无法访问
			return ""
		}

		// 更新当前值为字段值，继续下一层路径
		val = field

		// 处理字段值是指针的情况（如Name *Inner）
		for val.Kind() == reflect.Ptr {
			if val.IsNil() { // 字段指针为nil
				return ""
			}
			val = val.Elem()
		}
	}

	// 3. 将最终属性值转为字符串（支持任意类型：int、string、bool等）
	return fmt.Sprintf("%v", val.Interface())
}

// fieldValueToString 将反射值转为字符串（兼容常见类型）
func fieldValueToString(v reflect.Value) (string, error) {
	// 处理指针类型（如 *string、*int）
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", nil // 空指针返回空字符串
		}
		v = v.Elem() // 解指针
	}

	// 根据类型转换为字符串
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", v.Float()), nil
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool()), nil
	case reflect.Struct:
		// 特殊处理 time.Time 类型
		if t, ok := v.Interface().(time.Time); ok {
			return t.Format(time.RFC3339), nil
		}
		// 其他结构体默认调用 String() 方法（若实现）
		if strer, ok := v.Interface().(fmt.Stringer); ok {
			return strer.String(), nil
		}
		// 无 String() 方法则返回结构体字符串形式
		return fmt.Sprintf("%+v", v.Interface()), nil
	default:
		return "", fmt.Errorf("不支持的属性类型：%s", v.Kind())
	}
}

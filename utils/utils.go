package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/wonderivan/logger"
)

// RemoveDefaultValues 区域制定切片中的默认值
func RemoveDefaultValues(slice interface{}) interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("slice is not a slice")
	}
	result := reflect.MakeSlice(s.Type(), 0, s.Len())
	defaultValue := reflect.Zero(s.Type().Elem()).Interface()

	for i := 0; i < s.Len(); i++ {
		value := s.Index(i).Interface()
		if reflect.DeepEqual(value, defaultValue) {
			continue
		}
		result = reflect.Append(result, reflect.ValueOf(value))
	}

	return result.Interface()
}

// String 其他类型转string
func String(s interface{}) string {
	return fmt.Sprintf("%v", s)
}

// ImageBase64ToFile base64编码转图片
func ImageBase64ToFile(baseStr string, fileName string) error {
	parts := strings.Split(baseStr, ",")
	if len(parts) != 2 {
		return fmt.Errorf("base64编码异常")
	}
	// 解码 base64 字符串
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("base64解码失败:[%#v]", err)
	}

	// 将解码后的数据保存到文件
	err = ioutil.WriteFile(fileName, data, 0777)
	if err != nil {
		return fmt.Errorf("保存文件失败:[%#v]", err)
	}
	return nil
}

// Float64ToString float64 转 string保留一位小数
func Float64ToString(s interface{}) string {
	return fmt.Sprintf("%.1f", s)
}

// Int string 转 int64 慎用, 注意报错会返回0
func Int(s string) int64 {
	// 不判断错误, 错误时num=0
	num, _ := strconv.ParseInt(s, 10, 64)
	return num
}

// Float string 转 float64 慎用, 注意报错会返回0
func Float(s string) float64 {
	// 不判断错误, 错误时num=0
	num, _ := strconv.ParseFloat(s, 64)
	return num
}

// UInt string 转 uint 慎用, 注意报错会返回0
func UInt(s string) uint {
	// 不判断错误, 错误时num=0
	num, _ := strconv.ParseInt(s, 10, 64)
	return uint(num)
}

// Int0 string 转 int 慎用, 注意报错会返回0
func Int0(s string) int {
	// 不判断错误, 错误时num=0
	num, _ := strconv.ParseInt(s, 10, 64)
	return int(num)
}

// ParamEmptyCheck 判断结构体指定字段是否为空
func ParamEmptyCheck(cantBeEmpty []string, paramStruct interface{}) error {
	t := reflect.TypeOf(paramStruct)
	v := reflect.ValueOf(paramStruct)
	for i := 0; i < t.NumField(); i++ {
		// fmt.Println(t.SearchField(i).CarriageNames)
		for _, field := range cantBeEmpty {
			structFieldName := t.Field(i).Name
			if field == structFieldName {
				switch v.Field(i).Interface().(type) {
				case uint:
					if v.Field(i).Uint() == 0 {
						return errors.New("filed " + structFieldName + " can't be empty")
					}
				case int:
					if v.Field(i).Int() == 0 {
						return errors.New("filed " + structFieldName + " can't be empty")
					}
				case string:
					if v.Field(i).String() == "" {
						return errors.New("filed " + structFieldName + " can't be empty")
					}
				default:
					return nil
				}

			}
		}
	}
	return nil
}

func In[T int | int64 | int32 | string | uint](element T, list []T) bool {
	for _, e := range list {
		if element == e {
			return true
		}
	}
	return false
}

func IntInt32(element int32, list []int32) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}

func UIntIn(element uint, list []uint) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}
func StringIn(element string, list []string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}

func SafeSend(ch chan int, value int) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()
	ch <- value
	return false
}

// Minus 求两个数组的差集合 SA-SB
func Minus(SA, SB []string) []string {
	var SC []string
	saMap := make(map[string]string)
	for _, sa := range SA {
		saMap[sa] = sa
	}
	for _, sb := range SB {
		if _, ok := saMap[sb]; ok {
			delete(saMap, sb)
		}
	}
	for _, value := range saMap {
		SC = append(SC, value)
	}
	return SC
}

// StringInter 求两个字符串数组的交集
func StringInter(SA, SB []string) []string {
	var SC []string
	saMap := make(map[string]string)
	for _, sa := range SA {
		saMap[sa] = sa
	}
	for _, sb := range SB {
		if _, ok := saMap[sb]; ok {
			SC = append(SC, sb)
		}
	}
	return SC
}

// StringUnion 求两个字符串数组的并集
func StringUnion(SA, SB []string) []string {
	saMap := make(map[string]string)
	for _, sa := range SA {
		saMap[sa] = sa
	}
	for _, sb := range SB {
		if _, ok := saMap[sb]; ok {
			SA = append(SA, sb)
		}
	}
	return SA
}

// StringUnion 求两个字符串数组的差集
func StringDifferenceA(SA, SB []string) []string {
	var SC []string
	sbMap := make(map[string]int)
	for _, sb := range SB {
		sbMap[sb]++
	}
	for _, sa := range SA {
		times := sbMap[sa]
		if times == 0 {
			SC = append(SC, sa)
		}
	}
	return SC
}

func Deduplicate(SA []string) []string {
	var SC []string
	saMap := make(map[string]string)
	for _, sa := range SA {
		if _, ok := saMap[sa]; !ok {
			saMap[sa] = sa
		}
	}
	for _, value := range saMap {
		SC = append(SC, value)
	}
	return SC
}

func DeduplicateInt(SA []int) []int {
	var SC []int
	saMap := make(map[int]int)
	for _, sa := range SA {
		if _, ok := saMap[sa]; !ok {
			saMap[sa] = sa
		}
	}
	for _, value := range saMap {
		SC = append(SC, value)
	}
	return SC
}

// TimeCost @brief：耗时统计函数
func TimeCost() func(s string) {
	start := time.Now()
	return func(s string) {
		tc := time.Since(start)
		logger.Info("%s:%v", s, tc)
	}
}

// CaseToCamel 下划线转驼峰
func CaseToCamel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

func SaveROIFile(filePath string, roi *[]float64) error {
	var err error
	data, err := json.MarshalIndent(roi, "", "")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, data, 0777)
	// c := exec.Command("python3", "listToPly.py", filePath)
	// output, err := c.CombinedOutput()
	//fmt.Println(string(output))
	return err
}

func ReadROIFile(filePath string, roi *[]float64) error {
	var err error
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, roi)
	return err
}

func EncodeOptions(options *[]string) string {
	optionsStr := ""
	for _, option := range *options {
		optionsStr = optionsStr + "--" + option
	}
	// fmt.Println("已编码: " + optionsStr)
	return base64.StdEncoding.EncodeToString([]byte(optionsStr))
}

func DecodeOptions(optionsEncoded string) []string {
	decodedByte, _ := base64.StdEncoding.DecodeString(optionsEncoded)
	// fmt.Printf("已解码: " + string(decodedByte))
	return strings.Split(string(decodedByte), "--")
}

// LargeLetterIncrease 字母递增
func LargeLetterIncrease(s string) (error, string) {
	i := []rune(s)
	j := i[0]
	if j < 65 || j > 90 {
		return fmt.Errorf("传入的参数不是大写字母"), s
	}
	for index := range i {
		i[index] = i[index] + 1
	}
	return nil, string(i)
}

// Float64ArrToString float64列表转化为字符串
func Float64ArrToString(arr []float64) string {
	res := "["
	for i := 0; i < len(arr); i++ {
		res += strconv.FormatFloat(arr[i], 'f', 3, 32)
		if i != len(arr)-1 {
			res += ","
		}
	}
	return res + "]"
}

// GetDifference 比较对象修改前后的不同
func GetDifference(old, new interface{}) string {
	operationLog := ""
	var typeInfo1 = reflect.TypeOf(old)
	var valInfo1 = reflect.ValueOf(old)
	var valInfo2 = reflect.ValueOf(new)
	num := typeInfo1.NumField()
	for i := 0; i < num; i++ {
		key := typeInfo1.Field(i).Name
		if key == "Model" {
			continue
		}
		val1 := String(valInfo1.Field(i).Interface())
		val2 := String(valInfo2.Field(i).Interface())
		tmp1, tmp2 := []byte(val1), []byte(val2)
		sort.Slice(tmp1, func(i, j int) bool {
			return tmp1[i] < tmp1[j]
		})
		sort.Slice(tmp2, func(i, j int) bool {
			return tmp2[i] < tmp2[j]
		})
		if string(tmp1) != string(tmp2) { //记录改变的属性
			tmp := fmt.Sprintf("%s:%v -> %v;", key, val1, val2)
			operationLog = operationLog + tmp
		}
	}
	return operationLog
}

// GetDifference2 比较对象修改前后的不同
func GetDifference2(obj interface{}, fields map[string]interface{}) string {
	operationLog := ""
	var typeInfo1 = reflect.TypeOf(obj)
	var valInfo1 = reflect.ValueOf(obj)
	num := typeInfo1.NumField()
	for i := 0; i < num; i++ {
		key := typeInfo1.Field(i).Name
		val1 := valInfo1.Field(i).Interface()
		if val2, ok := fields[key]; ok {
			if String(val1) != String(val2) {
				tmp := fmt.Sprintf("%s:%v >> %v; ", key, val1, val2)
				operationLog = operationLog + tmp
			}
		}
	}
	return operationLog
}

// BBoxesOverlap 判断是否重叠
func BBoxesOverlap(box1 []float64, box2 []float64, kuobian float64) bool {
	if len(box1) > 4 {
		// 样本数据标注框可能是多边形，先转成外接矩形
		minX, minY, maxX, maxY := 9999999.0, 9999999.0, 0.0, 0.0
		for i := 0; i < len(box1); i += 2 {
			if box1[i] < minX {
				minX = box1[i]
			}
			if box1[i] > maxX {
				maxX = box1[i]
			}
			if box1[i+1] < minY {
				minY = box1[i+1]
			}
			if box1[i+1] > maxY {
				maxY = box1[i+1]
			}
		}
		box1 = []float64{minX, minY, maxX, maxY}
	}
	if box1[0]-kuobian > box2[2] {
		return false
	}
	if box1[1]-kuobian > box2[3] {
		return false
	}
	if box1[2]+kuobian < box2[0] {
		return false
	}
	if box1[3]+kuobian < box2[1] {
		return false
	}
	return true
}

// GetStringOfStruct 遍历结构体的属性和对应值，返回字符串
func GetStringOfStruct(obj interface{}) string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	res := ""
	for k := 0; k < t.NumField(); k++ {
		if t.Field(k).Name == "Model" {
			continue
		}
		tmp := fmt.Sprintf("%s: %v  ", t.Field(k).Name, v.Field(k).Interface())
		res += tmp
	}
	return res
}

// 字符串转日期
func ParseTime(layout string, timeStr string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}

// SliceIntIsRepeat int类型的slice是否有重复
func SliceIntIsRepeat(container []int) bool {
	if len(container) <= 1 {
		return false
	}
	set := make(map[int]bool)
	for _, item := range container {
		if set[item] {
			return true
		}
		set[item] = true
	}
	return false
}

func SliceContainsAny(array []string, item string) bool {
	for _, cur := range array {
		if cur == item {
			return true
		}
	}
	return false
}

// StringContainsAny 判断字符串中是否包含任意字符
func StringContainsAny(s string, strs []string) bool {
	for _, item := range strs {
		if strings.Contains(s, item) {
			return true
		}
	}
	return false
}

// StringIsNum 判断string是否是数字
func StringIsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// IndexSort // 只能作用于正整数列表
func IndexSort(array []int, max int) []int {
	if max == 0 {
		for _, value := range array {
			if value > max {
				max = value
			}
		}
	}
	tempListForSort := make([]int, max+1)
	for _, num := range array {
		tempListForSort[num]++
	}
	var sortedList []int
	for index, value := range tempListForSort {
		if value != 0 {
			sortedList = append(sortedList, index)
		}
	}
	return sortedList
}

// ConvertFileSize 将文件大小从字节转换为合适的单位
func ConvertFileSize(sizeInBytes float64) (float64, string) {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	var convertedSize float64
	var selectedUnit string

	for _, unit := range units {
		if sizeInBytes < 1024.0 {
			convertedSize = sizeInBytes
			selectedUnit = unit
			break
		}
		sizeInBytes /= 1024.0
	}

	return Decimal(convertedSize, 2), selectedUnit
}

// NumberToChinese 取[0,1999]的中文字符
func NumberToChinese(num int, preStr string, sufStr string) string {
	if num > 1999 || num < 0 {
		return "越界"
	}
	var units = [4]string{"", "十", "百", "千"}
	var numerals = [10]string{"", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	if num == 0 {
		return "零"
	}

	str := fmt.Sprintf("%d", num)
	length := len(str)

	var result strings.Builder
	for i, digit := range str {
		digitInt := int(digit - '0')
		unit := units[length-i-1]
		if digitInt == 0 {
			if strings.HasSuffix(result.String(), "零") || length-i == 1 {
				continue
			}
			result.WriteString("零")
			continue
		} else {
			result.WriteString(numerals[digitInt])
			result.WriteString(unit)
		}
	}

	// 处理十位数数字，“一十”转为“十”
	if strings.HasPrefix(result.String(), "一十") {
		return strings.TrimPrefix(result.String(), "一")
	}

	return preStr + result.String() + sufStr
}

// ChineseToNumber 中文字符转数字。
func ChineseToNumber(chineseNum string) int {
	numerals := map[rune]int{
		'零': 0,
		'一': 1,
		'二': 2,
		'三': 3,
		'四': 4,
		'五': 5,
		'六': 6,
		'七': 7,
		'八': 8,
		'九': 9,
	}

	units := map[rune]int{
		'十': 10,
		'百': 100,
		'千': 1000,
	}

	result := 0      // 最终结果
	tmp := 0         // 临时数字，用于存储个位数
	currentUnit := 1 // 当前的单位：个、十、百、千...

	// 将字符串转化为rune类型的切片，这样可以正确处理中文字符
	characters := []rune(chineseNum)

	for i := len(characters) - 1; i >= 0; i-- {
		number, exists := numerals[characters[i]]

		if exists {
			tmp = number * currentUnit
			if currentUnit > 1 && number == 0 { // 如果单位大于1且数字为0，表示零位，跳过
				continue
			}
			result += tmp
			if number > 0 {
				currentUnit = 1 // 重置单位
			}
		} else {
			currentUnit, exists = units[characters[i]]

			// 如果当前位是"十"且为首位，则视为10
			if currentUnit == 10 && i == 0 {
				result += currentUnit
			}
		}
	}

	return result
}

// CutStringBeforeBracket 截取以startChars中任意字符开头的前置字符
func CutStringBeforeBracket(s string, startChars []rune) string {
	m := make(map[rune]bool)
	for _, item := range startChars {
		m[item] = true
	}
	runes := []rune(s) // 转换为rune切片来正确处理中文以及其他Unicode字符

	// 遍历rune切片查找'('或'（'
	for i, r := range runes {
		if m[r] {
			return string(runes[:i])
		}
	}

	// 如果没有找到任何括号，返回原字符串
	return s
}

// DownLoadFile 文件下载
// url:文件下载地址
// fileName:文件名称(默认存储在程序运行的相对位置，后缀默认为原始后缀)
func DownLoadFile(url string, fileName string) error {
	// 发起GET请求
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("图像[%s]保存失败.%v", fileName, err)
	}
	defer response.Body.Close()

	// 确保我们得到了一个成功的响应
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("图像[%s]保存失败.%v", fileName, err)
	}
	fileName = fmt.Sprintf("%s%s", strings.Split(fileName, ".")[0], filepath.Ext(url))
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("图像[%s]保存失败.%v", fileName, err)
	}
	defer file.Close()

	// 将HTTP响应的Body内容写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		_ = os.RemoveAll(fileName)
		return fmt.Errorf("图像[%s]保存失败.%v", fileName, err)
	}
	return nil
}

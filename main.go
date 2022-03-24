package main

import (
	"errors"
	"fmt"
	"github.com/go-pdf/fpdf"
	"math/rand"
	"os/exec"
	"runtime"
	"time"
)

const Quantity = 50

type Formula func(int, int, int) (string, error)

type Course struct {
	name     string
	formulas []Formula
}

var (
	addition = func(n1, n2, n3 int) (string, error) {
		return fmt.Sprintf("%d+%d=", n1, n2), nil
	}
	subtraction = func(n1, n2, n3 int) (string, error) {
		if n1 < n2 {
			return "", errors.New("n1 > n2")
		}
		return fmt.Sprintf("%d-%d=", n1, n2), nil
	}
	additionAndSubtraction = func(n1, n2, n3 int) (string, error) {
		if n1+n2 < n3 {
			return "", errors.New("n1 > n2")
		}
		return fmt.Sprintf("%d+%d-%d=", n1, n2, n3), nil
	}
	subtractionAndAddition = func(n1, n2, n3 int) (string, error) {
		if n1 < n2 {
			return "", errors.New("n1 > n2")
		}
		return fmt.Sprintf("%d-%d+%d=", n1, n2, n3), nil
	}
	additionAndAddition = func(n1, n2, n3 int) (string, error) {
		return fmt.Sprintf("%d+%d+%d=", n1, n2, n3), nil
	}
	subtractionAndSubtract = func(n1, n2, n3 int) (string, error) {
		if n1 < n2 {
			return "", errors.New("n1 > n2")
		}
		return fmt.Sprintf("%d-%d+%d=", n1, n2, n3), nil
	}
	multiplication = func(n1, n2, n3 int) (string, error) {
		return fmt.Sprintf("%dx%d=", n1, n2), nil
	}
	division = func(n1, n2, n3 int) (string, error) {
		if n1%n2 > 0 {
			return "", errors.New("n1 > n2")
		}
		return fmt.Sprintf("%d÷%d=", n1, n2), nil
	}
)

var (
	algorithm = map[int]Course{
		1: Course{name: "加法", formulas: []Formula{
			addition,
		}},
		2: Course{name: "减法", formulas: []Formula{
			subtraction,
		}},
		3: Course{name: "加减法", formulas: []Formula{
			addition,
			subtraction,
			subtractionAndAddition,
			additionAndSubtraction,
			additionAndAddition,
			subtractionAndSubtract,
		}},
		4: Course{name: "乘法", formulas: []Formula{
			multiplication,
		}},
		5: Course{name: "除法", formulas: []Formula{
			division,
		}},
		6: Course{name: "乘除法", formulas: []Formula{
			multiplication,
			division,
		}},
		7: Course{name: "乘除加减", formulas: []Formula{
			multiplication,
			division,
		}},
	}
)

func generate(course Course, min, max int) []string {

	rand.Seed(time.Now().UnixNano())
	c := len(course.formulas)

	var results []string
	j := 0
	for i := 0; i < Quantity; {
		index := rand.Intn(c)
		callable := course.formulas[index]

		n1, n2, n3 := rand3Numbers(min, max, nil)
		row, err := callable(n1, n2, n3)
		j++

		if j >= 100000 {
			break
		}

		// 重复再跑一次
		if err != nil || isExist(results, row) {
			continue
		}
		results = append(results, row)
		i++
	}

	return results
}

func isExist(results []string, row string) bool {
	for _, v := range results {
		if v == row {
			return true
		}
	}

	return false
}

func rand3Numbers(min, max int, verify func(n1, n2, n3 int) bool) (int, int, int) {
	rand.Seed(time.Now().UnixNano() * time.Now().UnixNano())
	for {
		n1 := min + rand.Intn(max-min+1)
		n2 := min + rand.Intn(max-min+1)
		n3 := min + rand.Intn(max-min+1)
		if verify != nil && !verify(n1, n2, n3) {
			continue
		}

		return n1, n2, n3
	}
}

func isValidAlgorithm(values []int, i int) bool {
	for _, v := range values {
		if v == i {
			return true
		}
	}

	return false
}

func openFile(file string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", file)
	case "linux":
		cmd = exec.Command("gnome-open", file)
	default:
		cmd = exec.Command("open", file)
	}

	if err := cmd.Start(); err != nil {
		return errors.New("PDF 打开失败")
	}

	return nil
}

func main() {
	fmt.Println("数字算术出题")
	var t, min, max int
	var is []int

	for i := 1; i < len(algorithm)+1; i++ {
		fmt.Printf("%d.%s ", i, algorithm[i].name)
		is = append(is, i)
	}

	var course Course

	for {
		fmt.Print("\n请选择:")
		fmt.Scanf("%d:", &t)

		if !isValidAlgorithm(is, t) {
			fmt.Println("你的选择不正确")
			continue
		}

		course = algorithm[t]

		fmt.Printf("你已选择: %s", course.name)

		fmt.Print("\n请填最小数字:")
		fmt.Scan(&min)
		fmt.Print("请填最大数字(最小填10):")
		fmt.Scan(&max)

		if (min == 0 || max == 0) || min >= max {
			continue
		}

		if max < 10 {
			max = 10
		}

		if max < 11 || max-min < 10 {
			min = 1
		}

		fmt.Printf("\n\n你的选择: [%s] 最小数字：%d 最大数字：%d\n", course.name, min, max)
		break
	}

	fmt.Println("计算出以下题目：")
	results := generate(course, min, max)
	fmt.Printf("==共%d题==\n", len(results))

	fmt.Println(results)

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.AddUTF8Font("NotoSansSC-Regular", "", "NotoSansSC-Regular.ttf")
	pdf.SetFont("NotoSansSC-Regular", "", 16)

	var x float64 = 20
	var y float64
	// 一行显示多小个
	var rowCount = 4

	j := 0
	//line := ""
	_, lineHeight := pdf.GetFontSize()

	pdf.Write(lineHeight, fmt.Sprintf("%s (%d-%d)", course.name, min, max))
	pdf.Ln(20)
	y = pdf.GetY()

	for _, v := range results {

		x = float64(j)*40 + 20

		pdf.SetXY(x, y)
		pdf.Write(lineHeight, v)

		if j+1 == rowCount {
			pdf.Ln(15)
			y = pdf.GetY()
			x = 20
			//y += 20
			j = 0
			continue
		}

		j++
	}

	pdf.Ln(15)
	pdf.SetFont("NotoSansSC-Regular", "", 9)
	pdf.Write(lineHeight, fmt.Sprintf("Printed at: %s  %s", time.Now().Format(time.RFC850), "github.com/atans"))
	now := time.Now()
	filename := fmt.Sprintf("%s(%d-%d)_%d%d%d_%d%d%02d.pdf", course.name, min, max, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		fmt.Println("输出文件失败")
	} else {
		fmt.Printf("输出文件: %s\n", filename)
		if err := openFile(filename); err != nil {
			fmt.Println(err)
		}
	}

}

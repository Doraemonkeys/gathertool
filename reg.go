/*
	Description : 正则
	Author : ManGe
	Version : v0.2
	Date : 2021-05-15
*/

package gathertool

import (
	"regexp"
	"runtime"
	"strings"
)

// td:=`<td>(.*?)</td>`
// tdreg := regexp.MustCompile(td)
// tdList := tdreg.FindAllStringSubmatch(rest,-1)

func RegFindAll(regStr, rest string) [][]string {
	reg := regexp.MustCompile(regStr)
	List := reg.FindAllStringSubmatch(rest, -1)
	reg.FindStringSubmatch(rest)
	return List
}

var regMap = map[string]string{
	"RegHtmlA": `(?is:<a.*?</a>)`,
	"RegHtmlTitle": `(?is:<title.*?</title>)`,
	"RegHtmlTr": `(?is:<tr.*?</tr>)`,
	"RegHtmlInput": `(?is:<input.*?>)`,
	"RegHtmlTd": `(?is:<td.*?</td>)`,
	"RegHtmlP": `(?is:<p.*?</p>)`,
	"RegHtmlSpan": `(?is:<span.*?</span>)`,
	"RegHtmlSrc": `(?is:src=\".*?\")`,
	"RegHtmlHref": `(?is:href=\".*?\")`,
	"RegHtmlH1": `(?is:<h1.*?</h1>)`,
	"RegHtmlH2": `(?is:<h2.*?</h2>)`,
	"RegHtmlH3": `(?is:<h3.*?</h3>)`,
	"RegHtmlH4": `(?is:<h4.*?</h4>)`,
	"RegHtmlH5": `(?is:<h5.*?</h5>)`,
	"RegHtmlH6": `(?is:<h6.*?</h6>)`,

	// 取标签内容
	"RegHtmlSrcTxt": `(?is:src=\"(.*?)\")`,
	"RegHtmlHrefTxt": `(?is:href=\"(.*?)\")`,
	"RegHtmlHTxt1": `(?is:<h1.*?>(.*?)</h1>)`,
	"RegHtmlHTxt2": `(?is:<h2.*?>(.*?)</h2>)`,
	"RegHtmlHTxt3": `(?is:<h3.*?>(.*?)</h3>)`,
	"RegHtmlHTxt4": `(?is:<h4.*?>(.*?)</h4>)`,
	"RegHtmlHTxt5": `(?is:<h5.*?>(.*?)</h5>)`,
	"RegHtmlHTxt6": `(?is:<h6.*?>(.*?)</h6>)`,
}

// 获取正在运行的函数名
func runFuncName()string{
	pc := make([]uintptr,1)
	runtime.Callers(2,pc)
	f := runtime.FuncForPC(pc[0])
	fName := f.Name()
	fList := strings.Split(fName,".")
	return fList[len(fList)-1]
}

// 执行正则提取
func regFind(funcName, txt string) (dataList []string) {
	regStr, ok := regMap[funcName]
	if !ok{
		loger("reg func is not")
		return
	}
	reg := regexp.MustCompile(regStr)
	resList := reg.FindAllStringSubmatch(txt, -1)
	for _,v := range resList{
		loger(v)
		dataList = append(dataList, v[0])
	}
	return
}


func RegHtmlA(str string)[]string { return regFind(runFuncName(), str) }

func RegHtmlTitle(str string)[]string { return regFind(runFuncName(), str) }

func RegHtmlTr(str string)[]string { return regFind(runFuncName(), str) }

func RegHtmlInput(str string) []string { return regFind(runFuncName(), str) }

func RegHtmlTd(str string) []string { return regFind(runFuncName(), str) }

func RegHtmlP(str string) []string { return regFind(runFuncName(), str) }

func RegHtmlSpan(str string) []string { return regFind(runFuncName(), str) }

func RegHtmlSrc(str string) []string { return regFind(runFuncName(), str) }

func RegHtmlHref(str string) []string { return regFind(runFuncName(), str) }

func RegHtmlH(str, typeH string) []string {
	funcName := runFuncName()
	return regFind(funcName+typeH, str)
}



// 执行正则提取 只取内容
func regFindTxt(funcName, txt string) (dataList []string) {
	regStr, ok := regMap[funcName]
	if !ok{
		loger("reg func is not")
		return
	}
	reg := regexp.MustCompile(regStr)
	resList := reg.FindAllStringSubmatch(txt, -1)
	for _,v := range resList{
		if len(v) > 1{
			dataList = append(dataList, v[1])
		}
	}
	return
}

func RegHtmlATxt(str string)[]string { return regFindTxt(runFuncName(), str) }

func RegHtmlTitleTxt(str string)[]string { return regFindTxt(runFuncName(), str) }

func RegHtmlTrTxt(str string)[]string { return regFindTxt(runFuncName(), str) }

func RegHtmlInputTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlTdTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlPTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlSpanTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlSrcTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlHrefTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlHTxt(str, typeH string) []string {
	funcName := runFuncName()
	return regFindTxt(funcName+typeH, str)
}

package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/lib/github.com/antonmedv/expr"
	"strings"
)

type ExpressionEngineExpr struct {
	mapCache *ExpressionEngineLexerMapCache
}

//引擎名称
func (this *ExpressionEngineExpr) Name() string {
	return "ExpressionEngineExpr"
}

//编译一个表达式
//参数：lexerArg 表达式内容
//返回：interface{} 编译结果,error 错误
func (this *ExpressionEngineExpr) Lexer(expression string) (interface{}, error) {
	expression = this.repleaceExpression(expression)
	var result, err = expr.Parse(expression)
	return result, err
}

//执行一个表达式
//参数：lexerResult=编译结果，arg=参数
//返回：执行结果，错误
func (this *ExpressionEngineExpr) Eval(lexerResult interface{}, arg interface{}, operation int) (interface{}, error) {
	output, err := expr.Run(lexerResult.(expr.Node), arg)
	return output, err
}

//替换表达式中的值 and,or,参数 替换为实际值
func (this *ExpressionEngineExpr) repleaceExpression(expression string) string {
	if expression == "" {
		return expression
	}
	expression = strings.Replace(expression, ` and `, " && ", -1)
	expression = strings.Replace(expression, ` or `, " || ", -1)
	return expression
}

func (this *ExpressionEngineExpr) split(str string) (stringItems []string) {
	if str == "" {
		return nil
	}
	var andStrings = strings.Split(str, " && ")
	if andStrings == nil {
		return nil
	}
	var newStrings []string
	for _, v := range andStrings {
		var orStrings = strings.Split(v, " || ")
		if orStrings == nil {
			continue
		}
		for _, orStr := range orStrings {
			if newStrings == nil {
				newStrings = make([]string, 0)
			}
			if orStr == "" {
				continue
			}
			newStrings = append(newStrings, orStr)
		}
	}
	return newStrings
}

//Lexer缓存,可不提供。
func (this *ExpressionEngineExpr) LexerCache() ExpressionEngineLexerCache {
	if this.mapCache == nil {
		var cache = ExpressionEngineLexerMapCache{}.New()
		this.mapCache = &cache
	}
	return this.mapCache
}

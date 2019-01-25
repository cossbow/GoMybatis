package example

import (
	"bytes"
	"fmt"
	"github.com/zhuxiujia/GoMybatis"
	"github.com/zhuxiujia/GoMybatis/utils"
	"testing"
	"time"
)

type ExampleActivityGoMapper struct {
	SelectByCondition func(name *string, startTime *time.Time, endTime *time.Time, page *int, size *int) ([]Activity, error) `mapperParams:"name,startTime,endTime,page,size"`
}

func (this ExampleActivityGoMapper) New() ExampleActivityGoMapper {
	this.SelectByCondition = func(name *string, startTime *time.Time, endTime *time.Time, page *int, size *int) (activities []Activity, e error) {
		var buf bytes.Buffer
		var bind bytes.Buffer
		bind.WriteString(`%` + utils.GetValue(name, GoMybatis.PtrStringType) + `%`)
		buf.WriteString(`select * from biz_activity`)
		if name != nil {
			buf.WriteString(" and name like ")
			buf.WriteString(utils.GetValue(name, GoMybatis.PtrStringType))
		}
		if startTime != nil {
			buf.WriteString(" ")
			buf.WriteString(utils.GetValue(startTime, GoMybatis.PtrTimeType))
		}
		if endTime != nil {
			buf.WriteString(" ")
			buf.WriteString(utils.GetValue(endTime, GoMybatis.PtrTimeType))
		}
		if page != nil {
			buf.WriteString(" ")
			buf.WriteString(utils.GetValue(page, GoMybatis.PtrIntType))
		}
		if size != nil {
			buf.WriteString(" ")
			buf.WriteString(utils.GetValue(size, GoMybatis.PtrIntType))
		}
		//fmt.Println(buf.String())
		return nil, nil
	}
	return this
}

func Test(t *testing.T) {
	var name = "dsaf"
	var date = time.Now()
	var mapper = ExampleActivityGoMapper{}.New()

	var sql, _ = mapper.SelectByCondition(&name, &date, &date, nil, nil)
	fmt.Println(sql)
}

func Benchmark_sql_test(b *testing.B) {
	b.StopTimer()
	var name = "dsaf"
	var date = time.Now()
	var mapper = ExampleActivityGoMapper{}.New()
	var page = 1
	var size = 20
	defer utils.CountMethodTps(b.N, time.Now(), "sql_test")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var act, _ = mapper.SelectByCondition(&name, &date, &date, &page, &size)
		if act != nil {

		}

	}
}

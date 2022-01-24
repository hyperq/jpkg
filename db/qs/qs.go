package qs

import (
	"fmt"
	"github.com/hyperq/jpkg/tool"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var DefaultPageSize = 10

type QuerySet struct {
	where        map[string]interface{}
	sortarray    []string
	Having       string
	Offset       int
	Limit        int
	GroupBy      string
	OrderBy      string
	formatwhere  string
	formatother  string
	formatparams []interface{}
	Select       []string
}

func (q *QuerySet) Add(key string, value ...interface{}) *QuerySet {
	if _, ok := q.where[key]; !ok {
		q.sortarray = append(q.sortarray, key)
	}
	if len(value) == 1 {
		q.where[key] = value[0]
	} else {
		q.where[key] = value
	}
	return q
}

func (q *QuerySet) EQ(key string, value interface{}) *QuerySet {
	q.Add(key+"=?", value)
	return q
}

func (q *QuerySet) LT(key string, value interface{}) *QuerySet {
	q.Add(key+"<?", value)
	return q
}

func (q *QuerySet) GT(key string, value interface{}) *QuerySet {
	q.Add(key+">?", value)
	return q
}

func (q *QuerySet) NE(key string, value interface{}) *QuerySet {
	q.Add(key+"!=?", value)
	return q
}

func (q *QuerySet) GE(key string, value interface{}) *QuerySet {
	q.Add(key+">=?", value)
	return q
}

func (q *QuerySet) LE(key string, value interface{}) *QuerySet {
	q.Add(key+"<=?", value)
	return q
}

func (q *QuerySet) Like(key, value string) *QuerySet {
	q.Add(key+" LIKE ?", "%"+value+"%")
	return q
}

func (q *QuerySet) LikeLeft(key, value string) *QuerySet {
	q.Add(key+" LIKE ?", value+"%")
	return q
}

func (q *QuerySet) LikeRight(key, value string) *QuerySet {
	q.Add(key+" LIKE ?", "%"+value)
	return q
}

func (q *QuerySet) Between(key, left, right string) {
	q.GE(key, left)
	q.LE(key, right)
}

func (q *QuerySet) IN(key string, value ...interface{}) *QuerySet {
	vlength := len(value)
	if vlength == 0 {
		panic("value len must > 1")
	}
	if vlength == 1 {
		switch t := value[0].(type) {
		case []interface{}:
			q.Add(key+" in (?"+strings.Repeat(",?", len(t)-1)+")", t...)
		case []string:
			q.Add(key+" in (?"+strings.Repeat(",?", len(t)-1)+")", t)
		case []int:
			q.Add(key+" in (?"+strings.Repeat(",?", len(t)-1)+")", t)
		case []int8:
			q.Add(key+" in (?"+strings.Repeat(",?", len(t)-1)+")", t)
		case []int16:
			q.Add(key+" in (?"+strings.Repeat(",?", len(t)-1)+")", t)
		case []int32:
			q.Add(key+" in (?"+strings.Repeat(",?", len(t)-1)+")", t)
		case []int64:
			q.Add(key+" in (?"+strings.Repeat(",?", len(t)-1)+")", t)
		case []float32:
			q.Add(key+" in (?"+strings.Repeat(",?", len(t)-1)+")", t)
		case []float64:
			q.Add(key+" in (?"+strings.Repeat(",?", len(t)-1)+")", t)
		case string:
			q.Add(key+" in ("+value[0].(string)+")", false)
		default:
			q.Add(key+"=?", value[0])
		}
	} else {
		q.Add(key+" in (?"+strings.Repeat(",?", len(value)-1)+")", value...)
	}
	return q
}

func New(initquery ...interface{}) *QuerySet {
	qs := &QuerySet{where: make(map[string]interface{}), Limit: 1}
	qs.EQ("a.is_delete", 0)
	qs.init(initquery...)
	return qs
}

func (q *QuerySet) init(initquery ...interface{}) {
	if len(initquery)%2 == 1 {
		panic("传入参数必须是偶数个")
	}
	for i := range initquery {
		if i%2 == 0 {
			q.Add(initquery[i].(string), initquery[i+1])
		}
	}
}

func (q *QuerySet) Init(initquery ...interface{}) {
	q.where = make(map[string]interface{})
	q.sortarray = []string{}
	q.init(initquery...)
	q.ResetOther()
}

// Format format Qs struct
func (q *QuerySet) Format() (w string, p []interface{}, other string) {
	if q.formatwhere == "" {
		w, p = q.FormatWhere()
		other = q.FormatOther()
		q.formatother = other
		q.formatwhere = w
		q.formatparams = p
	} else {
		w = q.formatwhere
		p = q.formatparams
		other = q.formatother
	}
	return
}

func (q *QuerySet) FormatCache(key string) (res string) {
	w, p, o := q.Format()
	res = key + tool.ToMD5(w+o+fmt.Sprint(p))
	return
}

func (q *QuerySet) FormatWhere() (w string, p []interface{}) {
	var ws []string
	for i := range q.sortarray {
		v, ok := q.where[q.sortarray[i]]
		if ok {
			ws = append(ws, q.sortarray[i])
			if v != false {
				switch t := v.(type) {
				case []interface{}:
					p = append(p, t...)
				case []string:
					for _, tnode := range t {
						p = append(p, tnode)
					}
				case []int:
					for _, tnode := range t {
						p = append(p, tnode)
					}
				case []int8:
					for _, tnode := range t {
						p = append(p, tnode)
					}
				case []int16:
					for _, tnode := range t {
						p = append(p, tnode)
					}
				case []int32:
					for _, tnode := range t {
						p = append(p, tnode)
					}
				case []int64:
					for _, tnode := range t {
						p = append(p, tnode)
					}
				case []float32:
					for _, tnode := range t {
						p = append(p, tnode)
					}
				case []float64:
					for _, tnode := range t {
						p = append(p, tnode)
					}
				default:
					p = append(p, t)
				}
			}
		}
	}
	if len(ws) > 0 {
		w = " WHERE "
	}
	w += strings.Join(ws, " and ")
	return
}

func (q *QuerySet) FormatOther() (fs string) {
	if q.GroupBy != "" {
		fs += " GROUP BY " + q.GroupBy
	}
	if q.Having != "" {
		fs += " HAVING " + q.Having
	}
	if q.OrderBy != "" {
		fs += " ORDER BY " + q.OrderBy
	}
	fs += limit(q.Limit, q.Offset)
	return
}

func limit(limit, offset int) string {
	if offset >= 0 && limit > 0 {
		return " LIMIT " + strconv.Itoa(offset) + "," + strconv.Itoa(limit)
	}
	return ""
}

func (q *QuerySet) ResetOther() {
	q.Offset = 0
	q.Limit = 1
	q.GroupBy = ""
	q.Having = ""
	q.OrderBy = ""
}

func (q *QuerySet) LIMIT(i int) *QuerySet {
	q.Limit = i
	return q
}

func (q *QuerySet) OFFSET(i int) *QuerySet {
	q.Offset = i
	return q
}

func (q *QuerySet) GROUPBY(s string) *QuerySet {
	q.GroupBy = s
	return q
}

func (q *QuerySet) HAVING(s string) *QuerySet {
	q.Having = s
	return q
}

func (q *QuerySet) ORDERBY(s string) *QuerySet {
	q.OrderBy = s
	return q
}

type context interface {
	Query(key string) string
}

func (q *QuerySet) Paging(c context) *QuerySet {
	page := c.Query("current")
	pageSize := c.Query("pageSize")
	pageno, _ := strconv.Atoi(page)
	limit, _ := strconv.Atoi(pageSize)
	if limit < 1 {
		limit = 10
	}
	q.Offset = pagechange(pageno, limit)
	q.Limit = limit
	if orderby := c.Query("orderby"); orderby != "" {
		q.OrderBy = orderby
	} else {
		q.OrderBy = "a.id desc"
	}
	return q
}

type QueryParse struct {
	IN    map[string][]interface{} `json:"in"`
	GT    map[string]interface{}   `json:"gt"`
	EQ    map[string]interface{}   `json:"eq"`
	Like  map[string]string        `json:"like"`
	Time  map[string][]string      `json:"time"`
	LT    map[string]interface{}   `json:"lt"`
	NE    map[string]interface{}   `json:"ne"`
	GE    map[string]interface{}   `json:"ge"`
	LE    map[string]interface{}   `json:"le"`
	LikeL map[string]string        `json:"like_l"`
	LikeR map[string]string        `json:"like_r"`
}

func (q *QuerySet) ParseQuery(c context) *QuerySet {
	qps := c.Query("qp")
	qp := QueryParse{}
	_ = jsoniter.UnmarshalFromString(qps, &qp)
	for k, v := range qp.IN {
		if len(v) > 0 {
			q.IN("a."+k, v)
		}
	}
	for k, v := range qp.NE {
		switch v.(type) {
		case string:
			if v != "" {
				q.NE("a."+k, v)
			}
		default:
			q.NE("a."+k, v)
		}
	}
	for k, v := range qp.EQ {
		switch v.(type) {
		case string:
			if v != "" {
				q.EQ("a."+k, v)
			}
		default:
			q.EQ("a."+k, v)
		}
	}
	for k, v := range qp.GT {
		switch v.(type) {
		case string:
			if v != "" {
				q.GT("a."+k, v)
			}
		default:
			q.GT("a."+k, v)
		}
	}
	for k, v := range qp.GE {
		switch v.(type) {
		case string:
			if v != "" {
				q.GE("a."+k, v)
			}
		default:
			q.GE("a."+k, v)
		}
	}
	for k, v := range qp.LT {
		switch v.(type) {
		case string:
			if v != "" {
				q.LT("a."+k, v)
			}
		default:
			q.LT("a."+k, v)
		}
	}
	for k, v := range qp.LE {
		switch v.(type) {
		case string:
			if v != "" {
				q.LE("a."+k, v)
			}
		default:
			q.LE("a."+k, v)
		}
	}
	for k, v := range qp.Like {
		if v != "" {
			q.Like("a."+k, v)
		}
	}
	for k, v := range qp.LikeL {
		if v != "" {
			q.LikeLeft("a."+k, v)
		}
	}
	for k, v := range qp.LikeR {
		if v != "" {
			q.LikeRight("a."+k, v)
		}
	}
	for k, v := range qp.Time {
		if len(v) == 2 {
			right, _ := time.Parse("2006-01-02", v[1])
			right = right.AddDate(0, 0, 1)
			v[1] = right.Format("2006-01-02")
			q.Between(k, v[0], v[1])
		}
	}
	return q
}

func Auto(c context, initquery ...interface{}) *QuerySet {
	return New(initquery...).Paging(c).ParseQuery(c)
}

func (q *QuerySet) Auto(c context, initquery ...interface{}) *QuerySet {
	return New(initquery...).Paging(c).ParseQuery(c)
}

func (q *QuerySet) SetArray(c context, s ...string) {
	for _, v := range s {
		q.Set(c, "a."+v+"=?", v)
	}
}

func (q *QuerySet) SetInArray(c context, s ...string) {
	for _, v := range s {
		value := c.Query(v + "s")
		if value != "" {
			q.IN("a."+v, value)
		}
	}
}

func (q *QuerySet) SetMap(c context, m map[string]string) {
	for k, v := range m {
		q.Set(c, k, v)
	}
}

func (q *QuerySet) SetLike(c context, querykey, key string) {
	value := c.Query(key)
	if value != "" {
		q.Like(querykey, value)
	}
}

func (q *QuerySet) SetLikeArray(c context, s ...string) {
	for _, v := range s {
		q.SetLike(c, "a."+v, v)
	}
}

func (q *QuerySet) SetLikeL(c context, querykey, key string) {
	value := c.Query(key)
	if value != "" {
		q.LikeLeft(querykey, value)
	}
}

func (q *QuerySet) SetLikeLArray(c context, s ...string) {
	for _, v := range s {
		q.SetLikeL(c, "a."+v, v)
	}
}

func (q *QuerySet) SetBetween(c context, key, left, right string) {
	leftvalue := c.Query(left)
	rightvalue := c.Query(right)
	if leftvalue == "" || rightvalue == "" {
		return
	}
	q.GE(key, leftvalue)
	q.LE(key, rightvalue)
}

func (q *QuerySet) Set(c context, querykey, key string) {
	value := c.Query(key)
	if value != "" {
		q.Add(querykey, value)
	}
}

// pagechange 分页
func pagechange(now int, num int) (page int) {
	if now == 0 {
		now = 1
	}
	if num == 0 {
		num = DefaultPageSize
	}
	page = (now - 1) * num
	return
}

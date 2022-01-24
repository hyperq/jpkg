package ctx

import (
	"github.com/hyperq/jpkg/log"
	"github.com/hyperq/jpkg/tool"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"

	ginjwt "github.com/appleboy/gin-jwt/v2"
)

const (
	STATUSERROR = iota
	STATUSOK
)

// QueryInt 获取url里面的参数并转成int类型
func (c *Context) QueryInt(key string) int {
	q := c.Query(key)
	qi, _ := strconv.Atoi(q)
	return qi
}

// PostFormInt 获取form 里面的参数并转成int类型
func (c *Context) PostFormInt(key string) int {
	q := c.PostForm(key)
	qi, _ := strconv.Atoi(q)
	return qi
}

// UnmarshalFromString 把字符串解析到对象里面 o需要传入指针类型
func (c *Context) UnmarshalFromString(o interface{}) (err error) {
	tool.SetDCM(o)
	err = c.ShouldBind(o)
	if err != nil {
		return
	}
	return
}

// GetAdminID 根据jwt获取用户的id
func (c *Context) GetAdminID() int {
	claims := ginjwt.ExtractClaims(c.Context)
	id, ok := claims[jwt.Auth.IdentityKey]
	if !ok {
		return 0
	}
	return int(id.(float64))
}

// GetSession 获取session
func (c *Context) GetSession(key string) (user interface{}) {
	session := sessions.Default(c.Copy())
	user = session.Get(key)
	return
}

// GetUID 根据session获取用户的id
func (c *Context) GetUID(auth *ginjwt.GinJWTMiddleware) int {
	claims, err := auth.GetClaimsFromJWT(c.Context)
	if err != nil {
		return 0
	}
	// = ginjwt.ExtractClaims(c.Context)
	id, ok := claims[auth.IdentityKey]
	if !ok {
		return 0
	}
	return int(id.(float64))
}

// VersionCheck 数据版本检查
func (c *Context) VersionCheck(version, version1 int) bool {
	if version == version1 {
		return true
	}
	c.RespError("提交数据已被更新, 请刷新后重试")
	return false
}

// IsMobile 判断是否是前端请求
func (c *Context) IsMobile() bool {
	return strings.Split(c.Context.Request.RequestURI, "/")[3] == "m"
}

// R api接口返回值
type R struct {
	Status int         `json:"code"`
	Msg    string      `json:"message"`
	Data   interface{} `json:"result"`
}

func (c *Context) HandlerError(err error) bool {
	if err != nil {
		log.Error2(err)
		r := R{Msg: err.Error(), Status: STATUSERROR}
		c.AbortWithStatusJSON(200, r)
		return true
	}
	return false
}

func (c *Context) HandlerOk(data interface{}) {
	r := R{Status: STATUSOK, Data: data}
	c.JSON(200, r)
}

func (c *Context) RespError(err string) {
	log.Error2(err)
	r := R{Msg: err, Status: STATUSERROR}
	c.AbortWithStatusJSON(200, r)
}

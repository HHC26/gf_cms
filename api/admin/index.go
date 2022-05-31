package admin

import "github.com/gogf/gf/v2/frame/g"

type IndexReq struct {
	g.Meta `tags:"Admin" method:"get" summary:"后台首页"`
}
type IndexRes struct {
	g.Meta `mime:"text/html" example:"string"`
}

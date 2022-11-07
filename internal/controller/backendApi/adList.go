package backendApi

import (
	"context"
	"gf_cms/api/backendApi"
	"gf_cms/internal/dao"
	"gf_cms/internal/model"
	"github.com/gogf/gf/v2/os/gtime"
)

var (
	AdList = cAdList{}
)

type cAdList struct{}

// Index 广告列表
func (c *cAdList) Index(ctx context.Context, req *backendApi.AdListIndexReq) (res *backendApi.AdListIndexRes, err error) {
	var adList []*model.AdListItem
	m := dao.CmsAd.Ctx(ctx).As("ad")
	if req.ChannelId > 0 {
		m = m.Where("ad.channel_id", req.ChannelId)
	}
	err = m.LeftJoin(dao.CmsAdChannel.Table(), "ad_channel", "ad_channel.id=ad.Channel_id").
		Fields("ad.*", "ad_channel.channel_name").
		Page(req.Page, req.Size).Scan(&adList)
	if err != nil {
		return nil, err
	}
	total, _ := m.Count()
	for key, item := range adList {
		if item.Status == 0 {
			adList[key].StatusDesc = "已停用"
		} else if item.StartTime == item.EndTime {
			adList[key].StatusDesc = "长启用"
			adList[key].StartTime = "永久"
			adList[key].EndTime = "永久"
		} else if item.StartTime <= gtime.Datetime() && gtime.Datetime() <= item.EndTime {
			adList[key].StatusDesc = "显示中"
		} else if item.StartTime > gtime.Datetime() {
			adList[key].StatusDesc = "待生效"
		} else if item.EndTime < gtime.Datetime() {
			adList[key].StatusDesc = "已过期"
		}
	}
	res = &backendApi.AdListIndexRes{
		List:  adList,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}
	return
}

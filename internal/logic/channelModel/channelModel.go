package channelModel

import (
	"context"
	"gf_cms/api/backend"
	"gf_cms/internal/dao"
	"gf_cms/internal/model"
	"gf_cms/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sChannelModel struct{}

var (
	insChannelModel = sChannelModel{}
)

func init() {
	service.RegisterChannelModel(New())
}

func New() *sChannelModel {
	return &sChannelModel{}
}

func ChannelModel() *sChannelModel {
	return &insChannelModel
}

func (*sChannelModel) ModelArticle(ctx context.Context, in *backend.ChannelModelIndexReq) (out []*model.ChannelBackendApiListItem, err error) {
	if in.ChannelId > 0 {
		one, err := dao.CmsChannel.Ctx(ctx).Where(dao.CmsChannel.Columns().Id, in.ChannelId).One()
		if err != nil {
			return nil, err
		}
		if one == nil {
			return nil, gerror.New("频道不存在")
		}
	}
	channelTree, err := service.Channel().BackendChannelModelTree(ctx, in.Type, in.ChannelId)
	recycleBin := service.Util().GetSetting("recycle_bin")
	err = service.Response().View(ctx, "backend/channel_model/article/index.html", g.Map{
		"channelTree": channelTree,
		"modelType":   in.Type,
		"modelMap":    service.Channel().BackendModelMap(),
		"channelId":   in.ChannelId,
		"withTab":     in.WithTab,
		"deleteType":  recycleBin,
	})
	return
}

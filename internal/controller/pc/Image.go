package pc

import (
	"context"
	"gf_cms/api/pc"
	"gf_cms/internal/consts"
	"gf_cms/internal/dao"
	"gf_cms/internal/model"
	"gf_cms/internal/model/entity"
	"gf_cms/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	Image = cImage{}
)

type cImage struct{}

// List pc图集列表
func (c *cImage) List(ctx context.Context, req *pc.ImageListReq) (res *pc.ImageListRes, err error) {
	// 栏目详情
	channelInfo, err := Image.channelInfo(ctx, req.ChannelId)
	if err != nil {
		return nil, err
	}
	// 图集列表
	chImagePageList := make(chan *pc.ImageListRes, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		imagePageList, err := Image.imagePageList(ctx, req)
		if err != nil {
			return
		}
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc图集列表分页数据耗时"+gconv.String(endTime-startTime)+"毫秒")
		chImagePageList <- imagePageList
		defer close(chImagePageList)
	}()
	// 导航栏
	chNavigation := make(chan []*model.ChannelPcNavigationListItem, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		navigation, _ := service.Channel().PcNavigation(ctx, gconv.Int(channelInfo.Id))
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc导航耗时"+gconv.String(endTime-startTime)+"毫秒")
		chNavigation <- navigation
		defer close(chNavigation)
	}()
	// TKD
	chTDK := make(chan *model.ChannelTDK, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		pcTDK, _ := service.Channel().PcTDK(ctx, channelInfo.Id, 0)
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pcTDK耗时"+gconv.String(endTime-startTime)+"毫秒")
		chTDK <- pcTDK
		defer close(chTDK)
	}()
	// 面包屑导航
	chCrumbs := make(chan []*model.ChannelCrumbs, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		pcCrumbs, _ := service.Channel().PcCrumbs(ctx, channelInfo.Id)
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc面包屑导航耗时"+gconv.String(endTime-startTime)+"毫秒")
		chCrumbs <- pcCrumbs
		defer close(chCrumbs)
	}()
	// 产品中心栏目列表
	chGoodsChannelList := make(chan []*model.ChannelPcNavigationListItem, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		goodsChannelList, _ := service.Channel().PcHomeGoodsChannelList(ctx, consts.GoodsChannelTid)
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc栏目页产品中心栏目列表耗时"+gconv.String(endTime-startTime)+"毫秒")
		chGoodsChannelList <- goodsChannelList
		defer close(chGoodsChannelList)
	}()
	// 最新资讯-文字新闻
	chTextNewsList := make(chan []*model.ArticleListItem, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		textNewsList, _ := service.Article().PcHomeTextNewsList(ctx, consts.NewsChannelTid)
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc栏目页最新资讯文字新闻列表"+gconv.String(endTime-startTime)+"毫秒")
		chTextNewsList <- textNewsList
		defer close(chTextNewsList)
	}()
	// 获取模板
	chChannelTemplate := make(chan string, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		channelTemplate, _ := service.Channel().PcListTemplate(ctx, channelInfo)
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc获取栏目模板"+gconv.String(endTime-startTime)+"毫秒")
		chChannelTemplate <- channelTemplate
	}()
	imagePageList := <-chImagePageList
	pageInfo := service.PageInfo().PcPageInfo(ctx, imagePageList.Total, imagePageList.Size)
	err = service.Response().View(ctx, <-chChannelTemplate, g.Map{
		"channelInfo":      channelInfo,          // 栏目信息
		"navigation":       <-chNavigation,       // 导航
		"tdk":              <-chTDK,              // TDK
		"crumbs":           <-chCrumbs,           // 面包屑导航
		"goodsChannelList": <-chGoodsChannelList, // 产品中心栏目列表
		"textNewsList":     <-chTextNewsList,     // 最新资讯-文字新闻
		"imagePageList":    imagePageList,        // 图集列表
		"pageInfo":         pageInfo,             // 页码
	})
	if err != nil {
		service.Response().Status500(ctx)
		return nil, err
	}
	return
}

// Detail pc图集详情页面
func (c *cImage) Detail(ctx context.Context, req *pc.ImageDetailReq) (res *pc.ImageDetailRes, err error) {
	// 图集详情
	var imageInfo *model.ImageListItem
	err = dao.CmsImage.Ctx(ctx).
		Where(dao.CmsImage.Columns().Id, req.Id).
		Where(dao.CmsImage.Columns().Status, 1).
		Scan(&imageInfo)
	if err != nil {
		return nil, err
	}
	if imageInfo == nil {
		service.Response().Status404(ctx)
	}
	imageInfo.ClickNum++
	// 更新连击量
	go func() {
		_, err = dao.CmsImage.Ctx(ctx).Where(dao.CmsImage.Columns().Id, imageInfo.Id).Increment(dao.CmsImage.Columns().ClickNum, 1)
	}()
	// 栏目详情
	channelInfo, err := Image.channelInfo(ctx, imageInfo.ChannelId)
	if err != nil {
		return nil, err
	}
	// 导航栏
	chNavigation := make(chan []*model.ChannelPcNavigationListItem, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		navigation, _ := service.Channel().PcNavigation(ctx, gconv.Int(channelInfo.Id))
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc导航耗时"+gconv.String(endTime-startTime)+"毫秒")
		chNavigation <- navigation
		defer close(chNavigation)
	}()
	// TKD
	chTDK := make(chan *model.ChannelTDK, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		pcTDK, _ := service.Channel().PcTDK(ctx, gconv.Uint(imageInfo.ChannelId), gconv.Int64(imageInfo.Id))
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pcTDK耗时"+gconv.String(endTime-startTime)+"毫秒")
		chTDK <- pcTDK
		defer close(chTDK)
	}()
	// 面包屑导航
	chCrumbs := make(chan []*model.ChannelCrumbs, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		pcCrumbs, _ := service.Channel().PcCrumbs(ctx, channelInfo.Id)
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc面包屑导航耗时"+gconv.String(endTime-startTime)+"毫秒")
		chCrumbs <- pcCrumbs
		defer close(chCrumbs)
	}()
	// 产品中心栏目列表
	chGoodsChannelList := make(chan []*model.ChannelPcNavigationListItem, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		goodsChannelList, _ := service.Channel().PcHomeGoodsChannelList(ctx, consts.GoodsChannelTid)
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc文章详情页产品中心栏目列表耗时"+gconv.String(endTime-startTime)+"毫秒")
		chGoodsChannelList <- goodsChannelList
		defer close(chGoodsChannelList)
	}()
	// 最新资讯-文字新闻
	chTextNewsList := make(chan []*model.ArticleListItem, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		textNewsList, _ := service.Article().PcHomeTextNewsList(ctx, consts.NewsChannelTid)
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc文章详情页最新资讯文字新闻列表"+gconv.String(endTime-startTime)+"毫秒")
		chTextNewsList <- textNewsList
		defer close(chTextNewsList)
	}()
	// 上一篇
	chPrevImage := make(chan *model.ImageLink, 1)
	go func() {
		prevImage, _ := service.Image().PcPrevImage(ctx, imageInfo.ChannelId, imageInfo.Id)
		chPrevImage <- prevImage
	}()
	// 下一篇
	chNextImage := make(chan *model.ImageLink, 1)
	go func() {
		nextImage, _ := service.Image().PcNextImage(ctx, imageInfo.ChannelId, imageInfo.Id)
		chNextImage <- nextImage
	}()
	// 在线留言栏目链接
	chGuestbookUrl := make(chan string, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		guestbookUrl, _ := service.GenUrl().PcChannelUrl(ctx, consts.GuestbookChannelTid, "")
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc在线留言栏目url耗时"+gconv.String(endTime-startTime)+"毫秒")
		chGuestbookUrl <- guestbookUrl
		defer close(chGuestbookUrl)
	}()
	// 获取模板
	chChannelTemplate := make(chan string, 1)
	go func() {
		startTime := gtime.TimestampMilli()
		channelTemplate, _ := service.Channel().PcDetailTemplate(ctx, channelInfo)
		endTime := gtime.TimestampMilli()
		g.Log().Async().Info(ctx, "pc获取文章详情模板"+gconv.String(endTime-startTime)+"毫秒")
		chChannelTemplate <- channelTemplate
	}()
	err = service.Response().View(ctx, <-chChannelTemplate, g.Map{
		"navigation":       <-chNavigation,       // 导航
		"tdk":              <-chTDK,              // TDK
		"crumbs":           <-chCrumbs,           // 面包屑导航
		"goodsChannelList": <-chGoodsChannelList, // 产品中心栏目列表
		"textNewsList":     <-chTextNewsList,     // 最新资讯-文字新闻
		"imageInfo":        imageInfo,            // 图集详情
		"prevImage":        <-chPrevImage,        // 上一篇
		"nextImage":        <-chNextImage,        // 下一篇
		"guestbookUrl":     <-chGuestbookUrl,     // 在线留言栏目url
	})
	if err != nil {
		service.Response().Status500(ctx)
		return nil, err
	}
	return
}

// 栏目详情
func (c *cImage) channelInfo(ctx context.Context, channelId int) (out *entity.CmsChannel, err error) {
	err = dao.CmsChannel.Ctx(ctx).
		Where(dao.CmsChannel.Columns().Id, channelId).
		Where(dao.CmsChannel.Columns().Status, 1).
		Where(dao.CmsChannel.Columns().Type, 1).
		Scan(&out)
	if err != nil {
		return
	}
	// 栏目不存在，展示404
	if out == nil {
		service.Response().Status404(ctx)
	}
	return
}

// 获取文章列表分页数据
func (c *cImage) imagePageList(ctx context.Context, in *pc.ImageListReq) (res *pc.ImageListRes, err error) {
	// 当前栏目的所有级别的子栏目id们加自己
	childChannelIds, err := service.Channel().GetChildIds(ctx, in.ChannelId, true)
	if err != nil {
		return nil, err
	}
	m := dao.CmsImage.Ctx(ctx).
		WhereIn(dao.CmsImage.Columns().ChannelId, childChannelIds).
		Where(dao.CmsImage.Columns().Status, 1)
	count, err := m.Count()
	if err != nil {
		return
	}
	var imageList []*model.ImageListItem
	err = m.OrderAsc(dao.CmsImage.Columns().Sort).
		OrderDesc(dao.CmsImage.Columns().Id).
		Page(in.Page, in.Size).
		Scan(&imageList)
	if err != nil {
		return nil, err
	}
	for key, item := range imageList {
		url, err := service.GenUrl().PcDetailUrl(ctx, consts.ChannelModelImage, gconv.Int(item.Id))
		if err != nil {
			return nil, err
		}
		imageList[key].Router = url
		imageList[key], err = service.Image().BuildThumb(ctx, item)
		if err != nil {
			return nil, err
		}
	}
	res = &pc.ImageListRes{
		List:  imageList,
		Page:  in.Page,
		Size:  in.Size,
		Total: count,
	}
	return
}
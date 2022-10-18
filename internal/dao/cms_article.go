// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"gf_cms/internal/dao/internal"
)

// internalCmsArticleDao is internal type for wrapping internal DAO implements.
type internalCmsArticleDao = *internal.CmsArticleDao

// cmsArticleDao is the data access object for table cms_article.
// You can define custom methods on it to extend its functionality as you wish.
type cmsArticleDao struct {
	internalCmsArticleDao
}

var (
	// CmsArticle is globally public accessible object for table cms_article operations.
	CmsArticle = cmsArticleDao{
		internal.NewCmsArticleDao(),
	}
)

// Fill with you ideas below.

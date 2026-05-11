package service

import (
	"reflect"

	"socialai/backend"
	"socialai/constants"
	"socialai/model"

	"github.com/olivere/elastic/v7"
)

func SearchPostsByUser(user string) ([]model.Post, error) {

	// 创建 query
	query := elastic.NewTermQuery("user", user)

	// 调 backend
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)

	if err != nil {
		return nil, err
	}

	return getPostFromSearchResult(searchResult), nil
}

func SearchPostsByKeywords(keywords string) ([]model.Post, error) {

	// 创建全文搜索 query
	query := elastic.NewMatchQuery("message", keywords)

	query.Operator("AND")

	// 空关键词返回全部
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}

	// 调 backend
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)

	if err != nil {
		return nil, err
	}

	return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {

	var ptype model.Post

	var posts []model.Post

	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {

		p := item.(model.Post)

		posts = append(posts, p)
	}

	return posts
}

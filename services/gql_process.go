package services

import (
	"com.fractopus/fractopus-node/gql"
	"com.fractopus/fractopus-node/storage/db_dao"
	"com.fractopus/fractopus-node/storage/model"
	"log"
	"time"
)

func ProcessOnChainUri() {
	latestCursor := db_dao.GetLatestCursor()
	for {
		result, err := gql.GetUriOnChainList(latestCursor)
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Minute)
			continue
		}
		hasNextPage := result.Get("data.transactions.pageInfo.hasNextPage").Bool()
		edges := result.Get("data.transactions.edges").Array()
		log.Println("ProcessOnChainUri edges ", len(edges))
		urlMap := map[string]string{}
		for i, edge := range edges {
			tags := edge.Get("node.tags").Array()
			owner := edge.Get("node.owner.address").String()
			for _, tag := range tags {
				if tag.Get("name").String() == "uri" {
					urlMap[tag.Get("value").String()] = owner
					break
				}
			}

			if hasNextPage {
				if i == len(edges)-1 {
					latestCursor = edge.Get("cursor").String()
				}
			} else {
				if len(edge.Get("cursor").String()) > 0 {
					latestCursor = edge.Get("cursor").String()
				}
			}
		}
		if len(latestCursor) > 0 {
			db_dao.SaveLatestCursor(latestCursor)
		}

		saveUriListToDb(urlMap)
		if hasNextPage {
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(1 * time.Minute)
		}
	}
}
func ProcessWaitOnChainUri() {
	for {
		result, err := gql.GetUriWaitOnChainList()
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Minute)
			continue
		}

		edges := result.Get("data.transactions.edges").Array()

		log.Println("ProcessWaitOnChainUri edges ", len(edges))
		urlMap := map[string]string{}
		for _, edge := range edges {
			tags := edge.Get("node.tags").Array()

			if edge.Get("node.block").IsObject() {
				continue
			}

			owner := edge.Get("node.owner.address").String()
			for _, tag := range tags {
				if tag.Get("name").String() == "uri" {
					urlMap[tag.Get("value").String()] = owner
					break
				}
			}
		}
		saveUriListToDb(urlMap)
		time.Sleep(1 * time.Minute)
	}
}

func saveUriListToDb(urlMap map[string]string) {
	if len(urlMap) > 0 {
		for key := range urlMap {
			if db_dao.CheckUriExist(key) {
				delete(urlMap, key)
			}
		}
	}

	if len(urlMap) > 0 {
		var list []model.OpusUri
		// TODO 爬虫，验证uri是否是owner的
		for key, value := range urlMap {
			list = append(list, model.OpusUri{
				Uri:   key,
				Owner: value,
			})
		}
		db_dao.SaveUris(list)

		//TODO 获取上游详情，写入数据和redis
		for _, node := range list {
			txDetail, err := gql.GetLatestTxDetailByUri(node.Uri)
			log.Println(err)
			log.Println(txDetail.Raw)
		}

	}
}

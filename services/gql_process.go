package services

import (
	"com.fractopus/fractopus-node/gql"
	"com.fractopus/fractopus-node/storage/db_dao"
	"com.fractopus/fractopus-node/storage/model"
	"log"
	"time"
)

func ProcessUri() {
	latestCursor := db_dao.GetLatestCursor()
	for true {
		result, err := gql.GetUriList(latestCursor)
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Minute)
			continue
		}
		hasNextPage := result.Get("data.transactions.pageInfo.hasNextPage").Bool()
		edges := result.Get("data.transactions.edges").Array()
		urlMap := map[string]bool{}
		for i, edge := range edges {
			tags := edge.Get("node.tags").Array()
			for _, tag := range tags {
				if tag.Get("name").String() == "uri" {
					urlMap[tag.Get("value").String()] = true
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

		if len(urlMap) > 0 {
			for key, _ := range urlMap {
				if db_dao.CheckUriExist(key) {
					delete(urlMap, key)
				}
			}
		}

		if len(urlMap) > 0 {
			var list []model.OpusUri
			for key, _ := range urlMap {
				list = append(list, model.OpusUri{
					Uri: key,
				})
			}
			db_dao.SaveUris(list)
		}

		if hasNextPage {
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(2 * time.Minute)
		}
	}
}

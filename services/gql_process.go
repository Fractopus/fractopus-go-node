package services

import (
	"com.fractopus/fractopus-node/gql"
	"com.fractopus/fractopus-node/storage/db_dao"
	"com.fractopus/fractopus-node/storage/model"
	"github.com/tidwall/gjson"
	"log"
	"time"
)

type FractopusInfo struct {
	Uri       string
	Owner     string
	ShareList []gjson.Result
}

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
		urlMap := map[string]FractopusInfo{}
		for i, edge := range edges {
			getOpusInfo(edge, urlMap)
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
			saveUriListToDb(urlMap)
		}

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
		urlMap := map[string]FractopusInfo{}
		for _, edge := range edges {
			if edge.Get("node.block").IsObject() {
				continue
			}
			getOpusInfo(edge, urlMap)
		}
		if len(urlMap) > 0 {
			saveUriListToDb(urlMap)
		}
		time.Sleep(1 * time.Minute)
	}
}

func getOpusInfo(edge gjson.Result, urlMap map[string]FractopusInfo) {
	tags := edge.Get("node.tags").Array()
	opusUri := FractopusInfo{
		Owner: edge.Get("node.owner.address").String(),
	}
	for _, tag := range tags {
		if tag.Get("name").String() == "uri" {
			opusUri.Uri = tag.Get("value").String()
		}
		if tag.Get("name").String() == "shrL" {
			opusUri.ShareList = tag.Get("value").Array()
		}
	}
	if len(opusUri.Uri) > 0 {
		urlMap[opusUri.Uri] = opusUri
	}
}

/*
• “p": "fractopus" / 分形章鱼协议
• "uri" / URI
• shrL:[{"uri":"xxx","shr":"0.1"}] 如果上游的uri没有上链过，在分润的时候，在分润系统记账
*/
func saveUriListToDb(urlMap map[string]FractopusInfo) {
	if len(urlMap) > 0 {
		for key := range urlMap {
			if db_dao.CheckUriExist(key) {
				delete(urlMap, key)
			}
		}
	}

	if len(urlMap) > 0 {
		list := dealMainNodes(urlMap)
		dealUpstream(urlMap)

		//TODO 获取上游详情，写入到redis
		for _, node := range list {
			txDetail, err := gql.GetLatestTxDetailByUri(node.Uri)
			log.Println(err)
			log.Println(txDetail.Raw)
		}

	}
}

func dealMainNodes(urlMap map[string]FractopusInfo) []model.OpusNode {
	var list []model.OpusNode
	// TODO 爬虫，验证uri是否是owner的
	for _, value := range urlMap {
		list = append(list, model.OpusNode{
			Uri:   value.Uri,
			Owner: value.Owner,
		})
		if len(value.ShareList) > 0 {
			for _, shareItem := range value.ShareList {
				if shareItem.Get("uri").Exists() && len(shareItem.Get("uri").String()) > 0 {
					if !db_dao.CheckUriExist(shareItem.Get("uri").String()) {
						list = append(list, model.OpusNode{
							Uri:   shareItem.Get("uri").String(),
							Owner: "",
						})
					}
				}
			}
		}
	}
	if len(list) > 0 {
		db_dao.SaveUris(list)
	}
	return list
}

func dealUpstream(urlMap map[string]FractopusInfo) {
	var streamList []model.OpusStream
	for _, value := range urlMap {
		if len(value.ShareList) > 0 {
			currNodeInfo := db_dao.GetUriNodeByUri(value.Uri)
			if currNodeInfo != nil {
				for _, shareItem := range value.ShareList {
					if !shareItem.Get("uri").Exists() {
						continue
					}
					upstreamUri := shareItem.Get("uri").String()
					ratio := shareItem.Get("shr").Float()
					if ratio < 0 {
						continue
					}
					if len(upstreamUri) > 0 {
						upstreamNodeInfo := db_dao.GetUriNodeByUri(upstreamUri)
						if upstreamNodeInfo != nil {
							streamList = append(streamList, model.OpusStream{
								CurrUriId:     currNodeInfo.ID,
								UpstreamUriId: upstreamNodeInfo.ID,
								Ratio:         ratio,
							})
						}
					}
				}
				db_dao.DeleteUpstreamsByNodeId(currNodeInfo.ID)
			}
		}
	}
	if len(streamList) > 0 {
		db_dao.SaveUpstreams(streamList)
	}
}

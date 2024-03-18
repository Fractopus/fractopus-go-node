package gql

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
)

var baseGqlUrl = "https://knn3-gateway.knn3.xyz/arseeding/graphql"
var baseArSeedingUrl = "https://arseed.web3infra.dev/"

func init() {
	log.Println("begin read data from gql")
}

type GraphQLRequest struct {
	Query string `json:"query"`
}

// GetUriWaitOnChainList 获取等待上链的交易
func GetUriWaitOnChainList() (gjson.Result, error) {
	ql := `
query {
  transactions(
    first: 50
    sort: HEIGHT_ASC
    tags: [{ name: "p", values: ["fractopus"] }]
  ) {
    pageInfo {
      hasNextPage
    }
    edges {
      cursor
      node {
        id
        owner {
          address
        }
        tags {
          name
          value
        }
        block {
          height
          timestamp
        }
      }
    }
  }
}
			`
	return gqlHttpPost(ql)
}
func GetUriOnChainList(lastCursor string) (gjson.Result, error) {
	ql := `query {
			  transactions(
				first: 50
				after:"%v"
				sort: HEIGHT_ASC
				tags: [
				  { name: "p", values: ["fractopus"] }
				]
			  ) {
				pageInfo {
				  hasNextPage
				}
				edges {
				  cursor
				  node {
					id
					owner{
					  address
					}
					tags {
					  name
					  value
					}
				  }
				}
			  }
			}
			`
	ql = fmt.Sprintf(ql, lastCursor)
	return gqlHttpPost(ql)
}

func GetLatestTxDetailByUri(uri string) (gjson.Result, error) {
	ql := `query {
			  transactions(
				first: 1
				sort: HEIGHT_ASC
				tags: [
				  { name: "p", values: ["fractopus"] }
				  { name: "uri", values: ["%v"] }
				]
			  ) {
				edges {
				  node {
					id
					owner{
					  address
					}
					block {
					  height
					}
				  }
				}
			  }
			}`

	ql = fmt.Sprintf(ql, uri)
	result, err := gqlHttpPost(ql)
	if err == nil {
		edges := result.Get("data.transactions.edges").Array()
		if len(edges) > 0 {
			edge := edges[0]
			if !edge.Get("node.block").IsObject() {
				return txDetail(edge.Get("node.id").String())
			}
		}
	}

	ql =
		`query {
			  transactions(
				first: 1
				sort: HEIGHT_DESC
				tags: [
				  { name: "p", values: ["fractopus"] }
				  { name: "uri", values: ["%v"] }
				]
			  ) {
				edges {
				  node {
					id
					owner{
					  address
					}
				  }
				}
			  }
			}
			`
	ql = fmt.Sprintf(ql, uri)
	result, err = gqlHttpPost(ql)
	if err == nil {
		edges := result.Get("data.transactions.edges").Array()
		if len(edges) > 0 {
			edge := edges[0]
			return txDetail(edge.Get("node.id").String())
		}
	}
	return gjson.Result{}, errors.New("no data")
}

func gqlHttpPost(ql string) (gjson.Result, error) {
	data := GraphQLRequest{Query: ql}
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(baseGqlUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return gjson.Result{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return gjson.Result{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return gjson.Parse(string(body)), nil
	}
	return gjson.Result{}, errors.New("wrong get")
}

func txDetail(id string) (gjson.Result, error) {
	resp, err := http.Get(baseArSeedingUrl + id)
	if err != nil {
		log.Println(err)
		return gjson.Result{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return gjson.Result{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return gjson.Parse(string(body)), nil
	}
	return gjson.Result{}, errors.New("wrong get")
}

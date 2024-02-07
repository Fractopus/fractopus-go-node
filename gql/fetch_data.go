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

func Process() {
	post, err := gqlHttpPost(`query {
  transactions(
    # first: 10
    # after: "ODRKM00zMUQxRWtfbnhLVkkzeUI0WUtVSE9OaHd5cXY4ekl6Z0VhVDY2VQ=="
    tags: [
      { name: "Piece-Uuid", values: ["98bb10d0f62f46a3b32b592e21c7536e"] }
      # { name: "App-Name", values: ["Cascad3"] },
    ]
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
        block {
          id
          timestamp
          height
        }
      }
    }
  }
}`)

	if err == nil {
		log.Println(post.Get("data.transactions|@pretty"))
		log.Println(post.Get("data.transactions.pageInfo.hasNextPage").Bool())
	}

	detail, err := txDetail("-GpiL0yIJO5hlPJD2y42v3QfZVsoNRZ0HVXtAGFp3bE")
	if err == nil {
		log.Println(detail)
	}
}

func GetUriWaitOnChainList() (gjson.Result, error) {
	ql := `query {
  transactions(
    owners: ["Bdcp-GSeLfL5gsF19o4yf8jdyVAKn0UdZin1-sU28us"]
    first: 50
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
        tags {
          name
          value
        }
        block{
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
    owners: ["Bdcp-GSeLfL5gsF19o4yf8jdyVAKn0UdZin1-sU28us"]
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

func gqlHttpPost(ql string) (gjson.Result, error) {
	data := GraphQLRequest{Query: ql}
	jsonData, _ := json.Marshal(data)
	resp, _ := http.Post(baseGqlUrl, "application/json", bytes.NewBuffer(jsonData))
	body, _ := io.ReadAll(resp.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return gjson.Parse(string(body)), nil
	}
	return gjson.Result{}, errors.New("wrong get")
}

func txDetail(id string) (gjson.Result, error) {
	resp, _ := http.Get(baseArSeedingUrl + id)
	body, _ := io.ReadAll(resp.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return gjson.Parse(string(body)), nil
	}
	return gjson.Result{}, errors.New("wrong get")
}

package fetcher

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

const ConnectionApiCount = 2

func (f *fetcher) FetchConnections(address string) (results []ConnectionEntry, err error) {
	ch := make(chan ConnectionEntryList)

	// Part 1 - Demo data source
	// Context API
	go f.processContextConn(address, ch)
	// Rarible API
	go f.processRaribleConn(address, ch)
	// Part 2 - Add other data source here
	// TODO

	// Final Part - Aggregate all data & convert ens domain & filter out invalid connections
	for i := 0; i < ConnectionApiCount; i++ {
		entry := <-ch
		if entry.Err != nil {
			zap.L().With(zap.Error(entry.Err)).Error("connection api error: " + entry.msg)
			continue
		}
		results = append(results, entry.Conn...)
	}

	return
}

func (f *fetcher) getRaribleConnection(address string, isFollowing bool) ([]RaribleConnectionResp, error) {
	// Prepare request
	var url string
	if isFollowing {
		url = fmt.Sprintf(RaribleFollowingUrl, address)
	} else {
		url = fmt.Sprintf(RaribleFollowerUrl, address)
	}

	postBody, _ := json.Marshal(map[string]int{
		"size": 5000, // TODO
	})

	body, err := sendRequest(f.httpClient, RequestArgs{
		url:    url,
		method: "POST",
		body:   postBody,
	})

	var results []RaribleConnectionResp
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (f *fetcher) processRaribleConn(address string, ch chan<- ConnectionEntryList) {
	var rarTotal []RaribleConnectionResp
	result := ConnectionEntryList{}

	// Query Followings from Rarible
	rarFollowings, err := f.getRaribleConnection(address, true)
	if err != nil {
		result.Err = err
		result.msg = "[processRaribleConn] fetch Rarible followings failed"
		ch <- result
		return
	}

	// Query Followers from Rarible
	rarFollowers, err := f.getRaribleConnection(address, false)
	if err != nil {
		result.Err = err
		result.msg = "[processRaribleConn] fetch Rarible followers failed"
		ch <- result
		return
	}

	// Merge and printing out for Rarible followings
	rarTotal = append(rarFollowers, rarFollowings...)
	var results []ConnectionEntry
	for i := 0; i < len(rarTotal); i++ {
		if !addressFilter(rarTotal[i].Following.From) || !addressFilter(rarTotal[i].Following.To) {
			continue
		}
		result := ConnectionEntry{
			From:     rarTotal[i].Following.From,
			To:       rarTotal[i].Following.To,
			Platform: RARIBLE,
		}
		results = append(results, result)
	}

	result.Conn = append(result.Conn, results...)
	ch <- result
}

func (f *fetcher) getUserContextConnection(address string, isFollowing bool) (results []ConnectionEntry, err error) {
	var url string

	if isFollowing {
		url = fmt.Sprintf(ContextUrl, address+"/following")
	} else {
		url = fmt.Sprintf(ContextUrl, address+"/followers")
	}

	body, err := sendRequest(f.httpClient, RequestArgs{
		url:    url,
		method: "GET",
	})
	if err != nil {
		return nil, err
	}

	var contextRecord ContextConnection
	err = json.Unmarshal(body, &contextRecord)
	if err != nil {
		return nil, err
	}

	if isFollowing {
		for i := 0; i < len(contextRecord.Relationships); i++ {
			var toAddr string
			toActor := contextRecord.Relationships[i].Actor
			if isAddress(toActor) {
				toAddr = toActor
			} else if len(contextRecord.Profiles[toActor]) != 0 {
				toAddr = contextRecord.Profiles[toActor][0].Address
			} else {
				// Context.app lacks of data
				continue
			}
			if !addressFilter(toAddr) {
				continue
			}
			newContextRecord := ConnectionEntry{
				From:     address,
				To:       toAddr,
				Platform: CONTEXT,
			}
			results = append(results, newContextRecord)
		}
	} else {
		for i := 0; i < len(contextRecord.Relationships); i++ {
			var fromAddr string
			profileAcct := contextRecord.Relationships[i].Actor
			if len(contextRecord.Profiles[profileAcct]) != 0 {
				fromAddr = contextRecord.Profiles[profileAcct][0].Address
			} else {
				// Context.app lacks of data
				continue
			}
			if !addressFilter(fromAddr) {
				continue
			}
			newContextRecord := ConnectionEntry{
				From:     fromAddr,
				To:       address,
				Platform: CONTEXT,
			}
			results = append(results, newContextRecord)
		}
	}
	return results, nil
}

func (f *fetcher) processContextConn(address string, ch chan<- ConnectionEntryList) {
	result := ConnectionEntryList{}
	followingResults, err := f.getUserContextConnection(address, true)
	if err != nil {
		result.Err = err
		result.msg = "[processContextConn] fetch Context followings failed"
		ch <- result
		return
	}

	followerResults, err := f.getUserContextConnection(address, false)
	if err != nil {
		result.Err = err
		result.msg = "[processContextConn] fetch Context followers failed"
		ch <- result
		return
	}

	followingResults = append(followingResults, followerResults...)
	result.Conn = append(result.Conn, followingResults...)
	ch <- result
}

// return false if input is neither Ethereum address nor ENS
func addressFilter(addr string) bool {
	if isAddress(addr) {
		return true
	} else if len(addr) > 4 && addr[len(addr)-4:] == ".eth" {
		return true
	} else {
		return false
	}
}

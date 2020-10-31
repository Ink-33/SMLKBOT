package music

import (
	"SMLKBOT/data/botstruct"
	"SMLKBOT/utils/cqfunction"
	"fmt"
	"strings"
)

func (e *AMusicClient) sendMsg(FR *botstruct.FunctionRequest, d map[string]string) {

	xml := fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?><msg serviceID="2" templateID="1" action="web" brief="[Music] %s" sourceMsgId="0" url="%s" flag="0" adverSign="0" multiMsgFlag="0"><item layout="2"><audio cover="%s" src="%s"/><title>%s</title><summary>%s</summary></item><source name="iTunes" icon="https://www.apple.com/v/itunes/home/k/images/overview/itunes_logo__dwjkvx332d0m_large.png" url="" action="" a_actionData="" i_actionData="" appid="-1" /></msg>`, d["title"], d["url"], d["cover"], d["audio"], d["title"], d["content"])
	xml = strings.ReplaceAll(xml, "&", "&amp;")
	xml = strings.ReplaceAll(xml, "[", "&#91;")
	xml = strings.ReplaceAll(xml, "]", "&#93;")
	xml = strings.ReplaceAll(xml, ",", "&#44;")

	cqfunction.CQSendMsg(FR, "[CQ:xml,data:"+xml+"]")
}

type iTunesAPI struct {
	ResultCount int `json:"resultCount"`
	Results     []struct {
		ArtistName     string `json:"artistName"`
		CollectionName string `json:"collectionName"`
		TrackName      string `json:"trackName"`
		TrackViewURL   string `json:"trackViewUrl"`
		PreviewURL     string `json:"previewUrl"`
		ArtworkURL100  string `json:"artworkUrl100"`
	} `json:"results"`
}

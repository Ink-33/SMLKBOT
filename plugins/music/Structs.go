package music

import (
	"sync"

	"github.com/Ink-33/SMLKBOT/data/botstruct"
)

// Translate json to go by using https://www.sojson.com/json/json2go.html

// GetVTBMusicList also can be used as GetVTBHotMusicList
type GetVTBMusicList struct {
	Total int                   `json:"total"`
	Data  []GetVTBMusicListData `json:"data"`
}

// GetVTBMusicListData is used for json Unmarshal
type GetVTBMusicListData struct {
	ID         string `json:"id"`
	OriginName string `json:"originName"`
	VocalName  string `json:"vocalName"`
	CoverImg   string `json:"coverImg"`
	Music      string `json:"music"`
	CDN        string `json:"cdn"`
}

// GetVTBCDNList is used for json Unmarshal
type GetVTBCDNList struct {
	Data []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"data"`
}

// GetVTBMusicData is used for json Unmarshal
type GetVTBMusicData struct {
	*GetVTBMusicListData `json:"Data"`
}

// GetVTBVtbsData is used for json Unmarshal
type GetVTBVtbsData struct {
	ID           string `json:"id"`
	OriginalName string `json:"originalName"`
	ChineseName  string `json:"chineseName"`
	JapaneseName string `json:"japaneseName"`
	EnglistName  string `json:"englistName"`
	GroupsID     string `json:"groupsId"`
	AvatarImg    string `json:"avatarImg"`
	Bilibili     string `json:"bilibili"`
	YouTube      string `json:"youTube"`
	Twitter      string `json:"twitter"`
	Watch        int    `json:"watch"`
	Introduce    string `json:"introduce"`
}

// VTBMusicInfo includes the info of a music.
type VTBMusicInfo struct {
	MusicName  string
	MusicID    string
	MusicVocal string
	Cover      string
	MusicURL   string
	MusicCDN   string
}

// VTBMusicList includes the result of searching for musics.
type VTBMusicList struct {
	Total int
	Data  []GetVTBMusicListData
}

// VTBsList includes the result of searching for Vtbs.
type VTBsList struct {
	Total int
	Data  []GetVTBVtbsData
}

// EventBusSubcriber 点歌者信息
type EventBusSubcriber struct {
	run  func(*botstruct.FunctionRequest)
	done chan struct{}
}

// EventBus 存储点歌信息
type EventBus struct {
	subscribers map[string]*EventBusSubcriber
	lock        sync.RWMutex
}

// VTBMusicClient : VTBMusic点歌Client
type VTBMusicClient struct {
	MusicList []GetVTBMusicListData
}

// AMusicClient : AppleMusic点歌Client
type AMusicClient struct{}

// Client : 点歌Client
type Client interface {
	getMusicDetailandCQCode(int) string
	musicListLen() int
}

// MsgType
type msgType struct {
	// Search Keyword
	content string
	ctype   int8
}

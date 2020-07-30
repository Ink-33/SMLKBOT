package vtbmusic

//Translate json to go by using https://www.sojson.com/json/json2go.html

//GetMusicList also can be used as GetHotMusicList
type GetMusicList struct {
	Total     int                `json:"total"`
	Data      []GetMusicListData `json:"data"`
	Success   bool               `json:"success"`
	ErrorCode int                `json:"errorCode"`
	Msg       string             `json:"msg"`
}

//GetMusicListData is used for json Unmarshal
type GetMusicListData struct {
	ID              string `json:"id"`
	CreateTime      string `json:"createTime"`
	PublishTime     string `json:"publishTime"`
	CreatorID       string `json:"creatorId"`
	CreatorRealName string `json:"creatorRealName"`
	Deleted         bool   `json:"deleted"`
	OriginName      string `json:"originName"`
	VocalID         string `json:"vocalId"`
	VocalName       string `json:"vocalName"`
	CoverImg        string `json:"coverImg"`
	Music           string `json:"music"`
	Lyric           string `json:"lyric"`
	CDN             string `json:"cdn"`
	Source          string `json:"source"`
	BiliBili        string `json:"biliBili"`
	YouTube         string `json:"youTube"`
	Twitter         string `json:"twitter"`
	Likes           int    `json:"likes"`
	Length          int    `json:"length"`
	Label           string `json:"label"`
	VocalList       []struct {
		ID         string `json:"id"`
		Cn         string `json:"cn"`
		Jp         string `json:"jp"`
		En         string `json:"en"`
		Originlang string `json:"originlang"`
	} `json:"vocalList"`
}

//GetCDNList is used for json Unmarshal
type GetCDNList struct {
	Total int `json:"total"`
	Data  []struct {
		ID         string `json:"id"`
		CreateTime string `json:"createTime"`
		CreatorID  string `json:"creatorId"`
		Name       string `json:"name"`
		URL        string `json:"url"`
		Info       string `json:"info"`
	} `json:"data"`
	Success   bool   `json:"success"`
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
}

//GetMusicData is used for json Unmarshal
type GetMusicData struct {
	*GetMusicListData `json:"Data"`
	Success           bool   `json:"Success"`
	ErrorCode         int    `json:"ErrorCode"`
	Msg               string `json:"Msg"`
}

//GetVtbsList is used for json Unmarshal
type GetVtbsList struct {
	Total     int           `json:"total"`
	Data      []GetVtbsData `json:"data"`
	Success   bool          `json:"success"`
	ErrorCode int           `json:"errorCode"`
	Msg       string        `json:"msg"`
}

//GetVtbsData is used for json Unmarshal
type GetVtbsData struct {
	ID           string `json:"id"`
	CreateTime   string `json:"createTime"`
	CreatorID    string `json:"creatorId"`
	Deleted      bool   `json:"deleted"`
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

//MusicInfo includes the info of a music.
type MusicInfo struct {
	MusicName  string
	MusicID    string
	MusicVocal string
	Cover      string
	MusicURL   string
	MusicCDN   string
}

//MusicList includes the result of searching for musics.
type MusicList struct {
	Total int
	Data  []GetMusicListData
}

//VtbsList includes the result of searching for Vtbs.
type VtbsList struct {
	Total int
	Data  []GetVtbsData
}

type cdnResult interface {
	match(string) string
}

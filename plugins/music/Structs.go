package music

//Translate json to go by using https://www.sojson.com/json/json2go.html

//GetVTBMusicList also can be used as GetVTBHotMusicList
type GetVTBMusicList struct {
	Total     int                `json:"total"`
	Data      []GetVTBMusicListData `json:"data"`
	Success   bool               `json:"success"`
	ErrorCode int                `json:"errorCode"`
	Msg       string             `json:"msg"`
}

//GetVTBMusicListData is used for json Unmarshal
type GetVTBMusicListData struct {
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


//GetVTBCDNList is used for json Unmarshal
type GetVTBCDNList struct {
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

//GetVTBMusicData is used for json Unmarshal
type GetVTBMusicData struct {
	*GetVTBMusicListData `json:"Data"`
	Success           bool   `json:"Success"`
	ErrorCode         int    `json:"ErrorCode"`
	Msg               string `json:"Msg"`
}

//GetVTBVtbsList is used for json Unmarshal
type GetVTBVtbsList struct {
	Total     int           `json:"total"`
	Data      []GetVTBVtbsData `json:"data"`
	Success   bool          `json:"success"`
	ErrorCode int           `json:"errorCode"`
	Msg       string        `json:"msg"`
}

//GetVTBVtbsData is used for json Unmarshal
type GetVTBVtbsData struct {
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

//VTBMusicInfo includes the info of a music.
type VTBMusicInfo struct {
	MusicName  string
	MusicID    string
	MusicVocal string
	Cover      string
	MusicURL   string
	MusicCDN   string
}

//VTBMusicList includes the result of searching for musics.
type VTBMusicList struct {
	Total int
	Data  []GetVTBMusicListData
}

//VTBsList includes the result of searching for Vtbs.
type VTBsList struct {
	Total int
	Data  []GetVTBVtbsData
}

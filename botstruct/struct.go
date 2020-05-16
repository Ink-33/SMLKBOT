package botstruct

//MsgInfo includes some basic info about a message.
type MsgInfo struct {
	SenderID string
	GroupID  string
	Message  string
	MsgType  string
}

//Auinfo includes some basic info of a Au number.
type Auinfo struct {
	AuNumber   string
	AuStatus   bool
	AuMsg      string
	AuJumpURL  string
	AuCoverURL string
	AuURL      string
	AuTitle    string
	AuDesp     string
}

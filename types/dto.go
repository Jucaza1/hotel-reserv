package types

type MsgUpdated struct {
	Updated string `json:"updated"`
}

type MsgDeleted struct {
	Deleted string `json:"deleted"`
}

type MsgCancelled struct {
	Cancelled string `json:"cancelled"`
}

type MsgError struct {
	Error string `json:"error"`
}

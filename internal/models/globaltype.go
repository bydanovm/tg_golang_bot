package models

type StatusRetriever struct {
	MsgError error
}

type StatusChannel struct {
	Start bool
	Stop  bool
	Error error
	Data  interface{}
}

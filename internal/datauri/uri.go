package datauri

type URI struct {
	MediaType string
	Params    map[string]string
	Data      []byte
}

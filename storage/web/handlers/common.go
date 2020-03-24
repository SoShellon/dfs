package handlers

import "cmu.edu/dfs/common"

type pathParams struct {
	Path string `json:"path"`
}
type copyParams struct {
	pathParams
	common.StorageNode
}
type fileParams struct {
	Path   string `json:"path"`
	Offset int64  `json:"offset"`
	Length int64  `json:"length"`
	Data   []byte `json:"data"`
}

func illegalArgumentError(msg string) (int, map[string]string) {
	return 404, map[string]string{
		"exception_type": "IllegalArgumentException",
		"exception_info": msg,
	}
}

func fileNotFoundError(msg string) (int, map[string]string) {
	return 404, map[string]string{
		"exception_type": "FileNotFoundException",
		"exception_info": msg,
	}
}

func indexOutOfBoundsException(msg string) (int, map[string]string) {
	return 404, map[string]string{
		"exception_type": "IndexOutOfBoundsException",
		"exception_info": msg,
	}
}

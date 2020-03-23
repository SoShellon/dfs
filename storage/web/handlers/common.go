package handlers

type pathParams struct {
	Path string `json:"path"`
}

type fileParams struct {
	Path   string `json:"path"`
	Offset int64  `json:"offset"`
	Length int64  `json:"length"`
	Data   []byte `json:"data"`
}

func illegalArgumentError(msg string) (int, map[string]string) {
	return 404, map[string]string{
		"expection_type": "IllegalArgumentException",
		"exception_info": msg,
	}
}

func fileNotFoundError(msg string) (int, map[string]string) {
	return 404, map[string]string{
		"exception_type": "FileNotFoundException",
		"exception_info": msg,
	}
}

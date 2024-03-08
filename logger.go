package room

import (
	"bytes"
	"io"
	"strconv"
)

const (
	startTag       = "${"
	endTag         = "}"
	paramSeparator = ":"
)

const (
	TagTime        = "time"
	TagHeader      = "header:"
	TagStatus      = "status"
	TagOk          = "ok"
	TagQueryParams = "queryParams"
)

type Config struct {
	Format string
	Output io.Writer
}

// TODO
func createLogFuncMap() map[string]LogFunc {
	return map[string]LogFunc{
		TagHeader: func(output bytes.Buffer, response Response) (int, error) {
			return output.WriteString(response.Header.String())
		},
		//TagTime: func(output bytes.Buffer, response request.Response) (int, error) {
		//	return output.WriteString(response.Time().String())
		//},
		//TagStatus: func(output bytes.Buffer, response request.Response) (int, error) {
		//	return output.WriteString(response.StatusCode)
		//},
		TagOk: func(output bytes.Buffer, response Response) (int, error) {
			return output.WriteString("Ok:" + strconv.FormatBool(response.OK()))
		},
		TagQueryParams: func(output bytes.Buffer, response Response) (int, error) {
			return output.WriteString(response.RequestURI.Query())
		},
	}
}

type LogFunc func(output bytes.Buffer, response Response) (int, error)

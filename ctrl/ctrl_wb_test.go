package ctrl

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/test"
	"encoding/json"
)

func TestWriteError(t *testing.T) {
	// given
	ctrlErr := &ctrlError{errorMsg: "some error message", httpCode: 500}
	writer := &test.HttpWriterMock{}
	writer.Head = make(map[string][]string)

	// when
	writeError(writer, ctrlErr)

	// then
	var errMsg ErrorBody
	json.Unmarshal(writer.Data, &errMsg)
	if writer.Header().Get("Content-Type") != "application/json" {
		t.Error("Bad content type")
	} else if errMsg.Msg != "some error message" {
		t.Error("Bad error message")
	} else if writer.Code != 500 {
		t.Error("bad http code")
	}
}

package ctrl

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/utils"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)


func TestWriteError(t *testing.T) {
	// given
	ctrlErr := &ctrlError{errorMsg:"some error message", httpCode:500}
	writer := &utils.HttpWriterMock{}
	writer.Head = make(map[string][]string)

	// when
	writeError(writer, ctrlErr)

	// then
	var errMsg ErrorBody
	json.Unmarshal(writer.Data, &errMsg)
	assert.Equal(t, "application/json", writer.Header().Get("Content-Type"))
	assert.Equal(t, "some error message", errMsg.Msg)
	assert.Equal(t, 500, writer.Code)
}

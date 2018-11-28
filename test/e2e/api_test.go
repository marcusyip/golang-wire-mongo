package e2e

import (
	"bytes"
	"encoding/json"
	"io"
)

// func testAPI(engine *gin.Engine, w http.ResponseWriter, method, url string, body map[string]interface{}) {
// 	body := json.Marshal(body)
// 	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
// 	ctn.ApiRouter.With(ctn.Engine)
// 	ctn.Engine.ServeHTTP(w, req)
// }

func jsonBody(data interface{}) io.Reader {
	jsonStr, _ := json.Marshal(data)
	return bytes.NewBuffer(jsonStr)
}

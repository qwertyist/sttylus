package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func printRequestBody(r *http.Request) {
	if !debug {
		return 
	}
k$
	var bodyBytes []byte
	var err error

	if r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
			return
		}
		defer r.Body.Close()
	}

	fmt.Printf("Headers: %+v\n", r.Header)

	if len(bodyBytes) > 0 {
		var prettyJSON bytes.Buffer
		if err = json.Indent(&prettyJSON, bodyBytes, "", "\t"); err != nil {
			fmt.Printf("JSON parse error: %v", err)
			return
		}
		fmt.Println(string(prettyJSON.Bytes()))
	} else {
		fmt.Printf("Body: No Body Supplied\n")
	}

	// CORS headers
	//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	//w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	//w.Header().Set("Access-Control-Allow-Headers", "Accept-Encoding,Authorization,X-Forwarded-For,Content-Type,Origin,Server")
}

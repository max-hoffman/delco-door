package delcodoor

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", defaultHandler)
	// log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

const response = `
    <?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Say voice="woman">Please leave a message after the tone.</Say>
	</Response>
`

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	fmt.Fprint(w, response)
}

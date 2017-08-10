package delcodoor

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", defaultHandler)
	// log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

const (
	welcome = `<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Say voice="woman">It is Wednesday my dudes.</Say>
		<Record timeout="2"/>
	</Response>
`
	echo = `<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Play>%s</Play>
	</Response>
`
	accept = `<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Say voice="woman">Welcome to Delco.</Say>
		<Play digits="9"></Play>
	</Response>
`
	reject = `<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Say voice="woman">Incorrect password.</Say>
	</Response>
`
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")

	recording := r.FormValue("RecordingUrl")
	if recording == "" {
		fmt.Fprint(w, welcome)
		return
	}

	context := appengine.NewContext(r)
	text, err := transcribe(context, recording)
	if err != nil {
		http.Error(w, "could not transcribe recording", http.StatusInternalServerError)
		log.Errorf(context, "could not transcribe %v", recording)
		return
	}

	if text == "squidward" {
		fmt.Print(w, accept)
		return
	}
	fmt.Fprint(w, reject)
}

func transcribe(c context.Context, url string) (string, error) {
	bytes, err := fetchAudio(c, url)
	if err != nil {
		return "", err
	}

	return fetchTranscription(c, bytes)
}

func fetchAudio(c context.Context, url string) ([]byte, error) {
	//download url
	client := urlfetch.Client(c)
	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not download %v: %v", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status: %v", res.Status)
	}

	//base64 encode
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	return body, nil
}

var speechURL = " https://speech.googleapis.com/v1/speech:recognize?key=" + os.Getenv("SPEECH_API_KEY")

func fetchTranscription(c context.Context, bytes []byte) (string, error) {

}

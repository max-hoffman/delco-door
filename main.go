package delcodoor

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

var speechURL = " https://speech.googleapis.com/v1/speech:recognize?key=" +
	os.Getenv("SPEECH_API_KEY")

type speechRequest struct {
	Audio struct {
		Content string `json:"content"`
	} `json:"config"`
	Config struct {
		Encoding        string `json:"encoding"`
		SampleRateHertz int    `json:"sampleRateHertz"`
		LanguageCode    string `json:"languageCode"`
	} `json:"audio"`
}

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
	client := urlfetch.Client(c)
	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not download %v: %v", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status: %v", res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	return body, nil
}

func fetchTranscription(c context.Context, b []byte) (string, error) {
	var req speechRequest
	req.Audio.Content = base64.StdEncoding.EncodeToString(b)
	req.Config.Encoding = "LINEAR16"
	req.Config.SampleRateHertz = 8000
	req.Config.LanguageCode = "en-US"

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to encode speechrequest as json string: %v", err)
	}

	result, err := urlfetch.Client(c).Post(speechURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("speech api post failed with error: %v", err)
	}

	var data struct {
		Error struct {
			Code    int
			Message string
			Status  string
		}
		Results []struct {
			Alternatives []struct {
				Transcript string
				Confidence float64
			}
		}
	}

	defer result.Body.Close()
	if err := json.NewDecoder(result.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode speech response with error: %v", err)
	}
	if data.Error.Code != 0 {
		return "", fmt.Errorf("speech API failed: %d %s %s",
			data.Error.Code, data.Error.Status, data.Error.Message)
	}
	if len(data.Results) == 0 || len(data.Results[0].Alternatives) == 0 {
		return "", fmt.Errorf("no transcriptions found")
	}
	return data.Results[0].Alternatives[0].Transcript, nil
}

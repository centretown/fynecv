package ui

import (
	"fmt"
	"fynecv/appdata"
	"io"
	"log"
	"net/http"

	"fyne.io/fyne/v2/widget"
)

// req := fmt.Sprintf("http://192.168.0.7:9000/record?duration=%d", RecordDuration)

func NewRecordButton(data *appdata.AppData) *widget.Button {
	const RecordDuration = 300
	var (
		err          error
		resp         *http.Response
		isRecording  bool
		recordButton = widget.NewButtonWithIcon(msgStartRecording, MotionOnIcon, func() {})
	)

	recordButton.OnTapped = func() {
		if isRecording {
			recordButton.SetIcon(MotionOnIcon)
			recordButton.SetText(msgStartRecording)
			resp, err = http.Get("http://192.168.0.7:9000/record?duration=0")
		} else {
			recordButton.SetIcon(MotionOffIcon)
			recordButton.SetText(msgStopRecording)
			req := fmt.Sprintf("http://192.168.0.7:9000/record?duration=%d", RecordDuration)
			resp, err = http.Get(req)
		}

		if err != nil {
			log.Println("http.Get", err)
			return
		}

		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(string(buf))
			log.Println("ReadAll", err)
			return
		}

		isRecording = !isRecording
		recordButton.Refresh()
	}

	return recordButton
}

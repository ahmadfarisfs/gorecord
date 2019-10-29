package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	r "github.com/ahmadfarisfs/gorecord"

	"gocv.io/x/gocv"
)

type Config struct {
	Port     int `json:"port"`
	CameraID int `json:"camera_id"`
	Width    int `json:"width"`
	Height   int `json:"height"`
	FPS      int `json:"fps"`
}

func main() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	config := Config{}

	fmt.Println("Successfully Opened config.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	errorUnmarshal := json.Unmarshal(byteValue, &config)
	if errorUnmarshal != nil {
		log.Fatalln("error unmarshal " + errorUnmarshal.Error())
		return
	}

	dev := r.VideoRecorder{
		DeviceID: uint(config.CameraID),
		Width:    uint(config.Width),
		Height:   uint(config.Height),
		FPS:      uint(config.FPS),
	}
	errO := dev.Open()
	if errO != nil {
		log.Println(errO)
	}

	server := TheServer{
		Port:     config.Port,
		Recorder: dev,
	}

	server.Serve()

}

func mainvv() {
	// set to use a video capture device 0
	deviceID := 0
	videoWidth := 640
	videoHeight := 480
	videoFps := 60
	filename := "test"
	format := ".avi"

	//sleep time
	sleepTimeMs := float32(1.0/float32(videoFps)) * 1000.0

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	webcam.Set(gocv.VideoCaptureFrameWidth, float64(videoWidth))
	webcam.Set(gocv.VideoCaptureFrameHeight, float64(videoHeight))
	webcam.Set(gocv.VideoCaptureFPS, float64(videoFps))

	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	writer, errW := gocv.VideoWriterFile(filename+format, "XVID", float64(videoFps), videoWidth, videoHeight, true)
	if errW != nil {
		fmt.Println(errW)
		return
	}
	defer writer.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	fmt.Printf("Start reading camera device: %v\n", deviceID)
	startTime := time.Now()

	go func() {
		for {
			if ok := webcam.Read(&img); !ok {
				fmt.Printf("cannot read device %v\n", deviceID)
				return
			}
			if time.Now().After(startTime.Add(time.Duration(6) * time.Second)) {
				//stop recording
				break
			}
		}
	}()

	for {
		err = writer.Write(img)
		if err != nil {
			fmt.Println(err)
		}

		//calculate sleep time based on fps
		time.Sleep(time.Duration(sleepTimeMs) * time.Millisecond)

		if time.Now().After(startTime.Add(time.Duration(5) * time.Second)) {
			//stop recording
			break
		}
	}
}

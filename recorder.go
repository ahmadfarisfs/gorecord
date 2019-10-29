package recorder

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gocv.io/x/gocv"
)

type VideoRecorder struct {
	DeviceID    uint
	Width       uint
	Height      uint
	FPS         uint
	isOpened    bool
	isRecording bool
	camera      *gocv.VideoCapture
	rec         *gocv.VideoWriter
	img         *gocv.Mat
}

func (v *VideoRecorder) Open() error {
	if v.isOpened {
		log.Println("Device already opened")
		return nil
	}

	webcam, err := gocv.OpenVideoCapture(int(v.DeviceID))
	if err != nil {
		log.Println(err)
		return err
	}
	webcam.Set(gocv.VideoCaptureFrameWidth, float64(v.Width))
	webcam.Set(gocv.VideoCaptureFrameHeight, float64(v.Height))
	webcam.Set(gocv.VideoCaptureFPS, float64(v.FPS))

	imgs := gocv.NewMat()

	v.img = &imgs
	v.camera = webcam
	v.isOpened = true

	go func() {
		for {
			if ok := v.camera.Read(v.img); !ok {
				log.Println("Cannot read device ", v.DeviceID)
				return
			}
			if !v.isOpened {
				log.Println("Closing")
				return
			}
		}
	}()

	return nil

}

func (v *VideoRecorder) StartRecord() error {
	if v.isRecording {
		log.Println("Already recording")
		return nil
	}
	writer, errW := gocv.VideoWriterFile("temp_frs.avi", "XVID", float64(v.FPS), int(v.Width), int(v.Height), true)
	if errW != nil {
		fmt.Println(errW)
		return errW
	}
	v.rec = writer
	v.isRecording = true
	sleepTimeMs := float32(1.0/float32(v.FPS)) * 1000.0
	go func() {
		for {
			if !v.isRecording {
				log.Println("Stopping recording")
				v.rec.Close()
				return
			}
			err := writer.Write(*v.img)
			if err != nil {
				log.Println(err)
			}
			//calculate sleep time based on fps
			time.Sleep(time.Duration(sleepTimeMs) * time.Millisecond)
		}
	}()
	return nil
}

func (v *VideoRecorder) StopRecord(filename string) error {
	if !v.isRecording {
		log.Println("Already stopping")
		return nil
	}
	v.isRecording = false

	//Wait for recorder to finish its job
	timeStart := time.Now()
	for {
		if !v.rec.IsOpened() {
			break
		}
		if time.Now().After(timeStart.Add(time.Duration(5) * time.Second)) {
			//timeout
			return errors.New("Timeout waiting writer to be closed")
		}
	}

	//wait for unfinished process of compressing file
	timeStart = time.Now()
	for {
		errRename := os.Rename("temp_frs.avi", "result"+string(os.PathSeparator)+filename)
		if errRename == nil {
			break
		}
		if time.Now().After(timeStart.Add(time.Duration(5) * time.Second)) {
			//timeout
			return errors.New("Timeout waiting OS to release file")
		}

	}
	return nil
}

func (v *VideoRecorder) Close() error {
	if !v.isOpened {
		log.Println("Device already closed")
		return nil
	}

	v.isOpened = false
	v.camera.Close()
	v.rec.Close()
	return nil
}

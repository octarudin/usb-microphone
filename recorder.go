package usbmicrophone

import (
	"fmt"

	"github.com/gen2brain/malgo"
)

// Recorder mengelola status perekaman audio
type Recorder struct {
	ctx *malgo.AllocatedContext
}

// NewRecorder menginisialisasi konteks audio
func NewRecorder() (*Recorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal inisialisasi audio context: %w", err)
	}
	return &Recorder{ctx: ctx}, nil
}

// StartListening membaca input dari mikrofon USB/Default selama durasi tertentu
func (r *Recorder) StartListening(onDataCallback func(pSample []byte)) error {
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 44100
	deviceConfig.Alsa.NoMMap = 1

	// Callback saat data audio masuk dari mikrofon
	onRec := func(pOutputSample, pInputSample []byte, frameCount uint32) {
		if len(pInputSample) > 0 {
			onDataCallback(pInputSample)
		}
	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRec,
	}

	device, err := malgo.InitDevice(r.ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		return fmt.Errorf("gagal menginisialisasi perangkat audio: %w", err)
	}

	err = device.Start()
	if err != nil {
		return fmt.Errorf("gagal menjalankan mikrofon: %w", err)
	}

	fmt.Println("Mikrofon aktif! Dengarkan audio...")
	return nil
}

// Close membersihkan resource audio
func (r *Recorder) Close() {
	if r.ctx != nil {
		_ = r.ctx.Uninit()
		r.ctx.Free()
	}
}

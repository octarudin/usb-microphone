package usbmicrophone

import (
	"fmt"

	"github.com/gen2brain/malgo"
)

// Config baru di v2.0.0 untuk fleksibilitas pengaturan audio
type Config struct {
	SampleRate uint32 // contoh: 44100 atau 48000
	Channels   uint32 // contoh: 1 (Mono) atau 2 (Stereo)
}

type Recorder struct {
	ctx    *malgo.AllocatedContext
	config Config // menyimpan konfigurasi baru
}

// BREAKING CHANGE: NewRecorder sekarang wajib menerima parameter Config
func NewRecorder(cfg Config) (*Recorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal inisialisasi audio context: %w", err)
	}
	return &Recorder{
		ctx:    ctx,
		config: cfg,
	}, nil
}

func (r *Recorder) StartListening(onDataCallback func(pSample []byte)) error {
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16

	// Menggunakan konfigurasi dinamis dari struct v2
	deviceConfig.Capture.Channels = r.config.Channels
	deviceConfig.SampleRate = r.config.SampleRate
	deviceConfig.Alsa.NoMMap = 1

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

	return device.Start()
}

func (r *Recorder) Close() {
	if r.ctx != nil {
		_ = r.ctx.Uninit()
		r.ctx.Free()
	}
}

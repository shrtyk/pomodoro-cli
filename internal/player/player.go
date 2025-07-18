package player

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

const (
	sampleRate beep.SampleRate = 48000
)

// Player defines the interface for playing audio alerts.
type Player interface {
	PlaySoonAlert()
	PlayLoopAlert()
	PlayRoundAlert()
	Close()
}

// player handles the playback of different audio alerts.
type player struct {
	soonBuffer  *beep.Buffer
	loopBuffer  *beep.Buffer
	roundBuffer *beep.Buffer

	soonClose  func() error
	loopClose  func() error
	roundClose func() error
}

// NewPlayer creates a new Player instance.
// It decodes the audio files into buffers for playback.
func NewPlayer(soonAlertPath, loopAlertPath, roundAlertPath string) (Player, error) {
	soonBuffer, soonFile, err := bufferedAlert(soonAlertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to buffer file: %w", err)
	}

	loopBuffer, loopFile, err := bufferedAlert(loopAlertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to buffer file: %w", err)
	}

	roundBuffer, roundFile, err := bufferedAlert(roundAlertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to buffer file: %w", err)
	}

	// Initialize the speaker with the sample rate of the first audio file.
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	return &player{
		soonBuffer:  soonBuffer,
		loopBuffer:  loopBuffer,
		roundBuffer: roundBuffer,
		soonClose:   soonFile.Close,
		loopClose:   loopFile.Close,
		roundClose:  roundFile.Close,
	}, nil
}

// PlaySoonAlert plays the 'soon' alert sound.
func (p *player) PlaySoonAlert() {
	speaker.PlayAndWait(p.soonBuffer.Streamer(0, p.soonBuffer.Len()))
}

// PlayLoopAlert plays the 'done' alert sound.
func (p *player) PlayLoopAlert() {
	speaker.PlayAndWait(p.loopBuffer.Streamer(0, p.loopBuffer.Len()))
}

// PlayRoundAlert plays the 'round' alert sound.
func (p *player) PlayRoundAlert() {
	speaker.PlayAndWait(p.roundBuffer.Streamer(0, p.roundBuffer.Len()))
}

// Close closes the audio speaker.
func (p *player) Close() {
	speaker.Close()
	p.soonClose()
	p.loopClose()
	p.roundClose()
}

// bufferedAlert returns buffered alert file ready to be played.
func bufferedAlert(fileName string) (*beep.Buffer, *os.File, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open '%s': %w", fileName, err)
	}

	s, f, err := decodeFile(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode file '%s': %w", file.Name(), err)
	}
	defer s.Close()

	resampled := beep.Resample(4, f.SampleRate, sampleRate, s)
	buf := beep.NewBuffer(
		beep.Format{
			SampleRate:  sampleRate,
			NumChannels: f.NumChannels,
			Precision:   f.Precision,
		},
	)
	buf.Append(resampled)
	return buf, file, nil
}

// decodeFile decodes an audio file (mp3 or wav)
func decodeFile(file *os.File) (beep.StreamSeekCloser, beep.Format, error) {
	ext := filepath.Ext(file.Name())
	if len(ext) == 0 {
		file.Close()
		return nil, beep.Format{}, fmt.Errorf("wrong filename: %s", file.Name())
	}

	switch ext {
	case ".mp3":
		streamer, format, err := mp3.Decode(file)
		if err != nil {
			file.Close()
			return nil, beep.Format{}, fmt.Errorf("failed to decode mp3: %w", err)
		}
		return streamer, format, nil
	case ".wav":
		streamer, format, err := wav.Decode(file)
		if err != nil {
			file.Close()
			return nil, beep.Format{}, fmt.Errorf("failed to decode wav: %w", err)
		}
		return streamer, format, nil
	default:
		file.Close()
		return nil, beep.Format{}, fmt.Errorf("unknown file extension: %s", ext)
	}
}

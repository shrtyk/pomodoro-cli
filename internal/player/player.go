package player

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
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
	doneBuffer  *beep.Buffer
	roundBuffer *beep.Buffer
}

// NewPlayer creates a new Player instance.
// It decodes the audio files into buffers for playback.
func NewPlayer(soonAlertPath, loopAlertPath, roundAlertPath string) (Player, error) {
	soonBuffer, err := bufferedAlert(soonAlertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to buffer file: %w", err)
	}

	loopBuffer, err := bufferedAlert(loopAlertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to buffer file: %w", err)
	}

	roundBuffer, err := bufferedAlert(roundAlertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to buffer file: %w", err)
	}

	// Initialize the speaker with the sample rate of the first audio file.
	speaker.Init(
		soonBuffer.Format().SampleRate,
		soonBuffer.Format().SampleRate.N(time.Second/10),
	)

	return &player{
		soonBuffer:  soonBuffer,
		doneBuffer:  loopBuffer,
		roundBuffer: roundBuffer,
	}, nil
}

// PlaySoonAlert plays the 'soon' alert sound.
func (p *player) PlaySoonAlert() {
	speaker.PlayAndWait(p.soonBuffer.Streamer(0, p.soonBuffer.Len()))
}

// PlayLoopAlert plays the 'done' alert sound.
func (p *player) PlayLoopAlert() {
	speaker.PlayAndWait(p.doneBuffer.Streamer(0, p.doneBuffer.Len()))
}

// PlayRoundAlert plays the 'round' alert sound.
func (p *player) PlayRoundAlert() {
	speaker.PlayAndWait(p.roundBuffer.Streamer(0, p.roundBuffer.Len()))
}

// Close closes the audio speaker.
func (p *player) Close() {
	speaker.Close()
}

// bufferedAlert returns buffered alert file ready to be played.
func bufferedAlert(fileName string) (*beep.Buffer, error) {
	s, f, err := decodeFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file '%s': %w", fileName, err)
	}
	defer s.Close()
	buf := beep.NewBuffer(f)
	buf.Append(s)
	return buf, nil
}

// decodeFile decodes an audio file (mp3 or wav)
func decodeFile(fileName string) (beep.StreamSeekCloser, beep.Format, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("failed to open '%s': %w", fileName, err)
	}

	s := strings.Split(fileName, ".")
	if len(s) < 2 {
		file.Close()
		return nil, beep.Format{}, fmt.Errorf("wrong filename: %s", fileName)
	}

	ext := s[len(s)-1]
	switch ext {
	case "mp3":
		streamer, format, err := mp3.Decode(file)
		if err != nil {
			file.Close()
			return nil, beep.Format{}, fmt.Errorf("failed to decode mp3: %w", err)
		}
		return streamer, format, nil
	case "wav":
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

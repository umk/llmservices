package audio

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createValidWAV creates a minimal valid WAV file as a byte slice
func createValidWAV() []byte {
	// Simple mono 16-bit PCM at 44100Hz
	const (
		numChannels   = 1
		sampleRate    = 44100
		bitsPerSample = 16
		numSamples    = 10 // Just a short sample
	)

	bytesPerSample := bitsPerSample / 8
	dataSize := numSamples * numChannels * bytesPerSample
	fileSize := 36 + dataSize // 36 bytes header + data

	// Create buffer
	buf := make([]byte, 44+dataSize)

	// RIFF header
	copy(buf[0:4], "RIFF")
	binary.LittleEndian.PutUint32(buf[4:8], uint32(fileSize))
	copy(buf[8:12], "WAVE")

	// Format chunk
	copy(buf[12:16], "fmt ")
	binary.LittleEndian.PutUint32(buf[16:20], 16) // format chunk size
	binary.LittleEndian.PutUint16(buf[20:22], 1)  // PCM format
	binary.LittleEndian.PutUint16(buf[22:24], numChannels)
	binary.LittleEndian.PutUint32(buf[24:28], sampleRate)
	byteRate := sampleRate * uint32(numChannels) * uint32(bytesPerSample)
	binary.LittleEndian.PutUint32(buf[28:32], byteRate)
	blockAlign := numChannels * bytesPerSample
	binary.LittleEndian.PutUint16(buf[32:34], uint16(blockAlign))
	binary.LittleEndian.PutUint16(buf[34:36], bitsPerSample)

	// Data chunk
	copy(buf[36:40], "data")
	binary.LittleEndian.PutUint32(buf[40:44], uint32(dataSize))

	// Add some sample data
	for i := 0; i < numSamples; i++ {
		binary.LittleEndian.PutUint16(buf[44+i*bytesPerSample:], uint16(i*1000))
	}

	return buf
}

func TestParseWAV(t *testing.T) {
	t.Run("Valid WAV", func(t *testing.T) {
		wavData := createValidWAV()

		audio, err := ParseWAV(wavData)

		assert.NoError(t, err)
		assert.Equal(t, uint16(1), audio.Metadata.Format)         // PCM
		assert.Equal(t, uint16(1), audio.Metadata.Channels)       // Mono
		assert.Equal(t, uint32(44100), audio.Metadata.SampleRate) // 44.1kHz
		assert.Equal(t, uint16(16), audio.Metadata.BitsPerSample) // 16-bit
		assert.Equal(t, 20, len(audio.Data.Buf))                  // 10 samples * 2 bytes per sample
		assert.Equal(t, 2, audio.Data.BytesPerFrame)              // 16-bit mono = 2 bytes per frame
	})

	t.Run("Too Short", func(t *testing.T) {
		wavData := make([]byte, 40) // Less than minimum 44 bytes

		_, err := ParseWAV(wavData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file too short")
	})

	t.Run("Invalid RIFF Header", func(t *testing.T) {
		wavData := createValidWAV()
		copy(wavData[0:4], "ABCD") // Corrupt RIFF header

		_, err := ParseWAV(wavData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing RIFF header")
	})

	t.Run("Invalid WAVE Format", func(t *testing.T) {
		wavData := createValidWAV()
		copy(wavData[8:12], "ABCD") // Corrupt WAVE format

		_, err := ParseWAV(wavData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing WAVE format")
	})

	t.Run("Missing Format Chunk", func(t *testing.T) {
		wavData := createValidWAV()
		copy(wavData[12:16], "ABCD") // Corrupt format chunk ID

		_, err := ParseWAV(wavData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing \"fmt \" chunk")
	})

	t.Run("Format Chunk Too Small", func(t *testing.T) {
		wavData := createValidWAV()
		binary.LittleEndian.PutUint32(wavData[16:20], 8) // Set format chunk size too small

		_, err := ParseWAV(wavData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "format chunk too small")
	})

	t.Run("Unsupported WAV Format", func(t *testing.T) {
		wavData := createValidWAV()
		binary.LittleEndian.PutUint16(wavData[20:22], 2) // Change from PCM (1) to another format

		_, err := ParseWAV(wavData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported WAV format")
	})

	t.Run("Missing Data Chunk", func(t *testing.T) {
		wavData := createValidWAV()
		copy(wavData[36:40], "ABCD") // Corrupt data chunk ID

		_, err := ParseWAV(wavData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing \"data\" chunk")
	})

	t.Run("Data Chunk Too Large", func(t *testing.T) {
		wavData := createValidWAV()
		binary.LittleEndian.PutUint32(wavData[40:44], uint32(1000000)) // Set data chunk size too large

		_, err := ParseWAV(wavData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "data chunk exceeds file size")
	})

	t.Run("Zero Bytes Per Frame", func(t *testing.T) {
		wavData := createValidWAV()
		// Set both channels and bits per sample to 0 to make bytes per frame = 0
		binary.LittleEndian.PutUint16(wavData[22:24], 0) // channels
		binary.LittleEndian.PutUint16(wavData[34:36], 0) // bits per sample

		_, err := ParseWAV(wavData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "zero or negative bytes per frame")
	})
}

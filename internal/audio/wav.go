package audio

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/umk/llmservices/pkg/adapter"
)

// ParseWAV parses a WAV file from a byte slice and returns metadata and PCM data
func ParseWAV(buf []byte) (adapter.Audio, error) {
	if len(buf) < 44 { // Minimum WAV header size
		return adapter.Audio{}, errors.New("invalid WAV: file too short")
	}

	// Verify RIFF header
	if string(buf[0:4]) != "RIFF" {
		return adapter.Audio{}, errors.New("invalid WAV: missing RIFF header")
	}
	if string(buf[8:12]) != "WAVE" {
		return adapter.Audio{}, errors.New("invalid WAV: missing WAVE format")
	}

	// Find and process format chunk
	metadata, err := getWAVMetadata(buf)
	if err != nil {
		return adapter.Audio{}, err
	}

	// Find and process data chunk
	data, err := getWAVData(buf, metadata)
	if err != nil {
		return adapter.Audio{}, err
	}

	return adapter.Audio{
		Metadata: metadata,
		Data:     data,
	}, nil
}

// getWAVChunk searches for a chunk with the given ID in WAV data
// Returns the offset of the chunk and the chunk size, or an error if not found
func getWAVChunk(data []byte, chunkID string) (offset int, size int, err error) {
	offset = 12 // Start after the RIFF header and WAVE ID

	for {
		if offset+8 > len(data) {
			return 0, 0, fmt.Errorf("invalid WAV: missing %q chunk", chunkID)
		}

		currentChunkID := string(data[offset : offset+4])
		currentChunkSize := int(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))

		if currentChunkID == chunkID {
			return offset, currentChunkSize, nil
		}

		offset += 8 + currentChunkSize
	}
}

// getWAVMetadata finds and processes the format chunk and returns WAV metadata
func getWAVMetadata(data []byte) (adapter.AudioMetadata, error) {
	// Find 'fmt ' chunk
	fmtOffset, fmtChunkSize, err := getWAVChunk(data, "fmt ")
	if err != nil {
		return adapter.AudioMetadata{}, err
	}

	if fmtChunkSize < 16 || fmtOffset+8+fmtChunkSize > len(data) {
		return adapter.AudioMetadata{}, errors.New("invalid WAV: format chunk too small")
	}

	metadata := adapter.AudioMetadata{
		Format:        binary.LittleEndian.Uint16(data[fmtOffset+8 : fmtOffset+10]),
		Channels:      binary.LittleEndian.Uint16(data[fmtOffset+10 : fmtOffset+12]),
		SampleRate:    binary.LittleEndian.Uint32(data[fmtOffset+12 : fmtOffset+16]),
		ByteRate:      binary.LittleEndian.Uint32(data[fmtOffset+16 : fmtOffset+20]),
		BlockAlign:    binary.LittleEndian.Uint16(data[fmtOffset+20 : fmtOffset+22]),
		BitsPerSample: binary.LittleEndian.Uint16(data[fmtOffset+22 : fmtOffset+24]),
	}

	// Only support PCM format (1)
	if metadata.Format != 1 {
		return adapter.AudioMetadata{}, fmt.Errorf("unsupported WAV format: %d, only PCM (1) supported", metadata.Format)
	}

	return metadata, nil
}

// getWAVData finds and processes the data chunk and returns PCM data
func getWAVData(data []byte, metadata adapter.AudioMetadata) (adapter.AudioData, error) {
	// Find 'data' chunk
	dataOffset, dataChunkSize, err := getWAVChunk(data, "data")
	if err != nil {
		return adapter.AudioData{}, err
	}

	dataSize := uint32(dataChunkSize)
	if dataOffset+8+int(dataSize) > len(data) {
		return adapter.AudioData{}, errors.New("invalid WAV: data chunk exceeds file size")
	}

	bytesPerFrame := int(metadata.Channels) * int(metadata.BitsPerSample) / 8
	if bytesPerFrame <= 0 {
		return adapter.AudioData{}, errors.New("invalid WAV: zero or negative bytes per frame")
	}

	pcmData := data[dataOffset+8 : dataOffset+8+int(dataSize)]

	return adapter.AudioData{
		Data:          pcmData,
		BytesPerFrame: bytesPerFrame,
		Size:          dataSize,
	}, nil
}

package adapter

// Audio represents a parsed WAV file with metadata and raw PCM data
type Audio struct {
	Buf      []byte        `json:"-"`        // Raw byte buffer containing the entire WAV file content
	Metadata AudioMetadata `json:"metadata"` // Format information extracted from the WAV file header
	Data     AudioData     `json:"data"`     // Raw audio data and related frame information
}

// AudioMetadata holds information about the WAV file format
type AudioMetadata struct {
	Format        uint16 `json:"format"`          // Format of the audio (1 = PCM, other values indicate compression)
	Channels      uint16 `json:"channels"`        // Number of audio channels (1 for mono, 2 for stereo, etc.)
	SampleRate    uint32 `json:"sample_rate"`     // Number of samples per second (e.g., 44100, 48000)
	ByteRate      uint32 `json:"byte_rate"`       // Average bytes per second (SampleRate * NumChannels * BitsPerSample/8)
	BlockAlign    uint16 `json:"block_align"`     // Size of a single sample across all channels (NumChannels * BitsPerSample/8)
	BitsPerSample uint16 `json:"bits_per_sample"` // Bits per sample per channel (e.g., 8, 16, 24)
}

// AudioData holds PCM audio data and related properties
type AudioData struct {
	Buf           []byte `json:"data"`            // Raw PCM audio data bytes
	BytesPerFrame int    `json:"bytes_per_frame"` // Number of bytes in each audio frame (NumChannels * BitsPerSample/8)
	Size          uint32 `json:"size"`            // Total size of the PCM data in bytes
}

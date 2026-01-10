package objects

import (
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    string
	}{
		{
			name:    "empty content",
			content: []byte{},
			want:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:    "hello world",
			content: []byte("hello world"),
			want:    "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			name:    "binary data",
			content: []byte{0x00, 0x01, 0x02, 0xFF},
			want:    "3d1f57c984978ef98a18378c8166c1cb8ede02c03eeb6aee7e2f121dfeee3e56",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Hash(tt.content)
			if got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashConsistency(t *testing.T) {
	content := []byte("test content for consistency check")

	hash1 := Hash(content)
	hash2 := Hash(content)

	if hash1 != hash2 {
		t.Errorf("Hash() not consistent: %v != %v", hash1, hash2)
	}
}

func TestHashLength(t *testing.T) {
	hash := Hash([]byte("anything"))

	if len(hash) != HashLength {
		t.Errorf("Hash length = %d, want %d", len(hash), HashLength)
	}
}

func TestHashReader(t *testing.T) {
	content := "hello world"
	reader := strings.NewReader(content)

	got, err := HashReader(reader)
	if err != nil {
		t.Fatalf("HashReader() error = %v", err)
	}

	want := Hash([]byte(content))
	if got != want {
		t.Errorf("HashReader() = %v, want %v", got, want)
	}
}

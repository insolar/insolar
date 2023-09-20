package foundation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractCanonicalPublicKey(t *testing.T) {
	type args struct {
		pk string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "happy",
			args:    args{pk: "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEwDcgWZ1SbG+nbiXZkmYUZEfk2nkk\n1PEmEWoj4g6DLEkdaQVorOkqlloEz1zXclQaAE1S8i3F7OFNrNxLkm34ow==\n-----END PUBLIC KEY-----\n"},
			want:    "A8A3IFmdUmxvp24l2ZJmFGRH5Np5JNTxJhFqI-IOgyxJ",
			wantErr: false,
		},
		{
			name:    "wrong pk",
			args:    args{pk: "-----BEGIN PUBLIC KEY-----\nasdnjkDFHaldfjl==\n-----END PUBLIC KEY-----\n"},
			wantErr: true,
		},
		{
			name:    "wrong dsa key",
			args:    args{pk: "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCF5iQBRQAijZo6V83VpA+u6c1i\nfzOLYNFWYTk72VK+W/m9DAvyDe2LJCX7kRq3hUhkQpR+YyfMJuNmCCFpz4/IfrfN\n/GdtNHlcmJU6f0hHE+CzxbY2yptXBLZpyg7Ll4vXHGD4WEbRBTzc8CW6L5kS5kJ5\ni2pwohrbRVBgfXkkmQIDAQAB\n-----END PUBLIC KEY-----"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractCanonicalPublicKey(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractCanonicalPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractCanonicalPublicKey() got = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("ok with additional fields in pem", func(t *testing.T) {
		pk1 := "-----BEGIN PUBLIC KEY-----\nThisIsNewField:testvalue\nThisIsNewField:testvalue\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEwDcgWZ1SbG+nbiXZkmYUZEfk2nkk\n1PEmEWoj4g6DLEkdaQVorOkqlloEz1zXclQaAE1S8i3F7OFNrNxLkm34ow==\n-----END PUBLIC KEY-----\n"
		pk2 := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEwDcgWZ1SbG+nbiXZkmYUZEfk2nkk\n1PEmEWoj4g6DLEkdaQVorOkqlloEz1zXclQaAE1S8i3F7OFNrNxLkm34ow==\n-----END PUBLIC KEY-----\n"
		got1, _ := ExtractCanonicalPublicKey(pk1)
		got2, _ := ExtractCanonicalPublicKey(pk2)
		require.Equal(t, got1, got2)
	})

	t.Run("ok with compressed/uncompressed pk odd", func(t *testing.T) {
		pk1 := "-----BEGIN PUBLIC KEY-----\nMDYwEAYHKoZIzj0CAQYFK4EEAAoDIgAC45SRdXuMWUnPEHu0VcJlP4Ws6qj0rZzx\nDlv/xlyMdmo=\n-----END PUBLIC KEY-----\n"
		pk2 := "-----BEGIN PUBLIC KEY-----\nMFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAE45SRdXuMWUnPEHu0VcJlP4Ws6qj0rZzx\nDlv/xlyMdmqGLpPEuqecYA3Gw9EWQ1I7CulTYQL4tMx0+zh14Rc27g==\n-----END PUBLIC KEY-----\n"
		got1, _ := ExtractCanonicalPublicKey(pk1)
		got2, _ := ExtractCanonicalPublicKey(pk2)
		require.Equal(t, got1, got2)
	})

	t.Run("ok with compressed/uncompressed pk even", func(t *testing.T) {
		pk1 := "-----BEGIN PUBLIC KEY-----\nMDYwEAYHKoZIzj0CAQYFK4EEAAoDIgADHHub97QxRqqlcWExT+5IWBXpKQ8lE7ih\nHCsfdiEwR80=\n-----END PUBLIC KEY-----\n"
		pk2 := "-----BEGIN PUBLIC KEY-----\nMFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEHHub97QxRqqlcWExT+5IWBXpKQ8lE7ih\nHCsfdiEwR80q/bdILdePYFDTc/uRcQ7dwAcxZdVwE8XvJ1s6k1vHVQ==\n-----END PUBLIC KEY-----\n"
		got1, _ := ExtractCanonicalPublicKey(pk1)
		got2, _ := ExtractCanonicalPublicKey(pk2)
		require.Equal(t, got1, got2)
	})
}

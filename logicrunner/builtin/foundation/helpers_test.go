// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package foundation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrimPublicKey(t *testing.T) {
	for _, tc := range []struct {
		input  string
		result string
	}{
		{
			input:  "asDafasf",
			result: "asDafasf",
		},
		{
			input:  "-----BEGIN RSA PUBLIC KEY-----\naSDafasf\n-----END RSA PUBLIC KEY-----",
			result: "aSDafasf",
		},
	} {
		require.Equal(t, tc.result, TrimPublicKey(tc.input))
	}
}

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
			want:    "A8A3IFmdUmxvp24l2ZJmFGRH5Np5JNTxJhFqI+IOgyxJ",
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
		pk1 := "-----BEGIN PUBLIC KEY-----\nThisIsNewField:testvalue\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEwDcgWZ1SbG+nbiXZkmYUZEfk2nkk\n1PEmEWoj4g6DLEkdaQVorOkqlloEz1zXclQaAE1S8i3F7OFNrNxLkm34ow==\n-----END PUBLIC KEY-----\n"
		pk2 := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEwDcgWZ1SbG+nbiXZkmYUZEfk2nkk\n1PEmEWoj4g6DLEkdaQVorOkqlloEz1zXclQaAE1S8i3F7OFNrNxLkm34ow==\n-----END PUBLIC KEY-----\n"
		got1, _ := ExtractCanonicalPublicKey(pk1)
		got2, _ := ExtractCanonicalPublicKey(pk2)
		require.Equal(t, got1, got2)
	})
}

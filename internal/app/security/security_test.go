package security

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDecrypt(t *testing.T) {
	secretKey := []byte("G0pher")
	myUUID := uuid.New()
	hashedUUID := Encrypt(myUUID, secretKey)

	type args struct {
		hashString string
		secret     []byte
	}
	type want struct {
		uuid    uuid.UUID
		wantErr error
		equal   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test #1",
			args: args{
				hashString: hashedUUID,
				secret:     secretKey,
			},
			want: want{
				uuid:    myUUID,
				wantErr: nil,
				equal:   true,
			},
		},
		{
			name: "positive test #2",
			args: args{
				hashString: hashedUUID + "1",
				secret:     secretKey,
			},
			want: want{
				wantErr: ErrNotValidSing,
				equal:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, decrErr := Decrypt(tt.args.hashString, tt.args.secret)

			var equalState bool
			if got == myUUID {
				equalState = true
			}
			assert.Equal(t, tt.want.uuid, got, fmt.Errorf("expected UUID %s, got %s", tt.want.uuid, got))
			assert.Equal(t, tt.want.equal, equalState, fmt.Errorf("expected state %t, got %t", tt.want.equal, equalState))
			assert.Equal(t, tt.want.wantErr, decrErr, fmt.Errorf("expected error %s, got %s", tt.want.wantErr, decrErr))
		})
	}
}

func BenchmarkCrypt(b *testing.B) {
	secretKey := []byte("G0pher")
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()

	b.Run("encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Encrypt(uuid.New(), secretKey)
		}
	})

	b.Run("decrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer() // останавливаем таймер
			hashedUUID := Encrypt(uuid.New(), secretKey)
			b.StartTimer() // возобновляем таймер
			Decrypt(hashedUUID, secretKey)
		}
	})
}

package tokens

import (
	"reflect"
	"testing"
)

func TestChunk(t *testing.T) {
	type args struct {
		t            string
		maxTokenSize int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "chunking not needed",
			args: args{
				t:            "hello world",
				maxTokenSize: 100,
			},
			want: []string{
				"hello world",
			},
			wantErr: false,
		},
		{
			name: "chunking is needed",
			args: args{
				t:            "This is a really long sentence, that will need to be chunked. Here is some more text.",
				maxTokenSize: 15,
			},
			want: []string{
				"This is a really long sentence, that will need to be chunked.",
				" Here is some more text.",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Chunk(tt.args.t, tt.args.maxTokenSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chunk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCount(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "count tokens",
			args: args{
				t: "How many tokens is this?",
			},
			want:    6,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Count(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CountTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

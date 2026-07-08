package admin

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	apperrors "app/pkg/errors"
	adminv1 "app/proto/gen/admin/v1"
)

func TestAI_RemindSpeechTranscribe(t *testing.T) {
	validAudio := base64.StdEncoding.EncodeToString([]byte("audio"))

	tests := []struct {
		name       string
		audio      string
		text       string
		err        error
		wantText   string
		wantCode   int
		wantReason string
	}{
		{
			name:     "success",
			audio:    validAudio,
			text:     "明天上午九点提醒我喝水。",
			wantText: "明天上午九点提醒我喝水。",
		},
		{
			name:       "invalid base64",
			audio:      "not-base64",
			wantCode:   400,
			wantReason: "INVALID_AUDIO_BASE64",
		},
		{
			name:       "empty transcript",
			audio:      validAudio,
			text:       "  ",
			wantCode:   400,
			wantReason: "EMPTY_TRANSCRIPT",
		},
		{
			name:       "transcriber error",
			audio:      validAudio,
			err:        errors.New("doubao unavailable"),
			wantCode:   500,
			wantReason: "AI_REMIND_TRANSCRIBE_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transcriber := &fakeRemindSpeechTranscriber{text: tt.text, err: tt.err}
			svc := &AI{speechTranscriber: transcriber}

			resp, err := svc.RemindSpeechTranscribe(context.Background(),
				adminv1.RemindSpeechTranscribeRequest_builder{AudioBase64: tt.audio}.Build())

			if tt.wantReason != "" {
				if err == nil {
					t.Fatal("expected error")
				}
				if got := apperrors.Code(err); got != tt.wantCode {
					t.Fatalf("unexpected code: got %d want %d", got, tt.wantCode)
				}
				if got := apperrors.Reason(err); got != tt.wantReason {
					t.Fatalf("unexpected reason: got %q want %q", got, tt.wantReason)
				}
				return
			}

			if err != nil {
				t.Fatalf("RemindSpeechTranscribe() error = %v", err)
			}
			if resp.GetText() != tt.wantText {
				t.Fatalf("unexpected text: got %q want %q", resp.GetText(), tt.wantText)
			}
			if transcriber.gotAudio != tt.audio {
				t.Fatalf("unexpected transcriber audio: got %q want %q", transcriber.gotAudio, tt.audio)
			}
		})
	}
}

type fakeRemindSpeechTranscriber struct {
	text     string
	err      error
	gotAudio string
}

func (f *fakeRemindSpeechTranscriber) Transcribe(_ context.Context, audioBase64 string) (string, error) {
	f.gotAudio = audioBase64
	return f.text, f.err
}

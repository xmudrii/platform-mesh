package subroutines

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	tests := []struct {
		name            string
		result          Result
		wantContinue    bool
		wantPending     bool
		wantStopRequeue bool
		wantStop        bool
		wantRequeue     time.Duration
		wantMessage     string
	}{
		{
			name:         "zero value is continue",
			result:       Result{},
			wantContinue: true,
		},
		{
			name:         "OK",
			result:       OK(),
			wantContinue: true,
		},
		{
			name:         "OKWithRequeue",
			result:       OKWithRequeue(5 * time.Second),
			wantContinue: true,
			wantRequeue:  5 * time.Second,
		},
		{
			name:        "Pending",
			result:      Pending(10*time.Second, "waiting for dependency"),
			wantPending: true,
			wantRequeue: 10 * time.Second,
			wantMessage: "waiting for dependency",
		},
		{
			name:            "StopWithRequeue",
			result:          StopWithRequeue(30*time.Second, "rate limited"),
			wantStopRequeue: true,
			wantRequeue:     30 * time.Second,
			wantMessage:     "rate limited",
		},
		{
			name:        "Stop",
			result:      Stop("precondition failed"),
			wantStop:    true,
			wantMessage: "precondition failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantContinue, tt.result.IsContinue())
			assert.Equal(t, tt.wantPending, tt.result.IsPending())
			assert.Equal(t, tt.wantStopRequeue, tt.result.IsStopWithRequeue())
			assert.Equal(t, tt.wantStop, tt.result.IsStop())
			assert.Equal(t, tt.wantRequeue, tt.result.Requeue())
			assert.Equal(t, tt.wantMessage, tt.result.Message())
		})
	}
}

func TestPending_PanicsOnZeroDuration(t *testing.T) {
	assert.PanicsWithValue(t, "subroutines: Pending requires a positive requeue duration", func() {
		Pending(0, "bad")
	})
}

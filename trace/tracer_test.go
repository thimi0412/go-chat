package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("return value is nil for New!")
	} else {
		tracer.Trace("Hello trace package!")
		if buf.String() != "Hello trace package!\n" {
			t.Errorf("'%s' : output incorrect string!", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("data")
}

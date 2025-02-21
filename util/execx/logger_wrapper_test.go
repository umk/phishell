package execx

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

// errorWriter is an io.Writer that returns an error on every write.
type errorWriter struct {
	lastN int
}

func (e *errorWriter) Write(p []byte) (int, error) {
	e.lastN = len(p) / 2 // pretend it could only write half, then fail
	return e.lastN, fmt.Errorf("forced error in base writer")
}

// TestWriteSingleLine ensures writing a single line is tracked correctly.
func TestWriteSingleLine(t *testing.T) {
	var buf bytes.Buffer
	w := newOutputWrapper(&buf, 50)

	input := []byte("Hello world\n")
	n, err := w.Write(input)
	if err != nil {
		t.Fatalf("Write returned unexpected error: %v", err)
	}
	if n != len(input) {
		t.Errorf("Expected n=%d, got %d", len(input), n)
	}
	// Ensure base writer received the data.
	if got := buf.String(); got != string(input) {
		t.Errorf("Base writer mismatch. Got %q, want %q", got, input)
	}

	// Now retrieve from the wrapper
	output := w.Get()
	if output.err != nil {
		t.Errorf("Unexpected error in output: %v", output.err)
	}
	if len(output.lines) != 2 {
		t.Fatalf("Expected 2 lines, got %d lines: %+v", len(output.lines), output.lines)
	}
	if output.lines[0] != "Hello world" {
		t.Errorf("Line mismatch. Got %q, want %q", output.lines[0], "Hello world")
	}
	if output.lines[1] != "" {
		t.Errorf("Line mismatch. Got %q, want %q", output.lines[1], "Hello world")
	}
	if output.tail {
		t.Errorf("Expected tail=false, got true")
	}
}

// TestWriteMultiLine ensures multiple lines get split correctly.
func TestWriteMultiLine(t *testing.T) {
	var buf bytes.Buffer
	w := newOutputWrapper(&buf, 100)

	data := []byte("Line1\nLine2\r\nLine3\nLine4")
	_, err := w.Write(data)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	output := w.Get()
	if output.err != nil {
		t.Errorf("Unexpected error: %v", output.err)
	}

	wantLines := []string{"Line1", "Line2", "Line3", "Line4"}
	if !reflect.DeepEqual(output.lines, wantLines) {
		t.Errorf("Lines mismatch.\nGot  %#v\nWant %#v", output.lines, wantLines)
	}
	if output.tail {
		t.Errorf("Expected tail=false, got true")
	}
}

// TestExceedMaxLen writes a line that exceeds the maxLen, verifying discard behavior.
func TestExceedMaxLen(t *testing.T) {
	w := newOutputWrapper(nil, 5)

	// This line is 11 bytes total (including newline).
	// The line is longer than maxLen (5), so the line should be discarded
	// and w.lines should become nil, tail should be set to true.
	line := []byte("HelloWorld\n")
	_, err := w.Write(line)
	if err != nil {
		t.Fatalf("Unexpected write error: %v", err)
	}

	output := w.Get()
	if output.err != nil {
		t.Errorf("Unexpected error: %v", output.err)
	}
	if len(output.lines) != 0 {
		t.Fatalf("Expected 0 lines after exceeding maxLen, got %d", len(output.lines))
	}
	if !output.tail {
		t.Errorf("Expected tail=true after exceeding maxLen, got false")
	}
}

// TestMultipleLinesExceedingMaxLen tests multiple lines that together exceed the maxLen (aggregated).
func TestMultipleLinesExceedingMaxLen(t *testing.T) {
	w := newOutputWrapper(nil, 10)

	// We'll write 3 short lines, but the total (including newline overhead in linesLen counting)
	// will exceed 10 eventually. We want to see the discarding from the beginning.
	//
	// linesLen is computed as the sum of each line length + 1.
	// If we have lines "AB", "CDE", "FGHI", then linesLen = 2+1 + 3+1 + 4+1 = 12.
	// That exceeds 10, so we should discard from the start until linesLen <= 10.
	//
	// We'll break it down:
	// After first line "AB": linesLen=3 -> <=10, keep [AB].
	// After second line "CDE": linesLen=3+4=7 -> <=10, keep [AB, CDE].
	// After third line "FGHI": linesLen=7+(4+1)=12 -> exceeds 10 => discard from start
	// Discard "AB": linesLen=12 - (2+1)=9 => 9 <=10 now, so we keep [CDE, FGHI], tail=true
	data := []byte("AB\nCDE\nFGHI\n")
	_, err := w.Write(data)
	if err != nil {
		t.Fatalf("Unexpected write error: %v", err)
	}

	out := w.Get()
	if out.err != nil {
		t.Fatalf("Unexpected error: %v", out.err)
	}
	// We expect the final lines array to be [CDE, FGHI, ""], tail=true
	want := []string{"CDE", "FGHI", ""}
	if !reflect.DeepEqual(out.lines, want) {
		t.Errorf("Lines mismatch.\nGot  %#v\nWant %#v", out.lines, want)
	}
	if !out.tail {
		t.Errorf("Expected tail=true but got false")
	}
}

// TestPartialUTF8 ensures that partial (incomplete) multi-byte runes are handled.
func TestPartialUTF8(t *testing.T) {
	// A 3-byte UTF-8 character is "€": 0xE2 0x82 0xAC
	// We'll supply only two of its bytes first.
	w := newOutputWrapper(nil, 50)

	partial := []byte{0xE2, 0x82} // incomplete "€"
	_, err := w.Write(partial)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Now write the final byte to complete the character.
	rest := []byte{0xAC, '\n'} // now the "€" is complete, plus a newline
	_, err = w.Write(rest)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	out := w.Get()
	if out.err != nil {
		t.Errorf("Unexpected error after final UTF-8 completion: %v", out.err)
	}
	// We should have exactly two lines, first containing "€", second is empty.
	if len(out.lines) != 2 {
		t.Fatalf("Expected 1 line, got %d", len(out.lines))
	}
	if out.lines[0] != "€" {
		t.Errorf("Expected line to be '€', got %q", out.lines[0])
	}
	if out.lines[1] != "" {
		t.Errorf("Expected line to be empty, got %q", out.lines[1])
	}
	if out.tail {
		t.Errorf("Expected tail=false, got true")
	}
}

// TestFinalIncompleteUTF8 checks that finalizing with incomplete UTF-8 triggers an error.
func TestFinalIncompleteUTF8(t *testing.T) {
	w := newOutputWrapper(nil, 50)

	// Write partial bytes of "€" and do NOT complete it.
	partial := []byte{0xE2, 0x82}
	_, err := w.Write(partial)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Now get the output, which triggers finalize and
	// should detect incomplete UTF-8 => error
	out := w.Get()
	if out.err == nil {
		t.Error("Expected an error due to incomplete UTF-8, got nil")
	}
	if len(out.lines) != 0 {
		t.Errorf("Expected no lines on error, got %v", out.lines)
	}
	if out.tail {
		t.Errorf("Expected tail=false on error, got true")
	}
}

// TestInvalidUTF8 tests detection of invalid UTF-8 sequences (e.g., 0xC0, 0xAF).
func TestInvalidUTF8(t *testing.T) {
	w := newOutputWrapper(nil, 50)

	invalidSeq := []byte{0xC0, 0xAF} // not valid UTF-8
	_, err := w.Write(invalidSeq)
	if err != nil {
		t.Errorf("Write() itself should not return an error from w.base, got: %v", err)
	}

	// Now get the output, which should reflect the error from processOutput.
	out := w.Get()
	if out.err == nil {
		t.Error("Expected invalid UTF-8 error, got nil")
	}
	if len(out.lines) != 0 {
		t.Errorf("Expected no lines, got %v", out.lines)
	}
	if out.tail {
		t.Errorf("Expected tail=false on error, got true")
	}
}

// TestTrailingDataWithoutNewline ensures that finalize() captures trailing data
// that doesn't end with a newline.
func TestTrailingDataWithoutNewline(t *testing.T) {
	w := newOutputWrapper(nil, 50)
	_, err := w.Write([]byte("NoNewline"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	out := w.Get()
	if out.err != nil {
		t.Fatalf("Unexpected error: %v", out.err)
	}
	if len(out.lines) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(out.lines))
	}
	if out.lines[0] != "NoNewline" {
		t.Errorf("Line mismatch, got %q, want %q", out.lines[0], "NoNewline")
	}
	if out.tail {
		t.Errorf("Expected tail=false, got true")
	}
}

// TestErrorReset ensures that once an error happens in processOutput (like invalid UTF-8),
// subsequent writes do not affect lines, but still go to base writer.
func TestErrorReset(t *testing.T) {
	var base bytes.Buffer
	w := newOutputWrapper(&base, 50)

	// First, inject invalid UTF-8
	_, _ = w.Write([]byte{0xC0, 0xAF}) // sets w.err internally

	// Now attempt writing valid data
	validData := []byte("Hello\n")
	n, err := w.Write(validData)
	if err != nil {
		t.Errorf("Expected no error from Write (base writer is OK), got %v", err)
	}
	if n != len(validData) {
		t.Errorf("Write returned n=%d, expected %d", n, len(validData))
	}

	// The base writer should still get the second write:
	if got := base.Bytes(); !bytes.Equal(got, append([]byte{0xC0, 0xAF}, validData...)) {
		t.Errorf("Base writer mismatch. Got %q, want %q", got, append([]byte{0xC0, 0xAF}, validData...))
	}

	// But lines should not reflect the second write because error was set:
	out := w.Get()
	if out.err == nil {
		t.Fatalf("Expected an error to be retained, got nil")
	}
	if len(out.lines) != 0 {
		t.Errorf("Expected no lines, got %v", out.lines)
	}
	if out.tail {
		t.Errorf("Expected tail=false on error, got true")
	}
}

// TestWriteReturnValues ensures Write() always returns len(p), nil unless base writer fails.
func TestWriteReturnValues(t *testing.T) {
	t.Run("BaseWriterNoError", func(t *testing.T) {
		var base bytes.Buffer
		w := newOutputWrapper(&base, 10)
		data := []byte("abc")
		n, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if n != len(data) {
			t.Errorf("Expected n=%d, got %d", len(data), n)
		}
	})

	t.Run("BaseWriterError", func(t *testing.T) {
		// Create a base writer that always errors after N bytes or some condition
		errWriter := &errorWriter{}
		w := newOutputWrapper(errWriter, 10)
		data := []byte("abc")
		n, err := w.Write(data)
		if err == nil {
			t.Fatalf("Expected an error from base writer, got nil")
		}
		// The number of bytes that made it into the base writer is controlled by errorWriter.
		if n != errWriter.lastN {
			t.Errorf("Write returned n=%d, want %d", n, errWriter.lastN)
		}
	})
}

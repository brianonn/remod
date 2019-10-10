package remod

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
)

// Enable will enable remod dev replacements for the given modules.
func Enable(data []byte, modules []string, base string, version string) ([]byte, error) {

	if len(modules) == 0 {
		return data, nil
	}

	if bytes.Contains(data, []byte("replace ( // remod:replacements")) {
		return data, nil
	}

	if version != "" {
		version = " " + version
	}

	buf := bytes.NewBuffer(data)

	_, _ = buf.WriteString("\nreplace ( // remod:replacements\n")
	for _, m := range modules {
		_, _ = buf.WriteString(fmt.Sprintf("\t%s => %s%s%s\n", m, base, filepath.Base(m), version))
	}
	_, _ = buf.WriteString(")")

	return append(bytes.TrimSpace(buf.Bytes()), '\n'), nil
}

// Disable will disable remod dev replacements.
func Disable(data []byte) ([]byte, error) {

	scanner := bufio.NewScanner(bytes.NewBuffer(data))

	buf := bytes.NewBuffer(nil)

	start := []byte("replace ( // remod:replacements")
	end := []byte(")")

	var strip bool
	var last []byte
	for scanner.Scan() {

		line := scanner.Bytes()

		if bytes.Equal(line, start) {
			strip = true
			continue
		}

		if strip && bytes.HasPrefix(line, end) {
			strip = false
			continue
		}

		if strip {
			continue
		}

		_, _ = buf.Write(line)

		if !bytes.Equal(last, []byte("\n")) {
			_ = buf.WriteByte('\n')
		}

		last = line
	}

	return append(bytes.TrimSpace(buf.Bytes()), '\n'), nil
}

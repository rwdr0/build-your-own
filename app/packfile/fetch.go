// Package packfile fetches and unpacks the git wire protocol v2
package packfile

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Fetch(url string) ([]byte, error) {
	refs, err := fetchRefs(url)
	if err != nil {
		return nil, err
	}

	packfile, err := fetchPackfile(url, refs)
	if err != nil {
		return nil, err
	}

	return packfile, nil
}

func fetchRefs(url string) ([]string, error) {
	body := "0014command=ls-refs\n" +
		"0016object-format=sha1" +
		"0001" +
		"001bref-prefix refs/heads/\n" +
		"0014ref-prefix HEAD\n" +
		"0000"

	req, err := http.NewRequest("POST", url+"/git-upload-pack", strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-git-upload-pack-request")
	req.Header.Set("Accept", "application/x-git-upload-pack-result")
	req.Header.Set("Git-Protocol", "version=2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var refs []string
	data := string(respBody)
	for len(data) >= 4 {
		var pktLen int
		fmt.Sscanf(data[:4], "%x", &pktLen)
		data = data[4:]

		if pktLen == 0 { // flush packet
			break
		}
		if pktLen <= 4 { // special packets (delimiter, response-end) with no content
			continue
		}

		content := data[:pktLen-4]
		data = data[pktLen-4:]

		// ref line: "<40-char-hash> <refname>\n"
		if len(content) >= 41 && content[40] == ' ' {
			refs = append(refs, content[:40])
		}
	}
	return refs, nil
}

func fetchPackfile(url string, refs []string) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString("0011command=fetch")
	sb.WriteString("0016object-format=sha1")
	sb.WriteString("0001")
	sb.WriteString("000dofs-delta")

	for _, ref := range refs {
		sb.WriteString(fmt.Sprintf("0032want %s\n", ref))
	}

	sb.WriteString("0009done\n")
	sb.WriteString("0000")

	req, err := http.NewRequest("POST", url+"/git-upload-pack", strings.NewReader(sb.String()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-git-upload-pack-request")
	req.Header.Set("Accept", "application/x-git-upload-pack-result")
	req.Header.Set("Git-Protocol", "version=2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	dataChannel := make(chan []byte)
	errChannel := make(chan error, 1)

	go func() {
		defer close(dataChannel)

		// skip `packfile` header
		_, err := resp.Body.Read(make([]byte, 13))
		if err != nil {
			errChannel <- err
			return
		}

		buf := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])
				dataChannel <- chunk
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				errChannel <- err
				return
			}
		}
	}()

	var packData []byte
	var pending []byte

	for chunk := range dataChannel {
		pending = append(pending, chunk...)

		for len(pending) >= 4 {
			var pktLen int
			fmt.Sscanf(string(pending[:4]), "%x", &pktLen)

			if pktLen == 0 { // flush packet
				pending = pending[4:]
				continue
			}

			if len(pending) < pktLen {
				break // wait for more data
			}

			sideband := pending[4]
			if sideband == 1 {
				packData = append(packData, pending[5:pktLen]...)
			}

			pending = pending[pktLen:]
		}
	}

	select {
	case err := <-errChannel:
		return nil, err
	default:
	}

	return packData, nil
}

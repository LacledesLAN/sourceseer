package srcds

import (
	"bufio"
	"strings"
)

func (s *server) linkStdoutPipe(r *bufio.Reader) {
	go func(r *bufio.Reader) {
		for {
			outLine, _ := r.ReadString('\n')
			outLine = strings.Trim(strings.TrimSuffix(outLine, "\n"), "")

			if len(outLine) > 0 {

				Log(SRCDSOther, outLine)
			}
		}
	}(r)
}

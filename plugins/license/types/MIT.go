package types

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var MIT_LICENSE = new(mit)

type mit struct{}

func (mi *mit) GetLicense() string {
	author := os.Getenv("MC_LICENSE_AUTHOR")
	currentYear := time.Now().Year()

	cmd := exec.Command("git", "log", "--reverse", "--date=format:\"%Y\"", "--format=%ad")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error getting first commit year: %s", err)
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 || lines[0] == "" {
		log.Fatalf("Could not determine the first commit year.")
	}
	firstCommitYear := strings.TrimSpace(lines[0])
	firstCommitYear = strings.Trim(firstCommitYear, "\"")

	mit := fmt.Sprintf(`Copyright (c) %s-%d %s
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.`, firstCommitYear, currentYear, author)
	return mit
}

func (mi *mit) GetType() string {
	return "MIT"
}

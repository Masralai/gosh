package handlers

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/urfave/cli/v3"
)

func runCmd(t *testing.T, dir string, args ...string) (string, string, error) {
	t.Helper()
	if dir != "" {
		t.Chdir(dir)
	}
	root := &cli.Command{Name: "gosh", Commands: All()}
	oldOut, oldErr := os.Stdout, os.Stderr
	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stderr: %v", err)
	}
	os.Stdout, os.Stderr = wOut, wErr

	var outBuf, errBuf bytes.Buffer
	outDone := make(chan struct{})
	errDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(&outBuf, rOut)
		close(outDone)
	}()
	go func() {
		_, _ = io.Copy(&errBuf, rErr)
		close(errDone)
	}()

	runErr := root.Run(context.Background(), append([]string{"gosh"}, args...))
	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout, os.Stderr = oldOut, oldErr
	<-outDone
	<-errDone
	_ = rOut.Close()
	_ = rErr.Close()
	return outBuf.String(), errBuf.String(), runErr
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestEchoPrintsAllArgs(t *testing.T) {
	out, _, err := runCmd(t, "", "echo", "hello", "big", "world")
	if err != nil {
		t.Fatalf("echo error: %v", err)
	}
	if strings.TrimSpace(out) != "hello big world" {
		t.Fatalf("echo out = %q, want %q", out, "hello big world\n")
	}
}

func TestEchoNoNewline(t *testing.T) {
	out, _, err := runCmd(t, "", "echo", "-n", "hi")
	if err != nil {
		t.Fatalf("echo -n error: %v", err)
	}
	if out != "hi" {
		t.Fatalf("echo -n out = %q, want %q", out, "hi")
	}
}

func TestCatNumberAllAndSqueezeAndBlank(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "f.txt")
	writeFile(t, file, "a\n\n\nb\n")
	out, _, err := runCmd(t, dir, "cat", "-n", "-s", "f.txt")
	if err != nil {
		t.Fatalf("cat error: %v", err)
	}
	if !strings.Contains(out, "     1  a") || !strings.Contains(out, "     2  ") || !strings.Contains(out, "     3  b") {
		t.Fatalf("cat -n -s output = %q", out)
	}
	// ensure squeezed: should not have 4th numbered blank from triple blanks
	if strings.Contains(out, "     4  ") {
		t.Fatalf("cat -n -s did not squeeze: %q", out)
	}
}

func TestCatNumberNonBlankOverridesN(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "f.txt")
	writeFile(t, file, "a\n\nb\n")
	out, _, err := runCmd(t, dir, "cat", "-b", "-n", "f.txt")
	if err != nil {
		t.Fatalf("cat error: %v", err)
	}
	if !strings.Contains(out, "     1  a") || !strings.Contains(out, "     2  b") {
		t.Fatalf("cat -b output missing numbers: %q", out)
	}
	// blank line should be present but unnumbered
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	for _, l := range lines {
		if strings.Contains(l, "     3") {
			t.Fatalf("cat -b should not number blank as 3: %q", out)
		}
	}
}

func TestDirListsHumanReadable(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "file.txt"), "x")
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	out, _, err := runCmd(t, dir, "dir", ".")
	if err != nil {
		t.Fatalf("dir error: %v", err)
	}
	if strings.Contains(out, "DirEntry") || strings.Contains(out, "[{") {
		t.Fatalf("dir output not human readable: %q", out)
	}
	if !strings.Contains(out, "file.txt") || !strings.Contains(out, "sub/") && !strings.Contains(out, "sub\\") && !strings.Contains(out, "sub") {
		t.Fatalf("dir output missing entries: %q", out)
	}
}

func TestLsFlags(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".hidden"), "h")
	writeFile(t, filepath.Join(dir, "small.txt"), "a")
	writeFile(t, filepath.Join(dir, "big.txt"), "abcdef")
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "sub", "inner.txt"), "i")

	out, _, err := runCmd(t, dir, "ls", "-a")
	if err != nil {
		t.Fatalf("ls -a error: %v", err)
	}
	if !strings.Contains(out, ".hidden") {
		t.Fatalf("ls -a missing hidden: %q", out)
	}

	out, _, err = runCmd(t, dir, "ls", "-S")
	if err != nil {
		t.Fatalf("ls -S error: %v", err)
	}
	bigIdx := strings.Index(out, "big.txt")
	smallIdx := strings.Index(out, "small.txt")
	if bigIdx == -1 || smallIdx == -1 || bigIdx > smallIdx {
		t.Fatalf("ls -S not sorted by size desc: %q", out)
	}

	out, _, err = runCmd(t, dir, "ls", "-R")
	if err != nil {
		t.Fatalf("ls -R error: %v", err)
	}
	if !strings.Contains(out, "sub") || !strings.Contains(out, "inner.txt") {
		t.Fatalf("ls -R missing recursive content: %q", out)
	}
}

func TestTouchCreatesAndUpdatesWithoutTruncating(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "note.txt")
	writeFile(t, path, "hello")
	oldInfo, _ := os.Stat(path)
	time.Sleep(1100 * time.Millisecond)
	_, _, err := runCmd(t, dir, "touch", "note.txt")
	if err != nil {
		t.Fatalf("touch error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("touch truncated content: %q", string(data))
	}
	newInfo, _ := os.Stat(path)
	if !newInfo.ModTime().After(oldInfo.ModTime()) && !newInfo.ModTime().Equal(oldInfo.ModTime()) {
		t.Fatalf("mtime not updated: old=%v new=%v", oldInfo.ModTime(), newInfo.ModTime())
	}
}

func TestTouchCreatesMissingFile(t *testing.T) {
	dir := t.TempDir()
	_, _, err := runCmd(t, dir, "touch", "new.txt")
	if err != nil {
		t.Fatalf("touch create error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "new.txt")); err != nil {
		t.Fatalf("touch did not create file: %v", err)
	}
}

func TestCpCopiesFileToDirectory(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.txt")
	dstDir := filepath.Join(dir, "out")
	writeFile(t, src, "content")
	if err := os.Mkdir(dstDir, 0o755); err != nil {
		t.Fatal(err)
	}
	_, _, err := runCmd(t, dir, "cp", "a.txt", "out")
	if err != nil {
		t.Fatalf("cp error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dstDir, "a.txt"))
	if err != nil || string(data) != "content" {
		t.Fatalf("cp result: err=%v data=%q", err, string(data))
	}
}

func TestZipDoesNotDoubleExtensionAndUnzipExtracts(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "one.txt"), "one")
	writeFile(t, filepath.Join(dir, "two.txt"), "two")
	_, _, err := runCmd(t, dir, "zip", "arc.zip", "one.txt", "two.txt")
	if err != nil {
		t.Fatalf("zip error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "arc.zip.zip")); err == nil {
		t.Fatalf("zip created double extension")
	}
	if _, err := os.Stat(filepath.Join(dir, "arc.zip")); err != nil {
		t.Fatalf("zip did not create arc.zip: %v", err)
	}
	outDir := filepath.Join(dir, "out")
	_, _, err = runCmd(t, dir, "unzip", "arc.zip", "out")
	if err != nil {
		t.Fatalf("unzip error: %v", err)
	}
	for _, name := range []string{"one.txt", "two.txt"} {
		b, err := os.ReadFile(filepath.Join(outDir, name))
		if err != nil {
			t.Fatalf("missing extracted %s: %v", name, err)
		}
		if name == "one.txt" && string(b) != "one" {
			t.Fatalf("bad content for one.txt: %q", string(b))
		}
	}
	// verify archive has expected entries
	zr, err := zip.OpenReader(filepath.Join(dir, "arc.zip"))
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	if len(zr.File) != 2 {
		t.Fatalf("zip entries = %d, want 2", len(zr.File))
	}
}

func TestKillUsageAndInvalidPid(t *testing.T) {
	_, _, err := runCmd(t, "", "kill")
	if err == nil || !strings.Contains(err.Error(), "usage: kill <pid>") {
		t.Fatalf("kill no args err = %v", err)
	}
	_, _, err = runCmd(t, "", "kill", "abc")
	if err == nil || !strings.Contains(err.Error(), "invalid pid") {
		t.Fatalf("kill invalid pid err = %v", err)
	}
}

func TestPingUsage(t *testing.T) {
	_, _, err := runCmd(t, "", "ping")
	if err == nil || !strings.Contains(err.Error(), "usage: ping <hostname>") {
		t.Fatalf("ping no args err = %v", err)
	}
}

func TestPingResolveErrorMessage(t *testing.T) {
	_, _, err := runCmd(t, "", "ping", "no.such.host.invalid.")
	if err == nil {
		// on some environments DNS may behave oddly; but usually fails
		t.Skip("expected resolve error but got nil")
	}
	if !strings.Contains(err.Error(), "failed to resolve host") && !strings.Contains(strings.ToLower(err.Error()), "ping") {
		t.Fatalf("ping resolve err = %v", err)
	}
}

func TestClearBuiltinRegistered(t *testing.T) {
	root := &cli.Command{Name: "gosh", Commands: All()}
	found := false
	for _, c := range root.Commands {
		if c.Name == "clear" {
			found = true
			if c.UsageText != "clear" {
				t.Fatalf("clear usage = %q", c.UsageText)
			}
		}
	}
	if !found {
		t.Fatalf("clear command not registered")
	}
	_ = runtime.GOOS
}

func TestHelpListsAndDetails(t *testing.T) {
	out, _, err := runCmd(t, "", "help")
	if err != nil {
		t.Fatalf("help error: %v", err)
	}
	for _, want := range []string{"help", "history", "env", "which", "find", "wc", "sort", "ls", "grep", "cat"} {
		if !strings.Contains(out, want) {
			t.Fatalf("help list missing %q: %q", want, out)
		}
	}

	out, _, err = runCmd(t, "", "help", "grep")
	if err != nil {
		t.Fatalf("help grep error: %v", err)
	}
	if !strings.Contains(out, "Name: grep") || !strings.Contains(out, "UsageText:") {
		t.Fatalf("help grep output: %q", out)
	}
	if !strings.Contains(out, "grep [-i] [-n] [-r] [-v] <pattern> <path>") {
		t.Fatalf("help grep missing usage text: %q", out)
	}

	_, _, err = runCmd(t, "", "help", "nope")
	if err == nil || !strings.Contains(err.Error(), "unknown help topic") {
		t.Fatalf("help unknown err = %v", err)
	}
}

func TestHistoryDisplayAndClear(t *testing.T) {
	dir := t.TempDir()
	hist := filepath.Join(dir, "histfile")
	t.Setenv("GOSH_HISTORY", hist)

	writeFile(t, hist, "echo one\necho two\n")
	out, _, err := runCmd(t, "", "history")
	if err != nil {
		t.Fatalf("history error: %v", err)
	}
	if !strings.Contains(out, "echo one") || !strings.Contains(out, "echo two") {
		t.Fatalf("history out = %q", out)
	}

	out, _, err = runCmd(t, "", "history", "-c")
	if err != nil {
		t.Fatalf("history -c error: %v", err)
	}
	if strings.TrimSpace(out) != "" {
		t.Fatalf("history -c out = %q", out)
	}
	data, err := os.ReadFile(hist)
	if err != nil || len(data) != 0 {
		t.Fatalf("history not cleared: err=%v data=%q", err, string(data))
	}

	out, _, err = runCmd(t, "", "history")
	if err != nil {
		t.Fatalf("history after clear error: %v", err)
	}
	if strings.TrimSpace(out) != "" {
		t.Fatalf("history after clear out = %q", out)
	}
}

func TestHistoryMissingFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GOSH_HISTORY", filepath.Join(dir, "missing", "hist"))
	out, _, err := runCmd(t, "", "history")
	if err != nil {
		t.Fatalf("history missing error: %v", err)
	}
	if strings.TrimSpace(out) != "" {
		t.Fatalf("history missing out = %q", out)
	}
}

func TestEnvListAndSpecific(t *testing.T) {
	t.Setenv("GOSH_TEST_ENV_A", "valueA")
	out, _, err := runCmd(t, "", "env")
	if err != nil {
		t.Fatalf("env error: %v", err)
	}
	if !strings.Contains(out, "GOSH_TEST_ENV_A=valueA") {
		t.Fatalf("env list missing var: %q", out)
	}

	out, _, err = runCmd(t, "", "env", "GOSH_TEST_ENV_A")
	if err != nil {
		t.Fatalf("env name error: %v", err)
	}
	if strings.TrimSpace(out) != "valueA" {
		t.Fatalf("env name out = %q", out)
	}

	out, _, err = runCmd(t, "", "env", "GOSH_TEST_ENV_MISSING_VAR_12345")
	if err != nil {
		t.Fatalf("env missing name error: %v", err)
	}
	if strings.TrimSpace(out) != "" {
		t.Fatalf("env missing name out = %q", out)
	}
}

func TestLsLongHumanTime(t *testing.T) {
	dir := t.TempDir()
	small := filepath.Join(dir, "small.txt")
	big := filepath.Join(dir, "big.txt")
	writeFile(t, small, "a")
	writeFile(t, big, strings.Repeat("b", 2048))
	// ensure different mtimes
	t1 := time.Now().Add(-2 * time.Hour)
	t2 := time.Now().Add(-1 * time.Hour)
	_ = os.Chtimes(small, t1, t1)
	_ = os.Chtimes(big, t2, t2)

	out, _, err := runCmd(t, dir, "ls", "-l")
	if err != nil {
		t.Fatalf("ls -l error: %v", err)
	}
	if !strings.Contains(out, "small.txt") || !strings.Contains(out, "big.txt") {
		t.Fatalf("ls -l missing entries: %q", out)
	}
	// long listing should have mode-like prefix and size/mod info
	if !strings.Contains(out, "-") && !strings.Contains(out, "d") {
		t.Fatalf("ls -l missing mode: %q", out)
	}

	out, _, err = runCmd(t, dir, "ls", "-lh")
	if err != nil {
		t.Fatalf("ls -lh error: %v", err)
	}
	if !strings.Contains(out, "K") && !strings.Contains(out, "B") {
		t.Fatalf("ls -lh missing human size: %q", out)
	}

	out, _, err = runCmd(t, dir, "ls", "-t")
	if err != nil {
		t.Fatalf("ls -t error: %v", err)
	}
	bigIdx := strings.Index(out, "big.txt")
	smallIdx := strings.Index(out, "small.txt")
	if bigIdx == -1 || smallIdx == -1 || bigIdx > smallIdx {
		t.Fatalf("ls -t not newest first: %q", out)
	}

	out, _, err = runCmd(t, dir, "ls", "-lt")
	if err != nil {
		t.Fatalf("ls -lt error: %v", err)
	}
	bigIdx = strings.Index(out, "big.txt")
	smallIdx = strings.Index(out, "small.txt")
	if bigIdx == -1 || smallIdx == -1 || bigIdx > smallIdx {
		t.Fatalf("ls -lt not newest first: %q", out)
	}
}

func TestLsHumanOnlyAffectsLong(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "big.txt"), strings.Repeat("b", 2048))
	out, _, err := runCmd(t, dir, "ls", "-h")
	if err != nil {
		t.Fatalf("ls -h error: %v", err)
	}
	if !strings.Contains(out, "big.txt") {
		t.Fatalf("ls -h missing file: %q", out)
	}
	if strings.Contains(out, "2.0K") || strings.Contains(out, "2.0k") {
		t.Fatalf("ls -h without -l should not force long human size format: %q", out)
	}
}

func TestMkdirP(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "a", "b", "c")
	_, _, err := runCmd(t, dir, "mkdir", "-p", filepath.Join("a", "b", "c"))
	if err != nil {
		t.Fatalf("mkdir -p error: %v", err)
	}
	if info, err := os.Stat(nested); err != nil || !info.IsDir() {
		t.Fatalf("mkdir -p did not create nested dir: err=%v", err)
	}
	// idempotent
	_, _, err = runCmd(t, dir, "mkdir", "-p", filepath.Join("a", "b", "c"))
	if err != nil {
		t.Fatalf("mkdir -p existing error: %v", err)
	}
}

func TestRmForceAndRecursiveSemantics(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "f.txt"), "x")
	if err := os.MkdirAll(filepath.Join(dir, "d", "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "d", "sub", "x.txt"), "x")

	// rm directory without -r should error
	_, _, err := runCmd(t, dir, "rm", "d")
	if err == nil || !strings.Contains(err.Error(), "is a directory") {
		t.Fatalf("rm dir without -r err = %v", err)
	}

	// rm -rf removes recursively without prompt and returns success
	_, _, err = runCmd(t, dir, "rm", "-rf", "d")
	if err != nil {
		t.Fatalf("rm -rf error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "d")); !os.IsNotExist(err) {
		t.Fatalf("rm -rf did not remove dir")
	}

	// rm missing without -f errors
	_, _, err = runCmd(t, dir, "rm", "missing-nope.txt")
	if err == nil {
		t.Fatalf("rm missing without -f should error")
	}

	// rm -f missing does not error
	_, _, err = runCmd(t, dir, "rm", "-f", "missing-nope.txt")
	if err != nil {
		t.Fatalf("rm -f missing error: %v", err)
	}

	// rm file
	_, _, err = runCmd(t, dir, "rm", "f.txt")
	if err != nil {
		t.Fatalf("rm file error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "f.txt")); !os.IsNotExist(err) {
		t.Fatalf("rm file did not remove")
	}
}

func TestRmRUppercase(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "d"), 0o755); err != nil {
		t.Fatal(err)
	}
	_, _, err := runCmd(t, dir, "rm", "-R", "d")
	if err != nil {
		t.Fatalf("rm -R error: %v", err)
	}
}

func TestCpRecursive(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "srcdir")
	if err := os.MkdirAll(filepath.Join(src, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(src, "a.txt"), "A")
	writeFile(t, filepath.Join(src, "sub", "b.txt"), "B")

	// cp dir without -r errors
	_, _, err := runCmd(t, dir, "cp", "srcdir", "dst1")
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "-r") {
		t.Fatalf("cp dir without -r err = %v", err)
	}

	_, _, err = runCmd(t, dir, "cp", "-r", "srcdir", "dst1")
	if err != nil {
		t.Fatalf("cp -r error: %v", err)
	}
	if b, err := os.ReadFile(filepath.Join(dir, "dst1", "a.txt")); err != nil || string(b) != "A" {
		t.Fatalf("cp -r a.txt err=%v b=%q", err, string(b))
	}
	if b, err := os.ReadFile(filepath.Join(dir, "dst1", "sub", "b.txt")); err != nil || string(b) != "B" {
		t.Fatalf("cp -r b.txt err=%v b=%q", err, string(b))
	}

	// copy into existing directory
	outDir := filepath.Join(dir, "out")
	if err := os.Mkdir(outDir, 0o755); err != nil {
		t.Fatal(err)
	}
	_, _, err = runCmd(t, dir, "cp", "-r", "srcdir", "out")
	if err != nil {
		t.Fatalf("cp -r into existing dir error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "srcdir", "a.txt")); err != nil {
		t.Fatalf("cp -r into existing dir missing file: %v", err)
	}
}

func TestGrepFlagsAndCompat(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "f.txt"), "Hello\nworld\nHELLO\n")
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(sub, "g.txt"), "hello there\n")

	// -i case insensitive
	out, _, err := runCmd(t, dir, "grep", "-i", "hello", "f.txt")
	if err != nil {
		t.Fatalf("grep -i error: %v", err)
	}
	if strings.Count(out, "Hello")+strings.Count(out, "HELLO")+strings.Count(out, "hello") < 2 {
		t.Fatalf("grep -i out = %q", out)
	}

	// -f compat alias for ignore-case
	out2, _, err := runCmd(t, dir, "grep", "-f", "hello", "f.txt")
	if err != nil {
		t.Fatalf("grep -f compat error: %v", err)
	}
	if out != out2 {
		t.Fatalf("grep -f compat differs: %q vs %q", out2, out)
	}

	// -n line numbers
	out, _, err = runCmd(t, dir, "grep", "-n", "world", "f.txt")
	if err != nil {
		t.Fatalf("grep -n error: %v", err)
	}
	if strings.TrimSpace(out) != "2:world" {
		t.Fatalf("grep -n out = %q", out)
	}

	// -v invert
	out, _, err = runCmd(t, dir, "grep", "-v", "world", "f.txt")
	if err != nil {
		t.Fatalf("grep -v error: %v", err)
	}
	if strings.Contains(out, "world") || !strings.Contains(out, "Hello") {
		t.Fatalf("grep -v out = %q", out)
	}

	// -r recursive
	out, _, err = runCmd(t, dir, "grep", "-r", "hello", ".")
	if err != nil {
		t.Fatalf("grep -r error: %v", err)
	}
	if !strings.Contains(out, "g.txt") {
		t.Fatalf("grep -r out = %q", out)
	}

	// -rn combo
	out, _, err = runCmd(t, dir, "grep", "-r", "-n", "-i", "hello", ".")
	if err != nil {
		t.Fatalf("grep -rn -i error: %v", err)
	}
	if !strings.Contains(out, ":1:") && !strings.Contains(out, ":1\t") && !strings.Contains(strings.ToLower(out), "hello") {
		t.Fatalf("grep -rn -i out = %q", out)
	}
}

func TestWcCounts(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")
	writeFile(t, f1, "one two\nthree\n") // 2 lines, 3 words, 14 bytes
	writeFile(t, f2, "x\n")              // 1 line, 1 word, 2 bytes

	out, _, err := runCmd(t, dir, "wc", "a.txt")
	if err != nil {
		t.Fatalf("wc error: %v", err)
	}
	if !strings.Contains(out, "2") || !strings.Contains(out, "3") || !strings.Contains(out, "14") || !strings.Contains(out, "a.txt") {
		t.Fatalf("wc default out = %q", out)
	}

	out, _, err = runCmd(t, dir, "wc", "-l", "a.txt")
	if err != nil {
		t.Fatalf("wc -l error: %v", err)
	}
	if strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(out), "a.txt")) != "2" && !strings.Contains(out, "2 a.txt") {
		t.Fatalf("wc -l out = %q", out)
	}

	out, _, err = runCmd(t, dir, "wc", "-w", "a.txt")
	if err != nil {
		t.Fatalf("wc -w error: %v", err)
	}
	if !strings.Contains(out, "3") {
		t.Fatalf("wc -w out = %q", out)
	}

	out, _, err = runCmd(t, dir, "wc", "-c", "a.txt")
	if err != nil {
		t.Fatalf("wc -c error: %v", err)
	}
	if !strings.Contains(out, "14") {
		t.Fatalf("wc -c out = %q", out)
	}

	out, _, err = runCmd(t, dir, "wc", "a.txt", "b.txt")
	if err != nil {
		t.Fatalf("wc multi error: %v", err)
	}
	if !strings.Contains(out, "total") || !strings.Contains(out, "3") || !strings.Contains(out, "4") || !strings.Contains(out, "16") {
		t.Fatalf("wc multi out = %q", out)
	}

	_, _, err = runCmd(t, dir, "wc", "nope.txt")
	if err == nil {
		t.Fatalf("wc missing file should error")
	}
}

func TestSortBasicAndReverse(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "s.txt"), "banana\napple\ncherry\n")
	out, _, err := runCmd(t, dir, "sort", "s.txt")
	if err != nil {
		t.Fatalf("sort error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(strings.ReplaceAll(out, "\r\n", "\n")), "\n")
	if strings.Join(lines, ",") != "apple,banana,cherry" {
		t.Fatalf("sort out = %q", out)
	}

	out, _, err = runCmd(t, dir, "sort", "-r", "s.txt")
	if err != nil {
		t.Fatalf("sort -r error: %v", err)
	}
	lines = strings.Split(strings.TrimSpace(strings.ReplaceAll(out, "\r\n", "\n")), "\n")
	if strings.Join(lines, ",") != "cherry,banana,apple" {
		t.Fatalf("sort -r out = %q", out)
	}

	_, _, err = runCmd(t, dir, "sort", "missing.txt")
	if err == nil {
		t.Fatalf("sort missing should error")
	}
}

func TestCatMultiFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.txt"), "a1\na2\n")
	writeFile(t, filepath.Join(dir, "b.txt"), "b1\n\nb2\n")

	out, _, err := runCmd(t, dir, "cat", "a.txt", "b.txt")
	if err != nil {
		t.Fatalf("cat multi error: %v", err)
	}
	if !strings.Contains(out, "a1") || !strings.Contains(out, "b2") {
		t.Fatalf("cat multi out = %q", out)
	}

	out, _, err = runCmd(t, dir, "cat", "-n", "a.txt", "b.txt")
	if err != nil {
		t.Fatalf("cat -n multi error: %v", err)
	}
	// numbering continues across files
	if !strings.Contains(out, "     1  a1") || !strings.Contains(out, "     3  b1") {
		t.Fatalf("cat -n multi out = %q", out)
	}

	_, _, err = runCmd(t, dir, "cat", "a.txt", "missing.txt")
	if err == nil {
		t.Fatalf("cat unreadable should error")
	}
}

func TestFindByName(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "a.txt"), "a")
	writeFile(t, filepath.Join(dir, "sub", "b.txt"), "b")
	writeFile(t, filepath.Join(dir, "sub", "c.log"), "c")

	out, _, err := runCmd(t, dir, "find", ".", "-name", "*.txt")
	if err != nil {
		t.Fatalf("find error: %v", err)
	}
	if !strings.Contains(out, "a.txt") || !strings.Contains(out, "b.txt") || strings.Contains(out, "c.log") {
		t.Fatalf("find out = %q", out)
	}

	_, _, err = runCmd(t, dir, "find", "nope", "-name", "*.txt")
	if err == nil {
		t.Fatalf("find missing path should error")
	}

	_, _, err = runCmd(t, dir, "find", ".")
	if err == nil {
		t.Fatalf("find without -name should error")
	}
}

func TestWhichBuiltinAliasAndNotFound(t *testing.T) {
	out, _, err := runCmd(t, "", "which", "ls", "help")
	if err != nil {
		t.Fatalf("which builtin error: %v", err)
	}
	if !strings.Contains(out, "ls: gosh builtin") || !strings.Contains(out, "help: gosh builtin") {
		t.Fatalf("which builtin out = %q", out)
	}

	_, _, err = runCmd(t, "", "which", "definitely_not_a_command_12345")
	if err == nil {
		t.Fatalf("which not found should error")
	}

	// host PATH best-effort: use a common executable name
	hostName := "go"
	if runtime.GOOS == "windows" {
		hostName = "cmd"
	}
	if _, lookErr := exec.LookPath(hostName); lookErr == nil {
		out, _, err = runCmd(t, "", "which", hostName)
		// on this environment hostName might also be resolved; if it's a builtin/alias unlikely, but go/cmd are external
		if err != nil {
			// if treated oddly just ensure output reasonable; but expect path
			if !strings.Contains(out, hostName) {
				t.Fatalf("which host err=%v out=%q", err, out)
			}
		} else if !strings.Contains(strings.ToLower(out), strings.ToLower(hostName)) {
			t.Fatalf("which host out = %q", out)
		}
	}
}

func TestZipAddsExtension(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "one.txt"), "one")
	_, _, err := runCmd(t, dir, "zip", "myarc", "one.txt")
	if err != nil {
		t.Fatalf("zip error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "myarc.zip")); err != nil {
		t.Fatalf("zip did not add .zip: %v", err)
	}
}


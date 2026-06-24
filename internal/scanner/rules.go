package scanner

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var knownMaliciousHashes = map[string]string{
	"e99a18c428cb38d5f260853678922e03": "test-malware-eicar",
}

var suspiciousPatterns = []struct {
	Name     string
	Pattern  string
	Severity string
}{
	{Name: "crypto_miner", Pattern: "stratum+tcp://", Severity: "high"},
	{Name: "reverse_shell", Pattern: "/dev/tcp/", Severity: "critical"},
	{Name: "reverse_shell", Pattern: "bash -i", Severity: "critical"},
	{Name: "reverse_shell", Pattern: "nc -e", Severity: "critical"},
	{Name: "reverse_shell", Pattern: "ncat -e", Severity: "critical"},
	{Name: "privilege_escalation", Pattern: "chmod +s", Severity: "high"},
	{Name: "privilege_escalation", Pattern: "chown root", Severity: "medium"},
	{Name: "data_exfil", Pattern: "curl -F", Severity: "medium"},
	{Name: "data_exfil", Pattern: "wget --post", Severity: "medium"},
}

func (s *Scanner) scanFile(path string) ([]Finding, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() || info.Size() == 0 {
		return nil, nil
	}

	var findings []Finding

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hmd5 := md5.New()
	hsha := sha256.New()
	tr := io.MultiWriter(hmd5, hsha)

	if _, err := io.Copy(tr, f); err != nil {
		return nil, err
	}

	hashMD5 := fmt.Sprintf("%x", hmd5.Sum(nil))
	hashSHA := fmt.Sprintf("%x", hsha.Sum(nil))

	if _, ok := knownMaliciousHashes[hashMD5]; ok {
		findings = append(findings, Finding{
			File:       path,
			Threat:     "known_malware",
			Severity:   "critical",
			Md5:        hashMD5,
			Sha256:     hashSHA,
			Size:       info.Size(),
			Confidence: 1.0,
		})
	}

	f.Seek(0, io.SeekStart)
	content, err := io.ReadAll(f)
	if err == nil {
		text := string(content)
		for _, sp := range suspiciousPatterns {
			if strings.Contains(text, sp.Pattern) {
				findings = append(findings, Finding{
					File:       path,
					Threat:     sp.Name,
					Severity:   sp.Severity,
					Md5:        hashMD5,
					Sha256:     hashSHA,
					Size:       info.Size(),
					Confidence: 0.8,
					Rules:      []string{sp.Pattern},
				})
			}
		}
	}

	return findings, nil
}

func (s *Scanner) shouldExclude(path string) bool {
	for _, excl := range s.cfg.Exclude {
		matched, _ := filepath.Match(excl, path)
		if matched {
			return true
		}
	}
	excludeDirs := []string{"/proc", "/sys", "/dev", "/run", "/tmp/.X11-unix"}
	for _, ed := range excludeDirs {
		if strings.HasPrefix(path, ed) {
			return true
		}
	}
	return false
}

func (s *Scanner) scanDir(path string) ([]Finding, error) {
	var findings []Finding

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if s.shouldExclude(p) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		f, err := s.scanFile(p)
		if err != nil {
			return nil
		}
		findings = append(findings, f...)
		return nil
	})

	return findings, err
}

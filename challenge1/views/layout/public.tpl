<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Government Public Disclosure</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <header>
            <a href="/" class="logo">GOV.DISCLOSURE // PUBLIC ACCESS</a>
            <nav>
                <a href="/staff/login">[STAFF LOGIN]</a>
            </nav>
        </header>

        <main>
            {{.LayoutContent}}
        </main>

        <footer>
            <p>&copy; 2025 Government Disclosure Office. RESTRICTED ACCESS.</p>
        </footer>
    </div>
</body>
</html>

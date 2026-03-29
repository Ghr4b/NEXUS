<div class="auth-container" style="max-width: 400px; margin: 0 auto;">
    <h2>// STAFF ACCESS CONTROL</h2>
    
    {{if .Error}}
        <p class="alert-error">{{.Error}}</p>
    {{end}}

    <form action="/staff/login" method="post">
        {{.xsrfdata}}
        <label>AGENT ID (USERNAME)</label>
        <input type="text" name="username" required>
        
        <label>ACCESS CODE (PASSWORD)</label>
        <input type="password" name="password" required>
        
        <button type="submit">AUTHENTICATE</button>
    </form>
    
    <p><small style="text-align: center; display: block;">UNAUTHORIZED ACCESS IS A FELONY.</small></p>
    <p><a href="/staff/register">REQUEST NEW AGENT CLEARANCE</a></p>
</div>

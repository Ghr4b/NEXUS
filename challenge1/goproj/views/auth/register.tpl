<div class="auth-container" style="max-width: 500px; margin: 0 auto;">
    <h2>// NEW AGENT REGISTRATION</h2>
    
    {{if .Error}}
        <p class="alert-error">{{.Error}}</p>
    {{end}}

    <form action="/staff/register" method="post">
        {{.xsrfdata}}
        <label>USERNAME</label>
        <input type="text" name="username" required>
        
        <label>PASSWORD</label>
        <input type="password" name="password" required>
        
        <label>FIRST NAME</label>
        <input type="text" name="first_name" required>
        
        <label>LAST NAME</label>
        <input type="text" name="last_name" required>
        
        <label>EMAIL</label>
        <input type="email" name="email" required>
        
        <button type="submit">SUBMIT FOR CLEARANCE</button>
    </form>
</div>

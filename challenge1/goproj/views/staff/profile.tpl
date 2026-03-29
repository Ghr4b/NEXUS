<h2>// OPERATIVE PROFILE: {{.User.Username}}</h2>

{{if .Success}}
    <p class="alert-success" style="color: var(--secondary-color); border: 1px solid var(--secondary-color); padding: 10px;">{{.Success}}</p>
{{end}}

{{if .Error}}
    <p class="alert-error">{{.Error}}</p>
{{end}}

<div class="file-card">
    <form action="/staff/profile" method="post">
        {{.xsrfdata}}
        <label>FIRST NAME</label>
        <input type="text" name="first_name" value="{{.User.FirstName}}" required>
        
        <label>LAST NAME</label>
        <input type="text" name="last_name" value="{{.User.LastName}}" required>
        
        <label>EMAIL</label>
        <input type="email" name="email" value="{{.User.Email}}" required>
        
        <label>DEPARTMENT</label>
        <input type="text" value="{{if .Staff.Department}}{{.Staff.Department.Name}}{{else}}UNASSIGNED{{end}}" disabled style="opacity: 0.5;">
        
        <div style="margin-top: 20px; border-top: 1px solid #333; padding-top: 20px;">
            <h4>// SECURITY UPDATE</h4>
            <label>NEW PASSWORD (LEAVE BLANK TO KEEP CURRENT)</label>
            <input type="password" name="password" autocomplete="new-password">
        </div>
        
        <button type="submit" style="margin-top: 20px;">UPDATE RECORDS</button>
    </form>
</div>

<div class="actions">
    <a href="/staff/auditlog?Staff__User__Id={{.User.Id}}">[VIEW MY AUDIT HISTORY]</a>
</div>

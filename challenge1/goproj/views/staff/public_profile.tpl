<h2>// AGENT PROFILE</h2>

<div class="file-card">
    <div class="kv-row">
        <strong>AGENT ID:</strong> {{.Staff.User.Username}}
    </div>
    <div class="kv-row">
        <strong>NAME:</strong> {{.Staff.User.FirstName}} {{.Staff.User.LastName}}
    </div>
    <div class="kv-row">
        <strong>DEPARTMENT:</strong> {{if .Staff.Department}}{{.Staff.Department.Name}}{{else}}UNASSIGNED{{end}}
    </div>
    <div class="kv-row">
        <strong>STATUS:</strong> {{if .Staff.User.IsActive}}ACTIVE{{else}}INACTIVE{{end}}
    </div>
</div>

<div class="actions">
    <a href="/staff/auditlog?Staff__User__Id={{.Staff.User.Id}}"><button>VIEW ACTIVITY LOG</button></a>
</div>

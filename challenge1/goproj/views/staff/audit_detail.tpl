<h2>// AUDIT RECORD #{{.Audit.Id}}</h2>

{{if .Audit}}
<div class="file-card">
    <div class="kv-row">
        <strong>TIMESTAMP:</strong> {{date .Audit.Timestamp "Y-m-d H:i:s"}}
    </div>
    <div class="kv-row">
        <strong>AGENT:</strong> 
        <a href="/staff/profile/{{.Audit.Staff.User.Id}}">
            {{.Audit.Staff.User.FirstName}} {{.Audit.Staff.User.LastName}} ({{.Audit.Staff.User.Username}})
        </a>
    </div>
    <div class="kv-row">
        <strong>DEPARTMENT:</strong> {{if .Audit.Staff.Department}}{{.Audit.Staff.Department.Name}}{{else}}N/A{{end}}
    </div>
    <div class="kv-row">
        <strong>ACTION:</strong> {{.Audit.Action}}
    </div>
    <div class="kv-row">
        <strong>TARGET:</strong> {{.Audit.TargetType}} #{{.Audit.TargetId}}
    </div>
    
    <div class="content" style="margin-top: 20px; border-top: 1px dashed var(--secondary-color); padding-top: 10px;">
        <h3>// LOG MESSAGE</h3>
        <p style="font-family: monospace; white-space: pre-wrap;">{{.Audit.Message}}</p>
    </div>
</div>
{{else}}
<p class="alert-error">{{.Error}}</p>
{{end}}

<div class="actions">
    <a href="/staff/auditlog">[< BACK TO LOGS]</a>
</div>

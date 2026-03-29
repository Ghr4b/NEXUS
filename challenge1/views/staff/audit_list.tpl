<h2>// SYSTEM AUDIT LOGS</h2>

<div class="filter-panel" style="margin-bottom: 20px;">
    <form action="/staff/auditlog" method="get" style="display: flex; gap: 10px;">
        <input type="text" name="Staff__User__Username" placeholder="FILTER BY USERNAME" value="">
        <input type="text" name="Action" placeholder="FILTER BY ACTION" value="">
        <button type="submit">FILTER</button>
        <a href="/staff/auditlog"><button type="button" style="background: #333;">RESET</button></a>
    </form>
</div>

{{if .Logs}}
<table>
    <thead>
        <tr>
            <th>TIMESTAMP</th>
            <th>AGENT</th>
            <th>ACTION</th>
            <th>TARGET</th>
            <th>DETAILS</th>
        </tr>
    </thead>
    <tbody>
        {{range .Logs}}
        <tr>
            <td>{{date .Timestamp "Y-m-d H:i:s"}}</td>
            <td><a href="/staff/profile/{{.Staff.User.Id}}">{{.Staff.User.Username}}</a></td>
            <td>{{.Action}}</td>
            <td>{{.TargetType}} #{{.TargetId}}</td>
            <td><a href="/staff/auditlog/{{.Id}}">[VIEW]</a></td>
        </tr>
        {{end}}
    </tbody>
</table>
{{else}}
<p class="alert-error">NO RECORDED EVENTS MATCHING CRITERIA.</p>
{{end}}

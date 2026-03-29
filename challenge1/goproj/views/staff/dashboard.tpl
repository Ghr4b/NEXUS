<h2>// OPERATIVE DASHBOARD</h2>

<div class="actions">
    <a href="/staff/create"><button>+ INITIALIZE NEW CASE FILE</button></a>
</div>

<h3>// YOUR FILES</h3>

{{if .Files}}
<table>
    <thead>
        <tr>
            <th>UUID</th>
            <th>TITLE</th>
            <th>STATUS</th>
            <th>CREATED</th>
            <th>ACTIONS</th>
        </tr>
    </thead>
    <tbody>
        {{range .Files}}
        <tr>
            <td class="hash-code">{{.Uuid}}</td>
            <td>{{.Title}}</td>
            <td>{{if .IsPublished}}<span style="color:var(--secondary-color)">PUBLISHED</span>{{else}}<span style="color:var(--accent-color)">CLASSIFIED</span>{{end}}</td>
            <td>{{date .CreatedAt "Y-m-d"}}</td>
            <td>
                <a href="/staff/files/{{.Uuid}}">[MANAGE]</a>
            </td>
        </tr>
        {{end}}
    </tbody>
</table>
{{else}}
<p>NO FILES ASSIGNED.</p>
{{end}}

<div class="file-detail">
    <a href="/">[< BACK TO ARCHIVE]</a>
    
    <h1>CASE FILE: {{.File.Title}}</h1>
    <div class="meta-data">
        <p class="hash-code">UUID: {{.File.Uuid}}</p>
        <p>DATE: {{date .File.CreatedAt "Y-m-d H:i"}}</p>
    </div>

    <div class="content file-card">
        <h3>// DESCRIPTION</h3>
        <p>{{.File.Description}}</p>
    </div>

    <h3>// ATTACHMENTS</h3>
    {{if .File.Attachments}}
        <table>
            <thead>
                <tr>
                    <th>FILE NAME</th>
                    <th>SIZE</th>
                    <th>SHA-256 FINGERPRINT (VERIFY INTEGRITY)</th>
                    <th>ACTION</th>
                </tr>
            </thead>
            <tbody>
                {{range .File.Attachments}}
                <tr>
                    <td>{{.FileName}}</td>
                    <td>{{.FileSize}} B</td>
                    <td class="hash-code">{{.Sha256Hash}}</td>
                    <td><a href="/static/uploads/{{.FileName}}" download>DOWNLOAD</a></td>
                </tr>
                {{end}}
            </tbody>
        </table>
    {{else}}
        <p>NO ATTACHMENTS ON RECORD.</p>
    {{end}}
</div>

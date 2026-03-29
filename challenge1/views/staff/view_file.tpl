<div class="file-detail">
    <a href="/staff/dashboard">[< BACK TO DASHBOARD]</a>
    
    <h1>MANAGE CASE: {{.File.Title}}</h1>
    <div class="meta-data">
        <p class="hash-code">UUID: {{.File.Uuid}}</p>
        <p>STATUS: {{if .File.IsPublished}}PUBLIC{{else}}HIDDEN{{end}}</p>
    </div>

    <div class="content file-card">
        <h3>// DESCRIPTION</h3>
        <p>{{.File.Description}}</p>
    </div>

    <div class="actions" style="margin-bottom: 20px; display: flex; gap: 10px;">
        <a href="/staff/files/{{.File.Uuid}}/update"><button>EDIT CASE DETAILS</button></a>
        
        <a href="/staff/auditlog?TargetType=DisclosureFile&TargetId={{.File.Id}}"><button style="background-color: var(--secondary-color);">VIEW FILE AUDIT LOG</button></a>

        <form action="/staff/files/{{.File.Uuid}}/delete" method="post" onsubmit="return confirm('ARE YOU SURE YOU WANT TO DELETE THIS CASE? THIS ACTION CANNOT BE UNDONE.');">
            {{.xsrfdata}}
            <button type="submit" style="background-color: var(--accent-color);">DELETE CASE FILE</button>
        </form>
    </div>

    <h3>// EVIDENCE ATTACHMENTS</h3>
    
    <div class="upload-section" style="border: 1px dashed var(--secondary-color); padding: 20px; margin-bottom: 20px;">
        <h4>UPLOAD NEW EVIDENCE</h4>
        <form action="/staff/upload" method="post" enctype="multipart/form-data">
            {{.xsrfdata}}
            <input type="hidden" name="uuid" value="{{.File.Uuid}}">
            <label>SELECT FILE</label>
            <input type="file" name="attachment" required style="color: var(--text-color); margin-bottom: 10px;">
            <button type="submit">UPLOAD & FINGERPRINT</button>
        </form>
    </div>

    {{if .File.Attachments}}
        <table>
            <thead>
                <tr>
                    <th>FILE NAME</th>
                    <th>SIZE</th>
                    <th>SHA-256 FINGERPRINT</th>
                    <th>TIMESTAMP</th>
                    <th>ACTION</th>
                </tr>
            </thead>
            <tbody>
                {{range .File.Attachments}}
                <tr>
                    <td>{{.FileName}}</td>
                    <td>{{.FileSize}} B</td>
                    <td class="hash-code">{{.Sha256Hash}}</td>
                    <td>{{date .UploadedAt "Y-m-d H:i"}}</td>
                    <td>
                        <form action="/staff/attachments/{{.Id}}/delete" method="post" style="display:inline;" onsubmit="return confirm('DELETE EVIDENCE?');">
                            {{$.xsrfdata}}
                            <input type="hidden" name="file_uuid" value="{{$.File.Uuid}}">
                            <button type="submit" style="background-color: transparent; border: 1px solid var(--accent-color); color: var(--accent-color); padding: 5px;">[X]</button>
                        </form>
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>
    {{else}}
        <p>NO EVIDENCE ATTACHED.</p>
    {{end}}
</div>

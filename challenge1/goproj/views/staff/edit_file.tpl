<h2>// UPDATE CASE FILE</h2>

{{if .Error}}
    <p class="alert-error">{{.Error}}</p>
{{end}}

<form action="/staff/files/{{.File.Uuid}}/update" method="post">
    {{.xsrfdata}}
    <label>CASE TITLE</label>
    <input type="text" name="title" value="{{.File.Title}}" required>
    
    <label>DESCRIPTION (TEXT CONTENT)</label>
    <textarea name="description" rows="10" required>{{.File.Description}}</textarea>
    
    <label>STATUS</label>
    <select name="is_published">
        <option value="false" {{if not .File.IsPublished}}selected{{end}}>CLASSIFIED (HIDDEN)</option>
        <option value="true" {{if .File.IsPublished}}selected{{end}}>DECLASSIFIED (PUBLIC)</option>
    </select>
    
    <div style="margin-top: 20px;">
        <button type="submit">SAVE CHANGES</button>
        <a href="/staff/files/{{.File.Uuid}}" style="margin-left: 20px;">[CANCEL]</a>
    </div>
</form>

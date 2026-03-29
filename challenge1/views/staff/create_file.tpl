<h2>// NEW INVESTIGATION CASE</h2>

{{if .Error}}
    <p class="alert-error">{{.Error}}</p>
{{end}}

<form action="/staff/create" method="post">
    {{.xsrfdata}}
    <label>CASE TITLE</label>
    <input type="text" name="title" required>
    
    <label>DESCRIPTION (TEXT CONTENT)</label>
    <textarea name="description" rows="10" required></textarea>
    
    <label>STATUS</label>
    <select name="is_published">
        <option value="false">CLASSIFIED (HIDDEN)</option>
        <option value="true">DECLASSIFIED (PUBLIC)</option>
    </select>
    
    <button type="submit">CREATE RECORD</button>
</form>

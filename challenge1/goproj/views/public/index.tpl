<div class="search-box">
    <form action="/" method="get">
        <input type="text" name="search" placeholder="SEARCH DISCLOSURE ARCHIVE..." value="{{.Search}}">
        <button type="submit">SEARCH</button>
    </form>
</div>

<h2>// DISCLOSURE FILES</h2>

{{if .Files}}
    {{range .Files}}
    <div class="file-card">
        <h3><a href="/files/{{.Uuid}}">{{.Title}}</a></h3>
        <p class="hash-code">CASE ID: {{.Uuid}}</p>
        <p class="redacted-text">{{substr .Description 0 100}}...</p>
        <p><small>RELEASED: {{date .CreatedAt "Y-m-d H:i"}}</small></p>
    </div>
    {{end}}
{{else}}
    <p class="alert-error">NO RECORDS FOUND matching query.</p>
{{end}}

<h2>// SUBMIT REPORT</h2>

{{if .Error}}
    <p class="alert-error">{{.Error}}</p>
{{end}}

{{if .Success}}
    <p class="alert-success">{{.Success}}</p>
{{end}}

<form action="/staff/report" method="post" id="reportForm">

    <input type="hidden" name="uri" id="uri_field">

    <label>POST UUID</label>
    <input type="text" id="post_uuid" placeholder="Enter post UUID" required oninput="updateUrl()">

    <label>REPORT DESCRIPTION</label>
    <textarea name="content" rows="8" placeholder="Describe what is wrong with this post..." required></textarea>

    <button type="submit">SUBMIT REPORT</button>
</form>

<script>
function updateUrl() {
    var uuid = document.getElementById('post_uuid').value;
    document.getElementById('uri_field').value = 'staff/files/' + uuid;
}
</script>

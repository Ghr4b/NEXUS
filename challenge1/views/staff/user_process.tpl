<h2>// PROCESS USER: {{.User.Username}}</h2>

<div style="border: 1px solid var(--secondary-color); padding: 20px; margin-bottom: 20px;">
    <h3>USER DETAILS</h3>
    <p><strong>ID:</strong> {{.User.Id}}</p>
    <p><strong>Username:</strong> {{.User.Username}}</p>
    <p><strong>First Name:</strong> {{.User.FirstName}}</p>
    <p><strong>Last Name:</strong> {{.User.LastName}}</p>
    <p><strong>Email:</strong> {{.User.Email}}</p>
</div>

<div id="action-panel" style="border: 1px solid var(--accent-color); padding: 20px;">
    <h3>ACTIONS</h3>

    <form id="approve-form" action="/staff/management/approve" method="post" style="margin-bottom: 20px;">
        {{.xsrfdata}}
        <input type="hidden" name="user_id" value="{{.User.Id}}">
        <label>ASSIGN DEPARTMENT:</label>
        <select name="department_id" style="margin-right: 10px;">
            {{range .Departments}}
            <option value="{{.Id}}">{{.Name}}</option>
            {{end}}
        </select>
        <button type="submit">APPROVE & ASSIGN</button>
    </form>

    <form id="reject-form" action="/staff/management/reject" method="post">
        {{.xsrfdata}}
        <input type="hidden" name="user_id" value="{{.User.Id}}">
        <button type="submit" style="background-color: var(--accent-color);">REJECT & DELETE</button>
    </form>

    <br>
    <a href="/staff/management" style="color: white; text-decoration: underline;">← Back to List</a>
</div>

<h2>// STAFF REQUEST MANAGEMENT</h2>

<form style="display:none;">{{.xsrfdata}}</form>

<div class="filter-panel" style="border: 1px solid var(--secondary-color); padding: 15px; margin-bottom: 20px;">
    <h3>FILTER REQUESTS</h3>
    <div id="filter-rows">
        <div class="filter-row" style="margin-bottom: 10px; display: flex; gap: 10px;">
            <select class="filter-field">
                <option value="Username">USERNAME</option>
                <option value="Email">EMAIL</option>
                <option value="FirstName">FIRST NAME</option>
                <option value="LastName">LAST NAME</option>
            </select>
            <select class="filter-operator">
                <option value="__icontains">CONTAINS</option>
                <option value="__istartswith">STARTS WITH</option>
                <option value="__exact">EXACT MATCH</option>
            </select>
            <input type="text" class="filter-value" placeholder="VALUE">
            <button onclick="removeRow(this)" style="padding: 5px 10px; background: #333; color: white;">-</button>
        </div>
    </div>
    <button onclick="addRow()" style="background: #333; color: white; margin-right: 10px;">+ ADD FILTER</button>
    <button onclick="applyFilters()">APPLY FILTERS</button>
</div>

<div id="results-area">
    <table>
        <thead>
            <tr>
                <th>ID</th>
                <th>USERNAME</th>
                <th>FULL NAME</th>
                <th>EMAIL</th>
                <th>ACTIONS</th>
            </tr>
        </thead>
        <tbody id="results-body">
            </tbody>
    </table>
</div>

<script>
    // SECURE: Go's html/template engine detects this is inside a JavaScript string ("")
    // and will safely escape quotes and special characters to prevent breakout payloads.
    var savedFilter = "{{.SavedFilter}}";

    function addRow() {
        const row = document.querySelector('.filter-row').cloneNode(true);
        row.querySelector('.filter-value').value = '';
        document.getElementById('filter-rows').appendChild(row);
    }

    function removeRow(btn) {
        if(document.querySelectorAll('.filter-row').length > 1) {
            btn.parentElement.remove();
        }
    }

    function applyFilters() {
        const rows = document.querySelectorAll('.filter-row');
        const filters = {};

        rows.forEach(row => {
            const field = row.querySelector('.filter-field').value;
            const operator = row.querySelector('.filter-operator').value;
            const value = row.querySelector('.filter-value').value;
            if (value) {
                filters[field + operator] = value;
            }
        });

        const xsrfInput = document.querySelector('input[name="_xsrf"]');
        const xsrfToken = xsrfInput ? xsrfInput.value : '';

        fetch('/staff/management/search', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Xsrftoken': xsrfToken
            },
            body: JSON.stringify(filters)
        })
        .then(response => response.json())
        .then(data => {
            const tbody = document.getElementById('results-body');
            tbody.innerHTML = '';

            if (data && data.length > 0) {
                data.forEach(user => {
                    // SECURE: Building DOM elements directly avoids HTML string parsing vulnerabilities
                    const tr = document.createElement('tr');

                    const idTd = document.createElement('td');
                    idTd.textContent = user.Id;
                    tr.appendChild(idTd);

                    const usernameTd = document.createElement('td');
                    usernameTd.textContent = user.Username; // TextContent automatically neutralizes malicious scripts
                    tr.appendChild(usernameTd);

                    const nameTd = document.createElement('td');
                    nameTd.textContent = `${user.FirstName} ${user.LastName}`;
                    tr.appendChild(nameTd);

                    const emailTd = document.createElement('td');
                    emailTd.textContent = user.Email;
                    tr.appendChild(emailTd);

                    const actionTd = document.createElement('td');
                    const btn = document.createElement('button');
                    btn.textContent = '[PROCESS]';

                    // SECURE: Passing values natively in memory, no inline HTML attributes to break out of
                    btn.addEventListener('click', () => {
                        window.location.href = '/staff/management/user?id=' + encodeURIComponent(user.Id);
                    });
                    actionTd.appendChild(btn);

                    tr.appendChild(actionTd);
                    tbody.appendChild(tr);
                });
            } else {
                tbody.innerHTML = '<tr><td colspan="5">NO PENDING REQUESTS MATCHING CRITERIA.</td></tr>';
            }
        })
        .catch(error => {
            console.error("Fetch error:", error);
        });
    }

    if (savedFilter) {
        document.querySelector('.filter-value').value = savedFilter;
    }
    applyFilters();
</script>

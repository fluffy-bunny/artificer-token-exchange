{{define "views/accounts/index"}}
{{template "header" .}}
{{template "navbar" .}}

<body>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5">
        <h1>Accounts</h1>
        <button type="button" id="btnForceRefresh">Force Refresh</button>
    </div>
</div>
</body>

{{template "footer" .}}
     <script>
	    // get reference to button
	    var btn = document.getElementById("btnForceRefresh");
	    // add event listener for the button, for action "click"
	    btn.addEventListener("click", postAccountsForceRefresh);

    </script>
{{end}}
{{define "views/auth/unauthorized/index"}}
{{template "header" .}}
{{template "navbar" .}}
<body>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5">
        <h1>Unauthorized</h1>
    </div>
</div>
</body>
    
{{template "footer" .}}
{{end}}
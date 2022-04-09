{{define "views/deep/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}
<body>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5">
        <h1>Deep</h1>
        <div class="alert alert-success" role="alert">
        <div>ID:{{ .params.ID }}</div>
        <div>Name:{{ .params.Name }}</div>
        </div>
    </div>
</div>
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}
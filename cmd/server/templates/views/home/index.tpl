{{define "views/home/index"}}
{{template "header" .}}
{{template "navbar" .}}
<body>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5">
        <h1>A Bootstrap 5 Starter Template</h1>
        <p class="lead">A complete project boilerplate built with Bootstrap</p>
        <p>Bootstrap v5.1.3</p>
        <div class="alert alert-success" role="alert">
        {{range $claim := .claims}}
            <div>{{ $claim }}</div>
        {{end}}
        </div>
    </div>
</div>
</body>
    
{{template "footer" .}}
{{end}}
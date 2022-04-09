{{define "views/home/index"}}
{{template "html_begin" .}}
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
                 <table class="table table-striped">
                <thead>
                <tr>
                <th class="text-start" scope="col">#</th>
                <th class="text-start" scope="col">Type</th>
                <th class="text-start" scope="col">Value</th>
                </tr>
            </thead>
            <tbody>
            {{range $idx,$claim := .claims}}
                <tr>
                <th class="text-start" scope="row">{{$idx}}</th>
                <td class="text-start">{{$claim.Type}}</td>
                <td class="text-start">{{$claim.Value}}</td>
                </tr>
            {{end}}
            </tbody>
            </table>
        </div>
    </div>
</div>
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}
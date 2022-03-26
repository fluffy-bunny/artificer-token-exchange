{{define "content"}}
{{template "header" .}}

<!-- Responsive navbar-->
<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
    <div class="container">
        <a class="navbar-brand" href="{{ .paths.Home }}">Echo Starter</a>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation"><span class="navbar-toggler-icon"></span></button>
        <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <ul class="navbar-nav ms-auto mb-2 mb-lg-0">
                <li class="nav-item"><a class="nav-link active" aria-current="page" href="{{ .paths.Home }}">Home</a></li>
                <li class="nav-item dropdown">
                    <a class="nav-link dropdown-toggle" id="navbarDropdown" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">Account</a>
                    <ul class="dropdown-menu dropdown-menu-end" aria-labelledby="navbarDropdown">
                    {{ if .user }}
                        <li><a class="dropdown-item" href="{{ .paths.Logout }}">Logout</a></li>
                    {{ else }}
                        <li><a class="dropdown-item" href="{{ .paths.Login }}">Login</a></li>
                    {{end}}
                    </ul>
                </li>
            </ul>
        </div>
    </div>
</nav>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5">
        <h1>A Bootstrap 5 Starter Template</h1>
        <p class="lead">A complete project boilerplate built with Bootstrap</p>
        <p>Bootstrap v5.1.3</p>
        {{ if .user }}
        <div class="alert alert-success" role="alert">
            <div>ID: {{ .user.ID }}</div>
            <div>Email: {{ .user.Email }}</div>
            <div>Name: {{ .user.Name }}</div>
        </div>
        {{end}}
    </div>
</div>
 
    
{{template "footer" .}}
{{end}}
{{define "content"}}
{{template "header" .}}
<div>
<header>
  <!-- Fixed navbar -->
  <nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
    <div class="container-fluid">
      <a class="navbar-brand" href="/">HOME</a>
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarCollapse">
        <ul class="navbar-nav me-auto mb-2 mb-md-0">
          <li class="nav-item">
            <a class="nav-link active" aria-current="page" href="/">Home</a>
          </li>
          
        </ul>
        <form class="d-flex">
          {{ if .user }}
          <li class="nav-item">
            <a class="nav-link" href="/logout">Logout</a>
          </li>
          {{ else }}
           <li class="nav-item">
            <a class="nav-link" href="/login">Login</a>
          </li>
          {{end}}
        </form>
      </div>
    </div>
  </nav>
</header>
<main class="flex-shrink-0">
  <div class="container">
    {{ if .user }}
        <div class="alert alert-success" role="alert">
            <div>ID: {{ .user.ID }}</div>
            <div>Email: {{ .user.Email }}</div>
            <div>Name: {{ .user.Name }}</div>
        </div>
 
     
    {{end}}
  </div>
</main>

    
</div>
{{template "footer" .}}
{{end}}
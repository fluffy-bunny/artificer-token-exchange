{{define "views/auth/profiles/index"}}
{{template "html_begin" .}}

{{template "header" .}}
{{template "navbar" .}}
<body>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5">
        <h1>Switch Profiles</h1>
         
        <div class="alert alert-success" role="alert">
        {{ $csrf := .security.CSRF }}
        {{range $profile := .profiles}}
            <form id="account" method="post">
            <input type="hidden" name="csrf" value="{{ $csrf }}">
            <div>
               
                <input type="hidden" name="profile" value="{{ $profile }}">
            </div>
            <div>
            </div>
                <button id="submit" type="submit" class="w-100 btn btn-sm btn-primary">Switch to {{ $profile }}</button>
           
            </form>
        {{end}}
        </div>
    </div>
</div>
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}
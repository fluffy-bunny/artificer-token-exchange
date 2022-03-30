{{define "csfr"}}
 
   {{ $csrf := .security.CSRF }}
    <script>
    let csrf = '{{$csrf}}';
    </script>
{{end}}
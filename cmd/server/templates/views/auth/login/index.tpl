{{define "views/auth/login/index"}}
{{template "html_begin" .}}
{{template "header" .}}
 
<body>
<!-- Page content-->
<div class="container">
    <div class="text-left mt-5">
        <button   name="button" onclick="getMyRedirectURL()">Submit OIDC Request</button>
        <button   name="button" onclick="getSessionData()">GetSessionData</button>
    </div>
    <div class="text-left mt-5">
		<h2>Output:</h2>
		<div id="json"></div>
    </div>
</div>


 

 {{ $url := .url }}
<script>
    async function getSessionData() {
        let url = '/api/v1/dev?directive=session';  
        try {
            let res = await fetch(url,{
                method: 'GET',
                credentials: 'include'
            });
            payload =  await res.json();
            jsonViewer.showJSON(payload);

            return payload
        } catch (error) {
            console.log(error);
            alert(error)
        }
    }

    function getMyRedirectURL() {
       document.write("It will redirect within 1 seconds.....please wait...");//it will redirect after 3 seconds
       setTimeout(function() {
            window.location = {{$url}};
       }, 1000);
    }
</script>
</body>
{{template "footer" .}}
{{template "html_end" .}}
{{end}}
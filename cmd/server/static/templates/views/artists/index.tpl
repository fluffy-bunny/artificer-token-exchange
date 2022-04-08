{{define "views/artists/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}

<body>
<!-- Page content-->
<div class="container">
    <div class="text-left mt-5">
        <h1>Artists</h1>
        <button type="button" id="btnArtists">Artists</button>
        <button type="button" id="btnArtist">Artist</button>
        <button type="button" id="btnAlbums">Albums</button>
        <button type="button" id="btnPostArtist">Post Artist</button>
        <button type="button" id="btnPostArtistForgotCsrf">Post Artist - Forgot CSRF</button>
 
		
    </div>
	{{ $config := .config }}
 
	<div class="text-left mt-5">
		<h2>Output:</h2>
		<div id="json"></div>
    </div>

</div>
</body>

{{template "footer" .}}
     <script>

 
	    // get reference to button
	    var btn = document.getElementById("btnArtists");
	    // add event listener for the button, for action "click"
	    btn.addEventListener("click", getArtists);

         // get reference to button
	    var btn = document.getElementById("btnArtist");
	    // add event listener for the button, for action "click"
	    btn.addEventListener("click", getArtist);

         // get reference to button
	    var btn = document.getElementById("btnAlbums");
	    // add event listener for the button, for action "click"
	    btn.addEventListener("click", getAlbums);

         // get reference to button
	    var btn = document.getElementById("btnPostArtist");
	    // add event listener for the button, for action "click"
	    btn.addEventListener("click", postArtist);

         // get reference to button
	    var btn = document.getElementById("btnPostArtistForgotCsrf");
	    // add event listener for the button, for action "click"
	    btn.addEventListener("click", postArtistForgotCsrf);
 

  
    </script>
{{template "html_end" .}}
{{end}}
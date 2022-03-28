{{define "views/artists/index"}}
{{template "header" .}}
{{template "navbar" .}}
<body>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5">
        <h1>Artists</h1>
        <button type="button" id="btnArtists">Artists</button>
        <button type="button" id="btnArtist">Artist</button>
        <button type="button" id="btnAlbums">Albums</button>
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

    </script>
{{end}}
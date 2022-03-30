/*!
* Start Bootstrap - Bare v5.0.7 (https://startbootstrap.com/template/bare)
* Copyright 2013-2021 Start Bootstrap
* Licensed under MIT (https://github.com/StartBootstrap/startbootstrap-bare/blob/master/LICENSE)
*/
// This file is intentionally blank
// Use this file to add JavaScript to your project

 

async function getArtists() {
    let url = '/api/v1/artists';
    try {
        let res = await fetch(url,{
            method: 'GET',
            credentials: 'include'
          });
          payload =  await res.json();
          console.log(payload);
          alert(JSON.stringify(payload));
          return payload

    } catch (error) {
        console.log(error);
    }
}
async function postArtist() {
    let url = '/api/v1/artists/1';   
    try {
        let res = await fetch(url,{
            method: 'POST',
            credentials: 'include',
            headers: {"X-Csrf-Token": csrf },
            body: JSON.stringify({ name: 'test' }),
          });
        payload =  await res.json();
        console.log(payload);
        alert(JSON.stringify(payload));
        return payload
    } catch (error) {s
        console.log(error);
    }
}
async function getArtist() {
    let url = '/api/v1/artists/1';   
    try {
        let res = await fetch(url,{
            method: 'GET',
            credentials: 'include'
          });
        payload =  await res.json();
        console.log(payload);
        alert(JSON.stringify(payload));
       

        return payload
    } catch (error) {
        console.log(error);
    }
}
async function getAlbums() {
    let url = '/api/v1/artists/1/albums';  
    try {
        let res = await fetch(url,{
            method: 'GET',
            credentials: 'include'
          });
        payload =  await res.json();
        console.log(payload);
        alert(JSON.stringify(payload));

        return payload
    } catch (error) {
        console.log(error);
    }
}
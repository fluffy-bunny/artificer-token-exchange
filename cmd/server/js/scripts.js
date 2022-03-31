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
        alert(error)
    }
}
async function postArtistForgotCsrf() {
    let url = '/api/v1/artists/1';   
    try {
        let res = await fetch(url,{
            method: 'POST',
            credentials: 'include',
            headers: {
                "Content-Type": "application/json" 
            },
            body: JSON.stringify({ name: 'test' }),
          });
        payload =  await res.json();
        console.log(payload);
        alert(JSON.stringify(payload));
        return payload
    } catch (error) {s
        console.log(error);
        alert(error)
    }
}
async function postArtist() {
    let url = '/api/v1/artists/1';   
    try {
        let res = await fetch(url,{
            method: 'POST',
            credentials: 'include',
            headers: {
                "X-Csrf-Token": csrf,
                "Content-Type": "application/json" 
            },
            body: JSON.stringify({ name: 'test' }),
          });
        payload =  await res.json();
        console.log(payload);
        alert(JSON.stringify(payload));
        return payload
    } catch (error) {s
        console.log(error);
        alert(error)
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
        alert(error)
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
        alert(error)
    }
}
async function postAccountsForceRefresh() {
    let url = '/api/v1/accounts';   
    try {
        let res = await fetch(url,{
            method: 'POST',
            credentials: 'include',
            headers: {
                    "X-Csrf-Token": csrf,
                    "Content-Type": "application/json" 
                },
            body: JSON.stringify({ directive: 'force-refresh' }),
          });
        payload =  await res.json();
        console.log(payload);
        alert(JSON.stringify(payload));
        return payload
    } catch (error) {s
        console.log(error);
        alert(error)
    }
}
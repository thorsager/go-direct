function add_redirect(container_id, redirect, clear_input) {
    let out = document.getElementById('out');
    out.insertAdjacentHTML("beforeend", renderRedirect(redirect));
    if (clear_input) {
        let inputUrl = document.getElementById('input_url');
        inputUrl.value = ''
    }
}

function renderRedirect(reply) {
    return `
<div>
    <a id="${reply.id}" target="_blank" href="${reply.source_url}">${reply.source_url}</a>
    <img alt="Copy to Clipboard" class="small-icon cursor-click" src="/svg/copy-regular.svg" onclick="CopyToClipboard('${reply.id}');"/> -> <a href="${reply.target_url}">${reply.target_url}</a>
</div>
`
}

function CopyToClipboard(containerId) {
    try {
        let node = document.getElementById(containerId);
        console.log("select node:", node);
        let range = document.createRange();
        range.selectNode(node);
        console.log(range);
        window.getSelection().addRange(range);
        if (document.execCommand("copy")) {
            console.log("Copy failed?!")
        } else {
            alert("text copied")
        }
        // window.getSelection().removeAllRanges();
    } catch (err) {
        console.log(err)
    }

}

function create_redirect() {
    if (!validate_input()) {
        return
    }
    let xhr = new XMLHttpRequest();
    xhr.open('post', '/api');
    xhr.onload = function () {
        let rply;
        if (xhr.status !== 200) { // analyze HTTP status of the response
            console.log(`Error ${xhr.status}: ${xhr.statusText}`); // e.g. 404: Not Found
        } else { // show the result
            console.log(`Done, got ${xhr.response.length} bytes`); // responseText is the server
            rply = JSON.parse(xhr.responseText);
            console.log(rply);
            add_redirect('out', rply, true);
        }
    };

    xhr.onerror = function () {
        console.log(`Network Error`)
    };

    xhr.onprogress = function (event) { // triggers periodically
        // event.loaded - how many bytes downloaded
        // event.lengthComputable = true if the server sent Content-Length header
        // event.total - total number of bytes (if lengthComputable)
        console.log(`Received ${event.loaded} of ${event.total}`);
    };
    let inputUrl = document.getElementById('input_url');
    xhr.send(JSON.stringify({code: 307, url: inputUrl.value}))
}

function validate_input() {
    let out = document.getElementById('error');
    let url = document.getElementById('input_url');
    if (!url.checkValidity()) {
        out.innerHTML = url.validationMessage;
        return false
    } else {
        out.innerHTML = '';
        return true
    }

}
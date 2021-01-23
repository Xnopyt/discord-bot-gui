window.darwinMeta = false;
window.darwin = true;

document.addEventListener("keydown", async function(event) {
    if (event.code === "OSLeft") {
        window.darwinMeta = true;
        return
    }
    if (window.darwinMeta) {
        switch (event.code) {
            case "KeyV":
                if ((event.target.nodeName == "INPUT") && ((event.target.type == "password" || event.target.type == "text")) || (event.target.nodeName == "TEXTAREA")) {
                    var clip = await readClipboard();
			        var end = event.target.value.slice(event.target.selectionEnd);
			        var start = event.target.value.slice(0, event.target.selectionStart);
			        event.target.value = start + clip + end;
                }
                break;
            case "KeyC":
                var clip = window.getSelection().toString();
                if (clip != "") {
                    writeClipboard(clip);
                }
                break;
            case "KeyX":
                if (((event.target.nodeName == "INPUT") && ((event.target.type == "password" || event.target.type == "text")) || (event.target.nodeName == "TEXTAREA")) && (event.target.selectionStart != event.target.selectionEnd)) {
                    var end = event.target.value.slice(event.target.selectionEnd);
			        var start = event.target.value.slice(0, event.target.selectionStart);
			        var clip = event.target.value.slice(event.target.selectionStart, event.target.selectionEnd);
			        event.target.value = start + end;
			        writeClipboard(clip);
                }
                break;
            case "KeyA":
                if ((event.target.nodeName == "INPUT") && ((event.target.type == "password" || event.target.type == "text")) || (event.target.nodeName == "TEXTAREA")) {
                    event.target.selectionStart = 0;
                    event.target.selectionEnd = event.target.value.length;
                }
                break
        }
    }
})

document.addEventListener("keyup", function(event) {
    if (event.code === "OSLeft") {
        window.darwinMeta = false;
    }
})
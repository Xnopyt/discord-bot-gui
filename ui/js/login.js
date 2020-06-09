var loginInput = document.getElementById("token");
var button = document.getElementById("loginButton");
var loading = false;
loginInput.addEventListener("keyup", function(event) {
	if (loading) {
		return
	}
	if (event.keyCode === 13) {
		button.innerHTML = '<div class="lds-facebook"><div></div><div></div><div></div></div>';
		loading = true;
		event.preventDefault();
		var loginInput = document.getElementById("token");
		var returnMessage = {};
		returnMessage.type = "connect";
		returnMessage.content = loginInput.value
		wv(JSON.stringify(returnMessage));
	}
});
button.onclick = function() {
	if (loading) {
		return
	}
	loading = true;
	button.innerHTML = '<div class="lds-facebook"><div></div><div></div><div></div></div>';
	event.preventDefault();
	var loginInput = document.getElementById("token");
	var returnMessage = {};
	returnMessage.type = "connect";
	returnMessage.content = loginInput.value
	wv(JSON.stringify(returnMessage));
}

function fail() {
	var tok = document.getElementById("token");
	var lab = document.getElementById("tlabel");
	lab.style.color = "rgb(189, 53, 43)"
	tok.style.border = "1px solid rgb(189, 53, 43)"
	lab.innerHTML = "TOKEN - Login failed."
	button.innerHTML = "Login";
	loading = false;
}

window.ctrlHeld = false;

document.addEventListener("keyup", function(event) {
	if (event.keyCode === 17) {
		window.ctrlHeld = false;
	}
});

document.addEventListener("keydown", function(event) {
	if (event.keyCode === 17) {
		window.ctrlHeld = true;
	}
});

document.addEventListener('contextmenu', function(event) {
	if (window.ctrlHeld) {
		window.ctrlHeld = false;
		return
	}
	event.preventDefault();

	var copy = document.getElementById("copybutton");
	var clip = window.getSelection().toString();
	copy.addEventListener("click", function(event) {
		if (clip != "") {
			writeClipboard(clip);
		}
	}, {once: true})

	var paste = document.getElementById("pastebutton");
	var pasteTarget = event.target
	if ((pasteTarget.nodeName == "INPUT") && ((pasteTarget.type == "password" || pasteTarget.type == "text")) || (pasteTarget.nodeName == "TEXTAREA")) {
		paste.style.display = "block";
		paste.addEventListener("click", async function(event) {
			var clip = await readClipboard();
			var end = pasteTarget.value.slice(pasteTarget.selectionEnd);
			var start = pasteTarget.value.slice(0, pasteTarget.selectionStart);
			pasteTarget.value = start + clip + end;
		}, {once: true})
	} else {
		paste.style.display = "none";
	}

	var cut = document.getElementById("cutbutton");
	if (((pasteTarget.nodeName == "INPUT") && ((pasteTarget.type == "password" || pasteTarget.type == "text")) || (pasteTarget.nodeName == "TEXTAREA")) && (pasteTarget.selectionStart != pasteTarget.selectionEnd)) {
		cut.style.display = "block";
		cut.addEventListener("click", function(event) {
			var end = pasteTarget.value.slice(pasteTarget.selectionEnd);
			var start = pasteTarget.value.slice(0, pasteTarget.selectionStart);
			var clip = pasteTarget.value.slice(pasteTarget.selectionStart, pasteTarget.selectionEnd);
			pasteTarget.value = start + end;
			writeClipboard(clip);
		}, {once: true})
	} else {
		cut.style.display = "none";
	}

	var menu = document.getElementById("contextmenu");
	menu.style.visibility = "hidden";
	menu.style.display = "block";

	var rect = document.body.getBoundingClientRect();

	if ((event.clientY + menu.offsetHeight) > rect.bottom ) {
		menu.style.top = (event.clientY - menu.offsetHeight) + "px";
	} else {
		menu.style.top = event.clientY + "px";
	}

	if ((event.clientX + menu.offsetWidth) > rect.right) {
		menu.style.left = (event.clientX - menu.offsetWidth) + "px";
	} else {
		menu.style.left = event.clientX + "px";
	}

	menu.style.visibility = "visible";
})

document.addEventListener("click", function(event) {
	document.getElementById("contextmenu").style.display = "none";
})

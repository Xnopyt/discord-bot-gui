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

function createAlert(title, message) {
	document.getElementById("alerttitle").innerHTML = title;
	document.getElementById("alertmsg").innerHTML = message;
	document.getElementById("alertbox").style.display = "block";
}

var getClosest = function (elem, selector) {
	if (!Element.prototype.matches) {
		Element.prototype.matches = Element.prototype.msMatchesSelector ||
									Element.prototype.webkitMatchesSelector;
	  }
	for ( ; elem && elem !== document; elem = elem.parentNode ) {
		if ( elem.matches( selector ) ) return elem;
	}
	return null;
};

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

	var showContext = false;
	var menu = document.getElementById("contextmenu");
	menu.style.display = "none";

	var copy = document.getElementById("copybutton");
	var clip = window.getSelection().toString();
	if (clip != "") {
		copy.style.display = "block";
		showContext = true;
		copy.onclick = function(event) {
			writeClipboard(clip);
		}
	} else {
		copy.style.display = "none";
	}

	var paste = document.getElementById("pastebutton");
	var pasteTarget = event.target
	if ((pasteTarget.nodeName == "INPUT") && ((pasteTarget.type == "password" || pasteTarget.type == "text")) || (pasteTarget.nodeName == "TEXTAREA")) {
		paste.style.display = "block";
		showContext = true;
		paste.onclick = async function(event) {
			var clip = await readClipboard();
			var end = pasteTarget.value.slice(pasteTarget.selectionEnd);
			var start = pasteTarget.value.slice(0, pasteTarget.selectionStart);
			pasteTarget.value = start + clip + end;
		}
	} else {
		paste.style.display = "none";
	}

	var cut = document.getElementById("cutbutton");
	if (((pasteTarget.nodeName == "INPUT") && ((pasteTarget.type == "password" || pasteTarget.type == "text")) || (pasteTarget.nodeName == "TEXTAREA")) && (pasteTarget.selectionStart != pasteTarget.selectionEnd)) {
		cut.style.display = "block";
		showContext = true;
		cut.onclick = function(event) {
			var end = pasteTarget.value.slice(pasteTarget.selectionEnd);
			var start = pasteTarget.value.slice(0, pasteTarget.selectionStart);
			var clip = pasteTarget.value.slice(pasteTarget.selectionStart, pasteTarget.selectionEnd);
			pasteTarget.value = start + end;
			writeClipboard(clip);
		}
	} else {
		cut.style.display = "none";
	}

	var del = document.getElementById("deletebutton");
	var msg = getClosest(pasteTarget, ".message");
	if ( (msg != null) && msg.id != "" ) {
		del.style.display = "block";
		showContext = true;
		del.onclick = function(event) {
			document.getElementById("delconfirm").onclick = async function() {
				document.getElementById('deldialog').style.display = 'none';
				document.getElementById('confirmblock').style.display = 'none';
				var err = await deleteMessage(msg.id);
				if (err != "") {
					var x = err.split(",");
					x.shift();
					x = x.join(",");
					try {
						err = JSON.parse(x).message;
					} catch (e) {}
					createAlert("Failed to Delete Message", err);
				}
			};
			document.getElementById("confirmblock").style.display = "block";
			document.getElementById("deldialog").style.display = "block";
		}
	} else {
		del.style.display = "none";
	}

	if (showContext) {
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
	}
})

document.addEventListener("click", function(event) {
	document.getElementById("contextmenu").style.display = "none";
})
var loginInput = document.getElementById("token");
var button = document.getElementById("loginButton");
loginInput.addEventListener("keyup", function(event) {
	if (event.keyCode === 13) {
		event.preventDefault();
		var loginInput = document.getElementById("token");
		var returnMessage = {};
		returnMessage.type = "connect";
		returnMessage.content = loginInput.value
		wv(JSON.stringify(returnMessage));
	}
});
button.onclick = function() {
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
}
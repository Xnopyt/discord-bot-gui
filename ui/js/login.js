var loginInput = document.getElementById("token");
var button = document.getElementById("loginButton");
loginInput.addEventListener("keyup", function(event) {
	if (event.keyCode === 13) {
		event.preventDefault();
		var loginInput = document.getElementById("token");
		binder.connect(loginInput.value);
	}
});
button.onclick = function() {
	event.preventDefault();
	var loginInput = document.getElementById("token");
	binder.connect(loginInput.value);
}

function fail() {
	var tok = document.getElementById("token");
	var lab = document.getElementById("tlabel");
	lab.style.color = "rgb(189, 53, 43)"
	tok.style.border = "1px solid rgb(189, 53, 43)"
	lab.innerHTML = "TOKEN - Login failed."
}
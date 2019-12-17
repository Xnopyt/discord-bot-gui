var loginInput = document.getElementById("token");
var button = document.getElementById("loginButton");
loginInput.addEventListener("keyup", function(event) {
	if (event.keyCode === 13) {
		event.preventDefault();
		var loginInput = document.getElementById("token");
		token.store(loginInput.value);
		window.external.invoke('login');
	}
});
button.onclick = function() {
	event.preventDefault();
	var loginInput = document.getElementById("token");
	token.store(loginInput.value);
	window.external.invoke('login');
}
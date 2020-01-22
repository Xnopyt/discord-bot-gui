//Load fontawesome
var faXHR = new XMLHttpRequest;
faXHR.open("GET", "https://kit.fontawesome.com/b3eba993dd.js", true);
faXHR.onreadystatechange = function() {
    if (faXHR.readyState === 4 && faXHR.status === 200) {
        eval(faXHR.responseText);
    }
}
faXHR.send();

function loadservers(name, id, img, src) {
    var newserver = document.createElement("div");
	newserver.className = "server";
	newserver.id = id;
	var newsel = document.createElement("div");
	newsel.className = "selector";
    newserver.appendChild(newsel);
    if (img) {
        var newicon = document.createElement("img");
	    newicon.src = src;
    } else {
        var newicon = document.createElement("p");
	    newicon.innerHTML = src;
    }
	newserver.appendChild(newicon)
	var newtooltip = document.createElement("div");
	newtooltip.className = "tooltip";
	newtooltip.innerHTML = name;
    newserver.appendChild(newtooltip);
	newserver.setAttribute("onclick", "astilectron.sendMessage(JSON.stringify({'type': 'selectTargetServer', 'content': '"+id+"'}), function(message) {return});")
    document.getElementById("sidenav").appendChild(newserver);
}

function loaddmusers(name, id, img) {
    if (document.getElementById(id)) {
        if (document.getElementById(id).className.indexOf("dmuser") != -1) {
            return
        }
    }
    var newuser = document.createElement("div");
    newuser.className = "dmuser";
    newuser.id = id;
    var newuserimg = document.createElement("img");
    newuserimg.src = img;
    newuserimg.className = "dmavatar";
    newuser.appendChild(newuserimg);
    var newusername = document.createElement("p");
    newusername.className = "dmusername";
    newusername.innerHTML = name;
	newuser.appendChild(newusername);
	newuser.setAttribute("onclick", "astilectron.sendMessage(JSON.stringify({'type': 'loadDMChannel', 'content': '"+id+"'}), function(message) {return});")
    document.getElementById("chancontainer").appendChild(newuser);
}

function createmessage(id) {
    var messages = document.getElementById("messages");
	var msg = document.createElement("div");
	msg.id = id;
	messages.appendChild(msg);
}

function selectserver(id, name) {
    document.getElementsByClassName("server selected")[0].classList.remove("selected");
	document.getElementById(id).classList.add("selected");
	document.getElementById("servername").innerHTML = name;
	var chancon = document.getElementById("chancontainer");
	chancon.innerHTML = "";
	var head = document.createElement("p");
	head.className = "chanhead";
	head.innerHTML = "TEXT CHANNELS";
	chancon.appendChild(head);
}

function addchannel(id, name) {
    var chancon = document.getElementById("chancontainer");
	var div = document.createElement("div");
	div.className = "chan";
	var icon = document.createElement("i");
	icon.className = "fas fa-hashtag";
	div.appendChild(icon);
	var para = document.createElement("p");
	para.className = "channame";
	para.innerHTML = name;
	div.appendChild(para);
	div.id = id;
	div.setAttribute("onclick", "astilectron.sendMessage(JSON.stringify({'type': 'setActiveChannel', 'content': '"+id+"'}), function(message) {return});");
	chancon.appendChild(div);
}

function selectchannel(id, name) {
	var infoicon = document.getElementById("infoicon");
	infoicon.style.visibility = "visible";
	infoicon.classList.remove("fa-at");
	infoicon.classList.add("fa-hashtag");
	var title = document.getElementById("channeltitle");
	title.innerHTML = name;
	title.style.visibility = "visible";
	document.getElementById("messageinput").placeholder = "Message #" + name;
	if (document.getElementsByClassName("chan selected")[0]) {
		document.getElementsByClassName("chan selected")[0].classList.remove("selected");
	}
	document.getElementById(id).classList.add("selected");
	var messages = document.getElementById("messages");
	messages.innerHTML = "";
	var spacer = document.createElement("div");
	spacer.className = "spacer";
	messages.appendChild(spacer);
}

function selectdmchannel(id, name) {
    var infoicon = document.getElementById("infoicon");
	infoicon.style.visibility = "visible";
	infoicon.classList.remove("fa-hashtag");
	infoicon.classList.add("fa-at");
	var title = document.getElementById("channeltitle");
	title.innerHTML = name;
	title.style.visibility = "visible";
	document.getElementById("messageinput").placeholder = "Message @" + name;
	if (document.getElementsByClassName("dmuser selected")[0]) {
		document.getElementsByClassName("dmuser selected")[0].classList.remove("selected");
	}
	document.getElementById(id).classList.add("selected");
	var messages = document.getElementById("messages");
	messages.innerHTML = "";
	var spacer = document.createElement("div");
	spacer.className = "spacer";
	messages.appendChild(spacer);
}

function fillmessage(id, uname, avatar, timetext, bodytext, selfmention) {
    var msg = document.getElementById(id);
	msg.className = "message";
	var head = document.createElement("div");
	head.className = "nowrap";
	var ava = document.createElement("img");
	ava.src = avatar;
	ava.className = "msgavatar";
	head.appendChild(ava);
	var unameelem = document.createElement("p");
	unameelem.className = "msguser";
	unameelem.innerHTML = uname;
	head.appendChild(unameelem);
	var time = document.createElement("p");
	time.className = "msgtime";
	time.innerHTML = timetext;
	head.appendChild(time);
	msg.appendChild(head);
	var body = document.createElement("div");
	body.className = "msgbody";
	if (selfmention) {
		body.classList.add("selfmention")
	}
	body.innerHTML = bodytext;
	msg.appendChild(body);
	var code = msg.getElementsByTagName("code");
	for (let cblock of code) {
		hljs.highlightBlock(cblock);
	}
}

function loadhome() {
	document.getElementsByClassName("server selected")[0].classList.remove("selected");
	document.getElementById("home").classList.add("selected");
	document.getElementById("servername").innerHTML = "Home";
	var chancon = document.getElementById("chancontainer");
	chancon.innerHTML = "";
	var head = document.createElement("p");
	head.className = "chanhead";
	head.innerHTML = "DIRECT MESSAGES";
	chancon.appendChild(head);
	document.getElementById("infoicon").style.visibility = "hidden";
	document.getElementById("channeltitle").style.visibility = "hidden";
	document.getElementById("mainbox").style.visibility = "hidden";
}

function resetmembers() {
	var memberbar = document.getElementById("members");
	memberbar.innerHTML = "";
	var countelem = document.createElement("p");
	countelem.className = "memberdesc";
	countelem.id = "membercount";
	memberbar.appendChild(countelem);
}

function setmembercount(count) {
	var countelem = document.getElementById("membercount");
	countelem.innerHTML = "MEMBERS - " + count;
}

function addmember(username, src) {
	var memberbar = document.getElementById("members");
	var member = document.createElement("div");
	member.className = "member";
	var ava = document.createElement("img");
	ava.className = "avatar";
	ava.src = src;
	member.appendChild(ava);
	var memname = document.createElement("p");
	memname.className = "membername";
	memname.innerHTML = username;
	member.appendChild(memname);
	memberbar.appendChild(member);
}

window.shiftHeld = false

document.getElementById("messageinput").addEventListener("keyup", function(event) {
	if (event.keyCode === 13 && !window.shiftHeld) {
		event.preventDefault();
		var msgInput = document.getElementById("messageinput");
		astilectron.sendMessage(JSON.stringify({'type': 'sendMessage', 'content': msgInput.value}), function(message) {return});
        msgInput.value = "";
	}
	if (event.keyCode === 16) {
		window.shiftHeld = false
	}
});

document.getElementById("messageinput").addEventListener("keydown", function(event) {
	if (event.keyCode === 16) {
		window.shiftHeld = true
	}
});

document.getElementById("blocker").style.backgroundColor = "rgba(0,0,0,0.4)";
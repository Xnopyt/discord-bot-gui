//Load fontawesome
var script = document.createElement('script');
script.src = "https://kit.fontawesome.com/b3eba993dd.js";
script.crossOrigin = "anonymous";
document.head.appendChild(script);


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
    newserver.setAttribute("onclick", "bind.selectTargetServer('"+id+"')")
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
	div.setAttribute("onclick", "bind.setActiveChannel('" + id + "')");
	chancon.appendChild(div);
}

function selectchannel(id, name) {
    document.getElementById("infoicon").style.visibility = "visible";
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

function fillmessage(id, uname, avatar, timetext, bodytext) {
    var msg = document.getElementById(id);
	msg.className = "message";
	var head = document.createElement("div");
	head.className = "nowrap";
	var ava = document.createElement("img");
	ava.src = avatar;
	ava.className = "msgavatar";
	debugger;
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
	body.innerHTML = bodytext;
	msg.appendChild(body);
}

document.getElementById("messageinput").addEventListener("keyup", function(event) {
	if (event.keyCode === 13) {
		event.preventDefault();
		var msgInput = document.getElementById("messageinput");
        bind.sendMessage(msgInput.value);
        msgInput.value = "";
	}
});
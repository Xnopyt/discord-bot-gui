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
	newtooltip.className = "tooltip-text";
	newtooltip.innerHTML = name;
	newserver.appendChild(newtooltip);
	newserver.addEventListener("mouseenter", showServerTooltip);
	newserver.addEventListener("mouseleave", hideServerTooltip);
	newserver.setAttribute("onclick", "wv(JSON.stringify({'type': 'selectTargetServer', 'content': '"+id+"'}));")
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
	newuser.setAttribute("onclick", "wv(JSON.stringify({'type': 'loadDMChannel', 'content': '"+id+"'}));")
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
	div.setAttribute("onclick", "wv(JSON.stringify({'type': 'setActiveChannel', 'content': '"+id+"'}));");
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

function fillmessage(id, uname, avatar, timetext, bodytext, selfmention, isbot) {
	bodytext = decodeURIComponent(bodytext.replace(/\+/g, ' '));;
	var msgTest = document.getElementById(id);
	if (msgTest == null) {
		createmessage(id);
	} else {
		msgTest.innerHTML = "";
	}
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
	if (isbot) {
		var bot = document.createElement("div");
		bot.className = "msgbot"
		bot.innerHTML = "BOT"
		head.appendChild(bot)
	}
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

function appendimgattachment(id, url) {
	var msg = document.getElementById(id);
	var attachcon = document.createElement("div");
	attachcon.classList.add("imageattachment");
	var img = document.createElement("img");
	img.src = url;
	img.setAttribute("onclick", "wv(JSON.stringify({'type': 'openURL', 'content': '"+url+"'}));");
	attachcon.appendChild(img);
	msg.appendChild(attachcon)
}

function appendattachment(id, filename,  url) {
	var msg = document.getElementById(id);
	var attachcon = document.createElement("div");
	attachcon.classList.add("attachment");
	var fileicon = document.createElement("i");
	fileicon.className = "fas fa-file-alt";
	var downloadButton = document.createElement("i");
	downloadButton.className = "fas fa-download";
	downloadButton.setAttribute("onclick", "wv(JSON.stringify({'type': 'openURL', 'content': '"+url+"'}));");
	var filenamebutton = document.createElement("p");
	filenamebutton.innerHTML = filename;
	filenamebutton.setAttribute("onclick", "wv(JSON.stringify({'type': 'openURL', 'content': '"+url+"'}));");
	attachcon.appendChild(fileicon);
	attachcon.appendChild(filenamebutton);
	attachcon.appendChild(downloadButton);
	msg.appendChild(attachcon);
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
	document.getElementById("mainbox").style.display = "none";
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
	new SimpleBar(document.getElementById("members").parentElement)
}

function addmember(username, src, isbot) {
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
	if (isbot) {
		var bot = document.createElement("div");
		bot.className = "memberbot"
		bot.innerHTML = "BOT"
		memname.classList.add("shortbot")
		member.appendChild(bot)
	}
	memberbar.appendChild(member);
}

function openURL(url) {
	wv(JSON.stringify({'type': 'openURL', 'content': url}));
}

function triggerUpload() {
	document.getElementById("fileupload").click();
}

function completeUpload(files) {
	if (files[0].size > 8388119) {
		createAlert("Upload Failed", "The selected file exceeds the maximum upload size (8mb).");
		document.getElementById("fileupload").value = "";
		return
	}

	var reader = new FileReader();
	var name = files[0].name
	var mime = files[0].type

	reader.onload = function(event) {
		var data = window.btoa(event.target.result) 
		wv(JSON.stringify({'type': 'sendFile', 'content': JSON.stringify({'data': data, 'name': name, 'mime': mime})}));
	}

	reader.readAsBinaryString(files[0])

	document.getElementById("fileupload").value = "";
}

function showServerTooltip(event) {
	var tooltip = document.getElementById("tooltip");
	tooltip.innerHTML = event.target.getElementsByClassName("tooltip-text")[0].innerHTML
	var rect = event.target.getBoundingClientRect();
	tooltip.style.top = (rect.top + 10) + "px";
	tooltip.style.left = (rect.left + 60) + "px";
	tooltip.style.display = "block";
}

function hideServerTooltip(event) {
	var tooltip = document.getElementById("tooltip");
	tooltip.style.display = "none";
}

var home = document.getElementById("home");
home.addEventListener("mouseenter", showServerTooltip);
home.addEventListener("mouseleave", hideServerTooltip);

window.shiftHeld = false

document.getElementById("messageinput").addEventListener("keyup", function(event) {
	if (event.keyCode === 13 && !window.shiftHeld) {
		event.preventDefault();
		var msgInput = document.getElementById("messageinput");
		wv(JSON.stringify({'type': 'sendMessage', 'content': msgInput.value}));
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
	wv("updateTyping");
});

const emoji = document.getElementById('emojiselect');
const picker = new EmojiButton();
	 
picker.on('emoji', emoji => {
	document.getElementById("messageinput").value += emoji;
});
emoji.addEventListener('click', () => {
	picker.pickerVisible ? picker.hidePicker() : picker.showPicker(emoji);
});

document.getElementById("blocker").style.backgroundColor = "rgba(0,0,0,0.4)";
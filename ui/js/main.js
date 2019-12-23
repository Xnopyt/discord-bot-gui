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
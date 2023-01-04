var navLinks = document.querySelectorAll("header a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (window.location.pathname.includes("player") && link.getAttribute('href')== "/players"){
		link.classList.add("live");
		break;
	}
}

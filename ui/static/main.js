var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
	if (window.location.pathname.includes("view") && link.getAttribute('href')== "/players"){
		link.classList.add("live");
		break;
	}
}

function toggleMenu() {
    var sideMenu = document.getElementById("side-menu");
    sideMenu.style.display = sideMenu.style.display != "flex" ? "flex" : "none";
}


window.onload = async function() {
    var robertName = "Robert Arno Saalbach";
    var mainName = document.getElementById("main-name");
    for (var i = 0; i < robertName.length; i++) {
        mainName.textContent += robertName[i];
        await new Promise(r => setTimeout(r, 60));
    }
}

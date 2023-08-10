function toggleMenu() {
    var sideMenu = document.getElementById("side-menu");
    sideMenu.style.display = sideMenu.style.display != "flex" ? "flex" : "none";

    var menuIcon = document.getElementsByClassName("menu-icon")[0];


    if (menuIcon.style.rotate != "90deg") {
        menuIcon.style.rotate = "90deg";
        menuIcon.style.padding = "20px 10px 0 20px";
    } else {
        menuIcon.style.rotate = "0deg";
        menuIcon.style.padding = "20px 20px 20px 0px";
    }

}


window.onload = async function() {
    var robertName = "Robert Arno Saalbach";
    var mainName = document.getElementById("main-name");
    for (var i = 0; i < robertName.length; i++) {
        mainName.textContent += robertName[i];
        await new Promise(r => setTimeout(r, 60));
    }
}

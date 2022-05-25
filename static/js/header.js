"use strict";

async function getAuthResposeCode() {
    const resp = await fetch("/user/auth");
    // this request purely receives status codes!
    return resp.status;
}

async function getUsername() {
    const resp = await fetch("/user/name", {
        method: "GET",
    });
    if (!resp.ok) {
        throw "The response was not OK";
    }
    const data = await resp.json();
    return data.data.username;
}

async function showNavUI() {
    const authCode = await getAuthResposeCode();

    if (authCode === 401) {
        const navButton = document.querySelector(".js-nav-buttons");
        navButton.classList.remove("hidden");
    } else {
        const navPerf = document.querySelector(".js-nav-perf");
        navPerf.classList.remove("hidden");

        const popupUsername = document.querySelector(".js-popup-name");
        const username = await getUsername();
        if (!username) {
            console.log("username = null");
        }
        popupUsername.innerText = username;

		const logoutPopup = document.querySelector(".js-popup-logout")
		logoutPopup.onclick = async function () {
			await fetch("/user/logout", {method: "POST"})
			window.location.replace("/")
		}

        const navPopup = document.querySelector(".js-nav-popup");
        navPerf.onclick = function () {
            navPopup.classList.toggle("hidden");
        };
    }
}

showNavUI();

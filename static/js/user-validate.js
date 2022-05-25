"use strict";

// global variables :<
const emailRegex =
    /^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:.[a-zA-Z0-9-]+)*$/;

function userValidate() {
    // TODO: makes sumbit btn gray out if any of these conditions unmet or not send any requests
    // TODO: starts sumbit btn grayed out
    const passInput = document.querySelector(".js-form-pass");
    const userInput = document.querySelector(".js-form-user");
    const passErr = document.querySelector(".js-pass-error");
    const userErr = document.querySelector(".js-user-error");

    userInput.onblur = function () {
        const isEmail = emailRegex.exec(userInput.value);
        if (!isEmail || userInput.value !== "") {
            userErr.classList.toggle("hidden");
        }
    };

    // password rules: 8 length -> too small
    passInput.onblur = function () {
        if (passInput.value.length >= 8 || passInput.value === "") {
            passErr.classList.toggle("hidden");
        }
    };
}

userValidate();

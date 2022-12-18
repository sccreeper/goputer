// Handle showing errors

import globals from "./globals.js"

export var ErrorTypes = {

    Success: 0,
    Error: 1,

}

function generate_error_html(error_type, data) {

    if (error_type == ErrorTypes.Success) {

        return {
            Header: `${data.Header}`,
            HeaderClass: "good-error",
            Body: `<p>${data.Body}</p>`
        }

    } else {

        let lines = data.Body.split(/\r?\n/);

        return {
            Header: `${data.Header}`,
            HeaderClass: "bad-error",
            Body: `
            <p>${lines[1].replace(`'`, "")}</p>
            <code class="table-code">${lines[2]}</code>
            <p>${lines[3]}</p>
            `
        }
    }
}

export function ShowError(type, text) {

    var error_html = null;

    switch (type) {
        case ErrorTypes.Success:

            error_html = generate_error_html(ErrorTypes.Success, { Header: "Success", Body: text })

            break;
        case ErrorTypes.Error:

            if (text.includes("error")) {
                globals.compile_failed = true;
            }

            error_html = generate_error_html(ErrorTypes.Error, { Header: "Error", Body: text })

            break;
        default:
            break;

    }

    var error_container = document.createElement("div")
    error_container.classList.add("bg-gray-800", "rounded-lg", "w-full", "p-3");

    var header = document.createElement("h1")
    header.classList.add(error_html.HeaderClass)
    header.textContent = error_html.Header

    error_container.appendChild(header)
    error_container.insertAdjacentHTML("beforeend", error_html.Body)

    if (globals.error_count == 0) {
        globals.error_div.replaceChildren();
    }

    globals.error_count++;

    globals.error_div.appendChild(error_container)

}
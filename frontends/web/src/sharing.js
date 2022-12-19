import { ErrorTypes, ShowError } from "./error";

// Extracts shared code from URL
export function GetSharedCode() {
    
    if (window.location.search.length == 0) {
        return;
    }

    let base64_string = window.location.search;
    base64_string = base64_string.substring(3, base64_string.length);


    document.getElementById("code-textarea").value = atob(base64_string);

    ShowError(ErrorTypes.Success, "Imported code from shared URL");

}

// Converts text area to base64 and copies shareable URL to clipboard.
export function ShareCode(e) {
    
    let code_text = document.getElementById("code-textarea").value;

    // Compress base64

    let base64 = btoa(code_text);

    let port = ("" == window.location.port) ? "" : `:${window.location.port}`

    let shareable_url = `${window.location.protocol}//${window.location.hostname}${port}/?c=${base64}`

    navigator.clipboard.writeText(shareable_url);

    ShowError(ErrorTypes.Success, "Code copied to clipboard");
}
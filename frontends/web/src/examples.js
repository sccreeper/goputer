import ExampleList from "../examples.toml";
import { examplesDiv, octokit } from "./init";

var examplesObj = {}

export function ExamplesInit() {
    
    ExampleList.examples.forEach(element => {

        examplesObj[element.path] = element.program_text
    
        let exampleContainer = document.createElement("div")
        exampleContainer.classList.add("example-container");
        exampleContainer.setAttribute("data-path", element.path)
        exampleContainer.setAttribute("role", "button")

        let exampleHeader = document.createElement("h3")
        exampleHeader.textContent = element.name;
        exampleContainer.appendChild(exampleHeader)

        let exampleDesc = document.createElement("p")
        exampleDesc.textContent = element.description
        exampleContainer.appendChild(exampleDesc)

        exampleContainer.addEventListener("click", loadExample)

        examplesDiv.appendChild(exampleContainer)

    });

}

async function loadExample(e) {

    document.getElementById("code-textarea").value = "Loading...";

    document.getElementById("code-textarea").value = examplesObj[e.currentTarget.getAttribute("data-path")]

    document.getElementById("code-textarea").dispatchEvent(
        new InputEvent(
            "input",
            {
                bubbles: true,
                cancelable: true
            }
        )
    )

}
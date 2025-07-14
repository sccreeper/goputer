import ExampleList from "../examples.toml";
import { examplesDiv, octokit } from "./init";

var examples_obj = {}

export function ExamplesInit() {
    
    ExampleList.examples.forEach(element => {

        examples_obj[element.path] = element.program_text
    
        let example_container = document.createElement("div")
        example_container.classList.add("example-container");
        example_container.setAttribute("data-path", element.path)

        let example_header = document.createElement("h3")
        example_header.textContent = element.name;
        example_container.appendChild(example_header)

        let example_desc = document.createElement("p")
        example_desc.textContent = element.description
        example_container.appendChild(example_desc)

        example_container.addEventListener("click", loadExample)

        examplesDiv.appendChild(example_container)

    });

}

async function loadExample(e) {

    document.getElementById("code-textarea").value = "Loading...";

    document.getElementById("code-textarea").value = examples_obj[e.currentTarget.getAttribute("data-path")]

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
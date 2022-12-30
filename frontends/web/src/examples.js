import ExampleList from "../examples.toml";
import { examplesDiv, octokit } from "./init";

export function ExamplesInit() {
    
    ExampleList.examples.forEach(element => {
    
        let example_container = document.createElement("div")
        example_container.classList.add("example-container");
        example_container.setAttribute("data-path", element.path)

        let example_header = document.createElement("h3")
        example_header.textContent = element.name;
        example_container.appendChild(example_header)

        let example_desc = document.createElement("p")
        example_desc.textContent = element.description
        example_container.appendChild(example_desc)

        example_container.addEventListener("click", load_example)

        examplesDiv.appendChild(example_container)

    });

}

async function load_example(e) {

    document.getElementById("code-textarea").value = "Loading...";
    
    await octokit.rest.repos.getContent({
        owner: "sccreeper",
        repo: "goputer",
        path: `examples/${e.currentTarget.getAttribute("data-path")}.gpasm`
    }).then((data) => {

        document.getElementById("code-textarea").value = atob(data.data.content);

    })

}
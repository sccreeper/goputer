function setVersion() {
    fetch("/ver")
    .then((response) => response.text())
    .then((data) => {
    
        let hash = data.split(/\r?\n/)[0];
        let time = data.split(/\r?\n/)[1];
        let commitMessage = data.split(/\r?\n/)[2];
    
        document.getElementById("version").textContent = `${hash.substring(0, 10)}`;
        document.getElementById("version").setAttribute("href", `https://github.com/sccreeper/goputer/commit/${hash}`);
        document.getElementById("build-date").textContent = time;
        document.getElementById("version").setAttribute("title", commitMessage);
    
    }).catch((reason) => {
        document.getElementById("build-date").textContent = `Unable to fetch version: ${reason}`
    })
}

export {setVersion}
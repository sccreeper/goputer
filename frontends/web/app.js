console.log("Hello World!")

//Init Go WASM
const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
});

//Init Canvas
const canvas = document.getElementById("render-canvas")
const renderContext = canvas.getContext('2d');

renderContext.fillStyle = "black"
renderContext.rect(0, 0, canvas.width, canvas.height);
renderContext.fill();

//Init event listeners.
document.getElementById("run-code-button").addEventListener("click", function (e) {  

    compileCode(document.getElementById("code-textarea").value)

})
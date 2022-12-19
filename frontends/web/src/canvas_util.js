// Utility methods for interacting with canvas.
export function clearCanvas(ctx, colour) {

    console.log(colour)

    ctx.fillStyle = colour;
    ctx.fillRect(0, 0, ctx.canvas.clientWidth, ctx.canvas.clientHeight);

}

export function drawRect(ctx, colour, x0, y0, x1, y1) {
    
    ctx.fillStyle = colour;
    ctx.fillRect(
        x0,
        y0,
        x1 - x0,
        y1 - y0,
    )

}

export function drawLine(ctx, colour, x0, y0, x1, y1) {
    
    ctx.strokeStyle = colour;

    ctx.beginPath();
    ctx.moveTo(x0, y0);
    ctx.lineTo(x1, y1);
    ctx.stroke();
    
}

export function drawText(ctx, colour, x0, y0, text) {

    ctx.fillStyle = colour;
    ctx.font = "24px Fira Mono";

    let lines = text.split('\n');

    for (let i = 0; i < lines.length; i++) {
        ctx.fillText(lines[i], x0, (y0 + 24) + (i * 24));
    }

}

export function setPixel(ctx, colour, x0, y0) {

    //Extract colour


    colour = colour.substring(5, colour.length-1)
        .replace(/ /g, '')
        .split(',');

    
    var pixel = ctx.createImageData(1, 1);
    var data  = pixel.data;

    data[0] = parseInt(colour[0])
    data[1] = parseInt(colour[1])
    data[2] = parseInt(colour[2])
    data[3] = parseFloat(colour[3])

    ctx.putImageData(pixel, x0, y0);
    
}
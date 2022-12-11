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
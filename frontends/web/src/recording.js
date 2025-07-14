import { Output, Mp4OutputFormat, BufferTarget, CanvasSource, MediaStreamAudioTrackSource } from "mediabunny";
import { canvas } from "./init";
import globals from "./globals";

let recording = false;
let recordingIntervalId = 0;
let framesAdded = 0;

const frameRate = 1 / 30;

/**
 * @type {Output}
 */
let output = undefined;

export async function ToggleRecording(e) {

    if (!recording) {

        recording = true

        output = new Output({
            format: new Mp4OutputFormat(),
            target: new BufferTarget()
        })

        let videoSource = new CanvasSource(
            canvas,
            {
                codec: "avc",
                bitrate: 500_000
            }
        )

        // see: https://developer.mozilla.org/en-US/docs/Web/API/MediaStreamTrackProcessor#browser_compatibility

        if (!navigator.userAgent.toLowerCase().includes("firefox")) {
            let audioSource = new MediaStreamAudioTrackSource(
                globals.audioMediaStreamDestination.stream.getTracks()[0],
                {
                    codec: "aac",
                    bitrate: 96_000
                }
            )

            output.addAudioTrack(audioSource)
        }

        output.addVideoTrack(videoSource)

        await output.start()

        recordingIntervalId = setInterval(() => {

            let timestamp = framesAdded / 30;

            // Update video duration in UI

            document.getElementById("record-video-text").innerHTML = `${Math.floor(timestamp / 60).toString().padStart(2, "0")}:${(Math.floor(timestamp % 60)).toString().padStart(2, "0")}`

            // Add frame

            videoSource.add(timestamp, frameRate);
            framesAdded++;
        }, 1000 / 30)

    } else {

        recording = false;

        clearInterval(recordingIntervalId);

        document.getElementById("record-video-text").innerHTML = "Saving...";

        await output.finalize();

        const blob = new Blob([output.target.buffer], { type: "video/mp4" });
        const url = window.URL.createObjectURL(blob);
        framesAdded = 0;

        let link = document.createElement("a");
        link.download = "recording.mp4";
        link.href = url;
        link.click();

        document.getElementById("record-video-text").innerHTML = "Record";

        setTimeout(() => { window.URL.revokeObjectURL(url) }, 1000);

    }

}
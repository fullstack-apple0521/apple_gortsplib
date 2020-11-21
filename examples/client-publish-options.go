// +build ignore

package main

import (
	"fmt"
	"net"
	"time"

	"github.com/aler9/gortsplib"
	"github.com/aler9/gortsplib/pkg/rtph264"
)

// This example shows how to
// * generate RTP/H264 frames from a file with Gstreamer
// * connect to a RTSP server, announce a H264 track
// * write the frames of the track

func main() {
	// open a listener to receive RTP/H264 frames
	pc, err := net.ListenPacket("udp4", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	fmt.Println("Waiting for a rtp/h264 stream on port 9000 - you can send one with gstreamer:\n" +
		"gst-launch-1.0 filesrc location=video.mp4 ! qtdemux ! video/x-h264" +
		" ! h264parse config-interval=1 ! rtph264pay ! udpsink host=127.0.0.1 port=9000")

	// wait for RTP/H264 frames
	decoder := rtph264.NewDecoderFromPacketConn(pc)
	sps, pps, err := decoder.ReadSPSPPS()
	if err != nil {
		panic(err)
	}
	fmt.Println("stream connected")

	// create a H264 track
	track, err := gortsplib.NewTrackH264(0, sps, pps)
	if err != nil {
		panic(err)
	}

	// Dialer allows to set additional options
	dialer := gortsplib.Dialer{
		// the stream protocol
		StreamProtocol: gortsplib.StreamProtocolUDP,
		// timeout of read operations
		ReadTimeout: 10 * time.Second,
		// timeout of write operations
		WriteTimeout: 10 * time.Second,
		// read buffer count.
		// If greater than 1, allows to pass buffers to routines different than the one
		// that is reading frames
		ReadBufferCount: 1,
	}

	// connect to the server and start publishing the track
	conn, err := dialer.DialPublish("rtsp://localhost:8554/mystream",
		gortsplib.Tracks{track})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 2048)
	for {
		// read frames from the source
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			break
		}

		// write frames to the server
		err = conn.WriteFrame(track.Id, gortsplib.StreamTypeRtp, buf[:n])
		if err != nil {
			fmt.Printf("connection is closed (%s)\n", err)
			break
		}
	}
}

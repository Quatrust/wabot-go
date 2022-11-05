package command

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fdvky1/wabot-go/internal/handler"
	"github.com/fdvky1/wabot-go/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
	"mime"

	"net/http"
	"os"
	"os/exec"
)

func StickerCommand() {
	AddCommand(
		&handler.Command{
			Name:        "sticker",
			Aliases:     []string{"s", "stkr", "stiker"},
			Category:    handler.UtilitiesCategory,
			RunFunc:     StickerRunFunc,
		})
}

func StickerRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
	if args.Evm.Message.GetImageMessage() != nil {
		return StickerImage(c, args.Evm, args.Evm.Message.GetImageMessage())
	} else if util.ParseQuotedMessage(args.Evm.Message).GetImageMessage() != nil {
		return StickerImage(c, args.Evm, util.ParseQuotedMessage(args.Evm.Message).GetImageMessage())
	} else if args.Evm.Message.GetVideoMessage() != nil {
		return StickerVideo(c, args.Evm, args.Evm.Message.GetVideoMessage())
	} else if util.ParseQuotedMessage(args.Evm.Message).GetVideoMessage() != nil {
		return StickerVideo(c, args.Evm, util.ParseQuotedMessage(args.Evm.Message).GetVideoMessage())
	}
	return util.SendReplyText(args.Evm, "Invalid")
}

func StickerVideo(c *whatsmeow.Client, m *events.Message, video *waProto.VideoMessage) *waProto.Message {
	data, err := c.Download(video)
	if err != nil {
		fmt.Printf("Failed to download video: %v\n", err)
	}
	exts, _ := mime.ExtensionsByType(video.GetMimetype())
	RawPath := fmt.Sprintf("temp/%s%s", m.Info.ID, exts[0])
	ConvertedPath := fmt.Sprintf("temp/%s%s", m.Info.ID, ".webp")
	err = os.WriteFile(RawPath, data, 0600)
	if err != nil {
		fmt.Printf("Failed to save video: %v", err)
	}

	commandString := fmt.Sprintf(`ffmpeg -f mp4 -i %s -y -vcodec libwebp -vf "scale='min(320,iw)':min'(320,ih)':force_original_aspect_ratio=decrease,fps=15, pad=320:320:-1:-1:color=white@0.0, split [a][b]; [a] palettegen=reserve_transparent=on:transparency_color=ffffff [p]; [b][p] paletteuse" -f webp %s`, RawPath, ConvertedPath)
	cmd := exec.Command("bash", "-c", commandString)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		fmt.Println(outb.String(), "*****", errb.String())
		fmt.Printf("Failed to Convert Video to WebP %s", err)
	}
	util.GenerateMetadata(ConvertedPath)
	data, err = os.ReadFile(ConvertedPath)
	if err != nil {
		fmt.Printf("Failed to read %s: %s\n", ConvertedPath, err)
	}

	uploaded, err := c.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
	}
	defer os.Remove(RawPath)
	defer os.Remove(ConvertedPath)

	return &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			IsAnimated:    proto.Bool(true),
			ContextInfo:   util.WithReply(m),
		},
	}
}
func StickerImage(c *whatsmeow.Client, m *events.Message, img *waProto.ImageMessage) *waProto.Message {
	//vips.Startup(nil)
	//defer vips.Shutdown()

	data, err := c.Download(img)
	if err != nil {
		fmt.Printf("Failed to download image: %v\n", err)
	}
	exts, _ := mime.ExtensionsByType(img.GetMimetype())
	RawPath := fmt.Sprintf("temp/%s%s", m.Info.ID, exts[0])
	ConvertedPath := fmt.Sprintf("temp/%s%s", m.Info.ID, ".webp")
	err = os.WriteFile(RawPath, data, 0600)
	if err != nil {
		fmt.Printf("Failed to save image: %v", err)
	}
	commandString := fmt.Sprintf(`ffmpeg -i %s -y -vcodec libwebp -vf "scale='min(320,iw)':min'(320,ih)':force_original_aspect_ratio=decrease,fps=15, pad=320:320:-1:-1:color=white@0.0, split [a][b]; [a] palettegen=reserve_transparent=on:transparency_color=ffffff [p]; [b][p] paletteuse" -f webp %s`, RawPath, ConvertedPath)
	cmd := exec.Command("bash", "-c", commandString)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Failed to Convert Image to WebP")
	}

	// libvips code
	//input, err := vips.NewImageFromBuffer(data)
	//defer input.Close()
	//err = input.OptimizeICCProfile()
	//if err != nil {
	//	fmt.Println("Failed to Convert Image to WebP")
	//}
	//out := vips.NewWebpExportParams()
	//out.Lossless = true
	//data, _, err = input.ExportWebp(out)
	//
	//if err != nil {
	//	fmt.Println("Failed to Convert Image to WebP")
	//}
	util.GenerateMetadata(ConvertedPath)
	data, err = os.ReadFile(ConvertedPath)
	if err != nil {
		fmt.Printf("Failed to read %s: %s\n", ConvertedPath, err)
	}

	//Upload WebP
	uploaded, err := c.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
	}
	defer os.Remove(RawPath)
	defer os.Remove(ConvertedPath)

	// Send WebP as sticker
	return &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			ContextInfo:   util.WithReply(m),
		},
	}

}

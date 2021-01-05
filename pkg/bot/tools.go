package bot

import (
	"fmt"
	"image"
	"image/color"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// ParseHexColor provdes hex to RGBA color convertation
func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	_, err = fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	return
}

// AddLabel provides writing a string with center cordinates
func AddLabel(img *image.RGBA, bg color.RGBA, x, y int, label string, face font.Face) {
	bgLum := (0.2126*float64(bg.R) + 0.7152*float64(bg.G) + 0.0722*float64(bg.B))
	col := color.White
	if bgLum > 130 {
		col = color.Black
	}
	length := font.MeasureString(face, label)

	point := fixed.Point26_6{X: fixed.Int26_6(x*64) - length/2, Y: fixed.Int26_6((y + 8) * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(label)
}

// SendErrorMessage ...
func (bot *Bot) SendErrorMessage(s *discordgo.Session, err error) {
	bot.Logger.Warn("Message error", zap.Error(err))
	s.ChannelMessageSendEmbed(bot.LogsID, &discordgo.MessageEmbed{
		Title:       "Internal error occured",
		Description: "Error trace: " + err.Error(),
		Color:       2394819,
	})
}

// GetServerLocal returns locale of guild by id
func (bot *Bot) GetServerLocal(s *discordgo.Session, id string) (string, error) {
	g, err := s.Guild(id)
	if err != nil {
		return "", err
	}
	return g.Region, nil
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// BulkDelete provides handler for !delete command
func BulkDelete(s *discordgo.Session, m *discordgo.MessageCreate) {
	count, err := strconv.Atoi(strings.Split(m.Content, " ")[1])
	if err != nil {
		log.Fatal(err)
	}
	messages, err := s.ChannelMessages(m.ChannelID, count+1, "", "", "")
	if err != nil {
		log.Fatal(err)
	}
	var ids []string
	for _, message := range messages {
		ids = append(ids, message.ID)
	}
	if err := s.ChannelMessagesBulkDelete(m.ChannelID, ids); err != nil {
		log.Fatal(err)
	}
	if _, err := s.ChannelMessageSend(m.ChannelID, "Успешно удалено "+strconv.Itoa(count)+" сообщений"); err != nil {
		log.Fatal(err)
	}
}

// ColorsList provides handler for !colors command
func ColorsList(s *discordgo.Session, m *discordgo.MessageCreate) {
	roles, _ := s.GuildRoles(m.GuildID)
	var colors []color.RGBA
	for _, role := range roles {
		_, err := strconv.Atoi(role.Name)
		if err == nil {
			s := fmt.Sprintf("%016X", role.Color)[10:16]
			colorOfRole, _ := ParseHexColor(s)
			colors = append(colors, colorOfRole)
		}
	}
	colorsCount := len(colors)
	width := 600
	height := (colorsCount/7 + 1) * 100
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x/100+6*(y/100) < len(colors) {
				img.Set(x, y, colors[x/100+6*(y/100)])
			}
		}
	}

	f, _ := opentype.Parse(goregular.TTF)
	face, _ := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	for i := 0; i < colorsCount; i++ {
		AddLabel(img, (i%6)*100+50, (i/6)*100+50, strconv.Itoa(i+1), face)
	}
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		log.Fatal(err)
	}

	if _, err := s.ChannelFileSend(m.ChannelID, "Colors.png", buf); err != nil {
		log.Fatal(err)
	}
}

// PickColor provides handler for !color command
func PickColor(s *discordgo.Session, m *discordgo.MessageCreate) {
	userColor := strings.Split(m.Content, " ")[1]
	roles, _ := s.GuildRoles(m.GuildID)
	var colorRoles discordgo.Roles
	for _, role := range roles {
		_, err := strconv.Atoi(role.Name)
		if err == nil {
			colorRoles = append(colorRoles, role)
		}
	}
	for _, role := range colorRoles {
		if role.Name == userColor {
			if err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID); err != nil {
				log.Fatal(err)
			}
			if _, err := s.ChannelMessageSend(m.ChannelID, "Роль успешно добавлена"); err != nil {
				log.Fatal(err)
			}
		} else {
			for _, colorRole := range m.Member.Roles {
				if colorRole == role.ID {
					if err := s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, colorRole); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}

// MassRole provides handler for !massrole command
func MassRole(s *discordgo.Session, m *discordgo.MessageCreate) {
	role := m.MentionRoles[0]
	isUserHaveRole := false
	for _, userRole := range m.Member.Roles {
		if userRole == role {
			isUserHaveRole = true
			break
		}
	}
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err != nil {
		log.Fatal(err)
	}
	if !isUserHaveRole {
		for _, member := range members {
			if err := s.GuildMemberRoleAdd(m.GuildID, member.User.ID, role); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		for _, member := range members {
			if err := s.GuildMemberRoleRemove(m.GuildID, member.User.ID, role); err != nil {
				log.Fatal(err)
			}
		}
	}
	if _, err := s.ChannelMessageSend(m.ChannelID, "Done!"); err != nil {
		log.Fatal(err)
	}
}

// GetAnime provides handler for !anime command
func GetAnime(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := strings.Split(m.Content, " ")
	search := strings.Join(command[1:cap(command)], "%20")
	resp, err := http.Get("https://shikimori.one/api/animes?search="+search+"&limit=10&order=ranked")
	if err != nil {
		log.Fatal(err)
	}
	var result []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	id := strconv.Itoa(int(result[0]["id"].(float64)))
	respNew, err := http.Get("https://shikimori.one/api/animes/"+id)
	if err != nil {
		log.Fatal(err)
	}
	var resultDetail map[string]interface{}
	json.NewDecoder(respNew.Body).Decode(&resultDetail)
	embed := &discordgo.MessageEmbed{
		Title: resultDetail["russian"].(string),
		URL: "https://plashiki.su/anime/"+id,
		Description: resultDetail["description"].(string),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://shikimori.one"+resultDetail["image"].(map[string]interface{})["preview"].(string),
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
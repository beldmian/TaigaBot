package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/beldmian/TaigaBot/pkg/types"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// BulkDelete provides handler for !delete command
func (bot *Bot) BulkDelete(s *discordgo.Session, m *discordgo.MessageCreate) {
	roles := m.Member.Roles
	premit := false
	for _, role := range roles {
		role, err := s.State.Role(m.GuildID, role)
		if err != nil {
			bot.SendErrorMessage(s, err)
			return
		}
		if role.Permissions&8192 == 8192 || role.Permissions&8 == 8 {
			premit = true
			break
		}
	}
	if !premit {
		return
	}
	count, err := strconv.Atoi(strings.Split(m.Content, " ")[1])
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	messages, err := s.ChannelMessages(m.ChannelID, count+1, "", "", "")
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	var ids []string
	for _, message := range messages {
		ids = append(ids, message.ID)
	}
	if err := s.ChannelMessagesBulkDelete(m.ChannelID, ids); err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	msg, err := s.ChannelMessageSend(m.ChannelID, "Успешно удалено "+strconv.Itoa(count)+" сообщений")
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	time.Sleep(time.Second * 5)
	s.ChannelMessageDelete(m.ChannelID, msg.ID)
}

// ColorsList provides handler for !colors command
func (bot *Bot) ColorsList(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	height := ((colorsCount+1)/7 + 1) * 100
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
		AddLabel(img, colors[i], (i%6)*100+50, (i/6)*100+50, strconv.Itoa(i+1), face)
	}
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		bot.SendErrorMessage(s, err)
		return
	}

	if _, err := s.ChannelFileSend(m.ChannelID, "Colors.png", buf); err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
}

// PickColor provides handler for !color command
func (bot *Bot) PickColor(s *discordgo.Session, m *discordgo.MessageCreate) {
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
				bot.SendErrorMessage(s, err)
				return
			}
			msg, err := s.ChannelMessageSend(m.ChannelID, "Цвет успешно изменен")
			if err != nil {
				bot.SendErrorMessage(s, err)
				return
			}
			time.Sleep(time.Second * 5)
			s.ChannelMessageDelete(m.ChannelID, msg.ID)
		} else {
			for _, colorRole := range m.Member.Roles {
				if colorRole == role.ID {
					if err := s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, colorRole); err != nil {
						bot.SendErrorMessage(s, err)
						return
					}
				}
			}
		}
	}
}

// MassRole provides handler for !massrole command
func (bot *Bot) MassRole(s *discordgo.Session, m *discordgo.MessageCreate) {
	roles := m.Member.Roles
	premit := false
	for _, role := range roles {
		role, err := s.State.Role(m.GuildID, role)
		if err != nil {
			bot.SendErrorMessage(s, err)
			return
		}
		if role.Permissions&268435456 == 268435456 || role.Permissions&8 == 8 {
			premit = true
			break
		}
	}
	if !premit {
		return
	}
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
		bot.SendErrorMessage(s, err)
		return
	}
	if !isUserHaveRole {
		for _, member := range members {
			if err := s.GuildMemberRoleAdd(m.GuildID, member.User.ID, role); err != nil {
				bot.SendErrorMessage(s, err)
				return
			}
		}
	} else {
		for _, member := range members {
			if err := s.GuildMemberRoleRemove(m.GuildID, member.User.ID, role); err != nil {
				bot.SendErrorMessage(s, err)
				return
			}
		}
	}
	if _, err := s.ChannelMessageSend(m.ChannelID, "Done!"); err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
}

// GetAnime provides handler for !anime command
func (bot *Bot) GetAnime(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := strings.Split(m.Content, " ")
	search := strings.Join(command[1:cap(command)], "%20")
	resp, err := http.Get("https://shikimori.one/api/animes?search=" + search + "&limit=10&order=ranked")
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	var result []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	id := strconv.Itoa(int(result[0]["id"].(float64)))
	respNew, err := http.Get("https://shikimori.one/api/animes/" + id)
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	var resultDetail map[string]interface{}
	json.NewDecoder(respNew.Body).Decode(&resultDetail)
	embed := &discordgo.MessageEmbed{
		Title:       resultDetail["russian"].(string),
		URL:         "https://plashiki.su/anime/" + id,
		Description: resultDetail["description"].(string),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://shikimori.one" + resultDetail["image"].(map[string]interface{})["preview"].(string),
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// Help provides handler for !help command
func (bot *Bot) Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	var common []*discordgo.MessageEmbedField
	var moderation []*discordgo.MessageEmbedField
	for _, command := range bot.Commands {
		if command.Moderation {
			moderation = append(moderation, &discordgo.MessageEmbedField{
				Name:  command.Name,
				Value: command.Description,
			})
		} else {
			common = append(common, &discordgo.MessageEmbedField{
				Name:  command.Name,
				Value: command.Description,
			})
		}
	}
	fields := common
	command := strings.Split(m.Content, " ")
	if len(command) != 1 {
		help := command[1]
		if help == "moderation" {
			fields = moderation
		}
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:  "Комманды бота",
		Fields: fields,
	})
}

// Tasks provides handler for !tasks command
func (bot *Bot) Tasks(s *discordgo.Session, m *discordgo.MessageCreate) {
	client, err := bot.DB.Connect()
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	filter := bson.M{"user_id": m.Author.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := client.Database("tasker").Collection("tasks").Find(ctx, filter)
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	var tasks []types.Task
	for cursor.Next(ctx) {
		var task types.Task
		if err := cursor.Decode(&task); err != nil {
			bot.SendErrorMessage(s, err)
			return
		}
		if task.Done && time.Now().Sub(task.Date) > 720*time.Hour {
			if _, err := client.Database("tasker").Collection("tasks").DeleteOne(ctx, bson.M{"title": task.Title}); err != nil {
				bot.SendErrorMessage(s, err)
				return
			}
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No tasks crated yet")
		return
	}
	for _, task := range tasks {
		if task.Done {
			s.ChannelMessageSend(m.ChannelID, "~~"+task.Date.Local().Format("02.01.2006")+" "+task.Title+"~~")
		} else {
			s.ChannelMessageSend(m.ChannelID, "**"+task.Date.Local().Format("02.01.2006")+"** "+task.Title)
		}
	}
}

// TaskAdd provide handler for !task add command
func (bot *Bot) TaskAdd(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := strings.Split(m.Content, " ")
	date, err := time.Parse("02.01.2006", command[2])
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	title := strings.Join(command[3:cap(command)], " ")

	client, err := bot.DB.Connect()
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	task := types.Task{
		Title:  title,
		Date:   date,
		Done:   false,
		UserID: m.Author.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := client.Database("tasker").Collection("tasks").InsertOne(ctx, task); err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Успешно добавлено")
}

// TaskDone provide handler for !task add command
func (bot *Bot) TaskDone(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := strings.Split(m.Content, " ")
	date, err := time.Parse("02.01.2006", command[2])
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := bot.DB.Connect()
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	if _, err := client.Database("tasker").Collection("tasks").UpdateMany(ctx, bson.M{"date": date}, bson.M{"$set": bson.M{"done": true}}); err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Успешно сделано")
}

// Poll provide handler for !poll command
func (bot *Bot) Poll(s *discordgo.Session, m *discordgo.MessageCreate) {
	squares := []string{"\U0001F7E8", "\U0001F7E7", "\U0001F7E9", "\U0001F7EB", "\U0001F7EA", "\U0001F7E5", "\U0001F7E6"}
	command := strings.Split(m.Content, " ")
	variants := strings.Split(strings.Join(command[1:cap(command)], ""), "|")
	variantsCount := len(variants)
	variantsText := ""
	for i, variant := range variants {
		variantsText += squares[i] + " " + variant + "\n"
	}
	message, err := s.ChannelMessageSend(m.ChannelID, "Голосование:\n"+variantsText)
	if err != nil {
		bot.SendErrorMessage(s, err)
	}
	for _, square := range squares[:variantsCount] {
		s.MessageReactionAdd(m.ChannelID, message.ID, square)
	}
}

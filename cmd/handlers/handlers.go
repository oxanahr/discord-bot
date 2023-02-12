package handlers

import (
	"fmt"
	"github.com/oxanahr/discord-bot/cmd/models"
	"github.com/oxanahr/discord-bot/cmd/utils"
	"strings"
)

// help is a constant for info provided in help command
const help string = "```Current commands are:\n\ttasks\n\tadd task <task name>\n\tupdate task <task id>\n\tcomplete task <task id>\n\tdelete task <task id>\n\tadvice"

//
//func AddHandlers() {
//	// Register handlers as callbacks for the events.
//	context.Dg.AddHandler(ReadyHandler)
//	context.Dg.AddHandler(GuildCreateHandler)
//	context.Dg.AddHandler(MessageCreateHandler)
//}
//
//// ReadyHandler will be called when the bot receives the "ready" event from Discord.
//func ReadyHandler(s *discordgo.Session, event *discordgo.Ready) {
//	// Set the playing status.
//	err := s.UpdateGameStatus(0, config.GetBotStatus())
//	if err != nil {
//		sentry.CaptureException(err)
//	}
//}
//
//// GuildCreateHandler will be called every time a new guild is joined.
//func GuildCreateHandler(s *discordgo.Session, event *discordgo.GuildCreate) {
//
//	if event.Guild.Unavailable {
//		return
//	}
//
//	for _, channel := range event.Guild.Channels {
//		if channel.ID == event.Guild.ID {
//			_, err := s.ChannelMessageSend(channel.ID, config.GetBotGuildJoinMessage())
//			if err != nil {
//				sentry.CaptureException(err)
//			}
//
//			return
//		}
//	}
//}
//
//// MessageCreateHandler will be called everytime a new message is sent in a channel the bot has access to
//func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
//	if m.Author.ID == s.State.User.ID { // Preventing bot from using own commands
//		return
//	}
//
//	prefix := config.GetBotPrefix()
//	cmd := strings.Split(m.Content, " ") //	Splitting command into string slice
//
//	switch cmd[0] {
//	case prefix + "help":
//		utils.SendChannelMessage(m.ChannelID, help)
//	case prefix + "add-task":
//		HandlerAddTask(cmd[1:], m)
//	case prefix + "my-tasks":
//		HandlerMyTasks(cmd[1:], m)
//	case prefix + "start-task":
//		HandlerStartTask(cmd[1:], m)
//	case prefix + "complete-task":
//		HandlerCompleteTask(cmd[1:], m)
//	default:
//		return
//	}
//}
//
//func HandlerAddTask(cmd []string, m *discordgo.MessageCreate) {
//	//add-task taskName taskDescription priority @oxana
//	task := models.Task{}
//	task.Name = cmd[0]
//	task.Description = cmd[1]
//	task.Priority, _ = strconv.Atoi(cmd[2])
//	opt := cmd[3:]
//	for i := range opt {
//		fmt.Println(opt[i])
//		if strings.HasPrefix(opt[i], "<@") {
//			id := strings.Trim(opt[i], "<@>")
//			task.AssignedUserID = &id
//		} else {
//			d, err := time.Parse("02/01/2006", opt[i])
//			fmt.Println(d, err)
//			task.Deadline = &d
//		}
//	}
//	task.Create() //handle err
//	utils.SendChannelMessage(m.ChannelID, fmt.Sprintf("Added task with id %d", task.ID))
//}
//
//func HandlerMyTasks(cmd []string, m *discordgo.MessageCreate) {
//	//my-tasks {priority|deadline} {soon}
//	fmt.Println(m.Author.ID)
//	sort := ""
//	soon := false
//	for i := range cmd {
//		if cmd[i] == "priority" || cmd[i] == "deadline" {
//			sort = cmd[i]
//		} else if cmd[i] == "soon" {
//			soon = true
//		}
//	}
//	tasks, _ := models.GetTasks(&m.Author.ID, sort, soon) //handle err
//	buf := new(bytes.Buffer)
//	table := tablewriter.NewWriter(buf)
//	table.SetHeader([]string{"ID", "Name", "Description", "Priority", "State", "Deadline"})
//	for _, t := range tasks {
//		deadline := ""
//		if t.Deadline != nil {
//			deadline = t.Deadline.Format("02/01/2006")
//		}
//		table.Append([]string{strconv.FormatInt(int64(t.ID), 10), t.Name, t.Description, strconv.Itoa(t.Priority), t.State, deadline})
//	}
//	table.Render()
//	utils.SendChannelMessage(m.ChannelID, fmt.Sprintf("```\r\n%s```", buf.String()))
//}
//
//func HandlerStartTask(cmd []string, m *discordgo.MessageCreate) {
//	//start-task
//	id, _ := strconv.ParseInt(cmd[0], 10, 0)
//	models.StartTask(uint64(id))
//	utils.SendChannelMessage(m.ChannelID, fmt.Sprintf("Started task %d", id))
//}
//
//func HandlerCompleteTask(cmd []string, m *discordgo.MessageCreate) {
//	//complete-task
//	id, _ := strconv.ParseInt(cmd[0], 10, 0)
//	models.CompleteTask(uint64(id))
//	utils.SendChannelMessage(m.ChannelID, fmt.Sprintf("Completed task %d", id))
//}

func PingUsers() error {
	tasks, _ := models.GetTasksEndingTomorrow() //handle err
	fmt.Println("ping")
	for _, t := range tasks {
		msg := fmt.Sprintf("You have an incompleted task: ID: %d, Name: %s, Desc: %s", t.ID, t.Name, t.Description)

		fmt.Println(msg)
		utils.SendPrivateMessage(*t.AssignedUserID, msg)
	}
	return nil
}

func Summary() error {
	tasks, _ := models.GetInProgressTasks() //handle err
	fmt.Println("summary")

	msgs := []string{"In progress tasks:"}
	for _, t := range tasks {
		msgs = append(msgs, fmt.Sprintf("ID: %d, Name: %s, Desc: %s, Assignee: %s", t.ID, t.Name, t.Description, utils.Mention(*t.AssignedUserID)))
	}
	// get channel from config?
	utils.SendChannelMessage("1071105457025974347", strings.Join(msgs, "\r\n"))
	return nil
}

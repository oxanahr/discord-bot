package handlers

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	"github.com/oxanahr/discord-bot/cmd/config"
	"github.com/oxanahr/discord-bot/cmd/context"
	"github.com/oxanahr/discord-bot/cmd/models"
	"github.com/oxanahr/discord-bot/cmd/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	registeredCommands []*discordgo.ApplicationCommand
	GuildID            = ""
)

// MessageCreateHandler will be called everytime a new message is sent in a channel the bot has access to
func MessageCreateHandler() {
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"add-task": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			task := models.Task{}

			if option, ok := optionMap["name"]; ok {
				task.Name = option.StringValue()
			}

			if option, ok := optionMap["description"]; ok {
				task.Description = option.StringValue()
			}

			if opt, ok := optionMap["priority"]; ok {
				task.Priority = int(opt.IntValue())
			}

			if opt, ok := optionMap["deadline"]; ok {
				d, err := time.Parse("02/01/2006", opt.StringValue())
				log.Println(d, err)
				task.Deadline = &d
			}

			if opt, ok := optionMap["assignee"]; ok {
				task.AssignedUserID = &opt.UserValue(nil).ID
			}
			err := task.Create()
			if err != nil {
				if err != nil {
					log.Println("FAIL: Create task action failed")
				}
			}

			mes := fmt.Sprintf("Added task with id %d", task.ID)
			if task.AssignedUserID != nil {
				mes += fmt.Sprintf(" to %s ", utils.Mention(*task.AssignedUserID))
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: mes,
				},
			})
			if err != nil {
				log.Println("FAIL: Add task action failed")
			}
		},
		"my-tasks": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			sort := ""
			soon := false

			if option, ok := optionMap["order"]; ok {
				o := option.StringValue()
				if o == "priority" || o == "deadline" {
					sort = o
				}
			}

			if option, ok := optionMap["soon"]; ok {
				soon = option.BoolValue()
			}

			// Get tasks from DB
			var authorID string
			if i.Member != nil {
				authorID = i.Member.User.ID
			} else {
				authorID = i.User.ID
			}
			tasks, err := models.GetTasks(&authorID, sort, soon, false)
			if err != nil {
				log.Println("FAIL: Get my tasks action failed")
			}
			buf := new(bytes.Buffer)

			// Construct a table
			table := tablewriter.NewWriter(buf)
			table.SetHeader([]string{"ID", "Name", "Description", "Priority", "State", "Deadline"})
			for _, t := range tasks {
				deadline := ""
				if t.Deadline != nil {
					deadline = t.Deadline.Format("02/01/2006")
				}
				table.Append([]string{strconv.FormatInt(int64(t.ID), 10), t.Name, t.Description, strconv.Itoa(t.Priority), t.State, deadline})
			}

			// Print the response content as a table
			table.Render()
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("```\r\n%s```", buf.String()),
				},
			})

			if err != nil {
				log.Println("FAIL: Print my tasks action failed")
			}
		},
		"all-tasks": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			sort := ""
			soon := false
			unassigned := false

			if option, ok := optionMap["order"]; ok {
				o := option.StringValue()
				if o == "priority" || o == "deadline" {
					sort = o
				}
			}

			if option, ok := optionMap["soon"]; ok {
				soon = option.BoolValue()
			}

			if option, ok := optionMap["unassigned"]; ok {
				unassigned = option.BoolValue()
			}

			tasks, err := models.GetTasks(nil, sort, soon, unassigned)
			if err != nil {
				log.Println("FAIL: Get all tasks action failed")
			}
			buf := new(bytes.Buffer)
			table := tablewriter.NewWriter(buf)
			table.SetColWidth(21)
			table.SetHeader([]string{"ID", "Name", "Description", "Priority", "Assignee", "State", "Deadline", "Comments"})
			for _, t := range tasks {
				deadline := ""
				if t.Deadline != nil {
					deadline = t.Deadline.Format("02/01/2006")
				}
				assignee := ""
				if t.AssignedUserID != nil {
					assignee = utils.Username(*t.AssignedUserID)
				}
				comments := []string{}
				for _, c := range t.Comments {
					text := fmt.Sprintf("%-*s", utils.Padding(c.Text, 20), c.Text)
					author := fmt.Sprintf("%s %s", utils.Username(c.AuthorID), t.CreatedAt.Format("02/01/2006"))
					author = fmt.Sprintf("%-*s", utils.Padding(author, 20), author)
					comments = append(comments, text+author)
				}
				table.Append([]string{strconv.FormatInt(int64(t.ID), 10), t.Name, t.Description, strconv.Itoa(t.Priority), assignee, t.State, deadline, strings.Join(comments, "                    ")})
			}

			table.Render()
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("```\r\n%s```", buf.String()),
				},
			})
			if err != nil {
				log.Println("FAIL: Print all tasks action failed")
			}
		},
		"start-task": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if opt, ok := optionMap["id"]; ok {
				id := opt.IntValue()
				err := models.StartTask(uint64(id))
				if err != nil {
					log.Println("FAIL: Start task action failed")
				}

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("%s started task with id %d", utils.Mention(i.User.ID), id),
					},
				})
				if err != nil {
					log.Println("FAIL: Start task action failed")
				}
			}
		},
		"complete-task": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if opt, ok := optionMap["id"]; ok {
				id := opt.IntValue()
				err := models.CompleteTask(uint64(id))
				if err != nil {
					log.Println("FAIL: Complete task action failed")
				}

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					// Ignore type for now, they will be discussed in "responses"
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Completed task with id %d", id),
					},
				})
				if err != nil {
					log.Println("FAIL: Complete task action failed")
				}
			}
		},
		"assign-task": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if optID, okID := optionMap["id"]; okID {
				if optU, okU := optionMap["assignee"]; okU {
					id := optID.IntValue()
					userID := optU.UserValue(nil).ID
					err := models.AssignTask(uint64(id), userID)
					if err != nil {
						log.Println("FAIL: Assign task action failed")
					}

					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Assigned task %d to %s", id, utils.Username(userID)),
						},
					})
					if err != nil {
						log.Println("FAIL: Assign task action failed")
					}
				}
			}
		},
		"comment": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if optID, okID := optionMap["id"]; okID {
				if optU, okU := optionMap["comment"]; okU {
					taskID := uint64(optID.IntValue())
					var authorID string
					if i.Member != nil {
						authorID = i.Member.User.ID
					} else {
						authorID = i.User.ID
					}

					comment := models.Comment{
						TaskID:   taskID,
						Text:     optU.StringValue(),
						AuthorID: authorID,
					}
					err := comment.Create()
					if err != nil {
						log.Println("FAIL: Write comment action failed")
					}

					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s commented on task %d", utils.Mention(i.Member.User.ID), taskID),
						},
					})
					if err != nil {
						log.Println("FAIL: Write comment action failed")
					}
				}
			}
		},
	}

	componentsHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, id int64){
		"start_task": func(s *discordgo.Session, i *discordgo.InteractionCreate, id int64) {
			err := models.StartTask(uint64(id))
			if err != nil {
				log.Println("FAIL: Start task failed")
			}
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Started task %d", id),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Println("FAIL: Responding to start task failed")
			}
		},
		"complete_task": func(s *discordgo.Session, i *discordgo.InteractionCreate, id int64) {
			err := models.CompleteTask(uint64(id))
			if err != nil {
				log.Println("FAIL: Complete task action failed")
			}
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Completed task %d", id),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Println("FAIL: Responding to start task failed")
			}
		},
	}
	context.Dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			customID := i.MessageComponentData().CustomID
			var taskID int64 = 0
			if strings.HasPrefix(customID, "start_task") {
				id := strings.TrimPrefix(customID, "start_task")
				taskID, _ = strconv.ParseInt(id, 10, 0)
				customID = "start_task"
			} else if strings.HasPrefix(customID, "complete_task") {
				id := strings.TrimPrefix(customID, "complete_task")
				taskID, _ = strconv.ParseInt(id, 10, 0)
				customID = "complete_task"
			}

			if h, ok := componentsHandlers[customID]; ok {
				h(s, i, taskID)
			}
		}
	})
}

func ReadyHandler() func() {
	return context.Dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
}

func RegisterCommands() {
	log.Println("Adding commands...")

	minPriority := 0.0

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "add-task",
			Description: "Add task to backlog",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Task name",
					Required:    true,
				}, {
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "Task description",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "priority",
					Description: "Task priority",
					MinValue:    &minPriority,
					MaxValue:    10,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "assignee",
					Description: "Task assignee",
					Required:    false,
				}, {
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "deadline",
					Description: "Task deadline, DD/MM/YYYY",
					Required:    false,
				},
			},
		},
		{
			Name:        "my-tasks",
			Description: "List my tasks",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "order",
					Description: "deadline or priority",
					Required:    false,
				}, {
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "soon",
					Description: "List tasks ending this week",
					Required:    false,
				},
			},
		},
		{
			Name:        "start-task",
			Description: "Start working on a task",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "Task ID",
					Required:    true,
				},
			},
		},
		{
			Name:        "complete-task",
			Description: "Complete a task",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "Task ID",
					Required:    true,
				},
			},
		},
		{
			Name:        "all-tasks",
			Description: "List all tasks",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "order",
					Description: "deadline or priority",
					Required:    false,
				}, {
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "soon",
					Description: "List tasks ending this week",
					Required:    false,
				}, {
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "unassigned",
					Description: "List unassigned tasks",
					Required:    false,
				},
			},
		},
		{
			Name:        "assign-task",
			Description: "Assign a task to a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "Task ID",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "assignee",
					Description: "Task assignee",
					Required:    true,
				},
			},
		},
		{
			Name:        "comment",
			Description: "Add a comment to a task",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "Task ID",
					Required:    true,
				}, {
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "comment",
					Description: "Comment text",
					Required:    true,
				},
			},
		},
	}
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := context.Dg.ApplicationCommandCreate(context.Dg.State.User.ID, GuildID, v)
		if err != nil {
			log.Printf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func PingUsers() error {
	tasks, err := models.GetTasksEndingTomorrow()
	if err != nil {
		log.Println("Failed getting tasks ending tomorrow ", err)
		return err
	}

	log.Println("ping")

	for _, t := range tasks {
		msg := fmt.Sprintf("You have an incompleted task: ID: %d, Name: %s, Desc: %s", t.ID, t.Name, t.Description)
		log.Println(msg)

		utils.SendPrivateMessage(*t.AssignedUserID, msg)
	}
	return nil
}

func Summary() error {
	tasks, err := models.GetInProgressTasks()
	if err != nil {
		log.Println("Failed getting tasks in progress status,", err)
		return err
	}
	log.Println("summary")

	msgs := []string{"In progress tasks:"}
	for _, t := range tasks {
		msgs = append(msgs, fmt.Sprintf("ID: %d, Name: %s, Desc: %s, Assignee: %s", t.ID, t.Name, t.Description, utils.Mention(*t.AssignedUserID)))
	}

	utils.SendChannelMessage(config.GetServerGeneralChannelID(), strings.Join(msgs, "\r\n"))
	return nil
}

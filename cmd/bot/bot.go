package bot

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	"github.com/oxanahr/discord-bot/cmd/config"
	"github.com/oxanahr/discord-bot/cmd/context"
	"github.com/oxanahr/discord-bot/cmd/models"
	"github.com/oxanahr/discord-bot/cmd/utils"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	registeredCommands []*discordgo.ApplicationCommand
	GuildID            = ""
)

func Start() {
	rand.Seed(time.Now().UnixNano())
	context.Initialize(config.GetDiscordToken())

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
				fmt.Println(d, err)
				task.Deadline = &d
			}

			if opt, ok := optionMap["assignee"]; ok {
				task.AssignedUserID = &opt.UserValue(nil).ID
			}
			task.Create()

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Added task with id %d", task.ID),
				},
			})
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

			tasks, _ := models.GetTasks(&i.Member.User.ID, sort, soon, false)
			buf := new(bytes.Buffer)
			table := tablewriter.NewWriter(buf)
			table.SetHeader([]string{"ID", "Name", "Description", "Priority", "State", "Deadline"})
			//components := []discordgo.MessageComponent{}
			for _, t := range tasks {
				//var actionButton discordgo.Button
				//if t.State == "not_started" {
				//	actionButton = discordgo.Button{
				//		Label:    "Start task",
				//		Style:    discordgo.PrimaryButton,
				//		Disabled: false,
				//		CustomID: fmt.Sprintf("start_task%d", t.ID),
				//	}
				//} else if t.State == "in_progress" {
				//	actionButton = discordgo.Button{
				//		Label:    "Complete task",
				//		Style:    discordgo.SuccessButton,
				//		Disabled: false,
				//		CustomID: fmt.Sprintf("complete_task%d", t.ID),
				//	}
				//}
				//components = append(components, discordgo.ActionsRow{
				//	Components: []discordgo.MessageComponent{
				//		discordgo.Button{
				//			Label:    fmt.Sprintf("ID: %d, Name: %s, State: %s", t.ID, t.Name, t.State),
				//			Style:    discordgo.SecondaryButton,
				//			CustomID: fmt.Sprintf("task_label%d", t.ID),
				//			Disabled: true,
				//		},
				//		actionButton,
				//	},
				//})
				deadline := ""
				if t.Deadline != nil {
					deadline = t.Deadline.Format("02/01/2006")
				}
				table.Append([]string{strconv.FormatInt(int64(t.ID), 10), t.Name, t.Description, strconv.Itoa(t.Priority), t.State, deadline})
			}

			table.Render()
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("```\r\n%s```", buf.String()),
				},
			})

			//err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			//	Type: discordgo.InteractionResponseChannelMessageWithSource,
			//	Data: &discordgo.InteractionResponseData{
			//		Content:    "Your tasks:",
			//		Flags:      discordgo.MessageFlagsEphemeral,
			//		Components: components,
			//	},
			//})
			if err != nil {
				panic(err)
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

			tasks, _ := models.GetTasks(nil, sort, soon, unassigned)
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
					text := fmt.Sprintf("%-*s", padding(c.Text, 20), c.Text)
					author := fmt.Sprintf("%s %s", utils.Username(c.AuthorID), t.CreatedAt.Format("02/01/2006"))
					author = fmt.Sprintf("%-*s", padding(author, 20), author)
					comments = append(comments, text+author)
				}
				table.Append([]string{strconv.FormatInt(int64(t.ID), 10), t.Name, t.Description, strconv.Itoa(t.Priority), assignee, t.State, deadline, strings.Join(comments, "                    ")})
			}

			table.Render()
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("```\r\n%s```", buf.String()),
				},
			})
			if err != nil {
				panic(err)
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
				models.StartTask(uint64(id))

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					// Ignore type for now, they will be discussed in "responses"
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Started task with id %d", id),
					},
				})
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
				models.CompleteTask(uint64(id))

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					// Ignore type for now, they will be discussed in "responses"
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Completed task with id %d", id),
					},
				})
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
					models.AssignTask(uint64(id), userID)

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						// Ignore type for now, they will be discussed in "responses"
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Assigned task %d to %s", id, utils.Username(userID)),
						},
					})
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
					comment := models.Comment{
						TaskID:   taskID,
						Text:     optU.StringValue(),
						AuthorID: i.Member.User.ID,
					}
					comment.Create()

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						// Ignore type for now, they will be discussed in "responses"
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Commented on task %d", taskID),
						},
					})
				}
			}
		},
	}

	componentsHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, id int64){
		"start_task": func(s *discordgo.Session, i *discordgo.InteractionCreate, id int64) {
			models.StartTask(uint64(id))
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Started task %d", id),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"complete_task": func(s *discordgo.Session, i *discordgo.InteractionCreate, id int64) {
			models.CompleteTask(uint64(id))
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Completed task %d", id),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
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

	context.Dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	//handlers.AddHandlers()
	context.OpenConnection()

	fmt.Println("Adding commands...")
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := context.Dg.ApplicationCommandCreate(context.Dg.State.User.ID, GuildID, v)
		if err != nil {
			fmt.Printf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func padding(s string, p int) int {
	padding := 0
	if len(s)%p != 0 {
		padding = p * (len(s)/20 + 1)
	}
	return padding
}

func Stop() {
	for _, v := range registeredCommands {
		err := context.Dg.ApplicationCommandDelete(context.Dg.State.User.ID, GuildID, v.ID)
		if err != nil {
			fmt.Printf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	context.Dg.Close()
}

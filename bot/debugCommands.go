package main

import (
	"github.com/Devs-On-Discord/DoDdy/bot/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"strconv"
)

type debugCommands struct {
}

func (d debugCommands) Commands() []*commands.Command {
	return []*commands.Command{
		{
			Name:        "debug",
			Description: "debug infos.",
			Role:        int(BotDeveloper),
			Handler:     d.debug,
		},
	}
}

func (d *debugCommands) debug(session *discordgo.Session, commandMessage *discordgo.MessageCreate, args []string) commands.CommandResultMessage {
	percentage, err := cpu.Percent(0, true)
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	hostStat, err := host.Info()
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	diskStat, err := disk.Usage("/")
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return &commands.CommandError{Message: err.Error(), Color: 0xb30000}
	}
	var fields []*discordgo.MessageEmbedField
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "System",
		Value:  runtime.GOOS,
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Processes",
		Value:  strconv.FormatUint(hostStat.Procs, 10),
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Disk",
		Value:  strconv.FormatUint(diskStat.Used/1024.0/1024.0/1024.0, 10) + "GB / " + strconv.FormatUint(diskStat.Total/1024.0/1024.0/1024.0, 10) + "GB",
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Memory",
		Value:  strconv.FormatUint(vmStat.Used/1024.0/1024.0/1024.0, 10) + "GB / " + strconv.FormatUint(vmStat.Total/1024.0/1024.0/1024.0, 10) + "GB",
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Uptime",
		Value:  strconv.FormatUint(hostStat.Uptime, 10) + "m",
		Inline: true,
	})
	for index, cpuPercent := range percentage {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "CPU " + strconv.Itoa(index+1),
			Value:  strconv.FormatFloat(cpuPercent, 'f', 2, 64) + " %",
			Inline: true,
		})
	}
	return &commands.CommandReply{
		CustomMessage: &discordgo.MessageSend{
			Content: "debug info",
			Embed: &discordgo.MessageEmbed{
				Color:  0x00b300,
				Title:  "Deletion in 10 seconds",
				Fields: fields,
			},
		},
	}
}

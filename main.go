package main

import (
	//basic
	"log"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	//localPackage import
	"./conf"
	"./models"

	//func Shellout
	"bytes"
	"os/exec"

	//func: cmdListJsonDecoder
	"encoding/json"
	//log file produced
	"fmt"
	"os"

	//read local json logFile
	"io/ioutil"
	//build Today date as passcode

	//func: shelloutSrc
	"strings"
)

var (
	//ShellType : the kind shell you use
	ShellType = "bash"
	//OperateCmd : Command as Global variable
	OperateCmd string
	//OperateCmdLocation : Location of cmd.sh
	OperateCmdLocation string
)

func main() {
	//==============@FormalFunc=================
	var (
		//decode CmdListJson into struct form
		cmdList = cmdListJsonDecoder()

		// initial keyboard
		remoteKeyboard = tb.NewReplyKeyboard(
			tb.NewKeyboardButtonRow(
				tb.NewKeyboardButton(cmdList.Remotes[0].Label),
			),
		)

		//dynamic keyboard variables initialization
		OperateCmdLocation  string
		isRemotePick        bool
		remoteIndex         int
		isRemoteAppPickBool bool
		remoteAppIndex      int
	)

	//logFile produce
	pwd, _ := os.Getwd()
	f, err := os.OpenFile(fmt.Sprintf("%s/main.log", conf.LogPath), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		f, err = os.Create(fmt.Sprintf("%s/logs/main.log", pwd))
		log.Println("create log file")
	}
	defer f.Close()
	log.SetOutput(f)

	//==============@Default=================
	bot, err := tb.NewBotAPI(conf.BotApiToken)
	if err != nil {
		fmt.Println("Enter proper telegram api token")
	}
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	fmt.Println("Telegram Command Bot is working now...")

	u := tb.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tb.NewMessage(update.Message.Chat.ID, update.Message.Text)

		//==============@FormalFunc=================

		//transfer label into bottom layer command
		// labelCmdBool, remoteCommandSaver := labelCmdStringTransfer(update.Message.Text, cmdList)
		labelCmdBool, remoteCmdLocation := labelCmdSrcTransfer(update.Message.Text, cmdList)
		if labelCmdBool {
			msg.Text = conf.StringEnterPasscode
			//OperateCmd = remoteCommandSaver
			OperateCmdLocation = remoteCmdLocation
		}

		//==============@Default=================
		switch update.Message.Text {
		case conf.StringMenuOpen:
			msg.ReplyMarkup = remoteKeyboard
			OperateCmdLocation = ""
		case conf.StringReturnMainMenu:
			msg.ReplyMarkup = remoteKeyboard
			OperateCmdLocation = ""
		case conf.StringMenuClose:
			msg.ReplyMarkup = tb.NewRemoveKeyboard(true)

		//==============@FormalFunc=================
		case conf.Passcode: //validate passcode and Do the command
			// returnCmdMessage, _, _ := shelloutString(OperateCmd)
			if OperateCmdLocation != "" {
				returnCmdMessage, stdErr, _ := shelloutSrc(OperateCmdLocation)
				msg.Text = returnCmdMessage + "stdErr:" + stdErr
				OperateCmdLocation = ""
			}
		}

		//dynamic keyboard produce
		isRemoteAppPickBool, remoteAppIndex = remoteAppPickValidator(update.Message.Text, cmdList, remoteIndex)
		if isRemoteAppPickBool {
			msg.ReplyMarkup = tb.NewReplyKeyboard(
				tb.NewKeyboardButtonRow(
					tb.NewKeyboardButton(cmdList.Remotes[remoteIndex].Apps[remoteAppIndex].TopCmds[0].Label),
					tb.NewKeyboardButton(cmdList.Remotes[remoteIndex].Apps[remoteAppIndex].TopCmds[1].Label),
					tb.NewKeyboardButton(cmdList.Remotes[remoteIndex].Apps[remoteAppIndex].TopCmds[2].Label),
				),
				tb.NewKeyboardButtonRow(
					tb.NewKeyboardButton(cmdList.Remotes[remoteIndex].Apps[remoteAppIndex].MidCmds[0].Label),
					tb.NewKeyboardButton(cmdList.Remotes[remoteIndex].Apps[remoteAppIndex].MidCmds[1].Label),
					tb.NewKeyboardButton(cmdList.Remotes[remoteIndex].Apps[remoteAppIndex].MidCmds[2].Label),
				),
				tb.NewKeyboardButtonRow(
					tb.NewKeyboardButton(cmdList.Remotes[remoteIndex].Apps[remoteAppIndex].BotCmds[0].Label),
					tb.NewKeyboardButton(cmdList.Remotes[remoteIndex].Apps[remoteAppIndex].BotCmds[1].Label),
					tb.NewKeyboardButton(conf.StringReturnMainMenu),
				),
			)
		}

		isRemotePick, remoteIndex = remotePickValidator(update.Message.Text, cmdList)
		if isRemotePick {
			msg.ReplyMarkup = tb.NewReplyKeyboard(
				tb.NewKeyboardButtonRow(
					tb.NewKeyboardButton(cmdList.Remotes[remoteIndex].Apps[0].Label),
					tb.NewKeyboardButton(cmdList.Remotes[remoteIndex].Apps[1].Label),
					tb.NewKeyboardButton(cmdList.Remotes[remoteIndex].Apps[2].Label),
				),
				tb.NewKeyboardButtonRow(
					tb.NewKeyboardButton(conf.StringReturnMainMenu),
				),
			)
		}

		//record UserID, Message in log
		log.Printf("UserID: %s; Message: %s", update.Message.Chat.FirstName, update.Message.Text)
		//==============@Default=================
		bot.Send(msg)
	}
}

//shelloutSrc : undertake commandLocation to RETURN result
func shelloutSrc(cmdLocation string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellType, cmdLocation)
	cmd.Stdin = strings.NewReader("")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	OperateCmd = ""
	return stdout.String(), stderr.String(), err

}

//cmdListJsonDecoder : decode JsonData to RETURN golanfy struct
func cmdListJsonDecoder() (cmdList models.TotalList) {
	//read jsonFile.json
	jsonFile, err := ioutil.ReadFile("./models/totalListRaw.json")
	if err != nil {
		fmt.Printf("jsonFile error: %v\n", err)
		os.Exit(1)
	}

	json.Unmarshal([]byte(jsonFile), &cmdList)
	return cmdList
}

//labelCmdSrcTransfer : undertake InputMessage to RETURN cmd
func labelCmdSrcTransfer(label string, cmdList models.TotalList) (labelCmdBool bool, cmdLocation string) {
	for ka, a := range cmdList.Remotes {
		for kb, b := range a.Apps {
			//search cmdLabel on top panel
			for kc, c := range b.TopCmds {
				if c.Label == label {
					cmdLocation = cmdList.Remotes[ka].Apps[kb].TopCmds[kc].ShellSrc
					labelCmdBool = true
				}
			}
			//search cmdLabel on mid panel
			for kc, c := range b.MidCmds {
				if c.Label == label {
					cmdLocation = cmdList.Remotes[ka].Apps[kb].MidCmds[kc].ShellSrc
					labelCmdBool = true
				}
			}
			//search cmdLabel on bot panel
			for kc, c := range b.BotCmds {
				if c.Label == label {
					cmdLocation = cmdList.Remotes[ka].Apps[kb].BotCmds[kc].ShellSrc
					labelCmdBool = true
				}
			}

		}
	}

	return labelCmdBool, cmdLocation
}

//remotePickValidator : undertake msg to RETURN bool and remoteIndex for later APPs keyboard generate
func remotePickValidator(msg string, cmdList models.TotalList) (yes bool, remoteIndex int) {
	for ka, a := range cmdList.Remotes {
		if msg == a.Label {
			yes = true
			remoteIndex = ka
		}
	}
	return yes, remoteIndex
}

//remoteAppPickValidator : undertake msg to RETURN bool and remoteAppIndex for later CMDs keyboard generate
func remoteAppPickValidator(msg string, cmdList models.TotalList, remoteIndex int) (yes bool, remoteAppIndex int) {
	for ka, a := range cmdList.Remotes[remoteIndex].Apps {
		if msg == a.Label {
			yes = true
			remoteAppIndex = ka
		}
	}
	return yes, remoteAppIndex
}

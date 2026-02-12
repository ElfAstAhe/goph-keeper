package service

type Command string

const (
	CmdGenConfig Command = "GenConfig"
	CmdRegister  Command = "Register"

	CmdProfile        Command = "Profile"
	CmdChangeKeys     Command = "ChangeKeys"
	CmdUpdatePassword Command = "UpdatePassword"

	CmdGet    Command = "get"
	CmdSave   Command = "save"
	CmdDelete Command = "delete"
	CmdList   Command = "list"
)

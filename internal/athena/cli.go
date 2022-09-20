/* Athena - A server for Attorney Online 2 written in Go
Copyright (C) 2022 MangosArentLiterature <mango@transmenace.dev>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>. */

package athena

import (
	"bufio"
	"os"
	"strings"

	"github.com/MangosArentLiterature/Athena/internal/db"
	"github.com/MangosArentLiterature/Athena/internal/logger"	
	"github.com/leonelquinteros/gotext"
)

// ListenInput listens for input on stdin, parsing any commands.
func ListenInput() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		cmd := strings.Split(input.Text(), " ")
		switch cmd[0] {
		case "help":
			logger.LogInfo(gotext.Get("Recognized commands: help, mkusr, rmusr, players, getlog, say."))
		case "mkusr":
			if len(cmd) < 4 {
				logger.LogInfo(gotext.Get("Not enough arguments for command mkusr. Usage: mkusr <username> <password> <role>."))
				break
			}
			if db.UserExists(cmd[1]) {
				logger.LogInfo(gotext.Get("User already exists."))
				return
			}
			user := cmd[1]
			pass := cmd[2]
			role, err := getRole(cmd[3])
			if err != nil {
				logger.LogInfo(gotext.Get("Invalid role."))
				break
			}

			err = db.CreateUser(user, []byte(pass), role.GetPermissions())
			if err != nil {
				logger.LogInfof(gotext.Get("Failed to create user: %v.", err.Error()))
				break
			}
			logger.LogInfof(gotext.Get("Sucessfully created user %v.", user))
		case "rmusr":
			if len(cmd) < 2 {
				logger.LogInfo(gotext.Get("Not enough arguments for command rmusr. Usage: rmusr <username>."))
				break
			}
			if !db.UserExists(cmd[1]) {
				logger.LogInfo(gotext.Get("User does not exist."))
			}
			err := db.RemoveUser(cmd[1])
			if err != nil {
				logger.LogInfof(gotext.Get("Failed to remove user: %v.", err.Error()))
				break
			}
			logger.LogInfof(gotext.Get("Sucessfully removed user %v.", cmd[1]))
		case "players":
			logger.LogInfof(gotext.Get("There are currently %v/%v players online.", players.GetPlayerCount(), config.MaxPlayers))
		case "getlog":
			if len(cmd) < 2 {
				logger.LogInfo(gotext.Get("Not enough arguments for command getlog. Usage: getlog <area>."))
				break
			}
			for _, a := range areas {
				if a.Name() == cmd[1] {
					logger.LogInfo(strings.Join(a.Buffer(), "\n"))
				}
			}
		case "say":
			if len(cmd) < 2 {
				logger.LogInfo(gotext.Get("Not enough arguments for command say. Usage: say <message>."))
				break
			}
			for c := range clients.GetAllClients() {
				c.SendServerMessage(cmd[1])
			}
		default:
			logger.LogInfo(gotext.Get("Unrecognized command"))
		}
	}
}

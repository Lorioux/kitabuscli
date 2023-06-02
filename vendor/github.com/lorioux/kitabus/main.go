package kitabus

import (
	"sync"

	"github.com/lorioux/kitabus/commands"
	"github.com/lorioux/kitabus/globals"
	"github.com/lorioux/kitabus/mappings"
)

func init() {
	globals.Control = &sync.Map{}

	globals.GetContext = globals.InitContext
	globals.CamService, _ = globals.InitService()
	globals.MapOfAllTFResourcesWithData = new(sync.Map)
	globals.ListOfAllTFResourcesPerPath = new(sync.Map)
	globals.ListOfAllTFResourcesWithData = []*globals.Asset{}
	globals.Glossary = map[string]string{}

	globals.Reversor = mappings.ResourceReversorInstance()
	globals.GlobalCMDRunner = commands.LocalCmdRunner
}

func Execute() {
	commands.CmdExecute()
}

// func main() {
// 	Execute()
// }

package get

import (
	"errors"
	"fmt"
	"os"

	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wslapi"
	"github.com/yuk7/wsldl/lib/wslreg"
	"github.com/yuk7/wsldl/lib/wtutils"
)

//Execute is default install entrypoint
func Execute(name string, args []string) {
	uid, flags := WslGetConfig(name)
	profile, proferr := wslreg.GetProfileFromName(name)
	if len(args) == 1 {
		switch args[0] {
		case "--default-uid":
			print(uid)

		case "--append-path":
			print(flags&wslapi.FlagAppendNTPath == wslapi.FlagAppendNTPath)

		case "--mount-drive":
			print(flags&wslapi.FlagEnableDriveMounting == wslapi.FlagEnableDriveMounting)

		case "--wsl-version":
			if flags&wslapi.FlagEnableWsl2 == wslapi.FlagEnableWsl2 {
				print("2")
			} else {
				print("1")
			}

		case "--lxguid", "--lxuid":
			if profile.UUID == "" {
				if proferr != nil {
					utils.ErrorExit(proferr, true, true, false)
				}
				utils.ErrorExit(errors.New("lxguid get failed"), true, true, false)
			}
			print(profile.UUID)

		case "--default-term", "--default-terminal":
			switch profile.WsldlTerm {
			case wslreg.FlagWsldlTermWT:
				print("wt")
			case wslreg.FlagWsldlTermFlute:
				print("flute")
			default:
				print("default")
			}

		case "--wt-profile-name", "--wt-profilename", "--wt-pn":
			if profile.DistributionName != "" {
				name = profile.DistributionName
			}

			conf, err := wtutils.ReadParseWTConfig()
			if err != nil {
				utils.ErrorExit(err, true, true, false)
			}
			guid := "{" + wtutils.CreateProfileGUID(name) + "}"
			profileName := ""
			for _, profile := range conf.Profiles.ProfileList {
				if profile.GUID == guid {
					profileName = profile.Name
					break
				}
			}
			if profileName != "" {
				print(profileName)
			} else {
				utils.ErrorExit(errors.New("profile not found"), true, true, false)
			}

		case "--flags-val":
			print(flags)

		case "--flags-bits":
			fmt.Printf("%04b", flags)

		default:
			utils.ErrorExit(os.ErrInvalid, true, true, false)
		}
	} else {
		utils.ErrorExit(os.ErrInvalid, true, true, false)
	}
}

//WslGetConfig is getter of distribution configuration
func WslGetConfig(distributionName string) (uid uint64, flags uint32) {
	var err error
	_, uid, flags, err = wslapi.WslGetDistributionConfiguration(distributionName)
	if err != nil {
		utils.ErrorRedPrintln("ERR: Failed to GetDistributionConfiguration")
		utils.ErrorExit(err, true, true, false)
	}
	return
}

// ShowHelp shows help message
func ShowHelp(showTitle bool) {
	if showTitle {
		println("Usage:")
	}
	println("    get [setting [value]]")
	println("      - `--default-uid`: Get the default user uid in this instance.")
	println("      - `--append-path`: Get true/false status of Append Windows PATH to $PATH.")
	println("      - `--mount-drive`: Get true/false status of Mount drives.")
	println("      - `--wsl-version`: Get the version os the WSL (1/2) of this instance.")
	println("      - `--default-term`: Get Default Terminal type of this instance launcher.")
	println("      - `--wt-profile-name`: Get Profile Name from Windows Terminal")
	println("      - `--lxguid`: Get WSL GUID key for this instance.")
}

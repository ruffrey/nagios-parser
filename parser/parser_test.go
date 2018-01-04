package parser

import (
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
)

func Test_Parser(t *testing.T) {
	t.Run("it parses one line of config", func(t *testing.T) {
		example := `define host {
  host_name                      fleetm4vmprod
  alias                          FLEETM4VMPROD (fleet m4 prod App Server)
  address                        169.237.48.71
  use                            template-windows
  hostgroups                     windows-vm,windows-service-iis
}
`

		result := ParseConfText(example)

		expected := []*NagiosConfig{
			{
				HostName: "fleetm4vmprod",
				Address:  "169.237.48.71",
			},
		}

		equality := assert.ObjectsAreEqualValues(expected, result)
		if !equality {
			log.Fatal("Objects are not equal, got", result)
		}
	})
	t.Run("it parses two lines of config with a space in between", func(t *testing.T) {
		example := `define host {
  host_name                      fleetm4vmprod
  alias                          FLEETM4VMPROD (fleet m4 prod App Server)
  address                        169.237.48.71
  use                            template-windows
  hostgroups                     windows-vm,windows-service-iis
}

define host {
  host_name                      armcivil
  alias                          ARMCIVIL (W2K8 Facilities ArcGIS System)
  address                        169.237.206.206
  use                            template-windows
  hostgroups                     windows-disk-g,windows-disk-h,windows-service-iis,windows-service-mssql,windows-vm
}
`

		result := ParseConfText(example)

		expected := []*NagiosConfig{
			{
				HostName: "fleetm4vmprod",
				Address:  "169.237.48.71",
			},
			{
				HostName: "armcivil",
				Address:  "169.237.206.206",
			},
		}

		equality := assert.ObjectsAreEqualValues(expected, result)
		if !equality {
			log.Fatal("Objects are not equal, got", result)
		}
	})
}

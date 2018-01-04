package parser

import "strings"

// ParseConfText returns an array of NagiosConfig
func ParseConfText(text string) (configList []*NagiosConfig) {
	p := &parser{
		lines: strings.Split(text, "\n"),
	}

	var config *NagiosConfig
	var done bool
	for {
		config, done = p.getNextSection()
		if config != nil && config.HostName != "" {
			configList = append(configList, config)
		}
		if done {
			break
		}
	}

	return configList
}

// NagiosConfig is a thing
type NagiosConfig struct {
	HostName string
	Address  string
}

type parser struct {
	lines  []string
	cursor int
}

func (p *parser) getNextSection() (config *NagiosConfig, done bool) {
	var isBeginningOfSection, isEndOfSection, isConfigLine bool
	var configKey, configValue string
	config = &NagiosConfig{}
	for {
		isBeginningOfSection, isEndOfSection, isConfigLine, done, configKey, configValue = p.parseNextLine()

		if isBeginningOfSection {
			continue
		}
		if isConfigLine {
			switch configKey {
			case "host_name":
				config.HostName = configValue
				break
			case "address":
				config.Address = configValue
				break
			default:
				// ignored
			}
		}
		if isEndOfSection || done {
			break
		}
	}
	return config, done
}

func (p *parser) parseNextLine() (isBeginningOfSection, isEndOfSection, isConfigLine, done bool, configKey, configValue string) {
	p.cursor++
	done = p.cursor > len(p.lines)-1
	if done {
		return isBeginningOfSection, isEndOfSection, isConfigLine, done, configKey, configValue
	}

	currentLine := p.lines[p.cursor]

	isBeginningOfSection = strings.HasPrefix(currentLine, "define host {")
	isEndOfSection = currentLine == "}"
	emptyLine := strings.Replace(strings.TrimSpace(currentLine), "\n", "", -1) == ""
	isConfigLine = !isBeginningOfSection && !isEndOfSection && !emptyLine

	if isConfigLine {
		// parse the config line key and value
		s := strings.TrimSpace(currentLine)
		firstSpaceIndex := strings.Index(s, " ")
		configKey = s[0:firstSpaceIndex]
		configValue = strings.TrimSpace(s[firstSpaceIndex:])
	}

	return isBeginningOfSection, isEndOfSection, isConfigLine, done, configKey, configValue
}

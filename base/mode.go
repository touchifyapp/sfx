package main

import (
	"strconv"
	"strings"
)

const modOUTDATED = -1
const modEQUAL = 0
const modUPDATE = 1
const modINSTALL = 2

func getInstallMode(config *peConfig, destConfig *peConfig) int8 {
	if destConfig == nil {
		return modINSTALL
	}

	srcVer := strings.Split(config.Version, ".")
	destVer := strings.Split(destConfig.Version, ".")

	for i, part := range srcVer {
		srcPart, err := strconv.ParseInt(part, 10, 8)
		if err != nil {
			return modEQUAL
		}

		destPart, err := strconv.ParseInt(destVer[i], 10, 8)
		if err != nil || srcPart > destPart {
			return modUPDATE
		}

		if srcPart < destPart {
			return modOUTDATED
		}
	}

	return modEQUAL
}

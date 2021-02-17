package service

import (
	"fmt"
	"github.com/skeyic/ark-robot/utils"
	"testing"
)
import "github.com/grd/statistics"

func TestGetMagicPercent(t *testing.T) {
	utils.EnableGlogForTesting()

	var (
		pl = statistics.Float64{
			0.009707, 0.009688, 0.009703, 0.009712, 0.009716, 0.009711, 0.235008, 0.074526,
			0.009710, 0.009709, 0.009705, 0.009699, 0.404598, 0.009705, 0.009708, 0.009710,
		}
	)

	fmt.Println(PickAbnormalData(pl))
}

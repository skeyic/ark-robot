package utils

import (
	"strings"
	"testing"
)

func TestAlert(t *testing.T) {

	var content = `基于ARK基金公开的截止2021-08-06（不含）的持仓数据，以下个股在连续三个交易日发生同向变动：
RHHBY最近三日在ARKG中均被减持，分别是6.66%、4.90%以及2.20%，持股数分别为3721019股、3986520股和4191893股。
SEER最近三日在ARKG中均被减持，分别是5.27%、7.59%以及1.14%，持股数分别为1067058股、1126398股和1218960股。
NTDOY最近三日在ARKK中均被减持，分别是2.36%、1.29%以及2.54%，持股数分别为3313065股、3393200股和3437538股。
MKFG最近三日在ARKQ中都获得增持，分别是6.36%、14.30%以及14.37%，持股数分别为3661860股、3442924股和3012085股。
TCEHY最近三日在ARKQ中均被减持，分别是100.00%、79.47%以及88.73%，持股数分别为0股、39股和190股。
AVAV最近三日在ARKQ中都获得增持，分别是2.05%、5.32%以及4.88%，持股数分别为447052股、438079股和415966股。
MASS最近三日在ARKG中都获得增持，分别是2.61%、0.04%以及1.77%，持股数分别为2748517股、2678517股和2677432股。`
	SendAlertV2("连续变动股票2021-08-06", strings.ReplaceAll(content, "\n", "    "))
}

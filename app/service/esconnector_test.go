package service

import (
	"context"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"testing"
)
import "github.com/olivere/elastic/v7"

func TestESConnector(t *testing.T) {
	var (
		ctx = context.Background()
		url = config.Config.ESServer.URL
	)

	utils.EnableGlogForTesting()

	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	)
	if err != nil {
		glog.Errorf("New ES client to %s failed, err: %v", url, err)
		return
	}

	indexExistResp, err := client.IndexExists(holdingsIndex).Do(ctx)
	if err != nil {
		glog.Errorf("Check index %s failed, err: %v", holdingsIndex, err)
		return
	}

	if !indexExistResp {
		indexCreateResp, err := client.CreateIndex(holdingsIndex).BodyString(holdingsIndexSettings).Do(ctx)
		if err != nil || !indexCreateResp.Acknowledged {
			glog.Errorf("Create index %s failed, err: %v", holdingsIndex, err)
			return
		}
	}

	bulkProcessor, err := elastic.NewBulkProcessorService(client).Do(ctx)
	if err != nil {
		glog.Errorf("Create bulk processor failed, err: %v", err)
		return
	}
	defer bulkProcessor.Close()

	TheMaster.FreshInit()
	for _, fund := range allARKTypes {
		for _, holding := range TheLibrary.LatestStockHoldings.GetFundStockHoldings(fund).Holdings {
			glog.V(4).Infof("HOLDINGS: %v", holding)
			bulkProcessor.Add(elastic.NewBulkIndexRequest().
				Index(holdingsIndex).Id(holding.ESID()).Doc(holding.ESBody()))
		}
	}

	err = bulkProcessor.Start(ctx)
	if err != nil {
		glog.Errorf("Start bulk processor failed, err: %v", err)
		return
	}

	err = bulkProcessor.Flush()
	if err != nil {
		glog.Errorf("Flush bulk processor failed, err: %v", err)
		return
	}

	resp := bulkProcessor.Stats()
	glog.V(3).Infof("PROCESSOR STATS: %+v", resp)

}

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/adorigi/checkctl/pkg/config"
	"github.com/adorigi/checkctl/pkg/request"
	"github.com/adorigi/checkctl/pkg/types"
	"github.com/adorigi/checkctl/pkg/utils"
	"github.com/spf13/cobra"
)

// controlsCmd represents the controls command
var controlsCmd = &cobra.Command{
	Use:   "controls",
	Short: "Get information for all controls",
	Long: `Get information for all controls
	Example usage:
		checkctl get controls --page-number 1 --page-size 20 --output json
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := &http.Client{}
		configuration, err := config.ReadConfigFile()
		if err != nil {
			return err
		}

		outputFormat := utils.ReadStringFlag(cmd, "output")
		if outputFormat == "" {
			outputFormat = configuration.OutputFormat
		}

		benchmarkIDs, err := utils.ReadStringSliceFlag(cmd, "benchmark-id")
		if err != nil {
			return err
		}

		if _, ok := configuration.Benchmarks[benchmarkIDs[0]]; ok {
			fmt.Printf("Found stored Benchmark %s", benchmarkIDs[0])
			benchmarkIDs = configuration.Benchmarks[benchmarkIDs[0]]
		}

		requestPayload := types.GetControlsPayload{
			Cursor:  int(utils.ReadIntFlag(cmd, "page-number")),
			PerPage: int(utils.ReadIntFlag(cmd, "page-size")),
			FindingFilters: types.FindingFilters{
				BenchmarkID: benchmarkIDs,
			},
		}

		payload, err := json.Marshal(requestPayload)
		if err != nil {
			return err
		}

		request, err := request.GenerateRequest(
			configuration.ApiKey,
			configuration.ApiEndpoint,
			"POST",
			"main/compliance/api/v3/controls",
			payload,
		)
		if err != nil {
			return err
		}

		response, err := client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if response.StatusCode != 200 {
			fmt.Println(string(body))
			return nil
		}

		var getControlsResponse types.GetControlsResponse
		err = json.Unmarshal(body, &getControlsResponse)
		if err != nil {
			return err
		}

		// if outputFormat == "table" {
		// 	rows := utils.GenerateControlRows(getControlsResponse.Items)

		// 	tables.PrintControlsTable(rows)
		// } else {
		js, err := json.MarshalIndent(getControlsResponse.Items, "", "   ")
		if err != nil {
			return err
		}
		fmt.Print(string(js))
		// }

		fmt.Printf(
			"\n\n\n\nNext Page: \n\tcheckctl get controls --page-size %d --page-number %d --output %s\n",
			utils.ReadIntFlag(cmd, "page-size"),
			utils.ReadIntFlag(cmd, "page-number")+1,
			outputFormat,
		)

		return nil
	},
}

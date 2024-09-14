package run

import (
	"encoding/json"
	"fmt"
	"github.com/adorigi/opengovernance/pkg/request"
	"io"
	"net/http"

	"github.com/adorigi/opengovernance/pkg/config"
	"github.com/adorigi/opengovernance/pkg/types"
	"github.com/adorigi/opengovernance/pkg/utils"
	"github.com/spf13/cobra"
)

// benchmarksCmd represents the benchmarks command
var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "Run specified benchmark on given integrations",
	Long:  `Run specified benchmark on given integrations.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		client := &http.Client{}
		configuration, err := config.ReadConfigFile()
		if err != nil {
			return err
		}

		integrationsStr, err := utils.ReadStringArrayFlag(cmd, "integration-info")
		if err != nil {
			return err
		}
		resourceTypes, err := utils.ReadStringArrayFlag(cmd, "resource-type")
		if err != nil {
			return err
		}
		forceFull := utils.ReadBoolFlag(cmd, "force-full")

		var integrations []types.IntegrationFilterInfo
		for _, integrationStr := range integrationsStr {
			integrations = append(integrations, types.ParseIntegrationInfo(integrationStr))
		}
		req := types.RunDiscoveryRequest{
			IntegrationInfo: integrations,
			ResourceTypes:   resourceTypes,
			ForceFull:       forceFull,
		}

		payload, err := json.Marshal(req)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("main/schedule/api/v2/discovery/run")
		request, err := request.GenerateRequest(
			configuration.ApiKey,
			configuration.ApiEndpoint,
			"POST",
			url,
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

		var runDiscoveryResponse []types.RunDiscoveryResponse
		err = json.Unmarshal(body, &runDiscoveryResponse)
		if err != nil {
			return err
		}

		if configuration.OutputFormat == "table" {
			// TODO
		} else {
			js, err := json.MarshalIndent(runDiscoveryResponse, "", "   ")
			if err != nil {
				return err
			}
			fmt.Print(string(js))
		}

		return nil
	},
}

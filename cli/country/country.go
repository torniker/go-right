package country

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/torniker/go-right/app/country"
	"github.com/torniker/go-right/app/country/request"
	"github.com/torniker/go-right/pkg/logger"
	"github.com/torniker/go-right/pkg/repeat"
)

func Register(rootCmd *cobra.Command) {
	countries := &cobra.Command{
		Use:   "country",
		Short: "Import countries",
		Long:  "Fetch countries from restcountries.eu and import them into database",
		Run:   repeat.Tick(countriesHandler),
	}
	rootCmd.AddCommand(countries)
}

func countriesHandler(c context.Context) error {
	service := country.New()
	err := fetch(c, service)
	if err != nil {
		return fmt.Errorf("fetch countries error: %w", err)
	}
	return nil
}

func fetch(c context.Context, service *country.Service) error {
	client := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, "https://restcountries.eu/rest/v2/all", nil)
	resp, err := client.Do(req.WithContext(c))
	if err != nil {
		return fmt.Errorf("could not http.Get: %w", err)
	}
	defer resp.Body.Close()
	var countries []request.CountrySave
	err = json.NewDecoder(resp.Body).Decode(&countries)
	if err != nil {
		return err
	}
	for _, country := range countries {
		err = service.Save(c, country)
		if err != nil {
			logger.Errorf("country save err: %s", err)
		}
	}
	return nil
}
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand/v2"
	"os"
)

// EmailSubscription represents the structure of the JSON object that both Errata and Notifications agreed upon.
type EmailSubscription struct {
	Username        string   `json:"user_name"`
	OrgId           string   `json:"org_id"`
	UserPreferences []string `json:"preferences"`
}

// OutputFilename represents the filename where we are going to put the subscriptions in.
const OutputFilename = "subscriptions.json"

// SupportedSubscriptionTypes represents the email subscription types both Errata and Notifications support as of date.
var SupportedSubscriptionTypes = []string{"bugfix", "enhancement", "security"}

func main() {
	// Grab the number of email subscriptions we want to create.
	subscriptionsToGenerate := flag.Int("number", 0, "number of email subscriptions to generate")
	flag.Parse()

	// We should have at least one argument!
	if subscriptionsToGenerate == nil || *subscriptionsToGenerate == 0 {
		flag.Usage()
		os.Exit(0)
	}

	// Generate the email subscriptions.
	generatedEmailSubscriptions := make([]EmailSubscription, 0, *subscriptionsToGenerate)
	for i := 0; i < *subscriptionsToGenerate; i++ {
		id := uuid.New()

		// Pick a random number of user preferences to generate.
		numberSubscriptionsGenerate := rand.IntN(len(SupportedSubscriptionTypes))
		// Create a map to avoid duplicates.
		userPreferences := make([]string, 0, numberSubscriptionsGenerate)
		if numberSubscriptionsGenerate > 0 {
			var subscriptions = make(map[string]struct{})
			for len(subscriptions) < numberSubscriptionsGenerate {
				random := rand.IntN(len(SupportedSubscriptionTypes))

				// Take one subscription from the supported subscriptions' slice.
				sub := SupportedSubscriptionTypes[random]

				// When the subscription type is not present in the map add it, otherwise keep looping.
				if _, ok := subscriptions[sub]; !ok {
					subscriptions[sub] = struct{}{}
				}
			}

			// Populate the user preferences.
			for key, _ := range subscriptions {
				userPreferences = append(userPreferences, key)
			}
		}

		// Generate the email subscription.
		generatedEmailSubscriptions = append(generatedEmailSubscriptions, EmailSubscription{
			Username:        fmt.Sprintf("emt-%s", id),
			OrgId:           fmt.Sprintf("emt-%s", id),
			UserPreferences: userPreferences,
		})
	}

	// Marshal the file body.
	fileContents, err := json.Marshal(generatedEmailSubscriptions)
	if err != nil {
		log.Fatalf("Unable to marshal the file contents: %s", err)
	}

	// Write the contents to the file.
	err = os.WriteFile(OutputFilename, fileContents, 0644)
	if err != nil {
		log.Fatalf("Unable to write the contents to the file: %s", err)
	}

	fmt.Printf(`Successfully generated "%d" email subscriptions in the "%s" file`+"\n", *subscriptionsToGenerate, OutputFilename)
}

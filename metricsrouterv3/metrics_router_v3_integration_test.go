// +build integration

/**
 * (C) Copyright IBM Corp. 2022.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metricsrouterv3_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/metricsrouterv3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

/**
 * This file contains an integration test for the metricsrouterv3 package.
 *
 * Notes:
 *
 * The integration test will automatically skip tests if the required config file is not available.
 */

var _ = Describe(`MetricsRouterV3 Integration Tests`, func() {
	const externalConfigFile = "../metrics_router_v3.env"

	var (
		err          error
		metricsRouterService *metricsrouterv3.MetricsRouterV3
		serviceURL   string
		config       map[string]string

		// Variables to hold link values
		routeIDLink string
		targetIDLink string
	)

	var shouldSkipTest = func() {
		Skip("External configuration is not available, skipping tests...")
	}

	Describe(`External configuration`, func() {
		It("Successfully load the configuration", func() {
			_, err = os.Stat(externalConfigFile)
			if err != nil {
				Skip("External configuration file not found, skipping tests: " + err.Error())
			}

			os.Setenv("IBM_CREDENTIALS_FILE", externalConfigFile)
			config, err = core.GetServiceProperties(metricsrouterv3.DefaultServiceName)
			if err != nil {
				Skip("Error loading service properties, skipping tests: " + err.Error())
			}
			serviceURL = config["URL"]
			if serviceURL == "" {
				Skip("Unable to load service URL configuration property, skipping tests")
			}

			fmt.Fprintf(GinkgoWriter, "Service URL: %v\n", serviceURL)
			shouldSkipTest = func() {}
		})
	})

	Describe(`Client initialization`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It("Successfully construct the service client instance", func() {
			metricsRouterServiceOptions := &metricsrouterv3.MetricsRouterV3Options{}

			metricsRouterService, err = metricsrouterv3.NewMetricsRouterV3UsingExternalConfig(metricsRouterServiceOptions)
			Expect(err).To(BeNil())
			Expect(metricsRouterService).ToNot(BeNil())
			Expect(metricsRouterService.Service.Options.URL).To(Equal(serviceURL))

			core.SetLogger(core.NewLogger(core.LevelDebug, log.New(GinkgoWriter, "", log.LstdFlags), log.New(GinkgoWriter, "", log.LstdFlags)))
			metricsRouterService.EnableRetries(4, 30*time.Second)
		})
	})

	Describe(`CreateTarget - Create a target`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateTarget(createTargetOptions *CreateTargetOptions)`, func() {
			createTargetOptions := &metricsrouterv3.CreateTargetOptions{
				Name: core.StringPtr("my-mr-target"),
				DestinationCRN: core.StringPtr("crn:v1:bluemix:public:sysdig-monitor:us-south:a/0be5ad401ae913d8ff665d92680664ed:22222222-2222-2222-2222-222222222222::"),
				Region: core.StringPtr("us-south"),
			}

			target, response, err := metricsRouterService.CreateTarget(createTargetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(target).ToNot(BeNil())

			targetIDLink = *target.ID
			fmt.Fprintf(GinkgoWriter, "Saved targetIDLink value: %v\n", targetIDLink)
		})
	})

	Describe(`CreateRoute - Create a route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateRoute(createRouteOptions *CreateRouteOptions)`, func() {
			inclusionFilterModel := &metricsrouterv3.InclusionFilter{
				Operand: core.StringPtr("location"),
				Operator: core.StringPtr("is"),
				Value: []string{"testString"},
			}

			rulePrototypeModel := &metricsrouterv3.RulePrototype{
				TargetIds: []string{"c3af557f-fb0e-4476-85c3-0889e7fe7bc4"},
				InclusionFilters: []metricsrouterv3.InclusionFilter{*inclusionFilterModel},
			}

			createRouteOptions := &metricsrouterv3.CreateRouteOptions{
				Name: core.StringPtr("my-route"),
				Rules: []metricsrouterv3.RulePrototype{*rulePrototypeModel},
			}

			route, response, err := metricsRouterService.CreateRoute(createRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(route).ToNot(BeNil())

			routeIDLink = *route.ID
			fmt.Fprintf(GinkgoWriter, "Saved routeIDLink value: %v\n", routeIDLink)
		})
	})

	Describe(`ListTargets - List targets`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListTargets(listTargetsOptions *ListTargetsOptions)`, func() {
			listTargetsOptions := &metricsrouterv3.ListTargetsOptions{
			}

			targetList, response, err := metricsRouterService.ListTargets(listTargetsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(targetList).ToNot(BeNil())
		})
	})

	Describe(`GetTarget - Get details of a target`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetTarget(getTargetOptions *GetTargetOptions)`, func() {
			getTargetOptions := &metricsrouterv3.GetTargetOptions{
				ID: &targetIDLink,
			}

			target, response, err := metricsRouterService.GetTarget(getTargetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(target).ToNot(BeNil())
		})
	})

	Describe(`ReplaceTarget - Update a target`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ReplaceTarget(replaceTargetOptions *ReplaceTargetOptions)`, func() {
			replaceTargetOptions := &metricsrouterv3.ReplaceTargetOptions{
				ID: &targetIDLink,
				Name: core.StringPtr("my-mr-target"),
				DestinationCRN: core.StringPtr("crn:v1:bluemix:public:sysdig-monitor:us-south:a/0be5ad401ae913d8ff665d92680664ed:22222222-2222-2222-2222-222222222222::"),
			}

			target, response, err := metricsRouterService.ReplaceTarget(replaceTargetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(target).ToNot(BeNil())
		})
	})

	Describe(`ValidateTarget - Validate a target`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ValidateTarget(validateTargetOptions *ValidateTargetOptions)`, func() {
			validateTargetOptions := &metricsrouterv3.ValidateTargetOptions{
				ID: &targetIDLink,
			}

			target, response, err := metricsRouterService.ValidateTarget(validateTargetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(target).ToNot(BeNil())
		})
	})

	Describe(`ListRoutes - List routes`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListRoutes(listRoutesOptions *ListRoutesOptions)`, func() {
			listRoutesOptions := &metricsrouterv3.ListRoutesOptions{
			}

			routeList, response, err := metricsRouterService.ListRoutes(listRoutesOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(routeList).ToNot(BeNil())
		})
	})

	Describe(`GetRoute - Get details of a route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetRoute(getRouteOptions *GetRouteOptions)`, func() {
			getRouteOptions := &metricsrouterv3.GetRouteOptions{
				ID: &routeIDLink,
			}

			route, response, err := metricsRouterService.GetRoute(getRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(route).ToNot(BeNil())
		})
	})

	Describe(`ReplaceRoute - Update a route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ReplaceRoute(replaceRouteOptions *ReplaceRouteOptions)`, func() {
			inclusionFilterModel := &metricsrouterv3.InclusionFilter{
				Operand: core.StringPtr("location"),
				Operator: core.StringPtr("is"),
				Value: []string{"testString"},
			}

			rulePrototypeModel := &metricsrouterv3.RulePrototype{
				TargetIds: []string{"c3af557f-fb0e-4476-85c3-0889e7fe7bc4"},
				InclusionFilters: []metricsrouterv3.InclusionFilter{*inclusionFilterModel},
			}

			replaceRouteOptions := &metricsrouterv3.ReplaceRouteOptions{
				ID: &routeIDLink,
				Name: core.StringPtr("my-route"),
				Rules: []metricsrouterv3.RulePrototype{*rulePrototypeModel},
			}

			route, response, err := metricsRouterService.ReplaceRoute(replaceRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(route).ToNot(BeNil())
		})
	})

	Describe(`GetSettings - Get settings`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetSettings(getSettingsOptions *GetSettingsOptions)`, func() {
			getSettingsOptions := &metricsrouterv3.GetSettingsOptions{
			}

			settings, response, err := metricsRouterService.GetSettings(getSettingsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(settings).ToNot(BeNil())
		})
	})

	Describe(`ReplaceSettings - Modify settings`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ReplaceSettings(replaceSettingsOptions *ReplaceSettingsOptions)`, func() {
			replaceSettingsOptions := &metricsrouterv3.ReplaceSettingsOptions{
				MetadataRegionPrimary: core.StringPtr("us-south"),
				PrivateAPIEndpointOnly: core.BoolPtr(false),
				DefaultTargets: []string{"c3af557f-fb0e-4476-85c3-0889e7fe7bc4"},
				PermittedTargetRegions: []string{"us-south"},
			}

			settings, response, err := metricsRouterService.ReplaceSettings(replaceSettingsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(settings).ToNot(BeNil())
		})
	})

	Describe(`DeleteRoute - Delete a route`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteRoute(deleteRouteOptions *DeleteRouteOptions)`, func() {
			deleteRouteOptions := &metricsrouterv3.DeleteRouteOptions{
				ID: &routeIDLink,
			}

			response, err := metricsRouterService.DeleteRoute(deleteRouteOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteTarget - Delete a target`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteTarget(deleteTargetOptions *DeleteTargetOptions)`, func() {
			deleteTargetOptions := &metricsrouterv3.DeleteTargetOptions{
				ID: &targetIDLink,
			}

			warningReport, response, err := metricsRouterService.DeleteTarget(deleteTargetOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(warningReport).ToNot(BeNil())
		})
	})
})

//
// Utility functions are declared in the unit test file
//

package healthchecks

import (
	"fmt"
)

func Traverse(s []StatusEndpoint, dependencies []string, action string, protocol string, aboutFilePath string, versionFilePath string, customData map[string]interface{}) string {

	if action == "" {
		action = "about"
	}

	// base case
	if len(dependencies) == 0 {
		// run the action
		switch action {
		case "about":
			aboutResp, _ := About(s, protocol, aboutFilePath, versionFilePath, customData, APIV1, true)
			return aboutResp
		default:
			sl := StatusList{
				StatusList: []Status{
					{
						Description: "Unsupported action",
						Result:      CRITICAL,
						Details:     fmt.Sprintf("Unsupported traversal action '%s'", action),
					},
				},
			}

			return SerializeStatusList(sl, APIV1)
		}
	}

	headDependency := dependencies[0]
	headStatusEndpoint := FindStatusEndpoint(s, headDependency)

	if headStatusEndpoint == nil {
		sl := StatusList{
			StatusList: []Status{
				{
					Description: "Can't traverse",
					Result:      CRITICAL,
					Details:     fmt.Sprintf("Status path '%s' is not registered", headDependency),
				},
			},
		}

		return SerializeStatusList(sl, APIV1)
	}

	if !headStatusEndpoint.IsTraversable {
		sl := StatusList{
			StatusList: []Status{
				{
					Description: "Can't traverse",
					Result:      CRITICAL,
					Details:     fmt.Sprintf("%s is not traversable", headStatusEndpoint.Name),
				},
			},
		}

		return SerializeStatusList(sl, APIV1)
	}

	if headStatusEndpoint.TraverseCheck == nil {
		sl := StatusList{
			StatusList: []Status{
				{
					Description: "Can't traverse",
					Result:      CRITICAL,
					Details:     fmt.Sprintf("%s does not have a TraverseCheck() function defined", headStatusEndpoint.Name),
				},
			},
		}

		return SerializeStatusList(sl, APIV1)
	}

	// found dependency, continue to traverse with the tail of the dependencies
	tailDependencies := dependencies[1:]
	resp, err := headStatusEndpoint.TraverseCheck.Traverse(tailDependencies, action)
	if err != nil {
		sl := StatusList{
			StatusList: []Status{
				{
					Description: "Traverse",
					Result:      CRITICAL,
					Details:     err.Error(),
				},
			},
		}

		return SerializeStatusList(sl, APIV1)
	} else {
		return resp
	}
}

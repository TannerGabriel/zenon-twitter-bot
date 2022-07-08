package zenon

import (
	"fmt"
	"github.com/zenon-wiki/go-zdk/zdk"
	"math/big"
)

const (
	Accepted  = 1
	Paid      = 2
	Closed    = 3
	Completed = 4
)

func convertFundValues(value big.Int) string {
	return new(big.Int).Div(&value, big.NewInt(100000000)).String()
}

// HandleNewProject creates the Tweet for the project:new event
func HandleNewProject(project ProjectNew) string {
	tweetMessage := fmt.Sprintf(`New #AcceleratorZ project submission:  %s
								
								Requested funds: %d $ZNN  %d $QSR
								
								Project URL:
								%s`,
		project.Data.Name, int(project.Data.Znn), int(project.Data.Qsr), project.Data.Url,
	)

	return tweetMessage
}

// HandleProjectStatusUpdate fetches the project data and creates the Tweet for the project:status-update event
func HandleProjectStatusUpdate(statusUpdate ProjectStatusUpdate, zenon *zdk.Zdk) (string, error) {
	project, err := zenon.Embedded.Accelerator.GetProjectById(statusUpdate.Id)
	if err != nil {
		return "", err
	}

	// Check if project has been accepted
	if statusUpdate.NewStatus == Accepted {
		tweetMessage := fmt.Sprintf(`%s has been accepted into #AcceleratorZ
								
								Votes:
								Yes: %d 
								No: %d

								Funds Granted:
								%s $ZNN & %s $QSR

								%s`,
			project.Name, project.Votes.Yes, project.Votes.No, convertFundValues(*project.ZnnFundsNeeded), convertFundValues(*project.QsrFundsNeeded), project.Url,
		)

		return tweetMessage, nil
	}

	return "", nil
}

// HandlePhaseStatusUpdate fetches the phase data and creates the Tweet for the phase:status-update event
func HandlePhaseStatusUpdate(statusUpdate PhaseStatusUpdate, zenon *zdk.Zdk) (string, error) {
	phase, err := zenon.Embedded.Accelerator.GetPhaseById(statusUpdate.Id)
	if err != nil {
		return "", err
	}

	project, err := zenon.Embedded.Accelerator.GetProjectById(statusUpdate.Pid)
	if err != nil {
		return "", err
	}

	// Check if phase has been paid
	if statusUpdate.NewStatus == Paid {
		tweetMessage := fmt.Sprintf(`%s has been paid for %s

								Funds Granted:
								%s $ZNN & %s $QSR

								Phase URL:
								%s`,
			phase.Phase.Name, project.Name, convertFundValues(*phase.Phase.ZnnFundsNeeded), convertFundValues(*phase.Phase.QsrFundsNeeded), phase.Phase.Url,
		)

		return tweetMessage, err
	}
	return "", err
}

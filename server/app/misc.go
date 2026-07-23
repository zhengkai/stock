package app

import (
	"project/pb"
)

func normalizationAlert(a *pb.Alert) {
	if a.GetMax() == 0 && a.GetMin() > 0 {
		a.SetMax(a.GetMin() * 2)
	} else if a.GetMin() == 0 && a.GetMax() > 0 {
		a.SetMin(a.GetMax() / 2)
	}
}

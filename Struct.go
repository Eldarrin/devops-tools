type JobRequests struct {
	Count int `json:"count"`
	Value []struct {
		JobRequest struct {
			RequestId     int       `json:"requestId"`
			QueueTime     time.Time `json:"queueTime"`
			AssignTime    time.Time `json:"assignTime"`
			ReceiveTime   time.Time `json:"receiveTime"`
			LockedUntil   time.Time `json:"lockedUntil"`
			ServiceOwner  string    `json:"serviceOwner"`
			Result        *string   `json:"Result"`
			HostId        string    `json:"hostId"`
			ScopeId       string    `json:"scopeId"`
			PlanType      string    `json:"planType"`
			PlanId        string    `json:"planId"`
			JobId         string    `json:"jobId"`
			Demands       []string  `json:"demands"`
			MatchedAgents *[]struct {
				Links struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					Web struct {
						Href string `json:"href"`
					} `json:"web"`
				} `json:"_links"`
				Id                int    `json:"id"`
				Name              string `json:"name"`
				Version           string `json:"version"`
				Enabled           bool   `json:"enabled"`
				Status            string `json:"status"`
				ProvisioningState string `json:"provisioningState"`
			} `json:"matchedAgents"`
			ReservedAgent *struct {
				Links struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					Web struct {
						Href string `json:"href"`
					} `json:"web"`
				} `json:"_links"`
				Id                int    `json:"id"`
				Name              string `json:"name"`
				Version           string `json:"version"`
				OsDescription     string `json:"osDescription"`
				Enabled           bool   `json:"enabled"`
				Status            string `json:"status"`
				ProvisioningState string `json:"provisioningState"`
				AccessPoint       string `json:"accessPoint"`
			} `json:"reservedAgent"`
			Definition struct {
				Links struct {
					Web struct {
						Href string `json:"href"`
					} `json:"web"`
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"_links"`
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"definition"`
			Owner struct {
				Links struct {
					Web struct {
						Href string `json:"href"`
					} `json:"web"`
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"_links"`
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"owner"`
			Data struct {
				ParallelismTag string `json:"ParallelismTag"`
				IsScheduledKey string `json:"IsScheduledKey"`
			} `json:"data"`
			PoolId          int    `json:"poolId"`
			OrchestrationId string `json:"orchestrationId"`
			Priority        int    `json:"priority"`
		}
	} `json:"value"`
}

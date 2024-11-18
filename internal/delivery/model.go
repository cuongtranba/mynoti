package delivery

type Comic struct {
	Url         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Html        string `json:"html"`
	CronSpec    string `json:"cron_spec"`
}

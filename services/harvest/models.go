package harvest

import "time"

type Client struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	IsActive     *bool      `json:"is_active"`
	Address      *string    `json:"address"`
	StatementKey *string    `json:"statement_key"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	Currency     *string    `json:"currency"`
}

type Project struct {
	ID                               int        `json:"id"`
	Name                             string     `json:"name"`
	Code                             *string    `json:"code"`
	IsActive                         *bool      `json:"is_active"`
	BillBy                           *string    `json:"bill_by"`
	Budget                           *float64   `json:"budget"`
	BudgetBy                         *string    `json:"budget_by"`
	BudgetIsMonthly                  *bool      `json:"budget_is_monthly"`
	NotifyWhenOverBudget             *bool      `json:"notify_when_over_budget"`
	OverBudgetNotificationPercentage *float64   `json:"over_budget_notification_percentage"`
	OverBudgetNotificationDate       *time.Time `json:"over_budget_notification_date"`
	ShowBudgetToAll                  *bool      `json:"show_budget_to_all"`
	CreatedAt                        *time.Time `json:"created_at"`
	UpdatedAt                        *time.Time `json:"updated_at"`
	StartsOn                         *string    `json:"starts_on"`
	EndsOn                           *time.Time `json:"ends_on"`
	IsBillable                       *bool      `json:"is_billable"`
	IsFixedFee                       *bool      `json:"is_fixed_fee"`
	Notes                            *string    `json:"notes"`
	Client                           *Client    `json:"client"`
	CostBudget                       *float64   `json:"cost_budget"`
	CostBudgetIncludeExpenses        *bool      `json:"cost_budget_include_expenses"`
	HourlyRate                       *float64   `json:"hourly_rate"`
	Fee                              *float64   `json:"fee"`
}

type User struct {
	ID                           int        `json:"id"`
	FirstName                    string     `json:"first_name"`
	LastName                     *string    `json:"last_name"`
	Email                        *string    `json:"email"`
	Telephone                    *string    `json:"telephone"`
	Timezone                     *string    `json:"timezone"`
	HasAccessToAllFutureProjects *bool      `json:"has_access_to_all_future_projects"`
	IsContractor                 *bool      `json:"is_contractor"`
	IsActive                     *bool      `json:"is_active"`
	CreatedAt                    *time.Time `json:"created_at"`
	UpdatedAt                    *time.Time `json:"updated_at"`
	WeeklyCapacity               *int       `json:"weekly_capacity"`
	DefaultHourlyRate            *float64   `json:"default_hourly_rate"`
	CostRate                     *float64   `json:"cost_rate"`
	Roles                        []string   `json:"roles"`
	AvatarURL                    *string    `json:"avatar_url"`
}

type Task struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	BillableByDefault *bool      `json:"billable_by_default"`
	DefaultHourlyRate *float64   `json:"default_hourly_rate"`
	IsDefault         *bool      `json:"is_default"`
	IsActive          *bool      `json:"is_active"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

type ProjectUserAssignment struct {
	ID               int       `json:"id"`
	IsProjectManager bool      `json:"is_project_manager"`
	IsActive         bool      `json:"is_active"`
	UseDefaultRates  bool      `json:"use_default_rates"`
	Budget           *float64  `json:"budget"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	HourlyRate       float64   `json:"hourly_rate"`
	Project          *Project  `json:"project"`
	User             *User     `json:"user"`
}

type ProjectTaskAssignment struct {
	ID         int       `json:"id"`
	Billable   bool      `json:"billable"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	HourlyRate float64   `json:"hourly_rate"`
	Budget     *float64  `json:"budget"`
	Project    *Project  `json:"project"`
	Task       *Task     `json:"task"`
}

type LineItems struct {
	ID          int     `json:"id"`
	Kind        string  `json:"kind"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   int     `json:"unit_price"`
	Amount      int     `json:"amount"`
	Taxed       bool    `json:"taxed"`
	Taxed2      bool    `json:"taxed2"`
	Project     Project `json:"project"`
}

type Invoice struct {
	ID                 int         `json:"id"`
	ClientKey          *string     `json:"client_key"`
	Number             string      `json:"number"`
	PurchaseOrder      *string     `json:"purchase_order"`
	Amount             *float64    `json:"amount"`
	DueAmount          *float64    `json:"due_amount"`
	Tax                *int        `json:"tax"`
	TaxAmount          *float64    `json:"tax_amount"`
	Tax2               *int        `json:"tax2"`
	Tax2Amount         *float64    `json:"tax2_amount"`
	Discount           *int        `json:"discount"`
	DiscountAmount     *int        `json:"discount_amount"`
	Subject            *string     `json:"subject"`
	Notes              *string     `json:"notes"`
	State              *string     `json:"state"`
	PeriodStart        *string     `json:"period_start"`
	PeriodEnd          *string     `json:"period_end"`
	IssueDate          *string     `json:"issue_date"`
	DueDate            *string     `json:"due_date"`
	PaymentTerm        *string     `json:"payment_term"`
	SentAt             *time.Time  `json:"sent_at"`
	PaidAt             *time.Time  `json:"paid_at"`
	PaidDate           *time.Time  `json:"paid_date"`
	ClosedAt           *time.Time  `json:"closed_at"`
	RecurringInvoiceID *int        `json:"recurring_invoice_id"`
	CreatedAt          *time.Time  `json:"created_at"`
	UpdatedAt          *time.Time  `json:"updated_at"`
	Currency           *string     `json:"currency"`
	Client             *Client     `json:"client"`
	Estimate           *Estimate   `json:"estimate"`
	Retainer           interface{} `json:"retainer"`
	Creator            *User       `json:"creator"`
	LineItems          []LineItems `json:"line_items"`
}

type Estimate struct {
	ID             int         `json:"id"`
	ClientKey      string      `json:"client_key"`
	Number         string      `json:"number"`
	PurchaseOrder  string      `json:"purchase_order"`
	Amount         float64     `json:"amount"`
	Tax            float64     `json:"tax"`
	TaxAmount      float64     `json:"tax_amount"`
	Tax2           float64     `json:"tax2"`
	Tax2Amount     float64     `json:"tax2_amount"`
	Discount       float64     `json:"discount"`
	DiscountAmount float64     `json:"discount_amount"`
	Subject        string      `json:"subject"`
	Notes          string      `json:"notes"`
	State          string      `json:"state"`
	IssueDate      string      `json:"issue_date"`
	SentAt         time.Time   `json:"sent_at"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	AcceptedAt     interface{} `json:"accepted_at"`
	DeclinedAt     interface{} `json:"declined_at"`
	Currency       string      `json:"currency"`
	Client         Client      `json:"client"`
	Creator        User        `json:"creator"`
	LineItems      []LineItems `json:"line_items"`
}

type ExternalReference struct {
	ID             int         `json:"id"`
	GroupID        int         `json:"group_id"`
	AccountID      int         `json:"account_id"`
	Permalink      string      `json:"permalink"`
	Service        interface{} `json:"service"`
	ServiceIconUrl string      `json:"service_icon_url"`
}

type TimeEntry struct {
	ID                int                   `json:"id"`
	SpentDate         time.Time             `json:"spent_date"`
	User              User                  `json:"user"`
	Client            Client                `json:"client"`
	Project           Project               `json:"project"`
	Task              Task                  `json:"task"`
	UserAssignment    ProjectUserAssignment `json:"user_assignment"`
	TaskAssignment    ProjectTaskAssignment `json:"task_assignment"`
	Hours             float64               `json:"hours"`
	HoursWithoutTimer float64               `json:"hours_without_timer"`
	RoundedHours      float64               `json:"rounded_hours"`
	Notes             string                `json:"notes"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
	IsLocked          bool                  `json:"is_locked"`
	LockedReason      string                `json:"locked_reason"`
	IsClosed          bool                  `json:"is_closed"`
	IsBilled          bool                  `json:"is_billed"`
	TimerStartedAt    time.Time             `json:"timer_started_at"`
	StartedTime       string                `json:"started_time"`
	EndedTime         string                `json:"ended_time"`
	IsRunning         bool                  `json:"is_running"`
	Invoice           Invoice               `json:"invoice"`
	ExternalReference *ExternalReference    `json:"external_reference"`
	Billable          bool                  `json:"billable"`
	Budgeted          bool                  `json:"budgeted"`
	BillableRate      float64               `json:"billable_rate"`
	CostRate          float64               `json:"cost_rate"`
}

type TimeEntryResponse struct {
	TimeEntries  []TimeEntry `json:"time_entries"`
	PerPage      int         `json:"per_page"`
	TotalPages   int         `json:"total_pages"`
	TotalEntries int         `json:"total_entries"`
	NextPage     string      `json:"next_page"`
	PreviousPage string      `json:"previous_page"`
	Page         int         `json:"page"`
	Links        struct {
		First    string `json:"first"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Last     string `json:"last"`
	} `json:"links"`
}

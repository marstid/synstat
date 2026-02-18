package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	apiURL         = "https://api.synthetic.new/v2/quotas"
	updateInterval = time.Minute
)

type Theme struct {
	Name          string
	Primary       lipgloss.Color
	Secondary     lipgloss.Color
	Background    lipgloss.Color
	Surface       lipgloss.Color
	Text          lipgloss.Color
	TextMuted     lipgloss.Color
	TextDim       lipgloss.Color
	Border        lipgloss.Color
	Success       lipgloss.Color
	Warning       lipgloss.Color
	Danger        lipgloss.Color
	ProgressStart lipgloss.Color
	ProgressEnd   lipgloss.Color
}

var themes = []Theme{
	{
		Name:          "Cyan (Default)",
		Primary:       "#22D3EE",
		Secondary:     "#06B6D4",
		Background:    "#18181B",
		Surface:       "#27272A",
		Text:          "#FAFAFA",
		TextMuted:     "#A1A1AA",
		TextDim:       "#52525B",
		Border:        "#3F3F46",
		Success:       "#4ADE80",
		Warning:       "#FACC15",
		Danger:        "#F87171",
		ProgressStart: "#22D3EE",
		ProgressEnd:   "#0891B2",
	},
	{
		Name:          "Blue",
		Primary:       "#60A5FA",
		Secondary:     "#3B82F6",
		Background:    "#0F172A",
		Surface:       "#1E293B",
		Text:          "#F8FAFC",
		TextMuted:     "#94A3B8",
		TextDim:       "#475569",
		Border:        "#334155",
		Success:       "#4ADE80",
		Warning:       "#FACC15",
		Danger:        "#F87171",
		ProgressStart: "#60A5FA",
		ProgressEnd:   "#2563EB",
	},
	{
		Name:          "Green",
		Primary:       "#4ADE80",
		Secondary:     "#22C55E",
		Background:    "#052e16",
		Surface:       "#064e3b",
		Text:          "#F0FDF4",
		TextMuted:     "#86EFAC",
		TextDim:       "#166534",
		Border:        "#15803D",
		Success:       "#4ADE80",
		Warning:       "#FACC15",
		Danger:        "#F87171",
		ProgressStart: "#4ADE80",
		ProgressEnd:   "#16A34A",
	},
	{
		Name:          "Purple",
		Primary:       "#A78BFA",
		Secondary:     "#8B5CF6",
		Background:    "#18181B",
		Surface:       "#27272A",
		Text:          "#FAFAFA",
		TextMuted:     "#A1A1AA",
		TextDim:       "#52525B",
		Border:        "#3F3F46",
		Success:       "#4ADE80",
		Warning:       "#FACC15",
		Danger:        "#F87171",
		ProgressStart: "#A78BFA",
		ProgressEnd:   "#7C3AED",
	},
	{
		Name:          "Candy",
		Primary:       "#F472B6",
		Secondary:     "#EC4899",
		Background:    "#500724",
		Surface:       "#831843",
		Text:          "#FDF2F8",
		TextMuted:     "#F9A8D4",
		TextDim:       "#BE185D",
		Border:        "#9D174D",
		Success:       "#86EFAC",
		Warning:       "#FDE047",
		Danger:        "#FCA5A5",
		ProgressStart: "#F472B6",
		ProgressEnd:   "#DB2777",
	},
	{
		Name:          "Soda-pop",
		Primary:       "#22D3EE",
		Secondary:     "#06B6D4",
		Background:    "#083344",
		Surface:       "#164E63",
		Text:          "#ECFEFF",
		TextMuted:     "#67E8F9",
		TextDim:       "#0E7490",
		Border:        "#155E75",
		Success:       "#86EFAC",
		Warning:       "#FDE047",
		Danger:        "#FCA5A5",
		ProgressStart: "#67E8F9",
		ProgressEnd:   "#0891B2",
	},
	{
		Name:          "Hacker",
		Primary:       "#00FF00",
		Secondary:     "#00CC00",
		Background:    "#000000",
		Surface:       "#0D1117",
		Text:          "#00FF00",
		TextMuted:     "#00CC00",
		TextDim:       "#009900",
		Border:        "#003300",
		Success:       "#00FF00",
		Warning:       "#FFFF00",
		Danger:        "#FF0000",
		ProgressStart: "#00FF00",
		ProgressEnd:   "#003300",
	},
	{
		Name:          "Sunset",
		Primary:       "#FB923C",
		Secondary:     "#F97316",
		Background:    "#2A1810",
		Surface:       "#431407",
		Text:          "#FFF7ED",
		TextMuted:     "#FDBA74",
		TextDim:       "#9A3412",
		Border:        "#7C2D12",
		Success:       "#86EFAC",
		Warning:       "#FDE047",
		Danger:        "#FCA5A5",
		ProgressStart: "#FB923C",
		ProgressEnd:   "#C2410C",
	},
	{
		Name:          "Ocean",
		Primary:       "#38BDF8",
		Secondary:     "#0EA5E9",
		Background:    "#082F49",
		Surface:       "#0C4A6E",
		Text:          "#F0F9FF",
		TextMuted:     "#7DD3FC",
		TextDim:       "#0369A1",
		Border:        "#075985",
		Success:       "#86EFAC",
		Warning:       "#FDE047",
		Danger:        "#FCA5A5",
		ProgressStart: "#38BDF8",
		ProgressEnd:   "#0284C7",
	},
	{
		Name:          "Forest",
		Primary:       "#84CC16",
		Secondary:     "#65A30D",
		Background:    "#1A2E05",
		Surface:       "#365314",
		Text:          "#F7FEE7",
		TextMuted:     "#BEF264",
		TextDim:       "#4D7C0F",
		Border:        "#3F6212",
		Success:       "#86EFAC",
		Warning:       "#FDE047",
		Danger:        "#FCA5A5",
		ProgressStart: "#A3E635",
		ProgressEnd:   "#4D7C0F",
	},
	{
		Name:          "Monochrome",
		Primary:       "#E4E4E7",
		Secondary:     "#A1A1AA",
		Background:    "#18181B",
		Surface:       "#27272A",
		Text:          "#FAFAFA",
		TextMuted:     "#A1A1AA",
		TextDim:       "#52525B",
		Border:        "#3F3F46",
		Success:       "#A1A1AA",
		Warning:       "#A1A1AA",
		Danger:        "#A1A1AA",
		ProgressStart: "#A1A1AA",
		ProgressEnd:   "#52525B",
	},
}

func findThemeByName(name string) (Theme, int) {
	nameLower := strings.ToLower(strings.TrimSpace(name))
	for i, theme := range themes {
		if strings.ToLower(theme.Name) == nameLower {
			return theme, i
		}
	}
	// Try partial match
	for i, theme := range themes {
		if strings.Contains(strings.ToLower(theme.Name), nameLower) {
			return theme, i
		}
	}
	return themes[0], 0
}

type Styles struct {
	App         lipgloss.Style
	Header      lipgloss.Style
	HeaderTitle lipgloss.Style
	Row         lipgloss.Style
	RowTitle    lipgloss.Style
	RowContent  lipgloss.Style
	Icon        lipgloss.Style
	Label       lipgloss.Style
	Value       lipgloss.Style
	StatusOK    lipgloss.Style
	StatusWarn  lipgloss.Style
	StatusCrit  lipgloss.Style
	Footer      lipgloss.Style
	Help        lipgloss.Style
}

func NewStyles(theme Theme, width int) Styles {
	return Styles{
		App: lipgloss.NewStyle().
			Padding(1, 2),
		Header: lipgloss.NewStyle().
			Padding(1, 0).
			MarginBottom(1).
			Width(width - 4),
		HeaderTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Text).
			Align(lipgloss.Center).
			Width(width - 8),
		Row: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border).
			Padding(1, 2).
			MarginBottom(1).
			Width(width - 6),
		RowTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary),
		RowContent: lipgloss.NewStyle().
			Foreground(theme.TextMuted),
		Icon: lipgloss.NewStyle().
			Foreground(theme.Primary),
		Label: lipgloss.NewStyle().
			Foreground(theme.TextDim),
		Value: lipgloss.NewStyle().
			Foreground(theme.Text).
			Bold(true),
		StatusOK: lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true),
		StatusWarn: lipgloss.NewStyle().
			Foreground(theme.Warning).
			Bold(true),
		StatusCrit: lipgloss.NewStyle().
			Foreground(theme.Danger).
			Bold(true),
		Footer: lipgloss.NewStyle().
			Padding(0).
			MarginTop(1).
			Width(width - 4),
		Help: lipgloss.NewStyle().
			Foreground(theme.TextDim).
			Italic(true),
	}
}

type QuotaData struct {
	Subscription struct {
		Limit    float64 `json:"limit"`
		Requests float64 `json:"requests"`
		RenewsAt string  `json:"renewsAt"`
	} `json:"subscription"`
	Search struct {
		Hourly struct {
			Limit    float64 `json:"limit"`
			Requests float64 `json:"requests"`
			RenewsAt string  `json:"renewsAt"`
		} `json:"hourly"`
	} `json:"search"`
	FreeToolCalls struct {
		Limit    float64 `json:"limit"`
		Requests float64 `json:"requests"`
		RenewsAt string  `json:"renewsAt"`
	} `json:"freeToolCalls"`
}

type model struct {
	data       QuotaData
	styles     Styles
	theme      Theme
	themeIndex int
	progress   []progress.Model
	err        error
	lastUpdate time.Time
	loading    bool
	width      int
	height     int
}

type updateMsg struct {
	data QuotaData
	err  error
}

type tickMsg time.Time

func getTerminalWidth() int {
	// Try to get width from environment or use a reasonable default
	if w := os.Getenv("COLUMNS"); w != "" {
		var width int
		fmt.Sscanf(w, "%d", &width)
		if width > 0 {
			return width
		}
	}
	return 120
}

func initialModel() model {
	theme, themeIndex := themes[0], 0
	if envTheme := os.Getenv("SYNTHETIC_THEME"); envTheme != "" {
		theme, themeIndex = findThemeByName(envTheme)
	}
	width := getTerminalWidth()
	return model{
		styles:     NewStyles(theme, width),
		theme:      theme,
		themeIndex: themeIndex,
		progress:   createProgressBars(theme),
		loading:    true,
		width:      width,
	}
}

func createProgressBars(theme Theme) []progress.Model {
	bars := make([]progress.Model, 3)
	for i := range bars {
		bars[i] = progress.New(
			progress.WithGradient(string(theme.ProgressStart), string(theme.ProgressEnd)),
		)
		bars[i].ShowPercentage = false
	}
	return bars
}

func (m *model) updateTheme() {
	m.theme = themes[m.themeIndex]
	m.styles = NewStyles(m.theme, m.width)
	m.progress = createProgressBars(m.theme)
}

func (m *model) nextTheme() {
	m.themeIndex = (m.themeIndex + 1) % len(themes)
	m.updateTheme()
}

func (m *model) prevTheme() {
	m.themeIndex = (m.themeIndex - 1 + len(themes)) % len(themes)
	m.updateTheme()
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return fetchData() },
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Every(updateInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchData() tea.Msg {
	apiKey := os.Getenv("SYNTHETIC_API_KEY")
	if apiKey == "" {
		return updateMsg{err: fmt.Errorf("SYNTHETIC_API_KEY not set")}
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return updateMsg{err: err}
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return updateMsg{err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return updateMsg{err: fmt.Errorf("API returned status %d", resp.StatusCode)}
	}

	var data QuotaData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return updateMsg{err: err}
	}

	return updateMsg{data: data}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.styles = NewStyles(m.theme, m.width)
		for i := range m.progress {
			m.progress[i].Width = m.width - 40
			if m.progress[i].Width < 20 {
				m.progress[i].Width = 20
			}
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.loading = true
			return m, func() tea.Msg { return fetchData() }
		case "t", "right", "l":
			m.nextTheme()
			return m, nil
		case "left", "h":
			m.prevTheme()
			return m, nil
		}

	case updateMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.data = msg.data
			m.lastUpdate = time.Now()
			m.err = nil
		}
		return m, nil

	case tickMsg:
		m.loading = true
		return m, tea.Batch(
			func() tea.Msg { return fetchData() },
			tickCmd(),
		)
	}

	var cmds []tea.Cmd
	for i := range m.progress {
		newModel, cmd := m.progress[i].Update(msg)
		if p, ok := newModel.(progress.Model); ok {
			m.progress[i] = p
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return m.renderError()
	}

	if m.loading && m.lastUpdate.IsZero() {
		return m.renderLoading()
	}

	return m.renderDashboard()
}

func (m model) renderError() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.RowTitle.Render(" Error "),
		m.styles.RowContent.Render(m.err.Error()),
		"",
		m.styles.Help.Render("Press 'r' to retry, 'q' to quit"),
	)
	return m.styles.App.Render(m.styles.Row.Render(content))
}

func (m model) renderLoading() string {
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		m.styles.RowTitle.Render(" Loading... "),
		m.styles.RowContent.Render("Fetching your data"),
	)
	return m.styles.App.Render(m.styles.Row.Render(content))
}

func (m model) renderDashboard() string {
	header := m.renderHeader()
	rows := m.renderRows()
	footer := m.renderFooter()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		rows,
		footer,
	)

	return m.styles.App.Render(content)
}

func (m model) renderHeader() string {
	title := m.styles.HeaderTitle.Render("SYNTHETIC - USAGE DASHBOARD")
	return m.styles.Header.Render(title)
}

func (m model) renderRows() string {
	subRow := m.renderQuotaRow(
		"▶",
		"SUBSCRIPTION",
		m.data.Subscription.Requests,
		m.data.Subscription.Limit,
		m.data.Subscription.RenewsAt,
		0,
	)

	searchRow := m.renderQuotaRow(
		"🔍",
		"SEARCH",
		m.data.Search.Hourly.Requests,
		m.data.Search.Hourly.Limit,
		m.data.Search.Hourly.RenewsAt,
		1,
	)

	toolRow := m.renderQuotaRow(
		"⚡",
		"FREE TOOL CALLS",
		m.data.FreeToolCalls.Requests,
		m.data.FreeToolCalls.Limit,
		m.data.FreeToolCalls.RenewsAt,
		2,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		subRow,
		searchRow,
		toolRow,
	)
}

func (m model) renderQuotaRow(icon, title string, used, limit float64, renewsAt string, progressIndex int) string {
	pct := used / limit
	if limit == 0 {
		pct = 0
	}

	renewTime, _ := time.Parse(time.RFC3339, renewsAt)
	timeUntil := time.Until(renewTime)
	renewStr := formatDuration(timeUntil)

	percentStr := fmt.Sprintf("%.1f%%", pct*100)

	var statusStr string
	var statusStyle lipgloss.Style
	switch {
	case pct >= 0.9:
		statusStr = "✗ CRITICAL"
		statusStyle = m.styles.StatusCrit
	case pct >= 0.7:
		statusStr = "⚠ WARNING"
		statusStyle = m.styles.StatusWarn
	default:
		statusStr = "✓ NORMAL"
		statusStyle = m.styles.StatusOK
	}

	m.progress[progressIndex].Width = m.width - 45
	if m.progress[progressIndex].Width < 20 {
		m.progress[progressIndex].Width = 20
	}
	progressBar := m.progress[progressIndex].ViewAs(pct)

	line1 := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Icon.Render(icon),
		" ",
		m.styles.RowTitle.Render(title),
		" ",
		progressBar,
		" ",
		m.styles.Value.Render(percentStr),
	)

	line2 := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("Used: "),
		m.styles.Value.Render(fmt.Sprintf("%.1f", used)),
		m.styles.Label.Render(" / "),
		m.styles.Value.Render(fmt.Sprintf("%.0f", limit)),
		"    ",
		statusStyle.Render(statusStr),
		"    ",
		m.styles.Label.Render("Resets in: "),
		m.styles.Value.Render(renewStr),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		line1,
		line2,
	)

	return m.styles.Row.Render(content)
}

func (m model) renderFooter() string {
	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Help.Render("Theme: "),
		m.styles.Help.Render(m.theme.Name),
		lipgloss.NewStyle().Width(m.width-lipgloss.Width("Theme: ")-lipgloss.Width(m.theme.Name)-lipgloss.Width(" [←/→] Theme  [r] Refresh  [q] Quit")-12).Render(""),
		m.styles.Help.Render("[←/→] Theme  [r] Refresh  [q] Quit"),
	)

	return m.styles.Footer.Render(content)
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

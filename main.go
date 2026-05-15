package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

type LogEntry struct {
	IP         string
	Time       time.Time
	Method     string
	Path       string
	Protocol   string
	Status     int
	BodyBytes  int
	Referer    string
	UserAgent  string
}

type KeyCount struct {
	Key   string
	Count int
}

func sortMapByValueDesc(m map[string]int) []KeyCount {
	var list []KeyCount
	for k, v := range m {
		list = append(list, KeyCount{k, v})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Count > list[j].Count
	})
	return list
}

func topN(list []KeyCount, n int) []KeyCount {
	if len(list) < n {
		return list
	}
	return list[:n]
}

var logRegex = regexp.MustCompile(`^(\S+) - - \[([^\]]+)\] "(\S+) (\S+) (\S+)" (\d+) (\d+) "([^"]+)" "([^"]+)"$`)

func main() {
	rootCmd := &cobra.Command{
		Use:   "log-analyzer",
		Short: "Nginx log analyzer",
		Long:  "A CLI tool for analyzing Nginx access logs",
	}

	rootCmd.AddCommand(newGenSampleCmd())
	rootCmd.AddCommand(newStatsCmd())
	rootCmd.AddCommand(newFilterCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats <log-file>",
		Short: "Show log statistics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStats(args[0])
		},
	}

	return cmd
}

type FilterOptions struct {
	Status   string
	IP       string
	Path     string
	Method   string
	Since    string
	Until    string
}

func matchesStatus(entry LogEntry, statusFilter string) bool {
	if statusFilter == "" {
		return true
	}
	if len(statusFilter) == 3 && statusFilter[1] == 'x' && statusFilter[2] == 'x' {
		switch statusFilter[0] {
		case '2':
			return entry.Status >= 200 && entry.Status < 300
		case '3':
			return entry.Status >= 300 && entry.Status < 400
		case '4':
			return entry.Status >= 400 && entry.Status < 500
		case '5':
			return entry.Status >= 500 && entry.Status < 600
		}
	}
	s, err := strconv.Atoi(statusFilter)
	if err != nil {
		return false
	}
	return entry.Status == s
}

func matchesIP(entry LogEntry, ipFilter string) bool {
	if ipFilter == "" {
		return true
	}
	return entry.IP == ipFilter
}

func matchesPath(entry LogEntry, pathFilter string) bool {
	if pathFilter == "" {
		return true
	}
	return entry.Path == pathFilter
}

func matchesMethod(entry LogEntry, methodFilter string) bool {
	if methodFilter == "" {
		return true
	}
	return entry.Method == methodFilter
}

func newFilterCmd() *cobra.Command {
	var opts FilterOptions

	cmd := &cobra.Command{
		Use:   "filter <log-file>",
		Short: "Filter log entries",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFilter(args[0], opts)
		},
	}

	cmd.Flags().StringVar(&opts.Status, "status", "", "Filter by status code (e.g. 404 or 4xx)")
	cmd.Flags().StringVar(&opts.IP, "ip", "", "Filter by client IP")
	cmd.Flags().StringVar(&opts.Path, "path", "", "Filter by request path (exact match)")
	cmd.Flags().StringVar(&opts.Method, "method", "", "Filter by HTTP method (GET, POST, etc.)")

	return cmd
}

func runFilter(filename string, opts FilterOptions) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		entry, err := parseLogLine(line)
		if err != nil {
			continue
		}
		if !matchesStatus(entry, opts.Status) {
			continue
		}
		if !matchesIP(entry, opts.IP) {
			continue
		}
		if !matchesPath(entry, opts.Path) {
			continue
		}
		if !matchesMethod(entry, opts.Method) {
			continue
		}
		fmt.Println(line)
	}

	return scanner.Err()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func runStats(filename string) error {
	entries, err := parseLogFile(filename)
	if err != nil {
		return err
	}

	total := len(entries)
	statusCounts := make(map[int]int)
	categoryCounts := make(map[string]int)
	ipCounts := make(map[string]int)
	pathCounts := make(map[string]int)
	uaCounts := make(map[string]int)

	for _, e := range entries {
		statusCounts[e.Status]++
		ipCounts[e.IP]++
		pathCounts[e.Path]++
		uaCounts[e.UserAgent]++
		switch {
		case e.Status >= 200 && e.Status < 300:
			categoryCounts["2xx"]++
		case e.Status >= 300 && e.Status < 400:
			categoryCounts["3xx"]++
		case e.Status >= 400 && e.Status < 500:
			categoryCounts["4xx"]++
		case e.Status >= 500 && e.Status < 600:
			categoryCounts["5xx"]++
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "METRIC\tVALUE")
	fmt.Fprintf(w, "Total Requests\t%d\n", total)
	w.Flush()

	fmt.Println()
	fmt.Println("Status Code Distribution:")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CATEGORY\tCOUNT\tPERCENTAGE")
	for _, cat := range []string{"2xx", "3xx", "4xx", "5xx"} {
		count := categoryCounts[cat]
		pct := float64(count) / float64(total) * 100
		fmt.Fprintf(w, "%s\t%d\t%.1f%%\n", cat, count, pct)
	}
	w.Flush()

	fmt.Println()
	fmt.Println("Top 10 IPs:")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "IP\tREQUESTS")
	topIPs := topN(sortMapByValueDesc(ipCounts), 10)
	for _, kc := range topIPs {
		fmt.Fprintf(w, "%s\t%d\n", kc.Key, kc.Count)
	}
	w.Flush()

	fmt.Println()
	fmt.Println("Top 10 URLs:")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "URL\tREQUESTS")
	topPaths := topN(sortMapByValueDesc(pathCounts), 10)
	for _, kc := range topPaths {
		fmt.Fprintf(w, "%s\t%d\n", kc.Key, kc.Count)
	}
	w.Flush()

	fmt.Println()
	fmt.Println("Top 10 User Agents:")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "USER AGENT\tREQUESTS")
	topUAs := topN(sortMapByValueDesc(uaCounts), 10)
	for _, kc := range topUAs {
		fmt.Fprintf(w, "%s\t%d\n", truncate(kc.Key, 60), kc.Count)
	}
	w.Flush()

	return nil
}

func parseLogFile(filename string) ([]LogEntry, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		entry, err := parseLogLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: line %d: %v\n", lineNum, err)
			continue
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func parseLogLine(line string) (LogEntry, error) {
	matches := logRegex.FindStringSubmatch(line)
	if len(matches) != 10 {
		return LogEntry{}, fmt.Errorf("unable to parse line")
	}

	status, err := strconv.Atoi(matches[6])
	if err != nil {
		return LogEntry{}, fmt.Errorf("invalid status: %v", err)
	}

	bodyBytes, err := strconv.Atoi(matches[7])
	if err != nil {
		return LogEntry{}, fmt.Errorf("invalid body bytes: %v", err)
	}

	t, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		return LogEntry{}, fmt.Errorf("invalid time: %v", err)
	}

	return LogEntry{
		IP:         matches[1],
		Time:       t,
		Method:     matches[3],
		Path:       matches[4],
		Protocol:   matches[5],
		Status:     status,
		BodyBytes:  bodyBytes,
		Referer:    matches[8],
		UserAgent:  matches[9],
	}, nil
}

func newGenSampleCmd() *cobra.Command {
	var lines int
	var output string

	cmd := &cobra.Command{
		Use:   "gen-sample",
		Short: "Generate sample Nginx access logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateSample(lines, output)
		},
	}

	cmd.Flags().IntVarP(&lines, "lines", "n", 100, "Number of lines to generate")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file (default: stdout)")

	return cmd
}

func generateSample(lines int, output string) error {
	var w io.Writer = os.Stdout
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < lines; i++ {
		line := generateLogLine(r)
		fmt.Fprintln(w, line)
	}

	return nil
}

func generateLogLine(r *rand.Rand) string {
	ip := generateIP(r)
	ts := time.Now().Add(-time.Duration(r.Intn(86400)) * time.Second)
	method := pickMethod(r)
	path := pickPath(r)
	proto := "HTTP/1.1"
	status := pickStatus(r)
	bodyBytes := r.Intn(10000)
	referer := pickReferer(r)
	ua := pickUserAgent(r)

	return fmt.Sprintf(`%s - - [%s] "%s %s %s" %d %d "%s" "%s"`,
		ip,
		ts.Format("02/Jan/2006:15:04:05 -0700"),
		method,
		path,
		proto,
		status,
		bodyBytes,
		referer,
		ua,
	)
}

func generateIP(r *rand.Rand) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		r.Intn(256),
		r.Intn(256),
		r.Intn(256),
		r.Intn(256),
	)
}

func pickMethod(r *rand.Rand) string {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	return methods[r.Intn(len(methods))]
}

func pickPath(r *rand.Rand) string {
	paths := []string{
		"/",
		"/api/users",
		"/api/posts",
		"/static/css/style.css",
		"/static/js/app.js",
		"/images/logo.png",
		"/about",
		"/contact",
		"/blog/post-1",
		"/blog/post-2",
		"/api/products",
		"/cart",
		"/checkout",
		"/login",
		"/signup",
	}
	return paths[r.Intn(len(paths))]
}

func pickStatus(r *rand.Rand) int {
	statuses := []int{200, 200, 200, 200, 200, 201, 301, 302, 400, 401, 403, 404, 500}
	return statuses[r.Intn(len(statuses))]
}

func pickReferer(r *rand.Rand) string {
	if r.Float32() < 0.3 {
		return "-"
	}
	referers := []string{
		"https://google.com",
		"https://bing.com",
		"https://example.com",
		"https://example.com/blog",
		"https://example.com/about",
	}
	return referers[r.Intn(len(referers))]
}

func pickUserAgent(r *rand.Rand) string {
	uas := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"curl/7.68.0",
		"Go-http-client/1.1",
	}
	return uas[r.Intn(len(uas))]
}

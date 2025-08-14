package tests

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/cover"
	"html/template"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type FunctionCoverage struct {
	FunctionName string  `json:"function_name"`
	TotalLines   int64   `json:"total_lines"`
	CoveredLines int64   `json:"covered_lines"`
	Coverage     float64 `json:"coverage"`
}

type Row struct {
	FileName         string             `json:"file_name"`
	FunctionCoverage []FunctionCoverage `json:"function_name"`
	TotalLines       int64              `json:"total_lines"`
	CoveredLines     int64              `json:"covered_lines"`
	Coverage         float64            `json:"coverage"`
}

type CovReport struct {
	Module           string        `json:"module"`
	RepoTotalLines   int64         `json:"repo_total_lines"`
	RepoCoveredLines int64         `json:"repo_covered_lines"`
	OverallCoverage  float64       `json:"coverage"`
	Rows             []Row         `json:"rows"`
	CodeBranch       string        `json:"code_branch"`
	TotalTestTime    time.Duration `json:"total_test_time"`
}

// FuncExtent Taken from Official go cover package
// FuncExtent describes a function's extent in the source by file and position.
type FuncExtent struct {
	name      string
	startLine int
	startCol  int
	endLine   int
	endCol    int
}

// FuncVisitor Taken from Official go cover package
// FuncVisitor implements the visitor that builds the function position list for a file.
type FuncVisitor struct {
	fset    *token.FileSet
	name    string // Name of file.
	astFile *ast.File
	funcs   []*FuncExtent
}

func Start() {
	packages := []string{
		"./services/...",
				"./handlers/...",
						"./app/...",


	}
	// Test Start Time
	startTimeEpoch := int(time.Now().UnixNano() / 1e6)

	// Get Current Branch Name
	branchName := getCurrentBranchName()

	// Run tests and store the exit code
	exitCode := runTests(packages)

	// Test End Time
	endTimeEpoch := int(time.Now().UnixNano() / 1e6)

	// Generate the HTML coverage report
	err := generateHTMLCoverageReport("coverage.out", "coverage.html")
	if err != nil {
		log.Fatalf("Failed to generate HTML coverage report: %v", err)
	}

	startTime, endTime := time.Unix(int64(startTimeEpoch)/1000, 0), time.Unix(int64(endTimeEpoch)/1000, 0)
	fmt.Println(startTime, endTime)
	diff := endTime.Sub(startTime)
	report, err := generateDetailedReport("coverage.out", branchName, diff)
	if err != nil {
		return
	}
	fmt.Printf("Overall Coverage: %.2f%%\n", report.OverallCoverage)

	// Exit with the same exit code as the tests
	os.Exit(exitCode)
}

func runTests(packagesList []string) int {
	packages := strings.Join(packagesList, ",")
	cmd := exec.Command("go", "test", "-v", "-coverprofile=coverage.out", "-coverpkg="+packages, "./...", "-covermode=count")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Error running tests for package : %v", err)
		return 1
	}

	return 0
}

func generateHTMLCoverageReport(coverageFile, outputHTMLFile string) error {
	// Run "go tool cover" command to generate HTML coverage report
	cmd := exec.Command("go", "tool", "cover", "-html="+coverageFile, "-o="+outputHTMLFile)
	err := cmd.Run()
	if err != nil {
		log.Printf("Error generating HTML coverage report: %v\n", err)
		return err
	}

	return nil
}

func generateDetailedReport(file, branchName string, totalTime time.Duration) (CovReport, error) {
	rows := make([]Row, 0)

	profiles, err := cover.ParseProfiles(file)
	if err != nil {
		return CovReport{}, err
	}

	var repoTotalLines, repoCoveredLines int64

	for _, profile := range profiles {
		functionCoverages := make([]FunctionCoverage, 0)
		var total, covered int64
		fileName := strings.Split(profile.FileName, "aitring/")
		funcs, err := findFuncs(strings.TrimSpace(fileName[1]))
		if err != nil {
			return CovReport{}, err
		}
		//fmt.Println(funcs)
		for _, f := range funcs {
			c, t := f.coverage(profile)
			total += t
			covered += c

			functionCoverages = append(functionCoverages, FunctionCoverage{
				FunctionName: f.name,
				TotalLines:   t,
				CoveredLines: c,
				Coverage:     percent(c, t),
			})
		}

		sort.Slice(functionCoverages, func(i, j int) bool {
			return functionCoverages[i].Coverage > functionCoverages[j].Coverage
		})

		rows = append(rows, Row{
			FileName: profile.FileName, TotalLines: total, CoveredLines: covered, FunctionCoverage: functionCoverages,
			Coverage: percent(covered, total),
		})

		repoTotalLines += total
		repoCoveredLines += covered
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Coverage > rows[j].Coverage
	})
	covReport := CovReport{Module: "AI-TRING", RepoTotalLines: repoTotalLines,
		RepoCoveredLines: repoCoveredLines, OverallCoverage: percent(repoCoveredLines, repoTotalLines), Rows: rows,
		CodeBranch: branchName, TotalTestTime: totalTime,
	}

	html := `<!DOCTYPE html>
				<html>
				<head>
				  <title>{{.Module}} Coverage</title>
				  <style>
 						table td,
						table th {
						  text-overflow: ellipsis;
						  white-space: nowrap;
						  overflow: hidden;
						}
						.table-scroll {
						  border-radius: .5rem;
						}
						thead th {
						  color: #fff;
						}
						.table-scroll table thead th {
						  font-size: 1.25rem;
						}
						.font-12 {
							font-size: 12px;

						}
						body {
							font-family: Arial, sans-serif;
							background-color: #f1f1f1;
							padding: 20px;
						}
						td {
							cursor: pointer;
						}
						th, td {
							padding: 8px;
							text-align: left;
						}
						.table-arrow .arrow::before {
						  content: "\25B6";
						  font-weight: bold;
						  transition: transform 0.2s ease-in-out;
						}
					
						.table-arrow.collapsed .arrow::before {
						  transform: rotate(90deg);
						}
						.rows {
							background-color: antiquewhite
						}
						
						.sub-rows {
							background-color: floralwhite
						}
					</style>
					
				  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css">
   				  <script src="https://code.jquery.com/jquery-3.5.1.slim.min.js"></script>
  				  <script src="https://cdn.jsdelivr.net/npm/bootstrap@4.5.0/dist/js/bootstrap.min.js"></script>
	
				</head>
				<body>
					<h4 class="text-center">AI TRING Portal Unit Testing Coverage Report</h4>
					<br>
					<table class="table-bordered border-success" style="width:400px;margin-left:auto;margin-right:auto">
						<tr>
							<td style="background-color:seagreen;color:white;font-family:monospace">Total Lines</td>
							<td style="font-family:monospace">{{.RepoTotalLines}}</td>
						</tr>
						<tr>
							<td style="background-color:seagreen;color:white;font-family:monospace">Covered Lines</td>
							<td style="font-family:monospace">{{.RepoCoveredLines}}</td>
						</tr>
						<tr>
							<td style="background-color:seagreen;color:white;font-family:monospace">OverAll Coverage</td>
							<td style="font-family:monospace">{{.OverallCoverage}}%</td>
						</tr>
						<tr>
							<td style="background-color:seagreen;color:white;font-family:monospace">Branch</td>
							<td style="font-family:monospace">{{.CodeBranch}}</td>
						</tr>
						<tr>
							<td style="background-color:seagreen;color:white;font-family:monospace">Test Duration</td>
							<td style="font-family:monospace">{{.TotalTestTime}}</td>
						</tr>
					</table>
					
					<br>
					<div class="container table-responsive table-scroll" style="min-width:1250px;height:1000px">
						<table class="table table-hover table-bordered table-light border-success" style="font-family:monospace;font-size:14px">
							<thead style="background-color: #002d72;">
								<tr>
									<th scope="col">File</th>
									<th scope="col">Total Lines</th>
									<th scope="col">Covered Lines</th>
									<th scope="col">Coverage</th>
									<th scope="col">Details</th>
								</tr>
							</thead>
							<tbody>
								{{range $i, $row := .Rows}}
									<tr class="rows" data-toggle="collapse" data-target=".file{{$i}}-details" aria-expanded="false" aria-controls="file{{$i}}-details">
										<td>{{$row.FileName}}</td>
										<td>{{$row.TotalLines}}</td>
										<td>{{$row.CoveredLines}}</td>
										<td>{{$row.Coverage}}%</td>
										<td class="table-arrow"><span class="arrow"></span></td>
									</tr>
									{{range $j, $coverage := $row.FunctionCoverage}}
										<tr class="file{{$i}}-details collapse font-12 sub-rows">
											<td>{{$coverage.FunctionName}}</td>
											<td>{{$coverage.TotalLines}}</td>
											<td>{{$coverage.CoveredLines}}</td>
											<td colspan="2">{{$coverage.Coverage}}%</td>
										</tr>
									{{end}}
								{{end}}
							</tbody>
					   </table>
					</div>
				</body>
			</html>`

	// Create a new template from the HTML string
	t := template.Must(template.New("output").Parse(html))

	// Create the output file
	out, err := os.Create("AITringTestReport.html")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Execute the template and write to the file
	err = t.Execute(out, covReport)
	if err != nil {
		panic(err)
	}

	return covReport, nil
}

// Taken from Official go cover package
func findFuncs(name string) ([]*FuncExtent, error) {
	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		return nil, err
	}
	visitor := &FuncVisitor{
		fset:    fset,
		name:    name,
		astFile: parsedFile,
	}
	ast.Walk(visitor, visitor.astFile)
	return visitor.funcs, nil
}

func getCurrentBranchName() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	branch := strings.TrimSpace(string(output))
	return branch
}

// Visit Taken from Official go cover package
func (v *FuncVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Body == nil {
			// Do not count declarations of assembly functions.
			break
		}
		start := v.fset.Position(n.Pos())
		end := v.fset.Position(n.End())
		fe := &FuncExtent{
			name:      n.Name.Name,
			startLine: start.Line,
			startCol:  start.Column,
			endLine:   end.Line,
			endCol:    end.Column,
		}
		v.funcs = append(v.funcs, fe)
	}
	return v
}

func percent(covered, total int64) float64 {
	if total == 0 {
		total = 1 // Avoid zero denominator.
	}
	val := 100.0 * float64(covered) / float64(total)
	return float64(int(val*100)) / 100
}

// Taken from Official go cover package
// coverage returns the fraction of the statements in the function that were covered, as a numerator and denominator.
func (f *FuncExtent) coverage(profile *cover.Profile) (num, den int64) {
	// We could avoid making this n^2 overall by doing a single scan and annotating the functions,
	// but the sizes of the data structures is never very large and the scan is almost instantaneous.
	var covered, total int64
	// The blocks are sorted, so we can stop counting as soon as we reach the end of the relevant block.
	for _, b := range profile.Blocks {
		if b.StartLine > f.endLine || (b.StartLine == f.endLine && b.StartCol >= f.endCol) {
			// Past the end of the function.
			break
		}
		if b.EndLine < f.startLine || (b.EndLine == f.startLine && b.EndCol <= f.startCol) {
			// Before the beginning of the function
			continue
		}
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	return covered, total
}

package services

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

type Point struct {
	JudgeNumber  string
	TableNumber  string
	Points       float64
	Standardized float64
}

type Team struct {
	TableNumber string
	TeamName    string
}

type Aggregated struct {
	TableNumber        string
	StandardizedPoints float64
	JudgeCount         int
	TeamName           string
}

func aggregateAndSortPoints(pointsFile, teamsFile, outputFile string) error {
	// ---- Load points.csv ----
	pointsData, err := readCSV(pointsFile)
	if err != nil {
		return err
	}

	var points []Point
	for _, row := range pointsData[1:] { // skip header
		pVal, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}
		points = append(points, Point{
			JudgeNumber: row[0],
			TableNumber: row[1],
			Points:      pVal,
		})
	}

	// ---- Standardize points per judge ----
	judgeGroups := make(map[string][]float64)
	for _, p := range points {
		judgeGroups[p.JudgeNumber] = append(judgeGroups[p.JudgeNumber], p.Points)
	}

	judgeStats := make(map[string][2]float64) // mean, std
	for judge, vals := range judgeGroups {
		mean := 0.0
		for _, v := range vals {
			mean += v
		}
		mean /= float64(len(vals))

		std := 0.0
		for _, v := range vals {
			std += (v - mean) * (v - mean)
		}
		std = math.Sqrt(std / float64(len(vals)))
		judgeStats[judge] = [2]float64{mean, std}
	}

	for i, p := range points {
		stats := judgeStats[p.JudgeNumber]
		if stats[1] != 0 {
			points[i].Standardized = (p.Points - stats[0]) / stats[1]
		}
	}

	// ---- Group by tableNumber ----
	tableGroups := make(map[string][]Point)
	for _, p := range points {
		tableGroups[p.TableNumber] = append(tableGroups[p.TableNumber], p)
	}

	var agg []Aggregated
	for table, group := range tableGroups {
		sum := 0.0
		judges := make(map[string]bool)
		for _, p := range group {
			sum += p.Standardized
			judges[p.JudgeNumber] = true
		}
		mean := sum / float64(len(group))
		agg = append(agg, Aggregated{
			TableNumber:        table,
			StandardizedPoints: mean,
			JudgeCount:         len(judges),
		})
	}

	// ---- Load teams.csv ----
	teamsData, err := readCSV(teamsFile)
	if err != nil {
		return err
	}
	teamMap := make(map[string]string)
	for _, row := range teamsData[1:] {
		if len(row) < 2 {
			continue
		}
		teamMap[row[0]] = row[1]
	}

	// ---- Merge with team names ----
	for i := range agg {
		agg[i].TeamName = teamMap[agg[i].TableNumber]
	}

	// ---- Sort by standardized points descending ----
	sort.Slice(agg, func(i, j int) bool {
		return agg[i].StandardizedPoints > agg[j].StandardizedPoints
	})

	// ---- Write to CSV ----
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	writer.Write([]string{"tableNumber", "standardizedPoints", "judgeCount", "teamName"})
	for _, a := range agg {
		writer.Write([]string{
			a.TableNumber,
			fmt.Sprintf("%.4f", a.StandardizedPoints),
			strconv.Itoa(a.JudgeCount),
			a.TeamName,
		})
	}

	return nil
}

// Utility: read CSV file
func readCSV(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.ReuseRecord = true
	data, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func main() {
	err := aggregateAndSortPoints("points.csv", "teams.csv", "aggregated_and_sorted_points.csv")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println("✅ Aggregation complete. Output saved to aggregated_and_sorted_points.csv")
}

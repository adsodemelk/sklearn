package cluster

import (
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"time"

	"github.com/pa-m/sklearn/datasets"
	"github.com/pa-m/sklearn/preprocessing"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func ExampleDBSCAN() {
	// adapted from http://scikit-learn.org/stable/_downloads/plot_dbscan.ipynb
	// Generate sample data
	centers := mat.NewDense(3, 2, []float64{1, 1, -1, -1, 1, -1})
	NSamples := 750
	X, _ := datasets.MakeBlobs(&datasets.MakeBlobsConfig{NSamples: NSamples, Centers: centers, ClusterStd: .4}) //RandomState: rand.New(rand.NewSource(0)),

	X, _ = preprocessing.NewStandardScaler().FitTransform(X, nil)
	db := NewDBSCAN(&DBSCANConfig{Eps: .3, MinSamples: 10, Algorithm: "kd_tree"})
	db.Fit(X, nil)
	coreSampleMask := make([]bool, len(db.Labels))
	for sample := range db.CoreSampleIndices {
		coreSampleMask[sample] = true
	}
	labels := db.Labels
	labelsmap := make(map[int]int)
	for _, l := range labels {
		labelsmap[l] = l
	}
	nclusters := len(labelsmap)
	if _, ok := labelsmap[-1]; ok {
		nclusters--
	}
	fmt.Printf("Estimated number of clusters: %d\n", nclusters)

	canPlot := false
	if canPlot {

		// plot result
		p, err := plot.New()
		if err != nil {
			panic(err)
		}
		p.Title.Text = fmt.Sprintf("Estimated number of clusters: %d", nclusters)
		for cl := range labelsmap {
			var data plotter.XYs
			for sample := 0; sample < NSamples; sample++ {
				if labels[sample] == cl {
					data = append(data, struct{ X, Y float64 }{X.At(sample, 0), X.At(sample, 1)})
				}
			}
			s, err := plotter.NewScatter(data)
			if err != nil {
				panic(err)
			}
			var color0 color.RGBA
			switch cl {
			case -1:
				color0 = color.RGBA{0, 0, 0, 255}
			case 0:
				color0 = color.RGBA{176, 0, 0, 255}
			case 1:
				color0 = color.RGBA{0, 176, 0, 255}
			case 2:
				color0 = color.RGBA{0, 0, 176, 255}
			}
			s.GlyphStyle.Color = color0
			p.Add(s)
			//p.Legend.Add(fmt.Sprintf("scatter %d", cl), s)

		}
		// Save the plot to a PNG file.
		pngfile := "/tmp/plt.png"
		os.Remove(pngfile)
		if err := p.Save(4*vg.Inch, 4*vg.Inch, pngfile); err != nil {
			panic(err)
		}
		cmd := exec.Command("display", pngfile)
		err = cmd.Start()
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(200 * time.Millisecond)
		os.Remove(pngfile)
	}
	// Output:
	// Estimated number of clusters: 3
}

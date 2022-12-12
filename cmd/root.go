package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"xml2csv-parser/cfg"
	. "xml2csv-parser/internal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	path       string
	columnSet  string
	outFile    string
	nConsumers int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xml2cv-parser",
	Short: "xml2cv-parser creates list of orders in csv file from set of inbound Nestle xml files.",
	Long: `  The tool parses any set of fields from a list of xml files via xpath expressions in a yaml configuration file,
which is passed by -s <path/file_name> flag.
  By default the files are filtered for missing data, but you can configure mandatory output in csv result even if the data is missing.

  xml2cv-parser saves found data in csv, utf-8, comma delimiter.
  The xpath expression must return the result as the single value of a specific field. Array results and Node results are not supported.

  Example of yaml config for parsing:
includeFilename: true # filename will be added in last field
set:
- messageType: some_message_type
  columns:
  - name: some_column_name
    xpath: //SomeRoot/SomeNode/SomeElement
	optional: true # file will not be dropped if data is missing - csv field will be blank
  - name: another_column_name
    xpath: //SomeRoot/AnotherNode/AnotherElement[/@SomeAttribute='SomeValue']
- messageType: another_message_type
  columns:
  - name: some_column_name
    xpath: //AnotherRoot/SomeNode/SomeElement
  - name: another_column_name
    xpath: //AnotherRoot/AnotherNode/AnotherElement[/@SomeAttribute='SomeValue']`,
	Run: func(cmd *cobra.Command, args []string) {
		var setCfg cfg.XMLParser
		err := setCfg.Load(columnSet)
		if err != nil {
			log.Fatalf("can't load xpath config %s: %s\n", columnSet, err)
		}

		ctx, cancelFunc := context.WithCancel(context.Background())
		quit := make(chan interface{})

		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			close(quit)
			cancelFunc()
		}()

		t := time.Now()
		println("Started in path", path, "at", t.String())

		runtime.GOMAXPROCS(runtime.NumCPU())
		files := make(chan string, 4)

		p := NewProducer(path, ".xml", &files, quit)
		writer, err := NewCsvWriter(outFile)
		if err != nil {
			log.Fatalln(err)
		}
		c := NewConsumer(setCfg.CreateCompiled(), writer, &files, make(chan string, nConsumers))

		go p.Produce()
		go c.Consume(ctx)

		wg := &sync.WaitGroup{}
		wg.Add(nConsumers)
		for i := 0; i < nConsumers; i++ {
			go c.Work(wg)
		}

		wg.Wait()

		writer.Flush()
		writer.Close()

		println("Spent", time.Since(t).String())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't get current dir: %s", err)
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.xml2cv-parser.yaml)")
	rootCmd.Flags().StringVarP(&path, "path", "p", wd, "path for Nestle inbound xml's. Only files with 'xml' extensions will be processed")
	rootCmd.Flags().StringVarP(&columnSet, "column_set", "s", "", "config file for xml parsing in .yaml format")
	rootCmd.Flags().StringVarP(&outFile, "out_file", "o", "result.csv", "out filename for csv to where list of orders will be stored")
	rootCmd.Flags().IntVarP(&nConsumers, "threads", "t", 10, "number of consumers for concurrency executions")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".xml2cv-parser" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".xml2cv-parser")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

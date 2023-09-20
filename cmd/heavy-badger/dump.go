package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/insolar/insolar/insolar/record"
)

func (app *appCtx) valueHexDumpCommand() *cobra.Command {
	var dumpCmd = &cobra.Command{
		Use: "dump",
	}

	var key []byte
	keyArgCheck := func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a argument with hex key value")
		}
		var err error
		key, err = hex.DecodeString(args[0])
		if err != nil {
			return fmt.Errorf("provided argument should be valid hex string: %v", err)
		}
		return nil
	}

	var dumpBinaryCmd = &cobra.Command{
		Use:   "bin",
		Short: "dump binary value by key",
		Args:  keyArgCheck,
		Run: func(_ *cobra.Command, _ []string) {
			db, close := openDB(app.dataDir)
			defer close()
			value, err := readValueByKey(db, key)
			if err != nil {
				fatalf("failed to get key from badger: %v", err)
			}
			_, err = io.Copy(os.Stdout, bytes.NewReader(value))
			if err != nil {
				fatalf("failed copy to stdin: %v", err)
			}
		},
	}

	var dumpRecordCmd = &cobra.Command{
		Use:   "record",
		Short: "dump record protobuf by key",
		Args:  keyArgCheck,
		Run: func(_ *cobra.Command, _ []string) {
			db, close := openDB(app.dataDir)
			defer close()
			value, err := readValueByKey(db, key)
			if err != nil {
				fatalf("failed to get key from badger: %v", err)
			}
			var material record.Material
			err = material.Unmarshal(value)
			if err != nil {
				fatalf("failed to decode value to record.Material: %v", err)
			}

			printLine("-")
			fmt.Println("Material Record:")
			fmt.Println(pairFormatter{width: 20}.Pairs(
				"ID", material.ID.String(),
				"JetID", material.JetID.DebugString(),
				"ObjectID", material.ObjectID.String(),
			))

			fmt.Println("\nVirtual Record:")
			fmt.Println(prettyPrintVirtual(&material.Virtual))
			printLine("-")
		},
	}

	dumpCmd.AddCommand(
		dumpBinaryCmd,
		dumpRecordCmd,
	)

	return dumpCmd
}

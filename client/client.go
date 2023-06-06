package main

import (
	pb "Proj3/token/token"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

// Global definition of variables for tokens.

type gbtoken struct {
	Id           string   `yaml:"id"`
	Token_reader []string `yaml:"reader"`
	Token_writer []string `yaml:"writer"`
}

var stru = make(map[string]gbtoken)

func main() {

	filename, _ := filepath.Abs("config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	error1 := yaml.Unmarshal(yamlFile, &stru)
	if error1 != nil {
		panic(error1)
	}

	// Parse the flags according to the parameters set.
	var port string
	id := flag.String("id", "0", "Token ID")
	name := flag.String("name", "", "Value of Name")
	low := flag.Uint64("low", 0, "Low Value")
	mid := flag.Uint64("mid", 0, "Mid Value")
	high := flag.Uint64("high", 0, "High Value")
	read1 := flag.Bool("read", false, "Read Function")
	write1 := flag.Bool("write", false, "Write Function")
	flag.Parse()

	// Give the token as an input from cmd.
	// Then iterate through writer's address to get the writer's values.
	if *write1 {
		num_writers := len(stru[*id].Token_writer)
		for j := 0; j < num_writers; j++ {
			port = stru[*id].Token_writer[j]

			conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock()) // Port connection establishment.
			if err != nil {
				log.Fatalf("Connection Failed : %v", err)
			}
			defer conn.Close()
			c := pb.NewTokenServiceClient(conn)

			cont, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			// Verify for existing token IDs before new token creation or updating the existing tokens.
			_, inmap := stru[*id]
			if !inmap {
				var idvalue = *id
				r, err := c.Token_Create(cont, &pb.Create_ID{Cid: idvalue})
				if err != nil {
					log.Fatalf("New Token Creation not successful :%v", err)
				}
				fmt.Println(r.GetCinfo())

				// Write Token or Update existing Token.
				v, err := c.Token_Write(cont, &pb.Write_ID{Wid: *id, Wname: *name, Wlow: *low, Wmid: *mid, Whigh: *high})
				if err != nil {
					log.Fatalf("Write Token not successful  :%v", err)
				}
				log.Println(v.GetPvalue())
				log.Println(v.GetPinfo())
			} else {
				v, err := c.Token_Write(cont, &pb.Write_ID{Wid: *id, Wname: *name, Wlow: *low, Wmid: *mid, Whigh: *high})
				if err != nil {
					log.Fatalf("Write Token not successful :%v", err)
				}
				log.Println(v.GetPvalue())
				log.Println(v.GetPinfo())
			}
		}
	}
	var (
		pid   string
		pname string
		plow  uint64
		phigh uint64
		pmid  uint64
		pts   int64
	)
	var serverPt string
	var val_f uint64
	var info_f string
	var curentts int64
	var recentts int64

	// Read the token given as an input from cmd. Then iterate through the given port address of the reader and get values for the reader.
	if *read1 {
		num_readers := len(stru[*id].Token_reader)
		for i := 0; i < num_readers; i++ {
			port = stru[*id].Token_reader[i]
			conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("Error in connection : %v", err)
			}
			defer conn.Close()
			c := pb.NewTokenServiceClient(conn)

			cont, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			// Establish a new connection after closing the previous one and read function can now be activated.
			new, err := c.Token_Read(cont, &pb.Read_ID{Readid: *id})
			if err != nil {
				log.Fatalf("Token was not created :%v", err)
			}
			// Getting the most recent values as per the system's timestamp.
			if i == 0 {
				curentts = new.GetRts()
				recentts = curentts
				val_f = new.GetFinal()
				info_f = new.GetInfo()
				pid = new.GetRid()
				pname = new.GetRname()
				plow = new.GetLow()
				phigh = new.GetRhigh()
				pmid = new.GetRmid()
				pts = new.GetRts()
				serverPt = port
			} else {
				curentts = new.GetRts()
				if curentts >= recentts {
					recentts = curentts
					val_f = new.GetFinal()
					info_f = new.GetInfo()
					pid = new.GetRid()
					pname = new.GetRname()
					plow = new.GetLow()
					phigh = new.GetRhigh()
					pmid = new.GetRmid()
					pts = new.GetRts()
					serverPt = port
				}
			}

		}
		// Latest values of the expected content from the logs is printed here.
		log.Println(val_f)
		log.Println(info_f)
		log.Println(serverPt)
		log.Println(recentts)
		log.Println(pid)
		log.Println(pname)
		log.Println(plow)
		log.Println(pmid)
		log.Println(phigh)
		log.Println(pts)

		// The readers are kept up to date as per the latest values.
		for i := 0; i < num_readers; i++ {
			port = stru[*id].Token_reader[i]
			conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock()) // Establishing port connection.
			if err != nil {
				log.Fatalf("Connection not successful : %v", err)
			}
			defer conn.Close()
			c := pb.NewTokenServiceClient(conn)
			cont, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			//The current read member's values are now updated.
			v, err := c.Token_Write(cont, &pb.Write_ID{Wid: pid, Wname: pname, Wlow: plow, Wmid: pmid, Whigh: phigh})
			if err != nil {
				log.Fatalf("Write token not successful :%v", err)
			}
			log.Println(v.GetPvalue())
			log.Println(v.GetPinfo())
		}
	}

}

package main

import (
	pb "Proj3/token/token"

	"context"
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

//hash function given in the project requirement document

func Hash(name string, nonce uint64) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s %d", name, nonce)))
	return binary.BigEndian.Uint64(hasher.Sum(nil))
}

type tokenserver struct {
	pb.UnimplementedTokenServiceServer
}

type values struct {
	final_value   uint64
	partial_value uint64
}

type levels struct {
	low  uint64
	mid  uint64
	high uint64
}

type token_structure struct {
	id    string
	name  string
	times int64
	level levels
	res   values
}

var (
	token_1 = make(map[string]token_structure)
	mutex_1 sync.Mutex
)

//In case the token which is searched is not available, we need to create a new  token. Then, the token must be added to the existing records or database.

func (s *tokenserver) CreateToken(ctx context.Context, in *pb.Create_ID) (*pb.Token_Info, error) {
	log.Printf("Token received is : %v", in.GetCid())
	mutex_1.Lock()
	token_1[in.GetCid()] = token_structure{id: in.GetCid(), res: values{partial_value: 0, final_value: 0}}
	fmt.Printf("The current token has the following data -> id: %v, name: %v, low : %d, mid: %d, high: %d, partial_value: %d, final_value: %d", token_1[in.GetCid()].id, token_1[in.GetCid()].name, token_1[in.GetCid()].level.low, token_1[in.GetCid()].level.mid, token_1[in.GetCid()].level.high, token_1[in.GetCid()].res.partial_value, token_1[in.GetCid()].res.final_value)
	ids := make([]string, 0, len(token_1))
	j := 0
	for k := range token_1 {
		ids = append(ids, k)
		j++
	}
	fmt.Println("/n Following is the list of the currently existing token IDs : ", ids)
	var information string = "Token has been created successfully!"
	mutex_1.Unlock()
	return &pb.Token_Info{Cinfo: information}, nil

}

// Token Read Function to read the token requested by client.
func (s *tokenserver) ReadToken(ctx context.Context, in *pb.Read_ID) (*pb.Final, error) {
	mutex_1.Lock()
	var i uint64
	log.Printf("ReadToken")
	log.Printf("Recieved: %v", in.GetReadid())
	minhash := Hash(token_1[in.GetReadid()].name, token_1[in.GetReadid()].level.mid)
	min := token_1[in.GetReadid()].level.mid
	for i = token_1[in.GetReadid()].level.mid + 1; i < token_1[in.GetReadid()].level.high; i++ {
		if Hash(token_1[in.GetReadid()].name, i) < minhash {
			minhash = Hash(token_1[in.GetReadid()].name, i)
			min = i
		}
		if minhash > Hash(token_1[in.GetReadid()].name, token_1[in.GetReadid()].res.partial_value) {
			min = token_1[in.GetReadid()].res.partial_value
		}
		token_1[in.GetReadid()] = token_structure{id: in.GetReadid(), name: token_1[in.GetReadid()].name, res: values{partial_value: token_1[in.GetReadid()].res.partial_value, final_value: min}, level: levels{low: token_1[in.GetReadid()].level.low, mid: token_1[in.GetReadid()].level.mid, high: token_1[in.GetReadid()].level.high}, times: token_1[in.GetReadid()].times}
	}
	var finalread string = "Token has been successfully read!"
	fmt.Printf("The current token has the following data -> id: %v, name: %v, low : %d, mid: %d, high: %d, partial_value: %d, final_value: %d", token_1[in.GetReadid()].id, token_1[in.GetReadid()].name, token_1[in.GetReadid()].level.low, token_1[in.GetReadid()].level.mid, token_1[in.GetReadid()].level.high, token_1[in.GetReadid()].res.partial_value, token_1[in.GetReadid()].res.final_value)
	ids := make([]string, 0, len(token_1))
	// Printing all the updated values.
	j := 0
	for k := range token_1 {
		ids = append(ids, k)
		j++
	}
	fmt.Println("/n Following is the list of the currently existing token IDs : ", ids)
	mutex_1.Unlock()
	return &pb.Final{Final: min, Info: finalread, Ts: token_1[in.GetReadid()].times, Rid: token_1[in.GetReadid()].id, Rname: token_1[in.GetReadid()].name, Low: token_1[in.GetReadid()].level.low, Rhigh: token_1[in.GetReadid()].level.high, Rmid: token_1[in.GetReadid()].level.mid, Rts: token_1[in.GetReadid()].times}, nil
}

// The currently existing tokens in the system are updated through this function. In this, the input parameters are used for value calculation.

func (s *tokenserver) Token_Write(ctx context.Context, in *pb.Write_ID) (*pb.Partial, error) {
	mutex_1.Lock()
	var i uint64
	log.Printf("Write Token")
	log.Printf("The Token Recieved is : %v", in.GetWid())
	minhash := Hash(in.GetWname(), uint64(in.GetWlow()))
	min := in.GetWlow()

	if in.GetWlow() == in.GetWmid() {
		token_1[in.GetWid()] = token_structure{id: in.GetWid(), name: in.GetWname(), res: values{partial_value: min, final_value: 0}, level: levels{low: in.GetWlow(), mid: in.GetWmid(), high: in.GetWhigh()}, times: int64(time.Now().Minute())}
	} else {
		for i = in.GetWlow() + 1; i < in.GetWmid(); i++ {
			if Hash(in.GetWname(), i) < minhash {
				minhash = Hash(in.GetWname(), i)
				min = i
			}
			token_1[in.GetWid()] = token_structure{id: in.GetWid(), name: in.GetWname(), res: values{partial_value: min, final_value: 0}, level: levels{low: in.GetWlow(), mid: in.GetWmid(), high: in.GetWhigh()}, times: int64(time.Now().Minute())}
		}
	}
	fmt.Printf("The current token has the following data -> id: %v, name: %v, low : %d, mid: %d, high: %d, partial_value: %d, final_value: %d", token_1[in.GetWid()].id, token_1[in.GetWid()].name, token_1[in.GetWid()].level.low, token_1[in.GetWid()].level.mid, token_1[in.GetWid()].level.high, token_1[in.GetWid()].res.partial_value, token_1[in.GetWid()].res.final_value)
	ids := make([]string, 0, len(token_1))
	j := 0
	for k := range token_1 {
		ids = append(ids, k)
		j++
	}
	fmt.Println("/n Following is the list of the currently existing token IDs : ", ids) // These are the updated values/tokens.
	mutex_1.Unlock()
	return &pb.Partial{Pvalue: min, Pinfo: "Token has been successfully updated!"}, nil

}

// Token Drop Operation to delete the token requested by client.

func (s *tokenserver) DropToken(ctx context.Context, in *pb.Drop_ID) (*pb.Drop_Info, error) {
	log.Printf("Token to be dropped is : %v", in.GetDid())
	mutex_1.Lock()
	fmt.Printf("The current token has the following data -> id: %v, name: %v, low : %d, mid: %d, high: %d, partial_value: %d, final_value: %d", token_1[in.GetDid()].id, token_1[in.GetDid()].name, token_1[in.GetDid()].level.low, token_1[in.GetDid()].level.mid, token_1[in.GetDid()].level.high, token_1[in.GetDid()].res.partial_value, token_1[in.GetDid()].res.final_value)
	delete(token_1, in.Did)
	mutex_1.Unlock()
	ids := make([]string, 0, len(token_1))
	j := 0
	for k := range token_1 {
		ids = append(ids, k)
		j++
	}
	fmt.Println("/n Following is the list of the currently existing token IDs : ", ids)
	return &pb.Drop_Info{}, nil
}

// Main function
// Client-Server connection establishment, port accesses management and listen and error reporting.
func main() {

	port := flag.Int("port", 9000, "The server port") // Open port connection.
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTokenServiceServer(s, &tokenserver{}) // Connection establishment.
	log.Printf("server is listening at: %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve at the requested port : %v", err)
	}
}
